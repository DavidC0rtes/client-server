package server_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DavidC0rtes/client-server/server"
)

// TestBadListenRequest calls SplitRequest with a (listen) request that
// violates the custom protocol spec.
func TestBadListenRequest(t *testing.T) {
	ans, err := server.SplitRequest("listen")
	if err == nil {
		t.Errorf("SplitRequest(listen 2) = %v; want error", ans)
	}

	if len(ans) > 0 {
		t.Errorf("SplitRequest(listen 2) = %v; want nil", ans)
	}

}

// TestBadSendRequest calls SplitRequest with a (->) request that
// violates the custom protocol spec.
func TestBadSendRequest(t *testing.T) {
	ans2, err := server.SplitRequest("-> edji 232 deidje")
	if err == nil {
		t.Errorf("SplitRequest(-> edji 232 deidje) = %v; want error", ans2)
	}

	if len(ans2) > 0 {
		t.Errorf("SplitRequest(-> edji 232 deidje) = %v; want nil", ans2)
	}
}

// Make multiple calls to SplitRequest to test good and bad requests.
func TestSplitRequestTable(t *testing.T) {
	var tests = []struct {
		request string
		parts   []string
		err     error
	}{
		{"listen 1", []string{"listen", "1"}, nil},
		{"listen     	1", []string{"listen", "1"}, nil},
		{"-> 70 foo 0", []string{"->", "70", "foo", "0"}, nil},
		{"-> 70 foo	90", []string{"->", "70", "foo", "90"}, nil},
		{"-> 70 fADEOO90", nil, errors.New("Malformed request")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.request)
		t.Run(testname, func(t *testing.T) {
			ans, err := server.SplitRequest(tt.request)
			if !reflect.DeepEqual(ans, tt.parts) {
				t.Errorf("got %v want %v", ans, tt.parts)
			}

			if err != nil {
				if tt.err == nil {
					t.Errorf("got error %v want %v", err, tt.err)
				}
			}
		})
	}
}
