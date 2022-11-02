package address_loader_test

import (
	"testing"

	addressLoader "github.com/rarecircles/backend-challenge-go/internal/pkg/address_loader"
	"github.com/stretchr/testify/assert"
)

func TestAddressLoader(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     []string
		expErr   string
	}{
		{
			name:     "file open error",
			filePath: "wrong_file_path",
			expErr:   "failed to open a file",
		}, {
			name:     "success",
			filePath: "./testdata/addresses_test.jsonl",
			want: []string{
				"0x22f4a547ca569ae4dfee96c7aeff37884e25b1cf",
				"0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
				"0xe9c8934ebd00bf73b0e961d1ad0794fb22837206",
			},
			expErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan string, len(tt.want))
			loader := addressLoader.NewAddressLoader(nil, ch)

			err := loader.Load(tt.filePath)
			if err != nil {
				assert.ErrorContains(t, err, tt.expErr)
				return
			}

			var got []string
			for s := range ch {
				got = append(got, s)
			}

			assert.Equal(t, len(got), len(tt.want))
			for i := range got {
				assert.Equal(t, got[i], tt.want[i])
			}
		})
	}

}
