package types

import "time"

const EnableDangerousDevelopmentMode = "enable dangerous development mode"



type LocalConfig struct {
	
	
	
	
	BlockchainTimeout time.Duration

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ContractConfigConfirmations uint16

	
	
	
	ContractConfigTrackerPollInterval time.Duration

	
	
	ContractConfigTrackerSubscribeInterval time.Duration

	ContractTransmitterTransmitTimeout time.Duration

	
	
	
	DatabaseTimeout time.Duration

	
	
	
	DataSourceTimeout time.Duration

	
	
	DevelopmentMode string
}
