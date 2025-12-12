package managed

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/internal/util"
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
				return fmt.Errorf("ManagedOCR3_1Oracle: invalid ReportingPluginInfo"), false
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

			if prev, ok := sharedConfig.PublicConfig.GetPrevFields(); ok {
				err := tryCopyFromPrevInstance(
					ctx,
					sharedConfig.PublicConfig,
					logger,
					&devNullRegisterer{},
					registerer,
					reportingPluginInfo.Limits,
					keyValueDatabaseFactory,
					prev.PrevConfigDigest,
					prev.PrevSeqNr,
					sharedConfig.ConfigDigest,
				)
				if err != nil {
					return fmt.Errorf("ManagedOCR3_1Oracle: error during tryCopyFromPrevInstance: %w", err), true
				}
			}

			keyValueDatabase, err := keyValueDatabaseFactory.NewKeyValueDatabase(sharedConfig.ConfigDigest)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewKeyValueDatabase: %w", err), true
			}
			defer loghelper.CloseLogError(
				keyValueDatabase,
				logger,
				"ManagedOCR3_1Oracle: error during keyValueDatabase.Close()",
			)
			keyValueDatabaseWithMetrics := shim.NewKeyValueDatabaseWithMetrics(keyValueDatabase, registerer, logger)
			defer loghelper.CloseLogError(
				keyValueDatabaseWithMetrics,
				logger,
				"ManagedOCR3_1Oracle: error during keyValueDatabaseWithMetrics.Close()",
			)
			semanticOCR3_1KeyValueDatabase, err := shim.NewSemanticOCR3_1KeyValueDatabase(keyValueDatabaseWithMetrics, reportingPluginInfo.Limits, sharedConfig.PublicConfig, logger, registerer)
			if err != nil {
				return fmt.Errorf("ManagedOCR3_1Oracle: error during NewSemanticOCR3_1KeyValueDatabase: %w", err), true
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

func tryCopyFromPrevInstance(
	ctx context.Context,
	publicConfig ocr3_1config.PublicConfig,
	logger loghelper.LoggerWithContext,
	prevRegisterer prometheus.Registerer,
	nextRegisterer prometheus.Registerer,
	reportingPluginLimits ocr3_1types.ReportingPluginLimits,
	keyValueDatabaseFactory ocr3_1types.KeyValueDatabaseFactory,
	prevConfigDigest types.ConfigDigest,
	prevSeqNr uint64,
	nextConfigDigest types.ConfigDigest,
) error {
	nextKeyValueDatabase, err := keyValueDatabaseFactory.NewKeyValueDatabase(nextConfigDigest)
	if err != nil {
		return fmt.Errorf("error during NewKeyValueDatabase: %w", err)
	}
	defer nextKeyValueDatabase.Close()

	nextSemanticKeyValueDatabase, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		nextKeyValueDatabase,
		reportingPluginLimits,
		publicConfig,
		logger,
		nextRegisterer,
	)
	if err != nil {
		return fmt.Errorf("error during NewSemanticOCR3_1KeyValueDatabase: %w", err)
	}
	defer nextSemanticKeyValueDatabase.Close()

	nextHighestCommittedSeqNr, nextTreeSyncPhase, err := readHighestCommittedSeqNrAndTreeSyncPhase(nextSemanticKeyValueDatabase)
	if err != nil {
		return fmt.Errorf("error during readHighestCommittedSeqNrAndTreeSyncPhase: %w", err)
	}

	if nextHighestCommittedSeqNr >= prevSeqNr {
		logger.Info("⚙️ tryCopyFromPrevInstance: next instance kvdb is already synced, nothing to do", commontypes.LogFields{
			"nextHighestCommittedSeqNr": nextHighestCommittedSeqNr,
			"prevSeqNr":                 prevSeqNr,
		})
		return nil
	}

	err = ensureKeyValueDatabaseIsPristine(nextHighestCommittedSeqNr, nextTreeSyncPhase)
	if err != nil {
		logger.Warn("⚙️ tryCopyFromPrevInstance: next instance kvdb is not pristine, we won't copy from prev instance."+
			" If you would like to copy from prev instance, you should delete the next instance kvdb and restart.", commontypes.LogFields{
			"error": err,
		})
		return nil
	}

	prevKeyValueDatabase, err := keyValueDatabaseFactory.NewKeyValueDatabaseIfExists(prevConfigDigest)
	if err != nil {
		if errors.Is(err, ocr3_1types.ErrKeyValueDatabaseDoesNotExist) {
			logger.Warn("⚙️ tryCopyFromPrevInstance: prev instance kvdb does not exist, nothing to do", commontypes.LogFields{
				"prevConfigDigest": prevConfigDigest,
			})
			return nil
		}
		return fmt.Errorf("error during prev keyValueDatabaseFactory.NewKeyValueDatabaseIfExists: %w", err)
	}
	defer prevKeyValueDatabase.Close()

	// We do not have access to the public config of the prev instance, so we
	// reconstruct a good-enough public config for opening the semantic kvdb of
	// the prev instance and reading the tree sync chunks.
	prevPublicConfig := ocr3_1config.PublicConfig{
		// We do not know the snapshot interval of the prev instance, but
		// PrevSeqNr is guaranteed to be a snapshot sequence number in the prev
		// instance, and we will call ReadTreeSyncChunk with it.
		SnapshotInterval: util.PointerTo(uint64(1)),
		ConfigDigest:     prevConfigDigest,
		// Set all tree sync chunking parameters to be the same as the next
		// instance. They are guaranteed to fit a maximally sized key-value due
		// to checkPublicConfigParameters.
		MaxTreeSyncChunkKeys:                publicConfig.MaxTreeSyncChunkKeys,
		MaxTreeSyncChunkKeysPlusValuesBytes: publicConfig.MaxTreeSyncChunkKeysPlusValuesBytes,
	}
	prevSemanticKeyValueDatabase, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		prevKeyValueDatabase,
		reportingPluginLimits,
		prevPublicConfig,
		logger,
		prevRegisterer,
	)
	if err != nil {
		return fmt.Errorf("error during NewSemanticOCR3_1KeyValueDatabase: %w", err)
	}
	defer prevSemanticKeyValueDatabase.Close()

	prevTxn, err := prevSemanticKeyValueDatabase.NewReadTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("error during NewReadTransactionUnchecked: %w", err)
	}
	defer prevTxn.Discard()

	genesisBlock, err := findGenesisBlockInPrevKeyValueDatabase(logger, publicConfig, prevTxn, prevSeqNr)
	if err != nil {
		return fmt.Errorf("error during findGenesisBlockInPrevKeyValueDatabase: %w", err)
	}
	if genesisBlock == nil {
		logger.Warn("⚙️ tryCopyFromPrevInstance: prev seqnr block does not exist, nothing to do", commontypes.LogFields{
			"prevSeqNr":        prevSeqNr,
			"prevConfigDigest": prevConfigDigest,
		})
		return nil
	}

	startIndex := jmt.Digest{}
	for {
		logger.Info("⚙️ tryCopyFromPrevInstance: copying chunk from prev instance to next instance", commontypes.LogFields{
			"startIndex": fmt.Sprintf("%x", startIndex),
			"prevSeqNr":  prevSeqNr,
		})
		endInclIndex, done, err := copyChunkFromPrevInstance(
			publicConfig,
			prevSeqNr,
			genesisBlock.StateRootDigest,
			startIndex,
			prevTxn,
			nextSemanticKeyValueDatabase,
		)
		if err != nil {
			return fmt.Errorf("error during copyChunkFromPrevInstance: %w", err)
		}

		if done {
			break
		}

		var ok bool
		startIndex, ok = jmt.IncrementDigest(endInclIndex)
		if !ok {
			return fmt.Errorf("failed to increment endInclIndex even though we are not done copying chunks")
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	nextTxn, err := nextSemanticKeyValueDatabase.NewSerializedReadWriteTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("error during NewSerializedReadWriteTransactionUnchecked: %w", err)
	}
	defer nextTxn.Discard()

	err = nextTxn.WritePrevInstanceGenesisStateTransitionBlock(*genesisBlock)
	if err != nil {
		return fmt.Errorf("error during WritePrevInstanceGenesisStateTransitionBlock: %w", err)
	}

	err = nextTxn.WriteHighestCommittedSeqNr(prevSeqNr)
	if err != nil {
		return fmt.Errorf("error during WriteHighestCommittedSeqNr: %w", err)
	}

	err = nextTxn.WriteLowestPersistedSeqNr(prevSeqNr)
	if err != nil {
		return fmt.Errorf("error during WriteLowestPersistedSeqNr: %w", err)
	}

	err = nextTxn.Commit()
	if err != nil {
		return fmt.Errorf("error during Commit: %w", err)
	}

	return nil
}

