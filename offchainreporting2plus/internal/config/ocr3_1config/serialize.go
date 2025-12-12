package ocr3_1config

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/curve25519"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// Serialized configs must be no larger than this (arbitrary bound, to prevent
// resource exhaustion attacks)
var maxSerializedOffchainConfigSize = 2_000_000

// offchainConfig contains the contents of the oracle Config objects
// which need to be serialized
type offchainConfig struct {
	DeltaProgress                    time.Duration
	DeltaResend                      *time.Duration
	DeltaInitial                     *time.Duration
	DeltaRound                       time.Duration
	DeltaGrace                       time.Duration
	DeltaReportsPlusPrecursorRequest *time.Duration
	DeltaStage                       time.Duration

	// state sync
	DeltaStateSyncSummaryInterval *time.Duration

	// block sync
	DeltaBlockSyncMinRequestToSameOracleInterval *time.Duration
	DeltaBlockSyncResponseTimeout                *time.Duration
	MaxBlocksPerBlockSyncResponse                *int
	MaxParallelRequestedBlocks                   *uint64

	// tree sync
	DeltaTreeSyncMinRequestToSameOracleInterval *time.Duration
	DeltaTreeSyncResponseTimeout                *time.Duration
	MaxTreeSyncChunkKeys                        *int
	MaxTreeSyncChunkKeysPlusValuesBytes         *int
	MaxParallelTreeSyncChunkFetches             *int

	// snapshotting
	SnapshotInterval               *uint64
	MaxHistoricalSnapshotsRetained *uint64

	// blobs
	DeltaBlobOfferMinRequestToSameOracleInterval *time.Duration
	DeltaBlobOfferResponseTimeout                *time.Duration
	DeltaBlobBroadcastGrace                      *time.Duration
	DeltaBlobChunkMinRequestToSameOracleInterval *time.Duration
	DeltaBlobChunkResponseTimeout                *time.Duration
	BlobChunkBytes                               *int

	RMax                                    uint64
	S                                       []int
	OffchainPublicKeys                      []types.OffchainPublicKey
	PeerIDs                                 []string
	ReportingPluginConfig                   []byte
	MaxDurationInitialization               time.Duration
	WarnDurationQuery                       time.Duration
	WarnDurationObservation                 time.Duration
	WarnDurationValidateObservation         time.Duration
	WarnDurationObservationQuorum           time.Duration
	WarnDurationStateTransition             time.Duration
	WarnDurationCommitted                   time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration
	PrevConfigDigest                        *types.ConfigDigest
	PrevSeqNr                               *uint64
	PrevHistoryDigest                       *types.HistoryDigest
	SharedSecretEncryptions                 config.SharedSecretEncryptions
}

func checkSize(serializedOffchainConfig []byte) error {
	if len(serializedOffchainConfig) <= maxSerializedOffchainConfigSize {
		return nil
	} else {
		return fmt.Errorf("OffchainConfig length is %d bytes which is greater than the max %d",
			len(serializedOffchainConfig),
			maxSerializedOffchainConfigSize,
		)
	}
}

// serialize returns a binary serialization of o
func (o offchainConfig) serialize() []byte {
	offchainConfigProto := enprotoOffchainConfig(o)
	rv, err := proto.Marshal(&offchainConfigProto)
	if err != nil {
		panic(err)
	}
	if err := checkSize(rv); err != nil {
		panic(err.Error())
	}
	return rv
}

func deserializeOffchainConfig(
	b []byte,
) (offchainConfig, error) {
	if err := checkSize(b); err != nil {
		return offchainConfig{}, err
	}

	offchainConfigPB := OffchainConfigProto{}
	if err := proto.Unmarshal(b, &offchainConfigPB); err != nil {
		return offchainConfig{}, fmt.Errorf("could not unmarshal ContractConfig.OffchainConfig protobuf: %w", err)
	}

	return deprotoOffchainConfig(&offchainConfigPB)
}

