package types

import (
	"fmt"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypePhylaxSetUpdate           string = "PhylaxSetUpdate"
	ProposalTypeGovernanceWormholeMessage string = "GovernanceWormholeMessage"
)

func init() {
	gov.RegisterProposalType(ProposalTypePhylaxSetUpdate)
	gov.RegisterProposalTypeCodec(&PhylaxSetUpdateProposal{}, "wormhole/PhylaxSetUpdate")
	gov.RegisterProposalType(ProposalTypeGovernanceWormholeMessage)
	gov.RegisterProposalTypeCodec(&GovernanceWormholeMessageProposal{}, "wormhole/GovernanceWormholeMessage")
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

func NewGovernanceWormholeMessageProposal(title, description string, action uint8, targetChain uint16, module []byte, payload []byte) *GovernanceWormholeMessageProposal {
	return &GovernanceWormholeMessageProposal{
		Title:       title,
		Description: description,
		Module:      module,
		Action:      uint32(action),
		TargetChain: uint32(targetChain),
		Payload:     payload,
	}
}

func (sup *GovernanceWormholeMessageProposal) ProposalRoute() string { return RouterKey }
func (sup *GovernanceWormholeMessageProposal) ProposalType() string {
	return ProposalTypeGovernanceWormholeMessage
}
func (sup *GovernanceWormholeMessageProposal) ValidateBasic() error {
	if len(sup.Module) != 32 {
		return fmt.Errorf("invalid module length: %d != 32", len(sup.Module))
	}
	return gov.ValidateAbstract(sup)
}

func (sup *GovernanceWormholeMessageProposal) String() string {
	return fmt.Sprintf(`Governance Wormhole Message Proposal: 
  Title:       %s
  Description: %s
  Module: %x
  TargetChain: %d
  Payload: %x`, sup.Title, sup.Description, sup.Module, sup.TargetChain, sup.Payload)
}
