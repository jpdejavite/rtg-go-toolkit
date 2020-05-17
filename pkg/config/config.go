package config

import (
	"errors"
	"fmt"

	"github.com/jpdejavite/rtg-go-toolkit/pkg/firestore"
)

// IConfigs configs interface
type IConfigs interface {
	LoadConfig(app string, keys []string) error
	GetConfigAsInt(key string) int
	GetConfigAsStr(key string) string
}

// NewConfigs returns a new  interface
func NewConfigs(db firestore.IDBFirestore) IConfigs {
	return Configs{
		configs: make(map[string]interface{}),
		db:      db,
	}
}

// Configs implements IDBFirestore interface
type Configs struct {
	configs map[string]interface{}
	db      firestore.IDBFirestore
}

// LoadConfig load all  configs
func (c Configs) LoadConfig(app string, keys []string) error {
	ConfigData, err := c.db.GetDocumentData("configs", app)
	if err != nil {
		return err
	}
	if ConfigData == nil {
		return errors.New("no data in config")
	}

	for _, k := range keys {
		c.configs[k] = ConfigData[k]
		if c.configs[k] == nil || c.configs[k] == "" {
			return fmt.Errorf("missing config %s", k)
		}
	}
	return nil
}

// GetConfigAsInt get  config as int
func (c Configs) GetConfigAsInt(key string) int {
	val := c.configs[key]
	if val == nil {
		return 0
	}
	return val.(int)
}

// GetConfigAsStr get  config as string
func (c Configs) GetConfigAsStr(key string) string {
	val := c.configs[key]
	if val == nil {
		return ""
	}
	return val.(string)
}
