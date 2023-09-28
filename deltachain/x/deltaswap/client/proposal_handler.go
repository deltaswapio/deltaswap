package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/deltaswapio/deltachain/x/deltaswap/client/cli"
	"github.com/deltaswapio/deltachain/x/deltaswap/client/rest"
)

var PhylaxSetUpdateProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitPhylaxSetUpdateProposal, rest.ProposalPhylaxSetUpdateRESTHandler)
var WormholeGovernanceMessageProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitWormholeGovernanceMessageProposal, rest.ProposalWormholeGovernanceMessageRESTHandler)