func deprotoOffchainConfig(
	offchainConfigProto *OffchainConfigProto,
) (offchainConfig, error) {
	S := make([]int, 0, len(offchainConfigProto.GetS()))
	for _, elem := range offchainConfigProto.GetS() {
		S = append(S, int(elem))
	}

	offchainPublicKeys := make([]types.OffchainPublicKey, 0, len(offchainConfigProto.GetOffchainPublicKeys()))
	for _, ocpkRaw := range offchainConfigProto.GetOffchainPublicKeys() {
		var ocpk types.OffchainPublicKey
		if len(ocpkRaw) != len(ocpk) {
			return offchainConfig{}, fmt.Errorf("invalid offchain public key: %x", ocpkRaw)
		}
		copy(ocpk[:], ocpkRaw)
		offchainPublicKeys = append(offchainPublicKeys, ocpk)
	}

	sharedSecretEncryptions, err := deprotoSharedSecretEncryptions(offchainConfigProto.GetSharedSecretEncryptions())
	if err != nil {
		return offchainConfig{}, fmt.Errorf("could not unmarshal shared protobuf: %w", err)
	}

	var prevConfigDigest *types.ConfigDigest
	if len(offchainConfigProto.PrevConfigDigest) != 0 {
		d, err := types.BytesToConfigDigest(offchainConfigProto.PrevConfigDigest)
		if err != nil {
			return offchainConfig{}, fmt.Errorf("invalid PrevConfigDigest: %w", err)
		}
		prevConfigDigest = &d
	}

	var prevHistoryDigest *types.HistoryDigest
	if len(offchainConfigProto.PrevHistoryDigest) != 0 {
		d, err := types.BytesToHistoryDigest(offchainConfigProto.PrevHistoryDigest)
		if err != nil {
			return offchainConfig{}, fmt.Errorf("invalid PrevHistoryDigest: %w", err)
		}
		prevHistoryDigest = &d
	}

	return offchainConfig{
		time.Duration(offchainConfigProto.GetDeltaProgressNanoseconds()),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaResendNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaInitialNanoseconds),
		time.Duration(offchainConfigProto.GetDeltaRoundNanoseconds()),
		time.Duration(offchainConfigProto.GetDeltaGraceNanoseconds()),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaReportsPlusPrecursorRequestNanoseconds),
		time.Duration(offchainConfigProto.GetDeltaStageNanoseconds()),

		// state sync
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaStateSyncSummaryIntervalNanoseconds),

		// block sync
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlockSyncMinRequestToSameOracleIntervalNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlockSyncResponseTimeoutNanoseconds),
		util.PointerIntegerCast[int](offchainConfigProto.MaxBlocksPerBlockSyncResponse),
		offchainConfigProto.MaxParallelRequestedBlocks,

		// tree sync
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaTreeSyncMinRequestToSameOracleIntervalNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaTreeSyncResponseTimeoutNanoseconds),
		util.PointerIntegerCast[int](offchainConfigProto.MaxTreeSyncChunkKeys),
		util.PointerIntegerCast[int](offchainConfigProto.MaxTreeSyncChunkKeysPlusValuesBytes),
		util.PointerIntegerCast[int](offchainConfigProto.MaxParallelTreeSyncChunkFetches),

		// snapshotting
		offchainConfigProto.SnapshotInterval,
		offchainConfigProto.MaxHistoricalSnapshotsRetained,

		// blobs
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlobOfferMinRequestToSameOracleIntervalNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlobOfferResponseTimeoutNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlobBroadcastGraceNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlobChunkMinRequestToSameOracleIntervalNanoseconds),
		util.PointerIntegerCast[time.Duration](offchainConfigProto.DeltaBlobChunkResponseTimeoutNanoseconds),
		util.PointerIntegerCast[int](offchainConfigProto.BlobChunkBytes),

		offchainConfigProto.GetRMax(),
		S,
		offchainPublicKeys,
		offchainConfigProto.GetPeerIds(),
		offchainConfigProto.GetReportingPluginConfig(),
		time.Duration(offchainConfigProto.GetMaxDurationInitializationNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationQueryNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationObservationNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationValidateObservationNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationObservationQuorumNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationStateTransitionNanoseconds()),
		time.Duration(offchainConfigProto.GetWarnDurationCommittedNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationShouldAcceptAttestedReportNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationShouldTransmitAcceptedReportNanoseconds()),
		prevConfigDigest,
		offchainConfigProto.PrevSeqNr,
		prevHistoryDigest,
		sharedSecretEncryptions,
	}, nil
}

func deprotoSharedSecretEncryptions(sharedSecretEncryptionsProto *SharedSecretEncryptionsProto) (config.SharedSecretEncryptions, error) {
	var diffieHellmanPoint [curve25519.PointSize]byte
	if len(diffieHellmanPoint) != len(sharedSecretEncryptionsProto.GetDiffieHellmanPoint()) {
		return config.SharedSecretEncryptions{}, fmt.Errorf("DiffieHellmanPoint has wrong length. Expected %v bytes, got %v bytes", len(diffieHellmanPoint), len(sharedSecretEncryptionsProto.GetDiffieHellmanPoint()))
	}
	copy(diffieHellmanPoint[:], sharedSecretEncryptionsProto.GetDiffieHellmanPoint())

	var sharedSecretHash common.Hash
	if len(sharedSecretHash) != len(sharedSecretEncryptionsProto.GetSharedSecretHash()) {
		return config.SharedSecretEncryptions{}, fmt.Errorf("sharedSecretHash has wrong length. Expected %v bytes, got %v bytes", len(sharedSecretHash), len(sharedSecretEncryptionsProto.GetSharedSecretHash()))
	}
	copy(sharedSecretHash[:], sharedSecretEncryptionsProto.GetSharedSecretHash())

	encryptions := make([]config.EncryptedSharedSecret, 0, len(sharedSecretEncryptionsProto.GetEncryptions()))
	for i, encryptionRaw := range sharedSecretEncryptionsProto.GetEncryptions() {
		var encryption config.EncryptedSharedSecret
		if len(encryption) != len(encryptionRaw) {
			return config.SharedSecretEncryptions{}, fmt.Errorf("Encryptions[%v] has wrong length. Expected %v bytes, got %v bytes", i, len(encryption), len(encryptionRaw))
		}
		copy(encryption[:], encryptionRaw)
		encryptions = append(encryptions, encryption)
	}

	return config.SharedSecretEncryptions{
		diffieHellmanPoint,
		sharedSecretHash,
		encryptions,
	}, nil
}