func readHighestCommittedSeqNrAndTreeSyncPhase(semanticKeyValueDatabase protocol.KeyValueDatabase) (uint64, protocol.TreeSyncPhase, error) {
	txn, err := semanticKeyValueDatabase.NewReadTransactionUnchecked()
	if err != nil {
		return 0, protocol.TreeSyncPhaseInactive, fmt.Errorf("error during NewReadTransactionUnchecked: %w", err)
	}
	defer txn.Discard()

	highestCommittedSeqNr, err := txn.ReadHighestCommittedSeqNr()
	if err != nil {
		return 0, protocol.TreeSyncPhaseInactive, fmt.Errorf("error during ReadHighestCommittedSeqNr: %w", err)
	}

	treeSyncStatus, err := txn.ReadTreeSyncStatus()
	if err != nil {
		return 0, protocol.TreeSyncPhaseInactive, fmt.Errorf("error during ReadTreeSyncStatus: %w", err)
	}

	return highestCommittedSeqNr, treeSyncStatus.Phase, nil
}

func ensureKeyValueDatabaseIsPristine(highestCommittedSeqNr uint64, treeSyncPhase protocol.TreeSyncPhase) error {
	if highestCommittedSeqNr != 0 {
		return fmt.Errorf("highest committed sequence number is not 0")
	}
	if treeSyncPhase != protocol.TreeSyncPhaseInactive {
		return fmt.Errorf("tree sync phase is not inactive, it is %v", treeSyncPhase)
	}
	return nil
}

