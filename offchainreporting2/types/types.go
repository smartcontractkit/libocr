package types

import ocr2plustypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

// import (
// 	"bytes"
// 	"context"
// 	"crypto/ed25519"
// 	"time"

// 	"github.com/smartcontractkit/libocr/commontypes"
// 	"golang.org/x/crypto/curve25519"
// )

type BinaryNetworkEndpointLimits = ocr2plustypes.BinaryNetworkEndpointLimits

type BinaryNetworkEndpointFactory = ocr2plustypes.BinaryNetworkEndpointFactory

type BootstrapperFactory = ocr2plustypes.BootstrapperFactory

type Query = ocr2plustypes.Query

type Observation = ocr2plustypes.Observation

type AttributedObservation = ocr2plustypes.AttributedObservation

type ReportTimestamp = ocr2plustypes.ReportTimestamp

type ReportContext = ocr2plustypes.ReportContext
type Report = ocr2plustypes.Report

type AttributedOnchainSignature = ocr2plustypes.AttributedOnchainSignature

type ReportingPluginFactory = ocr2plustypes.ReportingPluginFactory

type ReportingPluginConfig = ocr2plustypes.ReportingPluginConfig

type ReportingPlugin = ocr2plustypes.ReportingPlugin

const (
	MaxMaxQueryLength       = ocr2plustypes.MaxMaxQueryLength
	MaxMaxObservationLength = ocr2plustypes.MaxMaxObservationLength
	MaxMaxReportLength      = ocr2plustypes.MaxMaxReportLength
)

type ReportingPluginLimits = ocr2plustypes.ReportingPluginLimits

type ReportingPluginInfo = ocr2plustypes.ReportingPluginInfo

type Account = ocr2plustypes.Account

type ContractTransmitter = ocr2plustypes.ContractTransmitter

type ContractConfigTracker = ocr2plustypes.ContractConfigTracker

type ContractConfig = ocr2plustypes.ContractConfig

type OffchainPublicKey = ocr2plustypes.OffchainPublicKey

type OnchainPublicKey = ocr2plustypes.OnchainPublicKey

type ConfigEncryptionPublicKey = ocr2plustypes.ConfigEncryptionPublicKey

type OffchainKeyring = ocr2plustypes.OffchainKeyring

type OnchainKeyring = ocr2plustypes.OnchainKeyring
