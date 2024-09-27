package notifier

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	v1 "github.com/reynn/notifier/gen/proto/notifier/v1"
	"github.com/reynn/notifier/internal/types"
)

func TestNotificationStatusToProto(t *testing.T) {
	tests := map[string]struct {
		s    types.NotificationStatus
		want v1.NotificationStatus
	}{
		"convert failed": {
			s:    "FAILED",
			want: v1.NotificationStatus_FAILED,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := NotificationStatusToProto(tt.s)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				// if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotificationStatusToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
