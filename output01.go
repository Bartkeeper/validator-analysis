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

// This is a struct that helps generate and assemble the output
type valAdd struct {
	moniker         string
	opAdd           string
	accAdd          string
	votingpower     float64
	selfDelegation  string
	totalDelegation string
}

// The slice in which the results are stored
var valAdds []valAdd

// An accumulator function that executes all functions for challenge 01 output 1 and stores variables
func output01(grpcConn *grpc.ClientConn) {
	var vals = getVals(grpcConn)
	var valAddresses = getValAddresses(vals)
	getSelfDelegation(grpcConn, valAddresses)

	sort.Slice(valAdds, func(i, j int) bool {
		return valAdds[i].votingpower > (valAdds[j].votingpower)
	})
	exportOP1(valAdds)
}

func getVals(grpcConn *grpc.ClientConn) []staking.Validator {

	// This creates a gRPC client to query the x/staking/types service.
	stakingClient := staking.NewQueryClient(grpcConn)
	stakingRes, err := stakingClient.Validators(
		context.Background(),
		&staking.QueryValidatorsRequest{
			Status:     "BOND_STATUS_BONDED", // Get all validators that are active and bonded
			Pagination: &query.PageRequest{Limit: 500, CountTotal: true}},
	)
	if err != nil {
		log.Fatal("error while querying validators, reason: ", err)
	}
	vals := stakingRes.Validators // storing the Validators in the global variable vals
	return vals
}

// This method contains multiple methods and eventually appends validator data to the validator variable
func getValAddresses(vals []staking.Validator) []valAdd {
	var ph sdk.Coin
	var totalConsensusPower = getTotalTokens(vals)

	for _, val := range vals {
		accAdd := deriveValAccAddress(val)
		vP := getVotingPower(val, totalConsensusPower)
		valAdds = append(valAdds, valAdd{val.GetMoniker(), val.OperatorAddress, accAdd.String(), vP, ph.String(), val.BondedTokens().String() + "uatom"})
	}

	return valAdds
}

// This method loops over the validator's bonded tokens and sums them up
func getTotalTokens(vals []staking.Validator) int64 {
	var total int64
	for _, val := range vals {
		total = total + val.GetBondedTokens().BigInt().Int64()
	}

	return total
}

// This method calculates the voting power of a validator
func getVotingPower(val staking.Validator, totalTokens int64) float64 {
	var vp = float64(val.BondedTokens().Int64()) / float64(totalTokens)
	return vp
}

func deriveValAccAddress(val staking.Validator) sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	if err != nil {
		log.Fatal("could not get validator address, reason: ", err)
	}
	accAddr, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
	if err != nil {
		log.Fatal("could not get account address, reason: ", err)
	}
	return accAddr
}

// This method queries the node to get a validator's self-delegation amount
func getSelfDelegation(grpcConn *grpc.ClientConn, valAdds []valAdd,
) {
	// This creates a gRPC client to query the x/staking service.
	stakingClient := staking.NewQueryClient(grpcConn)

	for i, val := range valAdds {
		// For each validator in the slice, we get their delegators
		stakingRes, err := stakingClient.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: val.accAdd,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
			},
		)
		if err != nil {
			log.Fatal("error while querying delegations, reason: ", err)
		}

		var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

		// for each validator, we look if their account address is in the set of delegators
		for _, del := range delRes {
			if del.Delegation.DelegatorAddress == val.accAdd {
				valAdds[i].selfDelegation = del.Balance.String() // append the self-delegation to the exportable slice
			} else {
				log.Fatalln("delegtor address does not equal validator address")
			}
		}
	}

}

// This function takes the slice and exports it in a CSV
func exportOP1(valAdds []valAdd) {
	file, err := os.Create("output_01_01.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{"Moniker", "Voting Power", "Self-Delegation", "Total Delegation"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, val := range valAdds {
		row := []string{val.moniker, fmt.Sprintf("%f", val.votingpower), val.selfDelegation, val.totalDelegation}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}
