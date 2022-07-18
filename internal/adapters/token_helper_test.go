package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/game"
	"testing"
	"time"
)

var (
	testUUID = uuid.New()
)

func getSecret() string {
	var secret string
	for i := 0; i < 32; i++ {
		secret += "a"
	}
	return secret
}

func getPlayer() *game.Player {
	return &game.Player{
		ID:        testUUID,
		Name:      "John",
		Email:     "john@test.com",
		Password:  "password",
		Points:    120,
		IsDeleted: false,
	}
}
func TestPASETOTokenHelper_GenerateDoesNotThrow(t *testing.T) {

	tokenHelper := NewPASETOTokenHelper(getSecret(), time.Minute)
	token, err := tokenHelper.Generate(getPlayer())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("Expected token, got empty string")
	}
	fmt.Println("Encoded token is ", token)
}

func generateToken(helper auth.TokenHelper, t *testing.T) auth.Token {
	res, err := helper.Generate(getPlayer())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	return res
}

func TestPasetoTokenHelper_Decode(t *testing.T) {

	tokenHelper := NewPASETOTokenHelper(getSecret(), time.Minute)
	token := generateToken(tokenHelper, t)

	var decodedPlayer game.Player
	err := tokenHelper.Decode(token, &decodedPlayer)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	fmt.Println("Decoded player is ", decodedPlayer)

	if decodedPlayer.Name != "John" {
		t.Errorf("Expected payload.Name to be 'John', got %v", decodedPlayer.Name)
	}
	if decodedPlayer.Email != "john@test.com" {
		t.Errorf("Expected payload.Age to be 30, got %T", decodedPlayer.Email)
	}
	if decodedPlayer.Points != 120 {
		t.Errorf("Expected payload.Points to be 120, got %v", decodedPlayer.Points)
	}
	if decodedPlayer.IsDeleted != false {
		t.Errorf("Expected payload.IsDeleted to be false, got %v", decodedPlayer.IsDeleted)
	}
	if decodedPlayer.Password != "" {
		t.Errorf("Expected payload.Password to be empty, got %v", decodedPlayer.Password)
	}
}

func TestPayload_UnmarshalJSON(t *testing.T) {
	var payload Payload
	js := `{
"Email": "john@test.com",
"ID": "e251a3b9-e792-41dd-ba70-ba036f14c88b",
"IsDeleted": false,
"Name":"John",
"Points":120,
"aud":"wordle-wf-players",
"exp":"2022-07-18T06:29:01+04:00",
"iat":"2022-07-18T06:28:01+04:00",
"iss":"wordle-wf",
"jti":"e251a3b9-e792-41dd-ba70-ba036f14c88b1658111281",
"nbf":"2022-07-18T06:28:01+04:00",
"sub":"e251a3b9-e792-41dd-ba70-ba036f14c88b"
}`
	err := json.Unmarshal([]byte(js), &payload)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if payload.Email != "john@test.com" {
		t.Errorf("Expected payload.Email to be '%s'", payload.Player.Email)
	}
	if payload.Issuer == "" {
		t.Errorf("Expected payload.Issuer to be '%s' got %v", "wordle-wf", payload.Issuer)
	}
}
