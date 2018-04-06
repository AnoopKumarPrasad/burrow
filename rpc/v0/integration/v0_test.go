// +build integration

// Space above here matters
// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"testing"

	"encoding/hex"

	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/rpc/v0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/hyperledger/burrow/binary"
)

func TestTransact(t *testing.T) {
	cli := v0.NewV0Client("http://localhost:1337/rpc")

	address := privateAccounts[1].Address()
	receipt, err := cli.Transact(v0.TransactParam{
		PrivKey:  privateAccounts[0].PrivateKey().RawBytes(),
		Address:  address.Bytes(),
		Data:     []byte{},
		Fee:      2,
		GasLimit: 10000,
	})
	require.NoError(t, err)
	assert.False(t, receipt.CreatesContract)
	assert.Equal(t, address, receipt.ContractAddress)
}

func TestTransactAndHold(t *testing.T) {
	cli := v0.NewV0Client("http://localhost:1337/rpc")

	bc, err := hex.DecodeString(strangeLoopBytecode)
	require.NoError(t, err)
	create, err := cli.TransactAndHold(v0.TransactParam{
		PrivKey:  privateAccounts[0].PrivateKey().RawBytes(),
		Address:  nil,
		Data:     bc,
		Fee:      2,
		GasLimit: 10000,
	})
	require.NoError(t, err)
	assert.Equal(t, 0, create.StackDepth)
	functionID := abi.FunctionID("UpsieDownsie()")
	call, err := cli.TransactAndHold2(v0.TransactParam{
		PrivKey:  privateAccounts[0].PrivateKey().RawBytes(),
		Address:  create.CallData.Callee.Bytes(),
		Data:     functionID[:],
		Fee:      2,
		GasLimit: 10000,
	})
	require.NoError(t, err)
	depth := binary.Uint64FromWord256(binary.LeftPadWord256(call.Return))
	// Would give 23 if taken from wrong frame
	assert.Equal(t, 18, int(depth))
}

