package internal

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Identifier ...
type Identifier interface {
	IdentifyRequest(*http.Request) (*string, error)
}

// IdentifierConfigJwtBearerHeader ...
type IdentifierConfigJwtBearerHeader struct {
	Algorithm string `yaml:"algorithm"`
	KeyID     string `yaml:"keyId"`
	Verifier  string `yaml:"verifier"`
	Claim     string `yaml:"claim"`
}

// IdentifierConfigJwtQueryParameter ...
type IdentifierConfigJwtQueryParameter struct {
	Algorithm string `yaml:"algorithm"`
	KeyID     string `yaml:"keyId"`
	Verifier  string `yaml:"verifier"`
	Claim     string `yaml:"claim"`
	Name      string `yaml:"name"`
}

// IdentifierConfig ...
type IdentifierConfig struct {
	JwtBearerHeader *IdentifierConfigJwtBearerHeader   `yaml:"jwtBearerHeader"`
	JwtQueryToken   *IdentifierConfigJwtQueryParameter `yaml:"jwtQueryParameter"`
}

// ExtractBearerToken ...
func ExtractBearerToken(req http.Request) (*string, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return nil, fmt.Errorf("bearer token header is missing")
	}
	tokenStr := authorizationHeader[7:]
	return &tokenStr, nil
}

// ExtractQueryParameter ...
func ExtractQueryParameter(parameterName string) func(http.Request) (*string, error) {
	return func(req http.Request) (*string, error) {
		parameterValue := req.URL.Query().Get(parameterName)
		if parameterValue == "" {
			return nil, fmt.Errorf("query parameter %s is missing", parameterName)
		}
		tokenStr := parameterValue
		return &tokenStr, nil
	}
}

// ExtractClaim ...
func ExtractClaim(claimName string) func(claims jwt.MapClaims) (*string, error) {
	return func(claims jwt.MapClaims) (*string, error) {
		claimValue, ok := claims[claimName].(string)
		if !ok {
			return nil, fmt.Errorf("jwt %s claim is missing", claimName)
		}
		return &claimValue, nil
	}
}

// JwtIdentifier ...
type JwtIdentifier struct {
	KeyID             string
	Algorithm         string
	Verifier          string
	TokenExtractor    func(http.Request) (*string, error)
	IdentityExtractor func(claims jwt.MapClaims) (*string, error)
}

// IdentifyRequest ...
func (id JwtIdentifier) IdentifyRequest(req *http.Request) (*string, error) {
	tokenStr, err := id.TokenExtractor(*req)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(*tokenStr, func(token *jwt.Token) (interface{}, error) {
		keyID, ok := token.Header["kid"].(string)
		if id.KeyID != "" {
			if !ok || id.KeyID != keyID {
				return nil, fmt.Errorf("jwt kid header %s does not match", keyID)
			}
		}
		if strings.HasPrefix(id.Algorithm, "HS") {
			return base64.StdEncoding.DecodeString(id.Verifier)
		} else if strings.HasPrefix(id.Algorithm, "RS") {
			return jwt.ParseRSAPublicKeyFromPEM([]byte(id.Verifier))
		} else if strings.HasPrefix(id.Algorithm, "ES") {
			return jwt.ParseECPublicKeyFromPEM([]byte(id.Verifier))
		} else {
			return nil, fmt.Errorf("unsupported algorithm %s", id.Algorithm)
		}
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("jwt is invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt is invalid")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix()+10, false) {
		return nil, fmt.Errorf("jwt has expired")
	}
	return id.IdentityExtractor(claims)
}
