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
