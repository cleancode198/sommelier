package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterInterfaces registers the oracle proto files
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSubmitCorkRequest{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&AddManagedCellarIDsProposal{},
		&RemoveManagedCellarIDsProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
