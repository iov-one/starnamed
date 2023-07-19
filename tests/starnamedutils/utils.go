package starnamedutils

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

func randomString(length int) string {
	rand.NewSource(time.Now().UnixNano())

	var alphabet string = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder

	l := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func IBCWalletFactory(t *testing.T, ctx context.Context, keyNamePrefix string, number_of_users int, token_ammount int64, chain ibc.Chain) (users [](ibc.Wallet)) {

	// Create the slice of users
	users = make([]ibc.Wallet, number_of_users)

	for i := 0; i < number_of_users; i++ {
		users[i] = interchaintest.GetAndFundTestUsers(t, ctx, keyNamePrefix, token_ammount, chain)[0]
	}

	return
}

func JsonUnmarshal(data string) (x map[string]interface{}, err error) {

	err = json.Unmarshal([]byte(data), &x)
	if err != nil {
		x = nil
		return
	}
	return
}
