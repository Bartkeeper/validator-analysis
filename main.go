package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/codec"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/nexidian/gocliselect"
)

// These are variables that can be called from different functions
var sfDelegators []staking.DelegationResponse
var validator staking.Validator
var GrpcConnC1 *grpc.ClientConn

func main() {

	startMenu()
}

// This method will start the CLI application. The user can choose which challenge should be performed
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
		GrpcConnC1 := getGrpcConnC1("cosmos-grpc.polkachu.com:14990")
		output01(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c1o2":
		fmt.Println("You just started Challenge #01 Output #2. Please hold on while we query the node and generate the csv dump.")
		GrpcConnC1 := getGrpcConnC1("cosmos-grpc.polkachu.com:14990")
		output02(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c1o3":
		fmt.Println("You just started Challenge #01 Output #3. Please hold on while we query the node and generate the csv dump.")
		GrpcConnC1 := getGrpcConnC1("cosmos-grpc.polkachu.com:14990")
		output03(GrpcConnC1)
		defer GrpcConnC1.Close()
		fmt.Println("CSV dump available")
	case "c2":
		fmt.Println("You just started Challenge #02. Please make sure the genesis file is in this repo.")
		Output04()
		fmt.Println("CSV dump available")
	}
}

// This method sets up a grpc connection to the Cosmos Hub and stores it in a global variable
func getGrpcConnC1(endpoint string) *grpc.ClientConn {
	grpcConnC1, err := grpc.Dial(
		endpoint, // Polkachu's gRPC server address for the Cosmos Hub
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		log.Panic("could not get grpc connection, reason: ", err)
	}

	return grpcConnC1
}
