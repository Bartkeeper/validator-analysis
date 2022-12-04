package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/codec"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/nexidian/gocliselect"
)

/*
type valAdd struct {
	moniker         string
	opAdd           string
	accAdd          string
	votingpower     float64
	selfDelegation  string
	totalDelegation string
}

var valAdds []valAdd

type delegators struct {
	address     string
	delegation  string
	votingpower float64
}

var delegator []delegators

type delegatorValidator struct {
	delegatorAddr string
	validatorAddr string
	bondedToken   string
}

var delVal []delegatorValidator
*/
var sfDelegators []staking.DelegationResponse
var validator staking.Validator
var GrpcConn *grpc.ClientConn

func main() {
	GrpcConn, _ = getGrpcConn()
	startMenu()
	/*

		Output01(GrpcConn)
		output02(GrpcConn)
		output03(GrpcConn)

		defer GrpcConn.Close() */
}

func startMenu() {
	menu := gocliselect.NewMenu("Choose a challenge")

	menu.AddItem("Challenge #01 Output #1", "c1o1")
	menu.AddItem("Challenge #01 Output #2", "c1o2")
	menu.AddItem("Challenge #01 Output #3", "c1o3")

	choice := menu.Display()

	fmt.Printf("Choice: %s\n", choice)

	switch choice {
	case "c1o1":
		fmt.Println("You just started Challenge #01 Output #1. Please hold on while we query the node and generate the csv dump.")
		output01(GrpcConn)
		fmt.Println("CSV dump available")
	case "c1o2":
		output02(GrpcConn)
	case "c1o3":
		output03(GrpcConn)
	}
}

func getGrpcConn() (*grpc.ClientConn, error) {
	grpcConn, err := grpc.Dial(
		"cosmos-grpc.polkachu.com:14990", // Polkachu's gRPC server address.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		return nil, err
	}

	return grpcConn, nil
}

// func output01(grpcConn *grpc.ClientConn) {
// 	var vals, _ = getVals(grpcConn)
// 	var valAddresses = getValAddresses(vals)
// 	getSelfDelegation(grpcConn, valAddresses)

// 	sort.Slice(valAdds, func(i, j int) bool {
// 		return valAdds[i].votingpower > (valAdds[j].votingpower)
// 	})
// 	exportOP1(valAdds)
// }

// func output02(grpcConn *grpc.ClientConn) {
// 	sfDelegators, _ = getSFDelegators(grpcConn)
// 	validator, _ = getValidator(grpcConn)
// 	appendDelegatorData(sfDelegators, validator)

// 	sort.Slice(delegator, func(i, j int) bool {
// 		return delegator[i].votingpower > delegator[j].votingpower
// 	})
// 	exportOP2(delegator)
// }

// func output03(grpcConn *grpc.ClientConn) {
// 	if len(sfDelegators) < 1 {
// 		sfDelegators, _ = getSFDelegators(grpcConn)
// 	}
// 	appendValidatorsOfDelegator(grpcConn, sfDelegators)
// 	sort.Slice(delVal, func(i, j int) bool {
// 		return delVal[i].delegatorAddr > delVal[j].delegatorAddr
// 	})
// 	exportOP3(delVal)
// }

// func getVals(grpcConn *grpc.ClientConn) ([]staking.Validator, error) {

// 	// This creates a gRPC client to query the x/staking service.
// 	stakingClient := staking.NewQueryClient(grpcConn)
// 	stakingRes, err := stakingClient.Validators(
// 		context.Background(),
// 		&staking.QueryValidatorsRequest{
// 			Status:     "BOND_STATUS_BONDED",
// 			Pagination: &query.PageRequest{Limit: 500, CountTotal: true}},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	vals := stakingRes.Validators // Here I only do get 100(vals) instead of 175(from stakingRes)
// 	return vals, nil
// }

// func getTotalTokens(vals []staking.Validator) int64 {
// 	var total int64
// 	for _, val := range vals {
// 		total = total + val.GetBondedTokens().BigInt().Int64()
// 	}

// 	return total
// }

// func getVotingPower(val staking.Validator, totalTokens int64) float64 {
// 	var vp = float64(val.BondedTokens().Int64()) / float64(totalTokens)
// 	return vp
// }

// func deriveValAccAddress(val staking.Validator) (sdk.AccAddress, error) {
// 	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
// 	if err != nil {
// 		return nil, err
// 	}
// 	accAddr, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return accAddr, nil
// }

// func getValAddresses(vals []staking.Validator) []valAdd {
// 	var ph sdk.Coin
// 	var totalConsensusPower = getTotalTokens(vals)

