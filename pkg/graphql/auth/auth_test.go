package auth_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/config"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/graphql/auth"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/graphql/errors"

	mock_config "github.com/jpdejavite/rtg-go-toolkit/mock/config"
)

const (
	gatewayPublicKey  = "-----BEGIN PUBLIC KEY-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAQXqqbEeiSK6d27LKcbNusbIUL+mn\nrMRbWx5ZzWLLJgSBUntTEb+GEDQB6vzjEEE4x033bbFMLv+eWFpbjJCwnIMBBpQO\nI9gO61dqPnaQLpnsFmHAeGRsBRif9zULvEbteTEstzMRKXP5eNzhPkmNfXT2sA6/\nOy7hTo82fAcNCEWK1uk=\n-----END PUBLIC KEY-----"
	gatewayPrivateKey = "-----BEGIN EC PRIVATE KEY-----\nMIHcAgEBBEIB7N7HkaB+pXzBlsSt+SIWd4IOpkT2ggax+rM7WqJqULBhjdU1LzSl\nzkrLMT9eWb0rI/urTZ/rh7aoYSKO0jgCe+GgBwYFK4EEACOhgYkDgYYABABBeqps\nR6JIrp3bsspxs26xshQv6aesxFtbHlnNYssmBIFSe1MRv4YQNAHq/OMQQTjHTfdt\nsUwu/55YWluMkLCcgwEGlA4j2A7rV2o+dpAumewWYcB4ZGwFGJ/3NQu8Ru15MSy3\nMxEpc/l43OE+SY19dPawDr87LuFOjzZ8Bw0IRYrW6Q==\n-----END EC PRIVATE KEY-----"
)

func TestValidateHasAllRolesNilRoles(t *testing.T) {
	claims := jwt.MapClaims{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.AuthorizationDataKey, claims)
	err := auth.ValidateHasAllRoles(ctx, []*string{})
	if diff := deep.Equal(err, errors.NotAuthorizedError); diff != nil {
		t.Error(diff)
	}
}

func TestValidateHasAllRolesEmptyRoles(t *testing.T) {
	claims := jwt.MapClaims{
		"roles": []interface{}{},
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.AuthorizationDataKey, claims)
	err := auth.ValidateHasAllRoles(ctx, []*string{})
	if diff := deep.Equal(err, errors.NotAuthorizedError); diff != nil {
		t.Error(diff)
	}
}

func TestValidateHasAllRolesClaimsWithoutAllRoles(t *testing.T) {
	adminRole := "admin"
	superAdmin := "super_admin"
	claims := jwt.MapClaims{
		"roles": []interface{}{adminRole},
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.AuthorizationDataKey, claims)
	err := auth.ValidateHasAllRoles(ctx, []*string{&adminRole, &superAdmin})
	if diff := deep.Equal(err, errors.NotAuthorizedError); diff != nil {
		t.Error(diff)
	}
}

func TestValidateHasAllRolesClaimsWithAllRoles(t *testing.T) {
	adminRole := "admin"
	superAdmin := "super_admin"
	claims := jwt.MapClaims{
		"roles": []interface{}{superAdmin, adminRole, "extra_role"},
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.AuthorizationDataKey, claims)
	err := auth.ValidateHasAllRoles(ctx, []*string{&adminRole, &superAdmin})
	if err != nil {
		t.Errorf("No error should be given, all roles should match: %v", err)
	}
}

func TestAddSecurityHandlerNoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	configMock := mock_config.NewMockIGlobalConfigs(ctrl)

	req, err := http.NewRequest("POST", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := auth.AddSecurityHandler(configMock)(nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := strings.ReplaceAll(rr.Body.String(), "\n", "")
	expect := strings.ReplaceAll(errors.NewGraphqlErrorToJSON("invalid_request", "Invalid request"), "\n", "")
	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestAddSecurityHandlerErrorValidatingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	configMock := mock_config.NewMockIGlobalConfigs(ctrl)

	configMock.EXPECT().
		GetGlobalConfigAsStr(config.GatewayPublicKey).
		Return("ble")

	req, err := http.NewRequest("POST", "/validate", nil)
	req.Header.Set(auth.GatewayTokenHeader, "bla")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := auth.AddSecurityHandler(configMock)(nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := strings.ReplaceAll(rr.Body.String(), "\n", "")
	expect := strings.ReplaceAll(errors.NewGraphqlErrorToJSON("invalid_gateway_token", "Invalid gateway token"), "\n", "")
	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestAddSecurityHandlerInvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	configMock := mock_config.NewMockIGlobalConfigs(ctrl)

	configMock.EXPECT().
		GetGlobalConfigAsStr(config.GatewayPublicKey).
		Return(gatewayPublicKey)

	req, err := http.NewRequest("POST", "/validate", nil)
	req.Header.Set(auth.GatewayTokenHeader, "bla")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := auth.AddSecurityHandler(configMock)(nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := strings.ReplaceAll(rr.Body.String(), "\n", "")
	expect := strings.ReplaceAll(errors.NewGraphqlErrorToJSON("invalid_gateway_token", "Invalid gateway token"), "\n", "")
	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestAddSecurityHandlerValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	configMock := mock_config.NewMockIGlobalConfigs(ctrl)

	configMock.EXPECT().
		GetGlobalConfigAsStr(config.GatewayPublicKey).
		Return(gatewayPublicKey)

	req, err := http.NewRequest("POST", "/validate", nil)
	req.Header.Set(auth.GatewayTokenHeader, generateToken())
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := auth.AddSecurityHandler(configMock)(oKHandler{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := strings.ReplaceAll(rr.Body.String(), "\n", "")
	expect := strings.ReplaceAll("OK", "\n", "")
	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

type oKHandler struct{}

func (okHandler oKHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "OK")
}

func generateToken() string {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.StandardClaims{})
	// Create the JWT string

	key, err := jwt.ParseECPrivateKeyFromPEM([]byte(gatewayPrivateKey))
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString
}
