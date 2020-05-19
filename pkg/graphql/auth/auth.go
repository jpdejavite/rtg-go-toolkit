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

type key int

const (
	// GatewayTokenHeader gateway token header
	GatewayTokenHeader = "gateway-token"
	// AppCoiHeader app coi header
	AppCoiHeader = "app-coi"
	// AuthorizationDataKey authorization data key in context
	AuthorizationDataKey key = iota
	// GlobalConfigsKey global config key in context
	GlobalConfigsKey key = iota
	// ConfigsKey config key in context
	ConfigsKey key = iota
	// AppCoi app coi key  in context
	AppCoi key = iota
)

// Data readed from gateway header
type Data struct {
	Service string
	Roles   []string
	Uid     string
	Iss     string
}

/*AddSecurityHandler extracts security credentials sent by gateway from header
and put into request context app using a standard struct */
func AddSecurityHandler(gc config.IGlobalConfigs, c config.IConfigs) func(http.Handler) http.Handler {

	addCtx := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(GatewayTokenHeader) == "" {
				http.Error(w, errors.NewGraphqlErrorToJSON("invalid_request", "Invalid request"), 400)
				return
			}

			gatewayPublicKey := gc.GetGlobalConfigAsStr(config.GatewayPublicKey)
			gatewayPublicKey = strings.ReplaceAll(gatewayPublicKey, "\\n", "\n")

			gatewayToken := r.Header.Get(GatewayTokenHeader)
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
			data := Data{}
			for key, val := range claims {
				switch key {
				case "service":
					data.Service = val.(string)
				case "iss":
					data.Iss = val.(string)
				case "uid":
					data.Uid = val.(string)
				case "roles":
					data.Roles = []string{}
					for _, v := range val.([]interface{}) {
						data.Roles = append(data.Roles, v.(string))
					}
				}
			}
			ctx := context.WithValue(r.Context(), AuthorizationDataKey, data)
			ctx = context.WithValue(ctx, GlobalConfigsKey, gc)
			ctx = context.WithValue(ctx, ConfigsKey, c)
			ctx = context.WithValue(ctx, AppCoi, r.Header.Get(AppCoiHeader))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return addCtx
}

// GetContextInfo extract info from context
func GetContextInfo(ctx context.Context) (*Data, *config.IGlobalConfigs, *config.IConfigs, string) {
	var data *Data
	if ctx.Value(AuthorizationDataKey) != nil {
		aux := ctx.Value(AuthorizationDataKey).(Data)
		data = &aux
	}
	var gc *config.IGlobalConfigs
	if ctx.Value(GlobalConfigsKey) != nil {
		aux := ctx.Value(GlobalConfigsKey).(config.IGlobalConfigs)
		gc = &aux
	}

	var c *config.IConfigs
	if ctx.Value(ConfigsKey) != nil {
		aux := ctx.Value(ConfigsKey).(config.IConfigs)
		c = &aux
	}

	var coi string
	if ctx.Value(AppCoi) != nil {
		coi = ctx.Value(AppCoi).(string)
	}
	return data, gc, c, coi
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

	authData := ctx.Value(AuthorizationDataKey).(Data)

	if authData.Roles == nil {
		return errors.NotAuthorizedError
	}

	if len(authData.Roles) == 0 {
		return errors.NotAuthorizedError
	}

	for _, reqRole := range roles {
		hasRole := false
		for _, userRole := range authData.Roles {
			if reqRole != nil && *reqRole == userRole {
				hasRole = true
			}
		}
		if !hasRole {
			return errors.NotAuthorizedError
		}
	}

	return nil
}
