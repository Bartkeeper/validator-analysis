package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	// auth staking "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Bear with me. I ran into an issue where I had to store the genesis file data in a struct.
// Alternative way (with unresolved bug) at the end of this file

type baseGenesis struct {
	Genesis_time     time.Time
	Chain_id         string
	Initial_height   string
	Consensus_Params ConsensusParam
	App_hash         string
	App_state        AppState
}

type ConsensusParam struct {
	Block     BlockStruct
	Evidence  EvidenceStruct
	Validator ValidatorStruct
	Version   any
}

type BlockStruct struct {
	Max_bytes    string
	Max_gas      string
	Time_iota_ms string
}

type EvidenceStruct struct {
	Max_age_num_blocks string
	Max_age_duration   string
	Max_bytes          string
}

type ValidatorStruct struct {
	Pub_key_types []string
}

type AppState struct {
	Auth         AuthStruct
	Authz        map[string]interface{}
	Bank         map[string]interface{}
	Capability   map[string]interface{}
	Crisis       map[string]interface{}
	Distribution map[string]interface{}
	Evidence     map[string]interface{}
	Feegrant     map[string]interface{}
	Genutil      map[string]interface{}
	Gov          map[string]interface{}
	Gravity      map[string]interface{}
	Ibc          map[string]interface{}
	Mint         map[string]interface{}
	Params       map[string]interface{}
	Slashing     map[string]interface{}
	Staking      map[string]interface{}
	Transfer     map[string]interface{}
	Upgrade      map[string]interface{}
	Vesting      map[string]interface{}
}

type AuthStruct struct {
	Params   map[string]interface{}
	Accounts []AccountStruct
}

type AccountStruct struct {
	Typ                  string
	Base_vesting_account BaseVestingAccount
}

type BaseVestingAccount struct {
	Base_account      BaseAccountStruct
	Original_vesting  []VestingStruct
	Delegated_free    []any
	Delegated_vesting []any
	End_time          string
}

type BaseAccountStruct struct {
	Address        string
	Pub_key        string
	Account_number string
	Sequence       string
}

type VestingStruct struct {
	Denom  string
	Amount string
}

// This is a struct that helps generate and assemble the output
type vestingAccs struct {
	vestAddress string
	vestTokens  string
	unlockDate  string
}

// This is a struct that helps generate and assemble the output
var vestingAnalysis []vestingAccs

// The genesisFile
var payload baseGenesis

// An accumulator function that executes all functions for challenge 01 output 1 and stores variables
func Output04() {
	getGenesis()
	appendVestingData(payload)
	exportOP4(vestingAnalysis)
}

// This function stores the genesisFile in the variable genesisFile
func getGenesis() {
	// make sure this file is in the root directory of the
	content, err := ioutil.ReadFile("./inputs/umee-genesis.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
}

// This function queries each account's end_time and formats it in a regular timestamp
// It appends a vesting account's data to the output slice
func appendVestingData(payload baseGenesis) {

	for _, acc := range payload.App_state.Auth.Accounts {
		var vestingAmount string
		var unlockTime time.Time
		if acc.Base_vesting_account.End_time != "" {
			i, err := strconv.ParseInt(acc.Base_vesting_account.End_time, 10, 64)
			if err != nil {
				panic(err)
			}

			unlockTime = time.Unix(i, 0)
		}
		for _, amount := range acc.Base_vesting_account.Original_vesting {
			vestingAmount = amount.Amount + amount.Denom
		}
		vestingAnalysis = append(vestingAnalysis, vestingAccs{
			acc.Base_vesting_account.Base_account.Address,
			vestingAmount,
			unlockTime.String(),
		})
	}

}

// This function takes the slice and exports it in a CSV
func exportOP4(vestingAnalysis []vestingAccs) {
	file, err := os.Create("output_04.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	header := []string{"Vesting Account", "Vested Tokens", "Unlock Date"}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	for _, van := range vestingAnalysis {
		row := []string{van.vestAddress, van.vestTokens, van.unlockDate}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}

/*
**** Legacy Code. ****

I tried to query the accounts from the node directly, but ran into a bug I couldn't fix in time.
This area is commented out because I want to show you guys my thought process and where I got stuck.

func getGrpcConnC2() (*grpc.ClientConn, error) {
	GrpcConnC2, err := grpc.Dial(
		"umee-grpc.polkachu.com:13690",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		return nil, err
	}

	return GrpcConnC2, nil
}

func getAccounts(grpcConn *grpc.ClientConn) {
	authClient := auth.NewQueryClient(grpcConn)

	authRes, err := authClient.Accounts{ //here is where the program paniced with "runtime error: invalid memory address or nil pointer dereference"
		context.Background(),
		&auth.QueryAccountsRequest{},
	}
	if err != nil {
		return nil, err
	}

}

// If this would have worked, I would have been able to query the vesting accounts that came from genesis by calling GetGenesisStateFromAppState.
// This would have avoided parsing the genesis file into a tree of structs.

*/
