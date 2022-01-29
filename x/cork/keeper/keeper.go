package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/x/cork/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
	gravityKeeper types.GravityKeeper
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper, gravityKeeper types.GravityKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		stakingKeeper: stakingKeeper,
		gravityKeeper: gravityKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

///////////////////////
// MsgSubmitCork //
///////////////////////

// SetCork sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetCork(ctx sdk.Context, val sdk.ValAddress, cork types.Cork) {
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetCorkForValidatorAddressKey(val, common.HexToAddress(cork.Address)), bz)
}

// GetCork gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetCork(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetCorkForValidatorAddressKey(val, cel))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteCork deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteCork(ctx sdk.Context, val sdk.ValAddress, cel common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetCorkForValidatorAddressKey(val, cel))
}

// HasCorkForContract gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasCorkForContract(ctx sdk.Context, val sdk.ValAddress, contract common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetCorkForValidatorAddressKey(val, contract))
}

// HasCork gets the existence of any commit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasCork(ctx sdk.Context, val sdk.ValAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetCorkValidatorKeyPrefix(val))
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	return iter.Valid()
}

// IterateCorks iterates over all votes in the store
func (k Keeper) IterateCorks(ctx sdk.Context, handler func(val sdk.ValAddress, cel common.Address, cork types.Cork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CorkForAddressKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.CorkForAddressKeyPrefix}))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		var cork types.Cork
		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if handler(val, cel, cork) {
			break
		}
	}
}

// IterateCorkAddresses iterates over all addresses who have committed corks
func (k Keeper) IterateCorkAddresses(ctx sdk.Context, handler func(addr common.Address) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CorkForAddressKeyPrefix})
	defer iter.Close()

	seenAddresses := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.CorkForAddressKeyPrefix}))
		keyPair.Next(20)
		address := common.BytesToAddress(keyPair.Bytes())

		// add seen address to set. if already in set, don't return to consumer
		if !seenAddresses.Add(address) {
			continue
		}

		if handler(address) {
			break
		}
	}
}

// IterateAddressCorks iterates over all corks for an address
func (k Keeper) IterateAddressCorks(ctx sdk.Context, handler func(addr common.Address) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CorkForAddressKeyPrefix})
	defer iter.Close()

	seenAddresses := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.CorkForAddressKeyPrefix}))
		keyPair.Next(20)
		address := common.BytesToAddress(keyPair.Bytes())

		// add seen address to set. if already in set, don't return to consumer
		if !seenAddresses.Add(address) {
			continue
		}

		if handler(address) {
			break
		}
	}
}

//////////////////
// CommitPeriod //
//////////////////

// SetCommitPeriodStart sets the current vote period start height
func (k Keeper) SetCommitPeriodStart(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.CommitPeriodStartKey}, sdk.Uint64ToBigEndian(uint64(height)))
}

// GetCommitPeriodStart returns the vote period start height
func (k Keeper) GetCommitPeriodStart(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.CommitPeriodStartKey})
	if len(bz) == 0 {
		return 0, false
	}

	return int64(sdk.BigEndianToUint64(bz)), true
}

// HasCommitPeriodStart returns true if the vote period start has been set
func (k Keeper) HasCommitPeriodStart(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has([]byte{types.CommitPeriodStartKey})
}

////////////
// Params //
////////////

// GetParamSet returns the vote period from the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

/////////////////////////
// Invalidation Nonces //
/////////////////////////

func (k Keeper) GetLatestInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LatestInvalidationNonceKey})
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestInvalidationNonce(ctx sdk.Context, invalidationNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(invalidationNonce))
}

func (k Keeper) IncrementInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextNonce := k.GetLatestInvalidationNonce(ctx) + 1
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(nextNonce))
	return nextNonce
}

///////////
// Votes //
///////////

func (k Keeper) GetWinningVotes(ctx sdk.Context, threshold sdk.Dec) (winningVotes []types.Cork) {

	var corks []types.Cork
	var corkPowers []int64

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateCorks(ctx, func(val sdk.ValAddress, addr common.Address, cork types.Cork) (stop bool) {
		validator := k.stakingKeeper.Validator(ctx, val)
		validatorPower := validator.GetConsensusPower(k.stakingKeeper.PowerReduction(ctx))

		found := false
		for i, rv := range corks {
			if rv.Equals(cork) {
				corkPowers[i] += validatorPower

				found = true
				break
			}
		}

		if !found {
			corks = append(corks, cork)
			corkPowers = append(corkPowers, validatorPower)
		}

		k.DeleteCork(ctx, val, addr)

		return false
	})

	var winningCorks []types.Cork

	for i, power := range corkPowers {
		quorumReached := sdk.NewDec(power).Quo(totalPower.ToDec()).GT(threshold)
		if quorumReached {
			winningCorks = append(winningCorks, corks[i])
		}
	}

	return winningCorks
}