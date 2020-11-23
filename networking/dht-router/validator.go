package dhtrouter


type AnnouncementValidator struct{}

const ValidatorNamespace = "peerinfo"


func (v AnnouncementValidator) Validate(key string, value []byte) error {
	
	peerId, err := dhtKeyToPeerId(key)
	if err != nil {
		return err
	}

	ann, err := deserializeSignedAnnouncement(value)
	if err != nil {
		return err
	}

	
	if !peerId.MatchesPublicKey(ann.PublicKey) {
		return InvalidDhtKey
	}

	
	return ann.verify()
}


func (v AnnouncementValidator) Select(_ string, values [][]byte) (int, error) {
	counter := announcementCounter{}
	latestRecord := 0
	for i := 0; i < len(values); i++ {
		ann, err := deserializeSignedAnnouncement(values[i])
		if err != nil {
			return 0, err
		}

		if ann.Counter.Gt(counter) {
			latestRecord = i
			counter = ann.Counter
		}
	}
	return latestRecord, nil
}
