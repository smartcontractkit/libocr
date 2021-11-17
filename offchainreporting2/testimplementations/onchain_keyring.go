package testimplementations

import (
	"bytes"
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"io"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type OnchainKeyring struct {
	privateKey *ecdsa.PrivateKey
}

var _ types.OnchainKeyring = OnchainKeyring{}

// NewOnchainKeyring returns an OnchainKeyring with the given keys.
//
// In any real implementation (except maybe in a test helper), this function
// should take no arguments, and use crypto/rand.{Read,Int}. It should return a
// pointer, so we aren't copying secret material willy-nilly, and have a way to
// destroy the secrets. Any persistence to disk should be encrypted, as in the
// chainlink keystores.
func NewOnchainKeyring(rand io.Reader) *OnchainKeyring {
	secret, err := cryptorand.Int(rand, ethcrypto.S256().Params().N)
	if err != nil {
		panic(err)
	}
	x, y := secp256k1.S256().ScalarBaseMult(secret.Bytes())
	return &OnchainKeyring{
		&ecdsa.PrivateKey{
			ecdsa.PublicKey{
				secp256k1.S256(),
				x, y,
			},
			secret,
		},
	}
}

func (ok OnchainKeyring) MaxSignatureLength() int {
	return 65
}

func (ok OnchainKeyring) Sign(repctx types.ReportContext, report types.Report) (signature []byte, err error) {
	rawReportContext := evmutil.RawReportContext(repctx)
	sigData := ethcrypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	return ethcrypto.Sign(ethcrypto.Keccak256(sigData), ok.privateKey)
}

func (ok OnchainKeyring) PublicKey() types.OnchainPublicKey {
	address := ethcrypto.PubkeyToAddress(ok.privateKey.PublicKey)
	return address[:]
}

func (ok OnchainKeyring) Verify(pubkey types.OnchainPublicKey, repctx types.ReportContext, report types.Report, sig []byte) bool {
	rawReportContext := evmutil.RawReportContext(repctx)
	sigData := ethcrypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	hash := ethcrypto.Keccak256(sigData)
	authorPubkey, err := ethcrypto.SigToPub(hash, sig)

	// fmt.Printf("author pubkey %x\n", authorPubkey)
	if err != nil {
		// fmt.Printf("error while doing SigToPub: %v\n", err)
		return false
	}
	authorAddress := ethcrypto.PubkeyToAddress(*authorPubkey)
	// fmt.Printf("author address %x\n", authorAddress)
	// fmt.Printf("expected address %x\n", common.BytesToAddress(pubkey))
	return bytes.Equal(pubkey[:], authorAddress[:])
}
