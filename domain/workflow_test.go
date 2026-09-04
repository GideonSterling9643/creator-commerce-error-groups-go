package domain

import (
	"reflect"
	"testing"
)

func TestProcessCapturesStableBusinessGroup(t *testing.T) {
	cases := []struct {
		name string
		job  Job
		want []string
	}{
		{"subscriber update", Job{"asset-42", "sub-7", "subscriber_update"}, []string{"creator-commerce", "subscriber_update"}},
		{"delivery", Job{"asset-9", "sub-2", "digital_asset_delivery"}, []string{"creator-commerce", "digital_asset_delivery"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got map[string]any
			status, err := Process(tc.job, func(payload map[string]any) error { got = payload; return nil })
			if err != nil || status != "captured" {
				t.Fatalf("process result: %q %v", status, err)
			}
			if !reflect.DeepEqual(got["fingerprint"], tc.want) {
				t.Fatalf("fingerprint = %#v, want %#v", got["fingerprint"], tc.want)
			}
		})
	}
}
