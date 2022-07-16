package auth

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
)

func TestLoginHandler_Handle(t *testing.T) {
	type fields struct {
		repository      *MockRepository
		tokenGenerator  *MockTokenHelper
		passwordChecker *MockPasswordChecker
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
				f.repository.EXPECT().FindByEmail("hello@gmail.com").Return(&User{id, "mathew", "hello@gmail.com", "passwordhash"}, nil)
				f.passwordChecker.EXPECT().Check("password", "passwordhash").Return(true)
				f.tokenGenerator.EXPECT().Generate(gomock.Any()).Return(Token("hello@gmail.compassword"))
			},
			"hello@gmail.compassword",
			false,
		},
		{"invalid password",
			args{LoginCommand{Email: "hello@gmail.com", Password: "passt"}}, func(f *fields) {
				f.repository.EXPECT().FindByEmail("hello@gmail.com").Return(&User{uuid.New(), "mathew", "hello@gmail.com", "password"}, nil)
				f.passwordChecker.EXPECT().Check("passt", "password").Return(false)
			}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			f := fields{
				repository:      NewMockRepository(ctrl),
				tokenGenerator:  NewMockTokenHelper(ctrl),
				passwordChecker: NewMockPasswordChecker(ctrl),
			}
			tt.prepare(&f)
			h := NewLoginHandler(f.repository, f.tokenGenerator, f.passwordChecker)
			got, err := h.Handle(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoginHandler.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
