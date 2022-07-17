package pg

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
)

func Test_userRepository_FindByID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Player
		wantErr bool
	}{
		{"some random id", args{uuid.New()}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockUserRepo.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
