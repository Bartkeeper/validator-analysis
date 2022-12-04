package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

type valAdd struct {
	moniker         string
	opAdd           string
	accAdd          string
	votingpower     float64
	selfDelegation  string
	totalDelegation string
}

var valAdds []valAdd

func output01(grpcConn *grpc.ClientConn) {
	var vals, _ = getVals(grpcConn)
	var valAddresses = getValAddresses(vals)
	getSelfDelegation(grpcConn, valAddresses)

	sort.Slice(valAdds, func(i, j int) bool {
		return valAdds[i].votingpower > (valAdds[j].votingpower)
	})
	exportOP1(valAdds)
}

func getVals(grpcConn *grpc.ClientConn) ([]staking.Validator, error) {

	// This creates a gRPC client to query the x/staking service.
	stakingClient := staking.NewQueryClient(grpcConn)
	stakingRes, err := stakingClient.Validators(
		context.Background(),
		&staking.QueryValidatorsRequest{
			Status:     "BOND_STATUS_BONDED",
			Pagination: &query.PageRequest{Limit: 500, CountTotal: true}},
	)
	if err != nil {
		return nil, err
	}
	vals := stakingRes.Validators // Here I only do get 100(vals) instead of 175(from stakingRes)
	return vals, nil
}

func getValAddresses(vals []staking.Validator) []valAdd {
	var ph sdk.Coin
	var totalConsensusPower = getTotalTokens(vals)

	for _, val := range vals {
		accAdd, _ := deriveValAccAddress(val)
		vP := getVotingPower(val, totalConsensusPower)
		valAdds = append(valAdds, valAdd{val.GetMoniker(), val.OperatorAddress, accAdd.String(), vP, ph.String(), val.BondedTokens().String() + "uatom"})
	}

	return valAdds
}

func getTotalTokens(vals []staking.Validator) int64 {
	var total int64
	for _, val := range vals {
		total = total + val.GetBondedTokens().BigInt().Int64()
	}

	return total
}

func getVotingPower(val staking.Validator, totalTokens int64) float64 {
	var vp = float64(val.BondedTokens().Int64()) / float64(totalTokens)
	return vp
}

func deriveValAccAddress(val staking.Validator) (sdk.AccAddress, error) {
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	if err != nil {
		return nil, err
	}
	accAddr, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
	if err != nil {
		return nil, err
	}
	return accAddr, nil
}

func getSelfDelegation(grpcConn *grpc.ClientConn, valAdds []valAdd,
) error {

	// This creates a gRPC client to query the x/staking service.
	stakingClient := staking.NewQueryClient(grpcConn)

	for i, val := range valAdds {
		stakingRes, err := stakingClient.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: val.accAdd,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
			},
		)
		if err != nil {
			return err
		}

		var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

		for _, del := range delRes {
			if del.Delegation.DelegatorAddress == val.accAdd {
				valAdds[i].selfDelegation = del.Balance.String()
			} else {
				return fmt.Errorf("something went wrong")
			}
		}
	}
	return nil
}

func exportOP1(valAdds []valAdd) {
	file, err := os.Create("output_01_01.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, val := range valAdds {
		row := []string{val.moniker, fmt.Sprintf("%f", val.votingpower), val.selfDelegation, val.totalDelegation}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		// fmt.Println(i, row)
	}
	defer file.Close()
}
