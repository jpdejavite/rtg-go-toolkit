package config

import (
	"errors"
	"fmt"

	"github.com/jpdejavite/rtg-go-toolkit/pkg/firestore"
)

const (
	// GatewayPublicKey gateway public key
	GatewayPublicKey = "gatewayPublicKey"
	// TokenExpirationInMinutes token expiration in minutes
	TokenExpirationInMinutes = "tokenExpirationInMinutes"
)

var globalConfigs map[string]interface{}

// GetGlobalKeys return list of global config keys
func GetGlobalKeys() []string {
	return []string{GatewayPublicKey, TokenExpirationInMinutes}
}

// LoadGlobalConfig load all global configs
func LoadGlobalConfig() error {
	globalConfigData, err := firestore.GetDocumentData("configs", "global")
	if err != nil {
		return nil
	}
	if globalConfigData == nil {
		return errors.New("no data in global config")
	}

	globalConfigs = make(map[string]interface{})
	for _, k := range GetGlobalKeys() {
		globalConfigs[k] = globalConfigData[k]
		if globalConfigs[k] == nil || globalConfigs[k] == "" {
			return fmt.Errorf("missing env var %s", k)
		}
	}
	return nil
}

// GetGlobalConfigAsInt get global config as int
func GetGlobalConfigAsInt(key string) int {
	return globalConfigs[key].(int)
}

// GetGlobalConfigAsStr get global config as string
func GetGlobalConfigAsStr(key string) string {
	return globalConfigs[key].(string)
}
