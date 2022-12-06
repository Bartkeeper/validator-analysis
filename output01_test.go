package main

import (
	"context"
	"fmt"
	"testing"

	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
)

func TestGetVals(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		want    string
		wantErr error
	}{
		{
			"PASS: Correct Validator",
			"cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p",
			"cosmosvaloper1x88j7vp2xnw3zec8ur3g4waxycyz7m0mahdv3p",
			nil,
		},
		{
			"FAIL: No validator found",
			"cosmosvaloper100000000000000000000000000000000000000",
			"",
			fmt.Errorf("Validator not found"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint("TestValidator#", i), func(t *testing.T) {
			var got = getGrpcConnC1("cosmos-grpc.polkachu.com:14990")
			stakingClient := staking.NewQueryClient(got)
			stakingRes, err := stakingClient.Validator(
				context.Background(),
				&staking.QueryValidatorRequest{ValidatorAddr: tt.target},
			)
			if err != nil {
				assert.Error(t, tt.wantErr)
				return
			}
			if stakingRes.GetValidator().OperatorAddress != tt.want {
				t.Errorf("GrpcClient Connection = %v, want %v", stakingRes.GetValidator().OperatorAddress, tt.want)
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
