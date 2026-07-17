package main

import (
	"flag"
	"log"
	"math"
	"os"
	"slices"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/keyvaluedatabase"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type keyValueDatabaseReadOnly interface {
	NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error)
	Close() error
}

type keyValueDatabaseNoWrite struct {
	kvdb keyValueDatabaseReadOnly
}

func (k keyValueDatabaseNoWrite) NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error) {
	return k.kvdb.NewReadTransaction()
}

// Sends reads/writes in a read-write transaction to /dev/null intentionally. We
// need this because otherwise the shim will try to initialize the kvdb schema
// version or migrate between schema versions, and we don't want to mutate the
// kvdb at all.
func (k keyValueDatabaseNoWrite) NewReadWriteTransaction() (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	txn, err := k.kvdb.NewReadTransaction()
	if err != nil {
		return nil, err
	}
	return keyValueDatabaseReadTransactionAsReadWriteNoWrite{txn}, nil
}

func (k keyValueDatabaseNoWrite) Close() error {
	return k.kvdb.Close()
}

type keyValueDatabaseReadTransactionAsReadWriteNoWrite struct {
	ocr3_1types.KeyValueDatabaseReadTransaction
}

func (keyValueDatabaseReadTransactionAsReadWriteNoWrite) Commit() error              { return nil }
func (keyValueDatabaseReadTransactionAsReadWriteNoWrite) Delete([]byte) error        { return nil }
func (keyValueDatabaseReadTransactionAsReadWriteNoWrite) Write([]byte, []byte) error { return nil }

func OpenPebbleKeyValueDatabaseNoWrite(path string) (ocr3_1types.KeyValueDatabase, error) {
	db, err := keyvaluedatabase.OpenPebbleKeyValueDatabaseReadOnlyForTooling(path)
	if err != nil {
		return nil, err
	}
	return keyValueDatabaseNoWrite{db}, nil
}

