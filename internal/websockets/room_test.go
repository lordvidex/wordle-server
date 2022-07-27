package websockets

import (
	"reflect"
	"testing"
)

func Test_convertToPayload(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		want    WSPayload
		wantErr bool
	}{
		{
			name: "convert PlayerJoined",
			data: []byte(`{"event":"PlayerJoined","data":{"player_id":"123","player_name":"test"}}`),
			want: WSPayload{
				Event: EventPlayerJoined,
				Data: map[string]any{
					"player_id":   "123",
					"player_name": "test",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid type",
			data:    []byte(`{"json121": "test"}`),
			want:    WSPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalPayload(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
