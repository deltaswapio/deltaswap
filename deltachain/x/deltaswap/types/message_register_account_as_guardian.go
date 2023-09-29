package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterAccountAsPhylax = "register_account_as_phylax"

var _ sdk.Msg = &MsgRegisterAccountAsPhylax{}

func NewMsgRegisterAccountAsPhylax(signer string, signature []byte) *MsgRegisterAccountAsPhylax {
	return &MsgRegisterAccountAsPhylax{
		Signer:    signer,
		Signature: signature,
	}
}

func (msg *MsgRegisterAccountAsPhylax) Route() string {
	return RouterKey
}

func (msg *MsgRegisterAccountAsPhylax) Type() string {
	return TypeMsgRegisterAccountAsPhylax
}

func (msg *MsgRegisterAccountAsPhylax) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgRegisterAccountAsPhylax) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterAccountAsPhylax) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer address (%s)", err)
	}
	return nil
}
