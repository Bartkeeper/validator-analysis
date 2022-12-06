package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/query"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

// This is a struct that helps generate and assemble the output
type delegators struct {
	address     string
	delegation  string
	votingpower float64
}

// The slice in which the results are stored
var delegator []delegators

// An accumulator function that executes all functions for challenge 01 output 1 and stores variables
func output02(grpcConn *grpc.ClientConn) {
	sfDelegators = getSFDelegators(grpcConn)
	validator = getValidator(grpcConn)
	appendDelegatorData(sfDelegators, validator)

	sort.Slice(delegator, func(i, j int) bool {
		return delegator[i].votingpower > delegator[j].votingpower
	})
	exportOP2(delegator)
}

// This function gets all delegators from the validator "Staking Facilities"
func getSFDelegators(grpcConn *grpc.ClientConn) []staking.DelegationResponse {
	stakingClient := staking.NewQueryClient(grpcConn)

	stakingRes, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&staking.QueryValidatorDelegationsRequest{
			ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p",
			Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
	)
	if err != nil {
		log.Fatal("could not get validator delegtions, reason: ", err)
	}
	var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

	return delRes
}

// This function returns a the validator object of "Staking Facilities"
func getValidator(grpcConn *grpc.ClientConn) staking.Validator {
	stakingClient := staking.NewQueryClient(grpcConn)

	stakingRes2, err := stakingClient.Validator(
		context.Background(),
		&staking.QueryValidatorRequest{ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p"},
	)
	if err != nil {
		log.Fatal("could not get validator, reason: ", err)
	}
	var val staking.Validator = stakingRes2.Validator

	return val
}

// This function appends delegator data and their voting power within one validator to the exportable object
func appendDelegatorData(delRes []staking.DelegationResponse, val staking.Validator) {
	var totalDelegation float64 = val.DelegatorShares.MustFloat64()
	for _, del := range delRes {
		delegator = append(delegator, delegators{del.Delegation.DelegatorAddress, del.Balance.String(), float64(del.Balance.Amount.Int64()) / totalDelegation})
	}
}

// This function takes the slice and exports it in a CSV
func exportOP2(delegator []delegators) {
	file, err := os.Create("output_01_02.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	header := []string{"Delegator Account", "Voting Power"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, del := range delegator {
		row := []string{del.address, fmt.Sprintf("%f", del.votingpower)}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}
