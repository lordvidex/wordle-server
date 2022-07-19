package auth

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
)

func TestLoginHandler_Handle(t *testing.T) {
	type fields struct {
		repository      *MockRepository
		tokenGenerator  *MockTokenHelper
		passwordChecker *MockPasswordHelper
	}

	type args struct {
		command LoginCommand
	}

	// test cases
	tests := []struct {
		name    string
		args    args
		prepare func(*fields)
		want    Token
		wantErr bool
	}{
		{"happy path",
			args{LoginCommand{Email: "hello@gmail.com", Password: "password"}},
			func(f *fields) {
				id := uuid.New()
				f.repository.EXPECT().FindByEmail("hello@gmail.com").Return(&game.Player{ID: id, Name: "mathew", Email: "hello@gmail.com", Password: "passwordhash"}, nil)
				f.passwordChecker.EXPECT().Validate("password", "passwordhash").Return(true)
				f.tokenGenerator.EXPECT().Generate(gomock.Any()).Return(Token("hello@gmail.compassword"), nil)
			},
			"hello@gmail.compassword",
			false,
		},
		{"invalid password",
			args{LoginCommand{Email: "hello@gmail.com", Password: "passt"}}, func(f *fields) {
				f.repository.EXPECT().FindByEmail("hello@gmail.com").Return(&game.Player{ID: uuid.New(), Name: "mathew", Email: "hello@gmail.com", Password: "password"}, nil)
				f.passwordChecker.EXPECT().Validate("passt", "password").Return(false)
			}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			f := fields{
				repository:      NewMockRepository(ctrl),
				tokenGenerator:  NewMockTokenHelper(ctrl),
				passwordChecker: NewMockPasswordHelper(ctrl),
			}
			tt.prepare(&f)
			h := NewLoginHandler(f.repository, f.tokenGenerator, f.passwordChecker)
			got, err := h.Handle(tt.args.command)
			// error cases
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Errorf("expected a token with a player, got nil")
			}
			if got != nil && got.Token != tt.want {
				t.Errorf("expected token %v, got %v", tt.want, got.Token)
			}
		})
	}
}
