package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/evmos/ethermint/server/config"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/UptickNetwork/uptick/contracts"
	"github.com/UptickNetwork/uptick/x/erc721/types"
)

// DeployERC721Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner.
func (k Keeper) DeployERC721Contract(
	ctx sdk.Context,
	msg *types.MsgConvertNFT,
) (common.Address, error) {

	class, err := k.nftKeeper.GetDenomInfo(ctx, msg.ClassId)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "nft class is invalid %s: %s", class.Id, err.Error())
	}

	ctorArgs, err := contracts.ERC721UpticksContract.ABI.Pack(
		"",
		class.Name,
		class.Symbol,
		class.Uri,
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "nft class is invalid %s: %s", class.Id, err.Error())
	}

	data := make([]byte, len(contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.ERC20MinterBurnerDecimalsContract.Bin)], contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	if _, err = k.CallEVMWithData(ctx, types.ModuleAddress, nil, data, true); err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy contract for %s", class.Id)
	}

	return contractAddr, nil
}

// QueryERC721 returns the data of a deployed ERC721 contract
func (k Keeper) QueryERC721(
	ctx sdk.Context,
	contract common.Address,
) (types.ERC721Data, error) {
	var (
		nameRes   types.ERC721StringResponse
		symbolRes types.ERC721StringResponse
	)

	erc721 := contracts.ERC721UpticksContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "name")
	if err != nil {
		return types.ERC721Data{}, err
	}

	if err := erc721.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ERC721Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack name: %s", err.Error(),
		)
	}

	// Symbol
	res, err = k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "symbol")
	if err != nil {
		return types.ERC721Data{}, err
	}

	if err := erc721.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ERC721Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack symbol: %s", err.Error(),
		)
	}

	return types.NewERC721Data(nameRes.Value, symbolRes.Value), nil
}

// QueryERC721Enhance returns the data of a deployed ERC721 contract
func (k Keeper) QueryERC721Enhance(
	ctx sdk.Context,
	contract common.Address,
	tokenID *big.Int,
) (types.ERC721Enhance, error) {

	fmt.Printf("xxl 01 QueryERC721Enhance start \n")
	erc721 := contracts.ERC721UpticksContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc721, types.ModuleAddress, contract, true, "getEnhanceInfo", tokenID)
	if err != nil {
		fmt.Printf("xxl 01 QueryERC721Enhance 1 \n")
		return types.ERC721Enhance{}, err
	}

	ret, err := erc721.Unpack("getEnhanceInfo", res.Ret)
	if err != nil {
		fmt.Printf("xxl 01 QueryERC721Enhance resRet %v \n", err)
	}
	fmt.Printf("xxl 01 QueryERC721Enhance resRet %v \n", ret)

	if len(ret) != 4 {
		return types.ERC721Enhance{}, nil
	}

	return types.NewERC721Enhance(ret[0].(string), ret[1].(string), ret[2].(string), ret[3].(string)), nil
}

// QueryERC721Token returns the data of a ERC721 token
func (k Keeper) QueryERC721Token(
	ctx sdk.Context,
	contract common.Address,
	tokenID *big.Int,
) (types.ERC721TokenData, error) {
	var (
		nameRes   types.ERC721TokenStringResponse
		symbolRes types.ERC721TokenStringResponse
		uriRes    types.ERC721TokenStringResponse
	)

	erc721 := contracts.ERC721UpticksContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "name")
	if err != nil {
		return types.ERC721TokenData{}, err
	}

	if err := erc721.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ERC721TokenData{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack name: %s", err.Error(),
		)
	}

	// Symbol
	res, err = k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "symbol")
	if err != nil {
		return types.ERC721TokenData{}, err
	}

	if err := erc721.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ERC721TokenData{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack symbol: %s", err.Error(),
		)
	}

	// Uri
	//res, err = k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "tokenURI", tokenID)
	//if err != nil {
	//	return types.ERC721TokenData{}, err
	//}

	if err := erc721.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ERC721TokenData{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack uri: %s", err.Error(),
		)
	}

	return types.NewERC721TokenData(nameRes.Value, symbolRes.Value, uriRes.Value), nil
}

// QueryERC721NextTokenID returns the next tokenID to mint
func (k Keeper) QueryERC721NextTokenID(
	ctx sdk.Context,
	contract common.Address,
) (*big.Int, error) {
	var idRes types.ERC721TokenIDResponse

	erc721 := contracts.ERC721UpticksContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "nextTokenId")
	if err != nil {
		return nil, err
	}

	if err := erc721.UnpackIntoInterface(&idRes, "nextTokenId", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack nextTokenId: %s", err.Error(),
		)
	}

	return idRes.Value, nil
}

// QueryERC721TokenOwner returns the owner of given tokenID
func (k Keeper) QueryERC721TokenOwner(
	ctx sdk.Context,
	contract common.Address,
	tokenID *big.Int,
) (common.Address, error) {
	var ownerRes types.ERC721TokenOwnerResponse

	erc721 := contracts.ERC721UpticksContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "ownerOf", tokenID)
	if err != nil {
		return common.Address{}, err
	}

	if err := erc721.UnpackIntoInterface(&ownerRes, "ownerOf", res.Ret); err != nil {
		return common.Address{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack owner: %s", err.Error(),
		)
	}

	return ownerRes.Value, nil
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	commit bool,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {

	fmt.Printf("xxl 01 CallEVM 001 k.CallEVM method %v \n", method)
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.CallEVMWithData(ctx, from, &contract, data, commit)
	if err != nil {
		fmt.Printf("xxl 01 CallEVM 002 k.CallEVM method %s-%s \n", method, contract)
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	gasCap := config.DefaultGasCap
	if commit {
		args, err := json.Marshal(evmtypes.TransactionArgs{
			From: &from,
			To:   contract,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		})
		if err != nil {
			return nil, err
		}
		gasCap = gasRes.Gas
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0), // amount
		gasCap,        // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		!commit,               // isFake
	)

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}

// monitorApprovalEvent returns an error if the given transactions logs include
// an unexpected `Approval` event
func (k Keeper) monitorApprovalEvent(res *evmtypes.MsgEthereumTxResponse) error {
	if res == nil || len(res.Logs) == 0 {
		return nil
	}

	//logApprovalSig := []byte("Approval(address,address,uint256)")
	//logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)
	//
	//for _, log := range res.Logs {
	//	if log.Topics[0] == logApprovalSigHash.Hex() {
	//		return sdkerrors.Wrapf(
	//			types.ErrUnexpectedEvent, "unexpected Approval event 1111" ,
	//		)
	//	}
	//}

	return nil
}
