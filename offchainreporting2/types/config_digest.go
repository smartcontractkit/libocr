package types

import ocr2plustypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

type ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefix

const (
	ConfigDigestPrefixEVM        ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixEVM
	ConfigDigestPrefixTerra      ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixTerra
	ConfigDigestPrefixSolana     ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixSolana
	ConfigDigestPrefixStarknet   ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixStarknet
	ConfigDigestPrefixMercuryV02 ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixMercuryV02
	ConfigDigestPrefixOCR1       ConfigDigestPrefix = ocr2plustypes.ConfigDigestPrefixOCR1
)

type ConfigDigest = ocr2plustypes.ConfigDigest

func BytesToConfigDigest(b []byte) (ConfigDigest, error) {
	return ocr2plustypes.BytesToConfigDigest(b)
}

type OffchainConfigDigester = ocr2plustypes.OffchainConfigDigester
