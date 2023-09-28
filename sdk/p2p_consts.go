package sdk

var (
	HeartbeatMessagePrefix = []byte("heartbeat|")

	SignedObservationRequestPrefix_old = []byte("signed_observation_request|")
	SignedObservationRequestPrefix     = []byte("signed_observation_request_000000|")
	SignedWormchainAddressPrefix       = []byte("signed_deltachain_address_00000000|")
)
