package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

var jwtapi *SJwtAPI
var token Token

var wg sync.WaitGroup

func mockServer() {
	// create mock server and endpoints for serving
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(APIResult{Data: "get-ok"})
	})
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Token{AccessToken: "access-test-token"})
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Token{AccessToken: "access-test-token"})
	})
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "access-test-token" {
			json.NewEncoder(w).Encode(APIResult{Data: "protected-ok"})
			return
		}
		json.NewEncoder(w).Encode(APIResult{Data: "protected-not-ok"})
	})
	wg.Add(1)
	http.ListenAndServe(":8080", nil)
	wg.Wait()
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "Authorization header format must be Bearer {token}"
	}
	return authHeaderParts[1]
}

func TestInit(t *testing.T) {
	go mockServer()
}
func TestAPIGet(t *testing.T) {
	url := "http://localhost:8080/get"
	api := SAPI{AllowInsecureSSL: true}
	res, err := api.Get(url)
	if err != nil {
		t.Errorf("APIGet error: %v\n", err)
	}
	exp := "get-ok"
	if res.Data != exp {
		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
	}
}

func TestGetToken(t *testing.T) {
	jwtapi = &SJwtAPI{}
	jwtapi.TokenURI = "http://localhost:8080/token"
	jwtapi.RefreshTokenURI = "http://localhost:8080/refresh"
	jwtapi.TokenRequestData.ClientID = "test0-client"
	jwtapi.TokenRequestData.ClientSecret = "123456678ABCCD"
	jwtapi.TokenRequestData.Scopes = "user"
	token, err := jwtapi.RequestTokenByCred()
	if err != nil {
		t.Errorf("Get-Token error: %v\n", err)
	}
	exp := "access-test-token"
	if token.AccessToken != exp {
		t.Errorf("Expected: %s,   Got: %s", exp, token.AccessToken)
	}
}

func TestRequestTokenByRefreshToken(t *testing.T) {
	rtoken := "test-rtoken"
	token, err := jwtapi.RequestTokenByRefreshToken(rtoken)
	if err != nil {
		fmt.Printf("Get-Refresh-Token error: %v\n", err)
	}
	exp := "access-test-token"
	if token.AccessToken != exp {
		t.Errorf("Expected: %s,   Got: %s", exp, token.AccessToken)
	}
}

func TestProtectedWithJWT(t *testing.T) {
	url := "http://localhost:8080/protected"
	res, err := jwtapi.Get(url)
	if err != nil {
		t.Errorf("APIGet error: %v\n", err)
	}
	exp := "protected-ok"
	if res.Data != exp {
		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
	}
}

func TestClose(t *testing.T) {
	wg.Done()
}
