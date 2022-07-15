package pg

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
	"testing"
)

func createGame(t *testing.T) *game.Game {
	g := &game.Game{
		ID:       uuid.New(),
		InviteID: "inviteid",
	}
	gm, err := mockGameRepo.Create(g)
	if err != nil {
		t.Errorf("Game Create() error = %v", err)
	}
	return gm
}
func Test_gameRepository_UpdateSettings(t *testing.T) {
	g := createGame(t)
	defer t.Cleanup(func() {
		_ = mockGameRepo.Delete(g.ID.String())
	})
	type args struct {
		settings *game.Settings
		gameID   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil settings", args{nil, g.ID.String()}, true},
		{"invalid uuid", args{&game.Settings{}, "invalid"}, true},
		{"not found", args{&game.Settings{}, uuid.New().String()}, true},
		{"no settings uses default values", args{func() *game.Settings { s := game.NewSettings(5); return &s }(), g.ID.String()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockGameRepo.UpdateSettings(tt.args.settings, tt.args.gameID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsEager(t *testing.T) {
	type args struct{}
	type args2 struct{}
	tests := []struct {
		name            string
		eagerType       [2]interface{}
		registeredTypes []interface{}
		want            bool
	}{
		{"nil",
			[2]interface{}{nil, nil},
			[]interface{}{nil},
			false,
		},
		{"empty",
			[2]interface{}{args{}, &args{}},
			[]interface{}{},
			false,
		},
		{"not found",
			[2]interface{}{args{}, &args{}},
			[]interface{}{args2{}},
			false,
		},
		{"found pointer",
			[2]interface{}{args{}, &args{}},
			[]interface{}{&args{}},
			true,
		},
		{
			"found struct",
			[2]interface{}{args{}, &args{}},
			[]interface{}{args{}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEager(tt.eagerType, tt.registeredTypes...); got != tt.want {
				t.Errorf("IsEager() = %v, want %v", got, tt.want)
			}
		})
	}
}
