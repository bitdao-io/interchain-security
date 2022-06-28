package types

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

func TestPendingClientKey(t *testing.T) {
	tests := []struct {
		name    string
		sec     int64
		nsec    int64
		chainId string
	}{
		{"malicious input", 83, 70, "foowiz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := time.Unix(tt.sec, tt.nsec)
			suffixKey := PendingClientKey(ts, tt.chainId)
			splitKey := strings.Split(string(suffixKey), "/")
			binary.BigEndian.Uint64([]byte(splitKey[1]))
		})
	}
}
