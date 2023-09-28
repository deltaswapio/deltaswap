package types

import (
	"fmt"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypePhylaxSetUpdate            string = "PhylaxSetUpdate"
	ProposalTypeGovernanceDeltaswapMessage string = "GovernanceDeltaswapMessage"
)

func init() {
	gov.RegisterProposalType(ProposalTypePhylaxSetUpdate)
	gov.RegisterProposalTypeCodec(&PhylaxSetUpdateProposal{}, "deltaswap/PhylaxSetUpdate")
	gov.RegisterProposalType(ProposalTypeGovernanceDeltaswapMessage)
	gov.RegisterProposalTypeCodec(&GovernanceDeltaswapMessageProposal{}, "deltaswap/GovernanceDeltaswapMessage")
}

func NewPhylaxSetUpdateProposal(title, description string, guardianSet PhylaxSet) *PhylaxSetUpdateProposal {
	return &PhylaxSetUpdateProposal{
		Title:        title,
		Description:  description,
		NewPhylaxSet: guardianSet,
	}
}

func (sup *PhylaxSetUpdateProposal) ProposalRoute() string { return RouterKey }
func (sup *PhylaxSetUpdateProposal) ProposalType() string  { return ProposalTypePhylaxSetUpdate }
func (sup *PhylaxSetUpdateProposal) ValidateBasic() error {
	if err := sup.NewPhylaxSet.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(sup)
}

func (sup *PhylaxSetUpdateProposal) String() string {
	return fmt.Sprintf(`Phylax Set Upgrade Proposal: 
  Title:       %s
  Description: %s
  PhylaxSet: %s`, sup.Title, sup.Description, sup.NewPhylaxSet.String())
}

func NewGovernanceDeltaswapMessageProposal(title, description string, action uint8, targetChain uint16, module []byte, payload []byte) *GovernanceDeltaswapMessageProposal {
	return &GovernanceDeltaswapMessageProposal{
		Title:       title,
		Description: description,
		Module:      module,
		Action:      uint32(action),
		TargetChain: uint32(targetChain),
		Payload:     payload,
	}
}

func (sup *GovernanceDeltaswapMessageProposal) ProposalRoute() string { return RouterKey }
func (sup *GovernanceDeltaswapMessageProposal) ProposalType() string {
	return ProposalTypeGovernanceDeltaswapMessage
}
func (sup *GovernanceDeltaswapMessageProposal) ValidateBasic() error {
	if len(sup.Module) != 32 {
		return fmt.Errorf("invalid module length: %d != 32", len(sup.Module))
	}
	return gov.ValidateAbstract(sup)
}

func (sup *GovernanceDeltaswapMessageProposal) String() string {
	return fmt.Sprintf(`Governance Deltaswap Message Proposal: 
  Title:       %s
  Description: %s
  Module: %x
  TargetChain: %d
  Payload: %x`, sup.Title, sup.Description, sup.Module, sup.TargetChain, sup.Payload)
}
