package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/query"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

// This is a struct that helps generate and assemble the output
type delegatorValidator struct {
	delegatorAddr string
	validatorAddr string
	bondedToken   string
}

// This is a struct that helps generate and assemble the output
var delVal []delegatorValidator

// An accumulator function that executes all functions for challenge 01 output 1 and stores variables
func output03(grpcConn *grpc.ClientConn) {
	// sFDelegators might already be filled if challenge 2 has been executed in the same process. Improves efficiency
	if len(sfDelegators) < 1 {
		sfDelegators = getSFDelegators(grpcConn)
	}
	appendValidatorsOfDelegator(grpcConn, sfDelegators)
	sort.Slice(delVal, func(i, j int) bool {
		return delVal[i].delegatorAddr > delVal[j].delegatorAddr
	})
	exportOP3(delVal)
}

// This function queries the addresses of delegators of "Staking Facilities"
// and returns all their delegations and appends it to the output slice
func appendValidatorsOfDelegator(grpcConn *grpc.ClientConn, delRes []staking.DelegationResponse) {
	stakingClient := staking.NewQueryClient(grpcConn)

	for _, del := range delRes {
		stakingRes2, err := stakingClient.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: del.Delegation.DelegatorAddress,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
		)
		if err != nil {
			log.Fatal("could not get delegator delegations, reason: ", err)
		}

		var delRes = stakingRes2.DelegationResponses

		for _, del := range delRes {
			delVal = append(delVal, delegatorValidator{del.Delegation.DelegatorAddress, del.Delegation.ValidatorAddress, del.Balance.String()})
		}
	}
}

// This function takes the slice and exports it in a CSV
func exportOP3(delVal []delegatorValidator) {
	file, err := os.Create("output_01_03.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	header := []string{"Delegator Account", "Delegated Validator", "Delegated Tokens"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, del := range delVal {
		row := []string{del.delegatorAddr, del.validatorAddr, del.bondedToken}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}
