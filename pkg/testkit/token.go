package testkit

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	clientID   = "5acaf304aa6cc50c77f7d228"
	clientName = "Moneylover Web"
	client     = "kHiZbFQOw5LV"
	userID     = "5c37cdd3fdc6be173620b466"
)

type requestClaim struct {
	Type       string    `json:"type"`
	Code       uuid.UUID `json:"code"`
	ClientID   string    `json:"clientId"`
	ClientName string    `json:"clientName"`
	Internal   bool      `json:"internal"`
	Client     string    `json:"client"`
	Iat        int64     `json:"iat"`
	Exp        int64     `json:"exp"`
}

type accessClaim struct {
	Type        string      `json:"type"`
	UserID      string      `json:"userId"`
	TokenDevice uuid.UUID   `json:"tokenDevice"`
	ClientID    string      `json:"clientId"`
	Client      string      `json:"client"`
	Scopes      interface{} `json:"scopes"`
	Iat         int64       `json:"iat"`
	Exp         int64       `json:"exp"`
}

// nolint: errcheck
func token(claim interface{}) string {
	signingKey := jose.SigningKey{Algorithm: jose.HS256, Key: []byte("")}
	signer, _ := jose.NewSigner(signingKey, (&jose.SignerOptions{}).WithType("JWT"))
	token, _ := jwt.Signed(signer).Claims(claim).CompactSerialize()

	return token
}

func requestToken() string {
	return token(requestClaim{
		Type:       "request-token",
		Code:       uuid.New(),
		ClientID:   clientID,
		ClientName: clientName,
		Internal:   true,
		Client:     client,
		Iat:        time.Now().Unix(),
		Exp:        time.Now().Add(time.Hour).Unix(),
	})
}

func accessToken() string {
	return token(accessClaim{
		Type:        "access-token",
		UserID:      userID,
		TokenDevice: uuid.New(),
		ClientID:    clientID,
		Client:      client,
		Scopes:      nil,
		Iat:         time.Now().Unix(),
		Exp:         time.Now().Add(time.Hour).Unix(),
	})
}

func refreshToken() string {
	return token(accessClaim{
		Type:        "refresh-token",
		UserID:      userID,
		TokenDevice: uuid.New(),
		ClientID:    clientID,
		Client:      client,
		Scopes:      nil,
		Iat:         time.Now().Unix(),
		Exp:         time.Now().Add(24 * time.Hour).Unix(),
	})
}
