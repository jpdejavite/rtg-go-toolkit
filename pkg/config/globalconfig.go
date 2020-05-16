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

// IGlobalConfigs global configs interface
type IGlobalConfigs interface {
	GetGlobalKeys() []string
	LoadGlobalConfig() error
	GetGlobalConfigAsInt(key string) int
	GetGlobalConfigAsStr(key string) string
}

// NewGlobalConfigs returns a new global interface
func NewGlobalConfigs(db firestore.IDBFirestore) IGlobalConfigs {
	return GlobalConfigs{
		configs: make(map[string]interface{}),
		db:      db,
	}
}

// GlobalConfigs implements IDBFirestore interface
type GlobalConfigs struct {
	configs map[string]interface{}
	db      firestore.IDBFirestore
}

// GetGlobalKeys return list of global config keys
func (gc GlobalConfigs) GetGlobalKeys() []string {
	return []string{GatewayPublicKey, TokenExpirationInMinutes}
}

// LoadGlobalConfig load all global configs
func (gc GlobalConfigs) LoadGlobalConfig() error {
	globalConfigData, err := gc.db.GetDocumentData("configs", "global")
	if err != nil {
		return err
	}
	if globalConfigData == nil {
		return errors.New("no data in global config")
	}

	for _, k := range gc.GetGlobalKeys() {
		gc.configs[k] = globalConfigData[k]
		if gc.configs[k] == nil || gc.configs[k] == "" {
			return fmt.Errorf("missing global config %s", k)
		}
	}
	return nil
}

// GetGlobalConfigAsInt get global config as int
func (gc GlobalConfigs) GetGlobalConfigAsInt(key string) int {
	val := gc.configs[key]
	if val == nil {
		return 0
	}
	return val.(int)
}

// GetGlobalConfigAsStr get global config as string
func (gc GlobalConfigs) GetGlobalConfigAsStr(key string) string {
	val := gc.configs[key]
	if val == nil {
		return ""
	}
	return val.(string)
}
