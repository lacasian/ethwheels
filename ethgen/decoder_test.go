package ethgen_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lacasian/ethwheels/ethgen/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tzapu/thelper"
)

var update = flag.Bool("update", false, "update golden files")

func TestUnpackLog(t *testing.T) {
	tests := map[string]struct {
		want    string
		wantErr bool
	}{
		"erc20-transfer-event": {
			wantErr: false,
		},
		// wrong Topic1
		"erc20-transfer-event-error": {
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			ef := fmt.Sprintf("testdata/%s.json", n)
			wf := fmt.Sprintf("testdata/%s.golden", n)

			var event types.Log
			thelper.Load(t, ef, &event)

			actual, err := testdata.ERC20.ERC20TransferEvent(event)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			thelper.SaveOnUpdate(t, update, wf, actual)
			var expected testdata.ERC20TransferEvent
			thelper.Load(t, wf, &expected)

			assert.Equal(t, expected, actual)
		})
	}
}
