package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jpdejavite/go-log/pkg/log"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/config"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/graphql/errors"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/model"
)

const (
	gatewayTokenHeader = "gateway-token"
)

/*AddSecurityHandler extracts security credentials sent by gateway from header
and put into request context app using a standard struct */
func AddSecurityHandler() func(http.Handler) http.Handler {

	addCtx := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(gatewayTokenHeader) == "" {
				http.Error(w, errors.NewGraphqlErrorToJSON("invalid_request", "Invalid request"), 400)
				return
			}

			gatewayPublicKey := config.GetGlobalConfigAsStr(config.GatewayPublicKey)
			gatewayPublicKey = strings.ReplaceAll(gatewayPublicKey, "\\n", "\n")

			gatewayToken := r.Header.Get(gatewayTokenHeader)
			gatewayToken = strings.ReplaceAll(gatewayToken, "Bearer ", "")

			val, err := verifyES512Token(gatewayToken, []byte(gatewayPublicKey))

			if err != nil {
				log.Error("AddSecurityHandler", "error verifiyng gateway token", model.NewMetaError(err), "coi")
				http.Error(w, errors.NewGraphqlErrorToJSON("invalid_gateway_token", "Invalid gateway token"), 400)
				return
			}
			if !val {
				log.Error("AddSecurityHandler", "invalid gateway token signature", nil, "coi")
				http.Error(w, errors.NewGraphqlErrorToJSON("invalid_gateway_token", "Invalid gateway token"), 400)
				return
			}

			tknPar, _ := jwt.Parse(gatewayToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(gatewayPublicKey), nil
			})
			claims := tknPar.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), "authorizationData", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return addCtx
}

func verifyES512Token(token string, publicKey []byte) (bool, error) {
	key, err := jwt.ParseECPublicKeyFromPEM(publicKey)
	if err != nil {
		return false, err
	}

	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		return false, nil
	}

	err = jwt.SigningMethodES512.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ValidateHasAllRoles validate has all roles directive
func ValidateHasAllRoles(ctx context.Context, roles []*string) error {

	reqCred := ctx.Value("authorizationData").(jwt.MapClaims)

	userRoles := reqCred["roles"].([]interface{})

	if userRoles == nil {
		return errors.New("Not_authorized", "Not authorized")
	}

	for _, reqRole := range roles {
		hasRole := false
		for _, userRole := range userRoles {
			if reqRole != nil && *reqRole == userRole {
				hasRole = true
			}
		}
		if !hasRole {
			return errors.New("Not_authorized", "Not authorized")
		}
	}

	return nil
}
