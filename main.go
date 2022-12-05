package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/codec"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/nexidian/gocliselect"
)

var sfDelegators []staking.DelegationResponse
var validator staking.Validator
var GrpcConnC1 *grpc.ClientConn
var GrpcConnC2 *grpc.ClientConn

// var UmeeGenesis string = "https://raw.githubusercontent.com/umee-network/mainnet/main/genesis.json"

// type baseGenesis struct {
// 	Genesis_time     time.Time
// 	Chain_id         string
// 	Initial_height   string
// 	Consensus_Params ConsensusParam
// 	App_hash         string
// 	App_state        AppState
// }

// type ConsensusParam struct {
// 	Block     BlockStruct
// 	Evidence  EvidenceStruct
// 	Validator ValidatorStruct
// 	Version   any
// }

// type BlockStruct struct {
// 	Max_bytes    string
// 	Max_gas      string
// 	Time_iota_ms string
// }

// type EvidenceStruct struct {
// 	Max_age_num_blocks string
// 	Max_age_duration   string
// 	Max_bytes          string
// }

// type ValidatorStruct struct {
// 	Pub_key_types []string
// }

// type AppState struct {
// 	Auth         AuthStruct
// 	Authz        map[string]interface{}
// 	Bank         map[string]interface{}
// 	Capability   map[string]interface{}
// 	Crisis       map[string]interface{}
// 	Distribution map[string]interface{}
// 	Evidence     map[string]interface{}
// 	Feegrant     map[string]interface{}
// 	Genutil      map[string]interface{}
// 	Gov          map[string]interface{}
// 	Gravity      map[string]interface{}
// 	Ibc          map[string]interface{}
// 	Mint         map[string]interface{}
// 	Params       map[string]interface{}
// 	Slashing     map[string]interface{}
// 	Staking      map[string]interface{}
// 	Transfer     map[string]interface{}
// 	Upgrade      map[string]interface{}
// 	Vesting      map[string]interface{}
// }

// type AuthStruct struct {
// 	Params   map[string]interface{}
// 	Accounts []AccountStruct
// }

// type AccountStruct struct {
// 	Typ                  string
// 	Base_vesting_account BaseVestingAccount
// }

// type BaseVestingAccount struct {
// 	Base_account      BaseAccountStruct
// 	Original_vesting  []VestingStruct
// 	Delegated_free    []any
// 	Delegated_vesting []any
// 	End_time          string
// }

// type BaseAccountStruct struct {
// 	Address        string
// 	Pub_key        string
// 	Account_number string
// 	Sequence       string
// }

// type VestingStruct struct {
// 	Denom  string
// 	Amount string
// }

// type vestingAccs struct {
// 	vestAddress string
// 	vestTokens  string
// 	unlockDate  string
// }

// var vestingAnalysis []vestingAccs
// var payload baseGenesis

func main() {

	startMenu()

}

func startMenu() {
	menu := gocliselect.NewMenu("Choose a challenge")

	menu.AddItem("Challenge #01 Output #1", "c1o1")
	menu.AddItem("Challenge #01 Output #2", "c1o2")
	menu.AddItem("Challenge #01 Output #3", "c1o3")
	menu.AddItem("Challenge #02", "c2")

	choice := menu.Display()

	fmt.Printf("Choice: %s\n", choice)

	switch choice {
	case "c1o1":
		fmt.Println("You just started Challenge #01 Output #1. Please hold on while we query the node and generate the csv dump.")
		GrpcConnC1, _ = getGrpcConnC1()
		output01(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c1o2":
		fmt.Println("You just started Challenge #01 Output #2. Please hold on while we query the node and generate the csv dump.")
		GrpcConnC1, _ = getGrpcConnC1()
		output02(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c1o3":
		fmt.Println("You just started Challenge #01 Output #3. Please hold on while we query the node and generate the csv dump.")
		GrpcConnC1, _ = getGrpcConnC1()
		output03(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c2":
		fmt.Println("You just started Challenge #02. Please make sure the genesis file is in this repo.")
		Output04()
		fmt.Println("CSV dump available")
	}
}

func getGrpcConnC1() (*grpc.ClientConn, error) {
	grpcConnC1, err := grpc.Dial(
		"cosmos-grpc.polkachu.com:14990", // Polkachu's gRPC server address.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		return nil, err
	}

	return grpcConnC1, nil
}

// func getGenesis(grpcConn *grpc.ClientConn) error {

// 	content, err := ioutil.ReadFile("./umee-genesis.json")
// 	if err != nil {
// 		log.Fatal("Error when opening file: ", err)
// 	}

// 	err = json.Unmarshal(content, &payload)
// 	if err != nil {
// 		log.Fatal("Error when opening file: ", err)
// 	}

// 	return err

// }

// func appendVestingData(payload baseGenesis) {

// 	for _, acc := range payload.App_state.Auth.Accounts {
// 		var vestingAmount string
// 		var unlockTime time.Time
// 		if acc.Base_vesting_account.End_time != "" {
// 			i, err := strconv.ParseInt(acc.Base_vesting_account.End_time, 10, 64)
// 			if err != nil {
// 				panic(err)
// 			}

// 			unlockTime = time.Unix(i, 0)
// 		}
// 		for _, amount := range acc.Base_vesting_account.Original_vesting {
// 			vestingAmount = amount.Amount + amount.Denom
// 		}
// 		vestingAnalysis = append(vestingAnalysis, vestingAccs{
// 			acc.Base_vesting_account.Base_account.Address,
// 			vestingAmount,
// 			unlockTime.String(),
// 		})
// 	}

// }

// func calcUnlockTime(unixTimestamp string) time.Time {
// 	if unixTimestamp != "" {
// 		i, err := strconv.ParseInt(unixTimestamp, 10, 64)
// 		if err != nil {
// 			panic(err)
// 		}

// 		var unlockTime = time.Unix(i, 0)

// 		return unlockTime
// 	}

// 	return payload.Genesis_time
// }
