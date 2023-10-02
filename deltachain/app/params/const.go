package params

const (
	// Name defines the application name of Deltachain.
	Name = "udelta"

	// BondDenom defines the native staking token denomination.
	BondDenom = "udelta"

	// DisplayDenom defines the name, symbol, and display value of the worm token.
	DisplayDenom = "WORM"

	// DefaultGasLimit - set to the same value as cosmos-sdk flags.DefaultGasLimit
	// this value is currently only used in tests.
	DefaultGasLimit = 200000
)
