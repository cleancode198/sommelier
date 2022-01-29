package keeper

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/x/cork/types"
	"github.com/stretchr/testify/require"
)

type ValCellar struct {
	Val  sdktypes.ValAddress
	Cork types.Cork
}

type VoteCalculatorTestCase struct {
	title        string
	description  string
	CellarID     common.Address
	ValCellars   []ValCellar
	WinningVotes []types.Cork
}

var (
	vallAddrA, _ = sdktypes.ValAddressFromHex("24ep6yqkhpwnfdrrapu6fzmjp3xrpsgca11ab1e")
	//vallAddrB, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hb01dface")
	//vallAddrC, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hdeadbeef")
	//vallAddrD, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hf005ba11")

	exampleAddrA = common.HexToAddress("0xc0ffee254729296a45a3885639AC7E10F9d54979")
)

func TestGetWinningVotes(t *testing.T) {
	testCases := []VoteCalculatorTestCase{
		{title: "Single voter",
			description: "Check that a single voter returns it's vote",
			CellarID:    exampleAddrA,
			ValCellars: []ValCellar{
				{Val: vallAddrA,
					Cork: types.Cork{
						Address: exampleAddrA.String(),
						Body:    []byte{33},
					},
				},
			},
			WinningVotes: []types.Cork{
				{
					Address: exampleAddrA.String(),
					Body:    []byte{33},
				},
			},
		},
	}

	for _, test := range testCases {
		input := CreateTestEnv(t)
		ctx := input.Context

		for _, vc := range test.ValCellars {
			commit := types.Cork{
				Address: exampleAddrA.String(),
				Body:    []byte{33},
			}

			input.corkKeeper.SetCork(ctx, vc.Val, commit)
		}

		winningVotes := input.corkKeeper.GetWinningVotes(ctx, sdktypes.MustNewDecFromStr("0.66"))
		require.Lenf(t, winningVotes, 1, "require that winning votes contains only one cellar")
	}
}