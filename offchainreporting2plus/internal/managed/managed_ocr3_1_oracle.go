package managed

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunManagedOCR3_1Oracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from types.BinaryNetworkEndpoint2 to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOCR3_1Oracle[RI any](
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database ocr3_1types.Database,
	keyValueDatabaseFactory ocr3_1types.KeyValueDatabaseFactory,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	messageNetEndpointFactory types.BinaryNetworkEndpoint2Factory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPluginFactory ocr3_1types.ReportingPluginFactory[RI],
) {
	subs := subprocesses.Subprocesses{}
	defer subs.Wait()

	var chTelemetrySend chan<- *serialization.TelemetryWrapper
	{
		chTelemetry := make(chan *serialization.TelemetryWrapper, 100)
		chTelemetrySend = chTelemetry
		subs.Go(func() {
			forwardTelemetry(ctx, logger, monitoringEndpoint, chTelemetry)
		})
	}

	metricsRegistererWrapper := metricshelper.NewPrometheusRegistererWrapper(metricsRegisterer, logger)

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, logger loghelper.LoggerWithContext, contractConfig types.ContractConfig) (err error, retry bool) {
			skipInsaneForProductionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount(ctx)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error getting FromAccount: %w", err), true
			}

			sharedConfig, oid, err := ocr3_1config.SharedConfigFromContractConfig(
				skipInsaneForProductionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				messageNetEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error while decoding ContractConfig: %w", err), false
			}

			registerer := prometheus.WrapRegistererWith(
				prometheus.Labels{
					// disambiguate different protocol instances by configDigest
					"config_digest": sharedConfig.ConfigDigest.String(),
					// disambiguate different oracle instances by offchainPublicKey
					"offchain_public_key": fmt.Sprintf("%x", offchainKeyring.OffchainPublicKey()),
				},
				metricsRegistererWrapper,
			)

			// Run with new config
			peerIDs := []string{}
			for _, identity := range sharedConfig.OracleIdentities {
				peerIDs = append(peerIDs, identity.PeerID)
			}

			childLogger := logger.MakeChild(commontypes.LogFields{
				"oid": oid,
			})

			blobEndpointWrapper := protocol.BlobEndpointWrapper{}

			maxDurationInitialization := sharedConfig.MaxDurationInitialization
			initCtx, initCancel := context.WithTimeout(ctx, maxDurationInitialization)
			defer initCancel()

			ins := loghelper.NewIfNotStopped(
				maxDurationInitialization+common.ReportingPluginTimeoutWarningGracePeriod,
				func() {
					logger.Error("ManagedOCR3_1Oracle: ReportingPluginFactory.NewReportingPlugin is taking too long", commontypes.LogFields{
						"maxDuration": maxDurationInitialization,
					})
				},
			)

			reportingPlugin, reportingPluginInfo_, err := reportingPluginFactory.NewReportingPlugin(initCtx, ocr3types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.WarnDurationQuery,
				sharedConfig.WarnDurationObservation,
				sharedConfig.MaxDurationShouldAcceptAttestedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			}, &blobEndpointWrapper)

			ins.Stop()

			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewReportingPlugin(): %w", err), true
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOCR3_1Oracle: error during reportingPlugin.Close()",
			)

			var reportingPluginInfo ocr3_1types.ReportingPluginInfo1
			switch rpi := reportingPluginInfo_.(type) {
			case ocr3_1types.ReportingPluginInfo1:
				reportingPluginInfo = rpi
			}

			if err := validateOCR3_1ReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOCR3_1Oracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return fmt.Errorf("ManagedOCR3_1Oracle: invalid MercuryPluginInfo"), false
			}

			maxSigLen := onchainKeyring.MaxSignatureLength()
			defaultLims, lowPriorityLimits, serializedLengthLimits, err := limits.OCR3_1Limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, maxSigLen)
			if err != nil {
				logger.Error("ManagedOCR3_1Oracle: error during limits", commontypes.LogFields{
					"error":                 err,
					"publicConfig":          sharedConfig.PublicConfig,
					"reportingPluginLimits": reportingPluginInfo.Limits,
					"maxSigLen":             maxSigLen,
				})
				return fmt.Errorf("ManagedOCR3_1Oracle: error during limits"), false
			}

			defaultPriorityConfig := types.BinaryNetworkEndpoint2Config{
				defaultLims,
				nil,
				nil,
			}
			lowPriorityConfig := types.BinaryNetworkEndpoint2Config{
				lowPriorityLimits,
				nil,
				nil,
			}

			binNetEndpoint, err := messageNetEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				defaultPriorityConfig,
				lowPriorityConfig,
			)
			if err != nil {
				logger.Error("ManagedOCR3_1Oracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewEndpoint"), true
			}
			defer loghelper.CloseLogError(
				binNetEndpoint,
				logger,
				"ManagedOCR3_1Oracle: error during BinaryNetworkEndpoint2.Close()",
			)

			netEndpoint := shim.NewOCR3_1SerializingEndpoint[RI](
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				maxSigLen,
				childLogger,
				registerer,
				reportingPluginInfo.Limits,
				sharedConfig.PublicConfig,
				serializedLengthLimits,
			)
			err = netEndpoint.Start()
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during netEndpoint.Start(): %w", err), true
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOCR3_1Oracle: error during netEndpoint.Close()",
			)

			keyValueDatabase, err := keyValueDatabaseFactory.NewKeyValueDatabase(sharedConfig.ConfigDigest)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewKeyValueDatabase: %w", err), false
			}
			defer loghelper.CloseLogError(
				keyValueDatabase,
				logger,
				"ManagedOCR3_1Oracle: error during keyValueDatabase.Close()",
			)
			keyValueDatabaseWithMetrics := shim.NewKeyValueDatabaseWithMetrics(keyValueDatabase, metricsRegisterer, logger)
			defer loghelper.CloseLogError(
				keyValueDatabaseWithMetrics,
				logger,
				"ManagedOCR3_1Oracle: error during keyValueDatabaseWithMetrics.Close()",
			)
			semanticOCR3_1KeyValueDatabase, err := shim.NewSemanticOCR3_1KeyValueDatabase(keyValueDatabaseWithMetrics, reportingPluginInfo.Limits, sharedConfig.PublicConfig, logger, metricsRegisterer)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewSemanticOCR3_1KeyValueDatabase: %w", err), false
			}
			defer loghelper.CloseLogError(
				semanticOCR3_1KeyValueDatabase,
				logger,
				"ManagedOCR3_1Oracle: error during semanticOCR3_1KeyValueDatabase.Close()",
			)

			protocol.RunOracle[RI](
				ctx,
				&blobEndpointWrapper,
				sharedConfig,
				contractTransmitter,
				&shim.SerializingOCR3_1Database{database},
				oid,
				semanticOCR3_1KeyValueDatabase,
				reportingPluginInfo.Limits,
				localConfig,
				childLogger,
				registerer,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckOCR3_1ReportingPlugin[RI]{reportingPlugin, reportingPluginInfo.Limits},
				shim.NewOCR3_1TelemetrySender(chTelemetrySend, childLogger),
			)

			return nil, false
		},
		localConfig,
		logger,
		offchainConfigDigester,
		defaultRetryParams(),
	)
}

