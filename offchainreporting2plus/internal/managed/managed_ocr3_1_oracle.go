package managed

import (
	"context"
	"errors"
	"fmt"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/RoSpaceDev/libocr/internal/util"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/shim"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultIncomingMessageBufferSize     = 10
	defaultOutgoingMessageBufferSize     = 10
	lowPriorityIncomingMessageBufferSize = 10
	lowPriorityOutgoingMessageBufferSize = 10
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
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount(ctx)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error getting FromAccount: %w", err), true
			}

			sharedConfig, oid, err := ocr3config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
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

			maxDurationInitialization := util.NilCoalesce(sharedConfig.MaxDurationInitialization, localConfig.DefaultMaxDurationInitialization)
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

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(initCtx, ocr3types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
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
				defaultIncomingMessageBufferSize,
				defaultOutgoingMessageBufferSize,
			}
			lowPriorityConfig := types.BinaryNetworkEndpoint2Config{
				lowPriorityLimits,
				lowPriorityIncomingMessageBufferSize,
				lowPriorityOutgoingMessageBufferSize,
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
			semanticOCR3_1KeyValueDatabase := shim.NewSemanticOCR3_1KeyValueDatabase(keyValueDatabase, reportingPluginInfo.Limits, logger, metricsRegisterer)

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
	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= ocr3_1types.MaxMaxQueryLength) {
		err = errors.Join(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, ocr3_1types.MaxMaxQueryLength))
	}
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= ocr3_1types.MaxMaxObservationLength) {
		err = errors.Join(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, ocr3_1types.MaxMaxObservationLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= ocr3_1types.MaxMaxReportLength) {
		err = errors.Join(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, ocr3_1types.MaxMaxReportLength))
	}
	if !(0 <= limits.MaxReportsPlusPrecursorLength && limits.MaxReportsPlusPrecursorLength <= ocr3_1types.MaxMaxReportsPlusPrecursorLength) {
		err = errors.Join(err, fmt.Errorf("MaxReportInfoLength (%v) out of range. Should be between 0 and %v", limits.MaxReportsPlusPrecursorLength, ocr3_1types.MaxMaxReportsPlusPrecursorLength))
	}
	if !(0 <= limits.MaxReportCount && limits.MaxReportCount <= ocr3_1types.MaxMaxReportCount) {
		err = errors.Join(err, fmt.Errorf("MaxReportCount (%v) out of range. Should be between 0 and %v", limits.MaxReportCount, ocr3_1types.MaxMaxReportCount))
	}

	if !(0 <= limits.MaxKeyValueModifiedKeysPlusValuesLength && limits.MaxKeyValueModifiedKeysPlusValuesLength <= ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesLength) {
		err = errors.Join(err, fmt.Errorf("MaxKeyValueModifiedKeysPlusValuesLength (%v) out of range. Should be between 0 and %v", limits.MaxKeyValueModifiedKeysPlusValuesLength, ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesLength))
	}
	if !(0 <= limits.MaxBlobPayloadLength && limits.MaxBlobPayloadLength <= ocr3_1types.MaxMaxBlobPayloadLength) {
		err = errors.Join(err, fmt.Errorf("MaxBlobPayloadLength (%v) out of range. Should be between 0 and %v", limits.MaxBlobPayloadLength, ocr3_1types.MaxMaxBlobPayloadLength))
	}
	return err
}
