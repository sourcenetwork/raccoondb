package marshal

import (
	"fmt"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/utils"
	"github.com/stretchr/testify/require"
)

func Test_Int64_EncDec_Inverses(t *testing.T) {
	vals := []int64{
		-10,
		1,
		10,
		-55,
		0,
		^0,
	}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v inverses", val), func(t *testing.T) {
			enc := EncodeInt64(val)
			dec := DecodeInt64(enc)
			require.Equal(t, val, dec)
		})
	}
}

func Test_Int64_EncDec_Orderable(t *testing.T) {
	vals := []int64{
		-10,
		1,
		10,
		-55,
		0,
		-1,
	}
	encodedVals := utils.MapSlice(vals, EncodeInt64)
	sortable := utils.FromComparator(encodedVals, utils.LexographicBytesComparator)
	sortedBytes := sortable.Sort()
	sorted := utils.MapSlice(sortedBytes, DecodeInt64)
	want := []int64{
		-55,
		-10,
		-1,
		0,
		1,
		10,
	}
	require.Equal(t, want, sorted)
}