func findGenesisBlockInPrevKeyValueDatabase(
	logger commontypes.Logger,
	publicConfig ocr3_1config.PublicConfig,
	prevTxn protocol.KeyValueDatabaseReadTransaction,
	prevSeqNr uint64,
) (*protocol.GenesisStateTransitionBlock, error) {
	prev, ok := publicConfig.GetPrevFields()
	if !ok {
		return nil, fmt.Errorf("previous instance is not specified in PublicConfig")
	}

	// ensure prevSeqNr is in range [lowest, highest]
	prevLowestPersistedSeqNr, err := prevTxn.ReadLowestPersistedSeqNr()
	if err != nil {
		return nil, fmt.Errorf("error during ReadLowestPersistedSeqNr: %w", err)
	}
	if prevSeqNr < prevLowestPersistedSeqNr {
		logger.Warn("⚙️ findGenesisBlockInPrevKeyValueDatabase: prevSeqNr is less than prevLowestPersistedSeqNr", commontypes.LogFields{
			"prevSeqNr":                prevSeqNr,
			"prevLowestPersistedSeqNr": prevLowestPersistedSeqNr,
		})
		return nil, nil
	}

	prevHighestCommittedSeqNr, err := prevTxn.ReadHighestCommittedSeqNr()
	if err != nil {
		return nil, fmt.Errorf("error during ReadHighestCommittedSeqNr: %w", err)
	}
	if prevSeqNr > prevHighestCommittedSeqNr {
		logger.Warn("⚙️ findGenesisBlockInPrevKeyValueDatabase: prevSeqNr is greater than prevHighestCommittedSeqNr", commontypes.LogFields{
			"prevSeqNr":                 prevSeqNr,
			"prevHighestCommittedSeqNr": prevHighestCommittedSeqNr,
		})
		return nil, nil
	}

	treeSyncStatus, err := prevTxn.ReadTreeSyncStatus()
	if err != nil {
		return nil, fmt.Errorf("error during ReadTreeSyncStatus: %w", err)
	}
	if treeSyncStatus.Phase != protocol.TreeSyncPhaseInactive {
		logger.Warn("⚙️ findGenesisBlockInPrevKeyValueDatabase: tree sync phase is not inactive, can't read from prev kvdb", commontypes.LogFields{
			"treeSyncStatus": treeSyncStatus,
		})
		return nil, nil
	}

	// check that history digest matches the block at prevSeqNr
	prevAstb, err := prevTxn.ReadAttestedStateTransitionBlock(prevSeqNr)
	if err != nil {
		return nil, fmt.Errorf("error during ReadAttestedStateTransitionBlock: %w", err)
	}
	if prevAstb.StateTransitionBlock.BlockSeqNr != prevSeqNr {
		// block does not exist
		return nil, nil
	}

	gstb := protocol.AttestedToGenesisStateTransitionBlock(prev.PrevConfigDigest, prevAstb)
	err = protocol.VerifyGenesisStateTransitionBlockFromPrevInstance(publicConfig, gstb)
	if err != nil {
		return nil, fmt.Errorf("error during VerifyGenesisStateTransitionBlockFromPrevInstance: %w", err)
	}
	return &gstb, nil
}

