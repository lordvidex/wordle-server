package adapters

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/o1egl/paseto"

	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/game"
)

var (
	payloadAudience = "wordle-wf-players"
	payloadIssuer   = "wordle-wf"
)

type Payload struct {
	game.Player
	paseto.JSONToken
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	var aggregator map[string]interface{}
	res, err := json.Marshal(p.JSONToken)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &aggregator)
	if err != nil {
		return nil, err
	}
	res, err = json.Marshal(p.Player)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &aggregator)
	if err != nil {
		return nil, err
	}
	return json.Marshal(aggregator)
}

func (p *Payload) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &p.Player)
	if err != nil {
		return err
	}
	// remove values that are not string because JSONToken uses a map[string]string
	temp := make(map[string]interface{})
	err = json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	for k, v := range temp {
		if _, ok := v.(string); !ok {
			delete(temp, k)
		}
	}
	b, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &p.JSONToken)
	if err != nil {
		return err
	}
	return nil
}

func NewPayload(player *game.Player, iat time.Duration) Payload {
	currentTime := time.Now()
	return Payload{
		Player: *player,
		JSONToken: paseto.JSONToken{
			Audience:   payloadAudience,
			IssuedAt:   currentTime,
			Subject:    player.ID.String(),
			Issuer:     payloadIssuer,
			Jti:        player.ID.String() + fmt.Sprintf("%d", time.Now().Unix()),
			NotBefore:  currentTime,
			Expiration: currentTime.Add(iat),
		},
	}
}

type pasetoTokenHelper struct {
	secret        string
	expiresBefore time.Duration
}

func (p *pasetoTokenHelper) Generate(payload *game.Player) (auth.Token, error) {
	tokenPayload := NewPayload(payload, p.expiresBefore)
	tokenString, err := paseto.NewV2().Encrypt([]byte(p.secret), tokenPayload, nil)
	if err != nil {
		return "", nil
	}
	return auth.Token(tokenString), nil
}

func (p *pasetoTokenHelper) Decode(token auth.Token, payload *game.Player) error {
	var tokenPayload Payload
	err := paseto.NewV2().Decrypt(string(token), []byte(p.secret), &tokenPayload, nil)
	if err != nil {
		return err
	}
	err = tokenPayload.Validate()
	if err != nil {
		return err
	}
	*payload = tokenPayload.Player
	return nil
}

func NewPASETOTokenHelper(secret string, expiresBefore time.Duration) auth.TokenHelper {
	if len([]byte(secret)) != 32 {
		panic("secret must be 32 bytes")
	}
	return &pasetoTokenHelper{secret, expiresBefore}
}
