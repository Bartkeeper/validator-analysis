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

type delegators struct {
	address     string
	delegation  string
	votingpower float64
}

var delegator []delegators

func output02(grpcConn *grpc.ClientConn) {
	sfDelegators, _ = getSFDelegators(grpcConn)
	validator, _ = getValidator(grpcConn)
	appendDelegatorData(sfDelegators, validator)

	sort.Slice(delegator, func(i, j int) bool {
		return delegator[i].votingpower > delegator[j].votingpower
	})
	exportOP2(delegator)
}

func getSFDelegators(grpcConn *grpc.ClientConn) ([]staking.DelegationResponse, error) {
	stakingClient := staking.NewQueryClient(grpcConn)

	stakingRes, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&staking.QueryValidatorDelegationsRequest{
			ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p",
			Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
	)
	if err != nil {
		return nil, err
	}
	var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

	return delRes, nil
}

func getValidator(grpcConn *grpc.ClientConn) (staking.Validator, error) {
	stakingClient := staking.NewQueryClient(grpcConn)

	stakingRes2, err := stakingClient.Validator(
		context.Background(),
		&staking.QueryValidatorRequest{ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p"},
	)
	if err != nil {
		return staking.Validator{}, err
	}
	var val staking.Validator = stakingRes2.Validator

	return val, nil
}

func appendDelegatorData(delRes []staking.DelegationResponse, val staking.Validator) {
	var totalDelegation float64 = val.DelegatorShares.MustFloat64()
	for _, del := range delRes {
		delegator = append(delegator, delegators{del.Delegation.DelegatorAddress, del.Balance.String(), float64(del.Balance.Amount.Int64()) / totalDelegation})
	}
}

func exportOP2(delegator []delegators) {
	file, err := os.Create("output_01_02.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, del := range delegator {
		row := []string{del.address, fmt.Sprintf("%f", del.votingpower)}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}
