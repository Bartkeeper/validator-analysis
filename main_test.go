package main

import (
	"context"
	"fmt"
	"testing"

	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
)

// This test is about making sure an error is thrown when a wrong grpc endpoint is provided
func TestGetGrpcConnC1(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr error
	}{
		{
			"PASS: Correct grpc endpoint",
			"cosmos-grpc.polkachu.com:14990",
			nil,
		},
		{
			"FAIL: Wrong grpc endpoint",
			"test",
			fmt.Errorf("Wrong endpoint"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint("TestValidator#", i), func(t *testing.T) {
			var got = getGrpcConnC1(tt.target)
			stakingClient := staking.NewQueryClient(got)
			stakingRes, err := stakingClient.Validator(
				context.Background(),
				&staking.QueryValidatorRequest{ValidatorAddr: "cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p"},
			)
			fmt.Println("Just a test to see if a ", stakingRes.Size(), " element is sent or not")
			if err != nil {
				assert.Error(t, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, tt.wantErr)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}

}