func copyChunkFromPrevInstance(
	publicConfig ocr3_1config.PublicConfig, // non essential
	prevSeqNr uint64,
	prevStateRootDigest protocol.StateRootDigest,
	startIndex jmt.Digest,
	prevTxn protocol.KeyValueDatabaseReadTransaction,
	nextSemanticKeyValueDatabase protocol.KeyValueDatabase,
) (endInclIndex jmt.Digest, done bool, err error) {
	nextTxn, err := nextSemanticKeyValueDatabase.NewSerializedReadWriteTransactionUnchecked()
	if err != nil {
		return jmt.Digest{}, false, fmt.Errorf("error during NewSerializedReadWriteTransactionUnchecked: %w", err)
	}
	defer nextTxn.Discard()

	requestEndInclIndex := jmt.MaxDigest
	endInclIndex, boundingLeaves, keyValues, err := prevTxn.ReadTreeSyncChunk(
		prevSeqNr,
		startIndex,
		requestEndInclIndex,
		publicConfig.GetMaxTreeSyncChunkKeysPlusValuesBytes(),
	)
	if err != nil {
		return jmt.Digest{}, false, fmt.Errorf("error during ReadTreeSyncChunk: %w", err)
	}

	writeAndVerifyTreeSyncChunkResult, err := nextTxn.VerifyAndWriteTreeSyncChunk(
		prevStateRootDigest,
		prevSeqNr,
		startIndex,
		requestEndInclIndex,
		endInclIndex,
		boundingLeaves,
		keyValues,
	)
	if err != nil {
		return jmt.Digest{}, false, fmt.Errorf("error during VerifyAndWriteTreeSyncChunk: %w", err)
	}

	switch writeAndVerifyTreeSyncChunkResult {
	case protocol.VerifyAndWriteTreeSyncChunkResultOkComplete:
		done = true
	case protocol.VerifyAndWriteTreeSyncChunkResultOkNeedMore:
		done = false
	case protocol.VerifyAndWriteTreeSyncChunkResultByzantine:
		return jmt.Digest{}, false, fmt.Errorf("byzantine error during VerifyAndWriteTreeSyncChunk")
	case protocol.VerifyAndWriteTreeSyncChunkResultUnrelatedError:
		return jmt.Digest{}, false, fmt.Errorf("unrelated error during VerifyAndWriteTreeSyncChunk")
	}

	err = nextTxn.Commit()
	if err != nil {
		return jmt.Digest{}, false, fmt.Errorf("error during Commit: %w", err)
	}

	return endInclIndex, done, nil
}

type devNullRegisterer struct{}

var _ prometheus.Registerer = (*devNullRegisterer)(nil)

func (d *devNullRegisterer) Register(collector prometheus.Collector) error {
	return nil
}

func (d *devNullRegisterer) MustRegister(collectors ...prometheus.Collector) {
}

func (d *devNullRegisterer) Unregister(collector prometheus.Collector) bool {
	return false
}
