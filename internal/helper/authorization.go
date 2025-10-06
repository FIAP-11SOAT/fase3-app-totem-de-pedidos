package helper

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CognitoClaims struct {
	Sub           string   `json:"sub,omitempty"`
	CognitoGroups []string `json:"cognito:groups,omitempty"`
	Iss           string   `json:"iss,omitempty"`
	ClientID      string   `json:"client_id,omitempty"`
	OriginJTI     string   `json:"origin_jti,omitempty"`
	EventID       string   `json:"event_id,omitempty"`
	TokenUse      string   `json:"token_use,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	AuthTime      int64    `json:"auth_time,omitempty"`
	Exp           int64    `json:"exp,omitempty"`
	Iat           int64    `json:"iat,omitempty"`
	JTI           string   `json:"jti,omitempty"`
	Username      string   `json:"username,omitempty"`
	jwt.RegisteredClaims
}

type Validator interface {
	Validate(tokenString string, keyMap map[string]*rsa.PublicKey) (*CognitoClaims, error)
}

type JWTValidator struct{}

func NewJWTValidator() Validator {
	return &JWTValidator{}
}

func (v *JWTValidator) Validate(tokenString string, keyMap map[string]*rsa.PublicKey) (*CognitoClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CognitoClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("token does not contain kid")
		}

		pubKey, exists := keyMap[kid]
		if !exists {
			return nil, fmt.Errorf("no key found for kid: %s", kid)
		}

		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate JWT: %w", err)
	}

	if claims, ok := token.Claims.(*CognitoClaims); ok && token.Valid {
		if claims.TokenUse != "access" {
			return nil, fmt.Errorf("invalid token_use: got %s, expected access", claims.TokenUse)
		}
		if claims.Exp < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
