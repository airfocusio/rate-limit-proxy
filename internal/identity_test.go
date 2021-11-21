package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
		IdentityExtractor: ExtractClaim("sub"),
	}
	identifierHS256UnknownKeyID := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "unknown",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim("sub"),
	}
	identifierHS256 := JwtIdentifier{
		Algorithm:         "HS256",
		KeyID:             "1",
		Verifier:          hs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim("sub"),
	}
	identifierRS256 := JwtIdentifier{
		Algorithm:         "RS256",
		KeyID:             "2",
		Verifier:          rs256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim("sub"),
	}
	identifierES256 := JwtIdentifier{
		Algorithm:         "ES256",
		KeyID:             "3",
		Verifier:          es256,
		TokenExtractor:    ExtractBearerToken,
		IdentityExtractor: ExtractClaim("sub"),
	}

	_, err := identifierHS256NoKeyID.IdentifyRequest(newTestRequest(""))
	if err == nil || err.Error() != "bearer token header is missing" {
		t.Errorf("expected error, got %v", err)
	}

	_, err = identifierHS256UnknownKeyID.IdentifyRequest(newTestRequest(jwtHs256User1))
	if err == nil || err.Error() != "jwt kid header 1 does not match" {
		t.Errorf("expected error, got %v", err)
	}

	subject1, err := identifierHS256.IdentifyRequest(newTestRequest(jwtHs256User1))
	if err != nil {
		t.Errorf("expected subject, got %v", err)
	} else if *subject1 != "user:1" {
		t.Errorf("expected subject, got %v", *subject1)
	}

	subject2, err := identifierRS256.IdentifyRequest(newTestRequest(jwtRs256User3))
	if err != nil {
		t.Errorf("expected subject, got %v", err)
	} else if *subject2 != "user:3" {
		t.Errorf("expected subject, got %v", *subject2)
	}

	subject3, err := identifierES256.IdentifyRequest(newTestRequest(jwtEs256User5))
	if err != nil {
		t.Errorf("expected subject, got %v", err)
	} else if *subject3 != "user:5" {
		t.Errorf("expected subject, got %v", *subject3)
	}

	_, err = identifierHS256.IdentifyRequest(newTestRequest(jwtHs256User7Expired))
	if err == nil || err.Error() != "Token is expired" {
		t.Errorf("expected error, got %v", err)
	}

	_, err = identifierHS256.IdentifyRequest(newTestRequest(jwtHs256User8Invalid))
	if err == nil || err.Error() != "signature is invalid" {
		t.Errorf("expected error, got %v", err)
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
		IdentityExtractor: ExtractClaim("sub"),
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
