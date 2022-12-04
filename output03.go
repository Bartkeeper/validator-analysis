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

type delegatorValidator struct {
	delegatorAddr string
	validatorAddr string
	bondedToken   string
}

var delVal []delegatorValidator

func output03(grpcConn *grpc.ClientConn) {
	if len(sfDelegators) < 1 {
		sfDelegators, _ = getSFDelegators(grpcConn)
	}
	appendValidatorsOfDelegator(grpcConn, sfDelegators)
	sort.Slice(delVal, func(i, j int) bool {
		return delVal[i].delegatorAddr > delVal[j].delegatorAddr
	})
	exportOP3(delVal)
}

func appendValidatorsOfDelegator(grpcConn *grpc.ClientConn, delRes []staking.DelegationResponse) error {
	stakingClient := staking.NewQueryClient(grpcConn)

	for _, del := range delRes {
		stakingRes2, err := stakingClient.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: del.Delegation.DelegatorAddress,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
		)
		if err != nil {
			return err
		}

		var delRes = stakingRes2.DelegationResponses

		for _, del := range delRes {
			delVal = append(delVal, delegatorValidator{del.Delegation.DelegatorAddress, del.Delegation.ValidatorAddress, del.Balance.String()})
		}

	}

	return nil
}

func exportOP3(delVal []delegatorValidator) {
	file, err := os.Create("output_01_03.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, del := range delVal {
		row := []string{del.delegatorAddr, del.validatorAddr, del.bondedToken}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		// fmt.Println(i, row)
	}
	defer file.Close()
}
