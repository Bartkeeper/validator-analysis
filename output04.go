package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var UmeeGenesis string = "https://raw.githubusercontent.com/umee-network/mainnet/main/genesis.json"

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

type vestingAccs struct {
	vestAddress string
	vestTokens  string
	unlockDate  string
}

var vestingAnalysis []vestingAccs
var payload baseGenesis

func Output04() {
	getGenesis()
	appendVestingData(payload)
	exportOP4(vestingAnalysis)
}

func getGenesis() error {

	content, err := ioutil.ReadFile("./umee-genesis.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	return err
}

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

func exportOP4(vestingAnalysis []vestingAccs) {
	file, err := os.Create("output_04.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, van := range vestingAnalysis {
		row := []string{van.vestAddress, van.vestTokens, van.unlockDate}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	defer file.Close()
}

/*
**** Legacy Code. I tried to query the accounts from the node directly, but ran into a bug I couldn't fix in time

// func getGrpcConnC2() (*grpc.ClientConn, error) {
// 	GrpcConnC2, err := grpc.Dial(
// 		"umee-grpc.polkachu.com:13690",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return GrpcConnC2, nil
// }

*/