// 	for _, val := range vals {
// 		accAdd, _ := deriveValAccAddress(val)
// 		vP := getVotingPower(val, totalConsensusPower)
// 		valAdds = append(valAdds, valAdd{val.GetMoniker(), val.OperatorAddress, accAdd.String(), vP, ph.String(), val.BondedTokens().String() + "uatom"})
// 	}

// 	return valAdds
// }

// func getSelfDelegation(grpcConn *grpc.ClientConn, valAdds []valAdd,
// ) error {

// 	// This creates a gRPC client to query the x/staking service.
// 	stakingClient := staking.NewQueryClient(grpcConn)

// 	for i, val := range valAdds {
// 		stakingRes, err := stakingClient.DelegatorDelegations(
// 			context.Background(),
// 			&staking.QueryDelegatorDelegationsRequest{
// 				DelegatorAddr: val.accAdd,
// 				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
// 			},
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

// 		for _, del := range delRes {
// 			if del.Delegation.DelegatorAddress == val.accAdd {
// 				valAdds[i].selfDelegation = del.Balance.String()
// 			} else {
// 				return fmt.Errorf("something went wrong")
// 			}
// 		}
// 	}
// 	return nil
// }

// func exportOP1(valAdds []valAdd) {
// 	file, err := os.Create("challenge_01.csv")
// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}
// 	w := csv.NewWriter(file)
// 	defer w.Flush()
// 	// Using Write
// 	for _, val := range valAdds {
// 		row := []string{val.moniker, fmt.Sprintf("%f", val.votingpower), val.selfDelegation, val.totalDelegation}
// 		if err := w.Write(row); err != nil {
// 			log.Fatalln("error writing record to file", err)
// 		}
// 		// fmt.Println(i, row)
// 	}
// 	defer file.Close()
// }

// func getSFDelegators(grpcConn *grpc.ClientConn) ([]staking.DelegationResponse, error) {
// 	stakingClient := staking.NewQueryClient(grpcConn)

// 	stakingRes, err := stakingClient.ValidatorDelegations(
// 		context.Background(),
// 		&staking.QueryValidatorDelegationsRequest{
// 			ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p",
// 			Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

// 	return delRes, nil
// }

// func appendValidatorsOfDelegator(grpcConn *grpc.ClientConn, delRes []staking.DelegationResponse) error {
// 	stakingClient := staking.NewQueryClient(grpcConn)

// 	for _, del := range delRes {
// 		stakingRes2, err := stakingClient.DelegatorDelegations(
// 			context.Background(),
// 			&staking.QueryDelegatorDelegationsRequest{
// 				DelegatorAddr: del.Delegation.DelegatorAddress,
// 				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true}},
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		var delRes = stakingRes2.DelegationResponses

// 		for _, del := range delRes {
// 			delVal = append(delVal, delegatorValidator{del.Delegation.DelegatorAddress, del.Delegation.ValidatorAddress, del.Balance.String()})
// 		}

// 	}

// 	return nil
// }

// func getValidator(grpcConn *grpc.ClientConn) (staking.Validator, error) {
// 	stakingClient := staking.NewQueryClient(grpcConn)

// 	stakingRes2, err := stakingClient.Validator(
// 		context.Background(),
// 		&staking.QueryValidatorRequest{ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p"},
// 	)
// 	if err != nil {
// 		return staking.Validator{}, err
// 	}
// 	var val staking.Validator = stakingRes2.Validator

// 	return val, nil
// }

// func appendDelegatorData(delRes []staking.DelegationResponse, val staking.Validator) {
// 	var totalDelegation float64 = val.DelegatorShares.MustFloat64()
// 	for _, del := range delRes {
// 		delegator = append(delegator, delegators{del.Delegation.DelegatorAddress, del.Balance.String(), float64(del.Balance.Amount.Int64()) / totalDelegation})
// 	}
// }

// func exportOP2(delegator []delegators) {
// 	file, err := os.Create("output_01_02.csv")
// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}
// 	w := csv.NewWriter(file)
// 	defer w.Flush()
// 	// Using Write
// 	for _, del := range delegator {
// 		row := []string{del.address, fmt.Sprintf("%f", del.votingpower)}
// 		if err := w.Write(row); err != nil {
// 			log.Fatalln("error writing record to file", err)
// 		}
// 		// fmt.Println(i, row)
// 	}
// 	defer file.Close()
// }

// func exportOP3(delVal []delegatorValidator) {
// 	file, err := os.Create("output_01_03.csv")
// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}
// 	w := csv.NewWriter(file)
// 	defer w.Flush()
// 	// Using Write
// 	for _, del := range delVal {
// 		row := []string{del.delegatorAddr, del.validatorAddr, del.bondedToken}
// 		if err := w.Write(row); err != nil {
// 			log.Fatalln("error writing record to file", err)
// 		}
// 		// fmt.Println(i, row)
// 	}
// 	defer file.Close()
// }