func main() {
	configDigestFlag := flag.String("config-digest", "", "OCR config digest hex string (32 bytes, with or without 0x prefix)")
	pathFlag := flag.String("path", "", "path to the kvdb directory (the <config-digest>.db directory itself)")
	contractConfigFlag := flag.String("contract-config", "", "path to JSON-encoded types.ContractConfig; see https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/types#ContractConfig for details. byte fields such as configDigest, signers, onchainConfig, and offchainConfig must be hex strings (with or without 0x)")
	skipExternalConsistencyChecksFlag := flag.Bool("skip-external-consistency-checks", false, "skip external consistency checks such as block signature validation, not recommended")
	logDespiteError := flag.Bool("log-despite-error", false, "log despite error, not recommended")
	flag.Parse()

	if flag.NArg() != 0 {
		log.Printf("error: unexpected positional arguments: %v", flag.Args())
		flag.Usage()
		os.Exit(2)
	}

	if *configDigestFlag == "" || *pathFlag == "" {
		flag.Usage()
		os.Exit(2)
	}

	if (*contractConfigFlag != "") == *skipExternalConsistencyChecksFlag {
		log.Printf("error: either -contract-config or -skip-external-consistency-checks must be passed, but not both")
		flag.Usage()
		os.Exit(2)
	}

	if *skipExternalConsistencyChecksFlag {
		log.Printf("WARNING: -skip-external-consistency-checks is set, so external consistency checks such as block signature validation will be skipped")
	}

	cdBytes, err := HexBytesFromString(*configDigestFlag)
	if err != nil {
		log.Fatalf("invalid config-digest: %v", err)
	}
	configDigest, err := types.BytesToConfigDigest(cdBytes)
	if err != nil {
		log.Fatalf("invalid config-digest: %v", err)
	}

	log.Printf("inspecting kvdb %q for config digest %s", *pathFlag, configDigest.Hex())

	var publicConfig ocr3_1config.PublicConfig
	if *contractConfigFlag != "" {
		ccBytes, err := os.ReadFile(*contractConfigFlag)
		if err != nil {
			log.Fatalf("error reading contract config file: %v", err)
		}
		contractConfig, err := parseContractConfigJSON(ccBytes)
		if err != nil {
			log.Fatalf("error parsing contract config JSON: %v", err)
		}
		if contractConfig.ConfigDigest != configDigest {
			log.Fatalf("error: config digest in contract config JSON (%s) does not match the CLI config-digest (%s)", contractConfig.ConfigDigest.Hex(), configDigest.Hex())
		}
		publicConfig, err = ocr3_1config.PublicConfigFromContractConfig(true, contractConfig)
		if err != nil {
			log.Fatalf("error constructing public config from contract config: %v", err)
		}
		if *skipExternalConsistencyChecksFlag { // defensive, as we enforce that exactly one of -contract-config and -skip-external-consistency-checks must be passed
			log.Printf("WARNING: -skip-external-consistency-checks is set, so external consistency checks will be skipped even though -contract-config was provided")
		} else {
			log.Printf("loaded contract config from %q; external consistency checks are enabled", *contractConfigFlag)
		}
	} else {
		log.Printf("WARNING: no contract config provided, external consistency checks are skipped. this is not recommended")
		publicConfig = ocr3_1config.PublicConfig{ConfigDigest: configDigest}
	}

	rawKVDB, err := OpenPebbleKeyValueDatabaseNoWrite(*pathFlag)
	if err != nil {
		log.Fatalf("error opening kvdb: %v", err)
	}
	defer rawKVDB.Close()

	semKVDB, err := shim.NewSemanticOCR3_1KeyValueDatabase(
		rawKVDB,
		ocr3_1types.ReportingPluginLimits{},
		publicConfig,
		DevnullLogger{},
		DevnullRegisterer{},
	)
	if err != nil {
		log.Fatalf("error opening semantic kvdb: %v", err)
	}
	defer semKVDB.Close()

	txn, err := semKVDB.NewReadTransactionUnchecked()
	if err != nil {
		log.Fatalf("error opening read transaction: %v", err)
	}
	defer txn.Discard()

	treeSyncStatus, err := txn.ReadTreeSyncStatus()
	if err != nil {
		log.Fatalf("error reading tree sync status: %v", err)
	}
	if treeSyncStatus.Phase != protocol.TreeSyncPhaseInactive {
		log.Fatalf("invalid kvdb backup: it was taken while tree sync was not inactive (phase=%s)", treeSyncStatus.Phase)
	}

	prevInstanceGenesisStateTransitionBlock, err := txn.ReadPrevInstanceGenesisStateTransitionBlock()
	if err != nil {
		log.Fatalf("error reading prev instance genesis state transition block: %v", err)
	}

	highestCommittedSeqNr, err := txn.ReadHighestCommittedSeqNr()
	if err != nil {
		log.Fatalf("error reading highest committed seq nr: %v", err)
	}
	log.Printf("highest committed seq nr: %d", highestCommittedSeqNr)

	versions, more, err := txn.ReadRootVersions(0, math.MaxInt)
	if err != nil {
		log.Fatalf("error reading tree roots: %v", err)
	}
	if more {
		log.Fatalf("error reading tree roots: unexpectedly truncated result")
	}

	slices.Sort(versions) // defensive

	if len(versions) == 0 {
		log.Printf("WARNING: found no tree roots")
	} else {
		log.Printf("found %d tree roots spanning versions %d..%d", len(versions), versions[0], versions[len(versions)-1])
	}

	var suitablePrevFields []ocr3_1config.PublicConfigPrevFields
	maxPrevSeqNr := uint64(0)
	corruptionIndicated := false

	for _, version := range versions {
		if prevInstanceGenesisStateTransitionBlock != nil && version <= prevInstanceGenesisStateTransitionBlock.SeqNr {
			// Special case: This root comes from a local tree sync due to
			// reconfiguration, so we should not expect to find an
			// AttestedStateTransitionBlock for it. It's fine to skip, because
			// we couldn't use it as a reconfiguration snapshot anyway.
			log.Printf("skipping tree root version %d because it is <= prev instance genesis state transition block seq nr %d", version, prevInstanceGenesisStateTransitionBlock.SeqNr)
			continue
		}

		prevFields, err := findAndCheckPrevFields(
			rawKVDB,
			highestCommittedSeqNr,
			version,
			publicConfig,
			*skipExternalConsistencyChecksFlag,
		)
		if err != nil {
			logfn := log.Fatalf
			if *logDespiteError {
				logfn = log.Printf
			}
			logfn("WARNING: we have indication that this backup is corrupted, integrity check failed for tree root version %d: %+v", version, err)
			corruptionIndicated = true
			continue
		}

		if version <= highestCommittedSeqNr {
			suitablePrevFields = append(suitablePrevFields, prevFields)
			if prevFields.PrevSeqNr > maxPrevSeqNr {
				maxPrevSeqNr = prevFields.PrevSeqNr
			}
		}
	}

	for _, prevFields := range suitablePrevFields {
		var suffix string
		if prevFields.PrevSeqNr == maxPrevSeqNr {
			suffix = " (👉 latest; this is the one you're most likely interested in)"
		}
		log.Printf("snapshot suitable for reconfiguration %+v%s", prevFields, suffix)
	}
	if len(suitablePrevFields) == 0 {
		log.Printf("WARNING: no intact snapshot suitable for reconfiguration was found")
	}

	if corruptionIndicated {
		log.Fatalf("This backup is likely corrupted. Do not use it for reconfiguration.")
	}
}
