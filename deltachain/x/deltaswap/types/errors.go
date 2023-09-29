package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/deltaswap module sentinel errors
var (
	ErrPhylaxSetNotFound                     = sdkerrors.Register(ModuleName, 1101, "phylax set not found")
	ErrSignaturesInvalid                     = sdkerrors.Register(ModuleName, 1102, "invalid signatures on VAA")
	ErrNoQuorum                              = sdkerrors.Register(ModuleName, 1103, "no quorum on VAA")
	ErrUnknownGovernanceModule               = sdkerrors.Register(ModuleName, 1105, "invalid governance module")
	ErrNoConfig                              = sdkerrors.Register(ModuleName, 1106, "config not set")
	ErrInvalidGovernanceEmitter              = sdkerrors.Register(ModuleName, 1107, "invalid governance emitter")
	ErrUnknownGovernanceAction               = sdkerrors.Register(ModuleName, 1108, "unknown governance action")
	ErrGovernanceHeaderTooShort              = sdkerrors.Register(ModuleName, 1109, "governance header too short")
	ErrInvalidGovernanceTargetChain          = sdkerrors.Register(ModuleName, 1110, "governance target chain does not match")
	ErrInvalidGovernancePayloadLength        = sdkerrors.Register(ModuleName, 1111, "governance payload has incorrect length")
	ErrPhylaxSetNotSequential                = sdkerrors.Register(ModuleName, 1112, "phylax set updates must be submitted sequentially")
	ErrVAAAlreadyExecuted                    = sdkerrors.Register(ModuleName, 1113, "VAA was already executed")
	ErrPhylaxSignatureMismatch               = sdkerrors.Register(ModuleName, 1114, "phylax signature mismatch")
	ErrSignerMismatch                        = sdkerrors.Register(ModuleName, 1115, "transaction signer doesn't match validator key")
	ErrPhylaxNotFound                        = sdkerrors.Register(ModuleName, 1116, "phylax not found in phylax set")
	ErrConsensusSetUndefined                 = sdkerrors.Register(ModuleName, 1117, "no consensus set defined")
	ErrPhylaxSetExpired                      = sdkerrors.Register(ModuleName, 1118, "phylax set expired")
	ErrNewPhylaxSetHasExpiry                 = sdkerrors.Register(ModuleName, 1119, "new phylax set should not have expiry time")
	ErrDuplicatePhylaxAddress                = sdkerrors.Register(ModuleName, 1120, "phylax set has duplicate addresses")
	ErrSignerAlreadyRegistered               = sdkerrors.Register(ModuleName, 1121, "transaction signer already registered as a phylax validator")
	ErrConsensusSetNotUpdatable              = sdkerrors.Register(ModuleName, 1122, "cannot make changes to active consensus phylax set")
	ErrInvalidHash                           = sdkerrors.Register(ModuleName, 1123, "could not verify the hash in governance action")
	ErrPhylaxIndexOutOfBounds                = sdkerrors.Register(ModuleName, 1124, "phylax index out of bounds for the phylax set")
	ErrInvalidAllowlistContractAddr          = sdkerrors.Register(ModuleName, 1125, "contract addresses in the wasm allowlist msg and vaa do not match")
	ErrInvalidAllowlistCodeId                = sdkerrors.Register(ModuleName, 1126, "code ids in the wasm allowlist msg and vaa do not match")
	ErrInvalidIbcComposabilityMwContractAddr = sdkerrors.Register(ModuleName, 1127, "contract addresses in the set ibc composability mw contract and vaa do not match")
)
