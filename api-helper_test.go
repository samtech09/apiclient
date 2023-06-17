package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

var jwtapi *JwtAPI
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
	http.HandleFunc("/protected-exp", func(w http.ResponseWriter, r *http.Request) {
		//token := extractToken(r)
		w.WriteHeader(http.StatusUnauthorized)
		//json.NewEncoder(w).Encode(APIResult{Data: "protected-not-ok"})
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
	url := "/get"
	api := API{AllowInsecureSSL: true, StructuredResponse: true}
	api.ResourceAPIBaseURL = "http://localhost:8080"
	res, err := api.Get(url)
	if err != nil {
		t.Errorf("APIGet error: %v\n", err)
	}
	exp := "get-ok"
	if res.Data != exp {
		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
	}
}

func TestAPIGetBasicAuth(t *testing.T) {
	url := "/orders/order_FgP6jhvCOWM1Hk/payments"
	api := API{AllowInsecureSSL: true, StructuredResponse: false, UseBasicAuth: true, BasicAuthUser: "rzp_test_ZlzD0ybRfBYmC5", BasicAuthPwd: "vicl3kdxVRmMaMj1w4qmsFNL"}
	api.ResourceAPIBaseURL = "https://api.razorpay.com/v1"
	res, err := api.Get(url)
	if err != nil {
		t.Errorf("APIGet error: %v\n", err)
	}
	exp := "get-ok"
	if res.Data != exp {
		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
	}
}

// func TestAPIPostForm(t *testing.T) {
// 	url := "https://api.secure.ebs.in/api/1_0/statusByRef"
// 	api := API{AllowInsecureSSL: true}
// 	api.ResourceAPIBaseURL = "https://api.secure.ebs.in/api/1_0"

// 	data := make(map[string]string)
// 	data["Action"] = "statusByRef"
// 	data["AccountID"] = "1111"
// 	data["SecretKey"] = "asasasasasasas"
// 	data["RefNo"] = "1212121212"

// 	res, err := api.APIPostForm(url, data)
// 	if err != nil {
// 		t.Errorf("APIPostForm error: %v\n", err)
// 	}

// 	fmt.Println("PostForm: ", res)

// 	exp := "<output"
// 	if !strings.Contains(res.Data, exp) {
// 		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
// 	}
// }

func TestGetToken(t *testing.T) {
	jwtapi = &JwtAPI{StructuredResponse: true}
	jwtapi.ResourceAPIBaseURL = "http://localhost:8080"
	jwtapi.TokenURI = "http://localhost:8080/token"
	jwtapi.RefreshTokenURI = "http://localhost:8080/refresh"
	jwtapi.TokenRequestData.ClientID = "test0-client"
	jwtapi.TokenRequestData.ClientSecret = "123456678ABCCD"
	jwtapi.TokenRequestData.Scopes = "user"
	jwtapi.Debug = true
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
	url := "/protected"
	res, err := jwtapi.Get(url)
	if err != nil {
		t.Errorf("APIGet error: %v\n", err)
	}
	exp := "protected-ok"
	if res.Data != exp {
		t.Errorf("Expected: %s,  Got: %s", exp, res.Data)
	}
}

func TestProtectedWithJWTExp(t *testing.T) {
	url := "/protected-exp"
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
