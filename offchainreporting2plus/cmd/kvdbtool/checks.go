package main

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/keyvaluedatabase"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type staticKeyValueDatabaseFactory struct {
	configDigest types.ConfigDigest
	database     ocr3_1types.KeyValueDatabase
}

var _ ocr3_1types.KeyValueDatabaseFactory = staticKeyValueDatabaseFactory{}

func (s staticKeyValueDatabaseFactory) NewKeyValueDatabase(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	if configDigest != s.configDigest {
		return nil, fmt.Errorf("unexpected config digest %s, expected %s", configDigest.Hex(), s.configDigest.Hex())
	}
	return nonClosingKeyValueDatabase{s.database}, nil
}

func (s staticKeyValueDatabaseFactory) NewKeyValueDatabaseIfExists(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return s.NewKeyValueDatabase(configDigest)
}

type nonClosingKeyValueDatabase struct {
	ocr3_1types.KeyValueDatabase
}

func (n nonClosingKeyValueDatabase) Close() error {
	return nil
}

type ExtendedPrevFields struct {
	ocr3_1config.PublicConfigPrevFields
	PrevStateRootDigest protocol.StateRootDigest
}

func findConsistentExtendedPrevFields(
	prevRawKVDB ocr3_1types.KeyValueDatabase,
	highestCommittedSeqNr uint64,
	treeRootVersion uint64,
	prevPublicConfig ocr3_1config.PublicConfig,
	skipBlockSignatureCheck bool,
) (ExtendedPrevFields, error) {
	prevSemKVDB, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		prevRawKVDB,
		ocr3_1types.ReportingPluginLimits{},
		ocr3_1config.PublicConfig{},
		DevnullLogger{},
		DevnullRegisterer{},
	)
	if err != nil {
		return ExtendedPrevFields{}, fmt.Errorf("failed to open source semantic kvdb: %w", err)
	}
	defer prevSemKVDB.Close()

	prevTxn, err := prevSemKVDB.NewReadTransactionUnchecked()
	if err != nil {
		return ExtendedPrevFields{}, fmt.Errorf("failed to create source read transaction: %w", err)
	}
	defer prevTxn.Discard()

	rootDigest, err := jmt.ReadRootDigest(prevTxn, prevTxn, treeRootVersion)
	if err != nil {
		return ExtendedPrevFields{}, fmt.Errorf("failed to read root digest for treeRootVersion %d: %w", treeRootVersion, err)
	}

	var blockSeqNr uint64
	if treeRootVersion > highestCommittedSeqNr {
		blockSeqNr = highestCommittedSeqNr
	} else {
		blockSeqNr = treeRootVersion
	}

	// Sanity check: Logic in the managed package agrees with our derived blockSeqNr <-> treeRootVersion mapping.
	{
		foundSnapshotVersion, foundStateRootDigest, err := managed.FindRoundedUpSnapshot(
			loghelper.MakeRootLoggerWithContext(DevnullLogger{}),
			DevnullRegisterer{},
			prevRawKVDB,
			prevPublicConfig.ConfigDigest,
			blockSeqNr,
		)
		if err != nil {
			return ExtendedPrevFields{}, fmt.Errorf("managed.FindRoundedUpSnapshot would have errored for PrevSeqNr/BlockSeqNr %v", blockSeqNr)
		}

		if foundSnapshotVersion != treeRootVersion {
			return ExtendedPrevFields{}, fmt.Errorf("managed.FindRoundedUpSnapshot would have returned snapshot version "+
				"%d for PrevSeqNr/BlockSeqNr %v, but we are considering tree root version %d", foundSnapshotVersion, blockSeqNr, treeRootVersion)
		}
		if foundStateRootDigest != rootDigest {
			return ExtendedPrevFields{}, fmt.Errorf("managed.FindRoundedUpSnapshot would have returned state root digest %x "+
				"for PrevSeqNr/BlockSeqNr %v, but we are considering tree root version %d with state root digest %x", foundStateRootDigest, blockSeqNr, treeRootVersion, rootDigest)
		}
	}

	prevAstb, err := prevTxn.ReadAttestedStateTransitionBlock(blockSeqNr)
	if err != nil {
		return ExtendedPrevFields{}, fmt.Errorf("failed to read attested state transition block %d: %w", blockSeqNr, err)
	}
	if prevAstb.StateTransitionBlock.BlockSeqNr != blockSeqNr {
		return ExtendedPrevFields{}, fmt.Errorf("attested state transition block %d does not exist", blockSeqNr)
	}
	if !skipBlockSignatureCheck {
		if err := prevAstb.Verify(prevPublicConfig); err != nil {
			return ExtendedPrevFields{}, fmt.Errorf("block signatures: did not verify: %w", err)
		}
	}
	if prevAstb.StateTransitionBlock.StateRootDigest != rootDigest {
		return ExtendedPrevFields{}, fmt.Errorf("state root digest mismatch, expected %x but attested block has %x", rootDigest, prevAstb.StateTransitionBlock.StateRootDigest)
	}

	certifiedCommit := prevAstb.ToCertifiedCommit(prevPublicConfig.ConfigDigest)
	prevHistoryDigest := certifiedCommit.HistoryDigest(prevPublicConfig.ConfigDigest)
	return ExtendedPrevFields{
		ocr3_1config.PublicConfigPrevFields{
			prevPublicConfig.ConfigDigest,
			blockSeqNr,
			prevHistoryDigest,
		},
		prevAstb.StateTransitionBlock.StateRootDigest,
	}, nil
}

