package cli

import (
	"io/ioutil"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/peggyjv/sommelier/v4/x/cork/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseAddManagedCellarsProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "deposit": "1000stake"
}
`)

	proposal := types.AddManagedCellarIDsProposalWithDeposit{}
	contents, err := ioutil.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Marshaler.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "1000stake", proposal.Deposit)
}

func TestParseRemoveManagedCellarsProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "deposit": "1000stake"
}
`)

	proposal := types.RemoveManagedCellarIDsProposalWithDeposit{}
	contents, err := ioutil.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Marshaler.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "1000stake", proposal.Deposit)
}