func enprotoOffchainConfig(o offchainConfig) OffchainConfigProto {
	s := make([]uint32, len(o.S))
	for i, d := range o.S {
		s[i] = uint32(d)
	}
	offchainPublicKeys := make([][]byte, 0, len(o.OffchainPublicKeys))
	for _, k := range o.OffchainPublicKeys {
		offchainPublicKeys = append(offchainPublicKeys, k[:])
	}
	sharedSecretEncryptions := enprotoSharedSecretEncryptions(o.SharedSecretEncryptions)

	var prevConfigDigestBytes []byte
	if o.PrevConfigDigest != nil {
		prevConfigDigestBytes = o.PrevConfigDigest[:]
	}
	var prevHistoryDigestBytes []byte
	if o.PrevHistoryDigest != nil {
		prevHistoryDigestBytes = o.PrevHistoryDigest[:]
	}

	return OffchainConfigProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint64(o.DeltaProgress),
		util.PointerIntegerCast[uint64](o.DeltaResend),
		util.PointerIntegerCast[uint64](o.DeltaInitial),
		uint64(o.DeltaRound),
		uint64(o.DeltaGrace),
		util.PointerIntegerCast[uint64](o.DeltaReportsPlusPrecursorRequest),
		uint64(o.DeltaStage),
		util.PointerIntegerCast[uint64](o.DeltaStateSyncSummaryInterval),
		util.PointerIntegerCast[uint64](o.DeltaBlockSyncMinRequestToSameOracleInterval),
		util.PointerIntegerCast[uint64](o.DeltaBlockSyncResponseTimeout),
		util.PointerIntegerCast[uint32](o.MaxBlocksPerBlockSyncResponse),
		util.PointerIntegerCast[uint64](o.MaxParallelRequestedBlocks),
		util.PointerIntegerCast[uint64](o.DeltaTreeSyncMinRequestToSameOracleInterval),
		util.PointerIntegerCast[uint64](o.DeltaTreeSyncResponseTimeout),
		util.PointerIntegerCast[uint32](o.MaxTreeSyncChunkKeys),
		util.PointerIntegerCast[uint32](o.MaxTreeSyncChunkKeysPlusValuesBytes),
		util.PointerIntegerCast[uint32](o.MaxParallelTreeSyncChunkFetches),
		util.PointerIntegerCast[uint64](o.SnapshotInterval),
		util.PointerIntegerCast[uint64](o.MaxHistoricalSnapshotsRetained),
		util.PointerIntegerCast[uint64](o.DeltaBlobOfferMinRequestToSameOracleInterval),
		util.PointerIntegerCast[uint64](o.DeltaBlobOfferResponseTimeout),
		util.PointerIntegerCast[uint64](o.DeltaBlobBroadcastGrace),
		util.PointerIntegerCast[uint64](o.DeltaBlobChunkMinRequestToSameOracleInterval),
		util.PointerIntegerCast[uint64](o.DeltaBlobChunkResponseTimeout),
		util.PointerIntegerCast[uint32](o.BlobChunkBytes),
		o.RMax,
		s,
		offchainPublicKeys,
		o.PeerIDs,
		o.ReportingPluginConfig,
		uint64(o.MaxDurationInitialization),
		uint64(o.WarnDurationQuery),
		uint64(o.WarnDurationObservation),
		uint64(o.WarnDurationValidateObservation),
		uint64(o.WarnDurationObservationQuorum),
		uint64(o.WarnDurationStateTransition),
		uint64(o.WarnDurationCommitted),
		uint64(o.MaxDurationShouldAcceptAttestedReport),
		uint64(o.MaxDurationShouldTransmitAcceptedReport),
		prevConfigDigestBytes,
		o.PrevSeqNr,
		prevHistoryDigestBytes,
		&sharedSecretEncryptions,
	}
}

func enprotoSharedSecretEncryptions(e config.SharedSecretEncryptions) SharedSecretEncryptionsProto {
	encs := make([][]byte, 0, len(e.Encryptions))
	for _, enc := range e.Encryptions {
		encs = append(encs, enc[:])
	}
	return SharedSecretEncryptionsProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		e.DiffieHellmanPoint[:],
		e.SharedSecretHash[:],
		encs,
	}
}