func checkTreeIntegrityWithTrialCopy(
	prevRawKVDB ocr3_1types.KeyValueDatabase,
	treeRootVersion uint64,
	prevPublicConfig ocr3_1config.PublicConfig,
	extendedPrevFields ExtendedPrevFields,
) error {
	tempDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("kvdbtool-tree-sync-%d-*", treeRootVersion))
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Permits copying of chunks from any seq nr including the highest committed seq nr
	prevPublicConfig.SnapshotInterval = util.PointerTo(treeRootVersion)

	prevSemKVDB, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		prevRawKVDB,
		ocr3_1types.ReportingPluginLimits{},
		prevPublicConfig,
		DevnullLogger{},
		DevnullRegisterer{},
	)
	if err != nil {
		return fmt.Errorf("failed to open source semantic kvdb: %w", err)
	}
	defer prevSemKVDB.Close()

	// Matches lib/offchainreporting2plus/internal/managed/managed_ocr3_1_oracle.go
	prevTxn, err := prevSemKVDB.NewReadTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create source read transaction: %w", err)
	}
	defer prevTxn.Discard()

	nextConfigDigest := types.ConfigDigest{}
	nextPublicConfig := prevPublicConfig
	nextPublicConfig.PrevConfigDigest = util.PointerTo(extendedPrevFields.PrevConfigDigest)
	nextPublicConfig.PrevSeqNr = util.PointerTo(extendedPrevFields.PrevSeqNr)
	nextPublicConfig.PrevHistoryDigest = util.PointerTo(extendedPrevFields.PrevHistoryDigest)
	nextPublicConfig.ConfigDigest = nextConfigDigest

	nextRawKVDB, err := keyvaluedatabase.NewPebbleKeyValueDatabaseFactory(tempDir).NewKeyValueDatabase(nextConfigDigest)
	if err != nil {
		return fmt.Errorf("failed to create target raw kvdb: %w", err)
	}
	defer nextRawKVDB.Close()

	nextSemKVDB, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		nextRawKVDB,
		ocr3_1types.ReportingPluginLimits{},
		nextPublicConfig,
		DevnullLogger{},
		DevnullRegisterer{},
	)
	if err != nil {
		return fmt.Errorf("failed to open target semantic kvdb: %w", err)
	}
	defer nextSemKVDB.Close()

	err = managed.CopyAllTreeSyncChunksFromPrevInstance(
		context.Background(),
		loghelper.MakeRootLoggerWithContext(DevnullLogger{}),
		nextPublicConfig,
		extendedPrevFields.PrevSeqNr,
		extendedPrevFields.PrevStateRootDigest,
		prevTxn,
		nextSemKVDB,
	)
	if err != nil {
		return fmt.Errorf("failed to trial-copy from prev instance: %w", err)
	}
	return nil
}

func findAndCheckPrevFields(
	prevRawKVDB ocr3_1types.KeyValueDatabase,
	highestCommittedSeqNr uint64,
	treeRootVersion uint64,
	prevPublicConfig ocr3_1config.PublicConfig,
	skipBlockSignatureCheck bool,
) (ocr3_1config.PublicConfigPrevFields, error) {
	extendedPrevFields, err := findConsistentExtendedPrevFields(
		prevRawKVDB,
		highestCommittedSeqNr,
		treeRootVersion,
		prevPublicConfig,
		skipBlockSignatureCheck,
	)
	if err != nil {
		return ocr3_1config.PublicConfigPrevFields{}, fmt.Errorf("failed to find consistent prev fields: %w", err)
	}

	err = checkTreeIntegrityWithTrialCopy(
		prevRawKVDB,
		treeRootVersion,
		prevPublicConfig,
		extendedPrevFields,
	)
	if err != nil {
		return ocr3_1config.PublicConfigPrevFields{}, fmt.Errorf("failed to check tree integrity with trial copy: %w", err)
	}

	return extendedPrevFields.PublicConfigPrevFields, nil
}
