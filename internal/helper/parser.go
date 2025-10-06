package helper

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
)

// JWK representa uma chave no JWKS.
type JWK struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

// JWKS representa o conjunto de chaves.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

func ParseJWKS(jwksJSON string) (map[string]*rsa.PublicKey, error) {
	var jwks JWKS
	if err := json.Unmarshal([]byte(jwksJSON), &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS JSON: %w", err)
	}

	keyMap := make(map[string]*rsa.PublicKey)
	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" || jwk.Use != "sig" || jwk.Alg != "RS256" {
			continue
		}
		pubKey, err := jwkToPublicKey(jwk)
		if err != nil {
			continue
		}
		keyMap[jwk.Kid] = pubKey
	}

	if len(keyMap) == 0 {
		return nil, fmt.Errorf("no valid RSA keys found in JWKS")
	}

	return keyMap, nil
}

func jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus for kid %s: %w", jwk.Kid, err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent for kid %s: %w", jwk.Kid, err)
	}

	e := 0
	for _, b := range eBytes {
		e = (e << 8) + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}
