//go:generate go run ../main.go --abi-folder ../ethgen/testdata/_source --package-path ../ethgen/testdata

package ethgen_test

import (
	"io/ioutil"
	"testing"

	"github.com/lacasian/ethwheels/ethgen"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	// make sure we are using the latest generated version
	raw, err := ioutil.ReadFile("testdata/erc20.go")
	require.NoError(t, err)

	cd, err := ethgen.ProcessFile("testdata/_source", "ERC20.json")
	require.NoError(t, err)

	code, err := ethgen.GenerateCode("ERC20", "testdata", *cd)
	require.NoError(t, err)

	require.Equal(t, string(raw), string(code))
}
