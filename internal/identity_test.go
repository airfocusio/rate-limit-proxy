package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtBearerTokenIdentifier(t *testing.T) {
	newTestRequest := func(token string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		}
		return req
	}

	identifierHS256NoKeyID := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}
	identifierHS256UnknownKeyID := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "unknown",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}
	identifierHS256 := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "1",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}
	identifierHS256WithClientID := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "1",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub", "client_id"}),
	}
	identifierRS256 := JwtIdentifier{
		Algorithm:         "RS256",
		KeyID:             "2",
		Verifier:          rs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}
	identifierES256 := JwtIdentifier{
		Algorithm:         "ES256",
		KeyID:             "3",
		Verifier:          es256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}

	if _, err := identifierHS256NoKeyID.IdentifyRequest(newTestRequest("")); assert.Error(t, err, "bearer token header is missing") {
	}

	if _, err := identifierHS256UnknownKeyID.IdentifyRequest(newTestRequest(jwtHs256User1)); assert.Error(t, err, "jwt kid header 1 does not match") {
	}

	if id, err := identifierHS256.IdentifyRequest(newTestRequest(jwtHs256User1)); assert.NoError(t, err) {
		assert.Equal(t, "user:1", *id)
	}

	if id, err := identifierRS256.IdentifyRequest(newTestRequest(jwtRs256User3)); assert.NoError(t, err) {
		assert.Equal(t, "user:3", *id)
	}

	if id, err := identifierES256.IdentifyRequest(newTestRequest(jwtEs256User5)); assert.NoError(t, err) {
		assert.Equal(t, "user:5", *id)
	}

	if _, err := identifierHS256NoKeyID.IdentifyRequest(newTestRequest(jwtHs256User7Expired)); assert.Error(t, err, "Token is expired") {
	}

	if _, err := identifierHS256NoKeyID.IdentifyRequest(newTestRequest(jwtHs256User8Invalid)); assert.Error(t, err, "signature is invalid") {
	}

	if id, err := identifierHS256.IdentifyRequest(newTestRequest(jwtHs256User9WithClientId)); assert.NoError(t, err) {
		assert.Equal(t, "user:9", *id)
	}
	if id, err := identifierHS256WithClientID.IdentifyRequest(newTestRequest(jwtHs256User9WithClientId)); assert.NoError(t, err) {
		assert.Equal(t, "user:9:client:1", *id)
	}
}

func TestJwtQueryParameterIdentifier(t *testing.T) {
	newTestRequest := func(token string) *http.Request {
		url := "/"
		if token != "" {
			url = fmt.Sprintf("%s?token=%s", url, token)
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		return req
	}

	identifier := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "1",
		Verifier:          hs256,
		TokenExtractor:    ExtractQueryParameter("token"),
		IdentityExtractor: ExtractClaim([]string{"sub"}),
	}

	_, err := identifier.IdentifyRequest(newTestRequest(""))
	if err == nil || err.Error() != "query parameter token is missing" {
		t.Errorf("expected error, got %v", err)
	}

	subject1, err := identifier.IdentifyRequest(newTestRequest(jwtHs256User1))
	if err != nil {
		t.Errorf("expected subject, got %v", err)
	} else if *subject1 != "user:1" {
		t.Errorf("expected subject, got %v", *subject1)
	}
}