func validateOCR3_1ReportingPluginLimits(limits ocr3_1types.ReportingPluginLimits) error {
	var err error
	if !(0 <= limits.MaxQueryBytes && limits.MaxQueryBytes <= ocr3_1types.MaxMaxQueryBytes) {
		err = errors.Join(err, fmt.Errorf("MaxQueryBytes (%v) out of range. Should be between 0 and %v", limits.MaxQueryBytes, ocr3_1types.MaxMaxQueryBytes))
	}
	if !(0 <= limits.MaxObservationBytes && limits.MaxObservationBytes <= ocr3_1types.MaxMaxObservationBytes) {
		err = errors.Join(err, fmt.Errorf("MaxObservationBytes (%v) out of range. Should be between 0 and %v", limits.MaxObservationBytes, ocr3_1types.MaxMaxObservationBytes))
	}
	if !(0 <= limits.MaxReportBytes && limits.MaxReportBytes <= ocr3_1types.MaxMaxReportBytes) {
		err = errors.Join(err, fmt.Errorf("MaxReportBytes (%v) out of range. Should be between 0 and %v", limits.MaxReportBytes, ocr3_1types.MaxMaxReportBytes))
	}
	if !(0 <= limits.MaxReportsPlusPrecursorBytes && limits.MaxReportsPlusPrecursorBytes <= ocr3_1types.MaxMaxReportsPlusPrecursorBytes) {
		err = errors.Join(err, fmt.Errorf("MaxReportsPlusPrecursorBytes (%v) out of range. Should be between 0 and %v", limits.MaxReportsPlusPrecursorBytes, ocr3_1types.MaxMaxReportsPlusPrecursorBytes))
	}
	if !(0 <= limits.MaxReportCount && limits.MaxReportCount <= ocr3_1types.MaxMaxReportCount) {
		err = errors.Join(err, fmt.Errorf("MaxReportCount (%v) out of range. Should be between 0 and %v", limits.MaxReportCount, ocr3_1types.MaxMaxReportCount))
	}

	if !(0 <= limits.MaxKeyValueModifiedKeys && limits.MaxKeyValueModifiedKeys <= ocr3_1types.MaxMaxKeyValueModifiedKeys) {
		err = errors.Join(err, fmt.Errorf("MaxKeyValueModifiedKeys (%v) out of range. Should be between 0 and %v", limits.MaxKeyValueModifiedKeys, ocr3_1types.MaxMaxKeyValueModifiedKeys))
	}
	if !(0 <= limits.MaxKeyValueModifiedKeysPlusValuesBytes && limits.MaxKeyValueModifiedKeysPlusValuesBytes <= ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes) {
		err = errors.Join(err, fmt.Errorf("MaxKeyValueModifiedKeysPlusValuesBytes (%v) out of range. Should be between 0 and %v", limits.MaxKeyValueModifiedKeysPlusValuesBytes, ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes))
	}
	if !(0 <= limits.MaxBlobPayloadBytes && limits.MaxBlobPayloadBytes <= ocr3_1types.MaxMaxBlobPayloadBytes) {
		err = errors.Join(err, fmt.Errorf("MaxBlobPayloadBytes (%v) out of range. Should be between 0 and %v", limits.MaxBlobPayloadBytes, ocr3_1types.MaxMaxBlobPayloadBytes))
	}
	return err
}
