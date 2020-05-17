package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jpdejavite/go-log/pkg/log"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/firestore"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/model"
)

const (
	// DefaultRefreshTimeoutInSeconds default config timetout to refresh configs from db (in seconds)
	DefaultRefreshTimeoutInSeconds = 30
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
		if os.Getenv(k) != "" {
			if c.configs[k] == nil {
				log.Info("config", "override config", OverrideMeta{k}, log.GenerateCoi(nil))
			}
			c.configs[k] = os.Getenv(k)
			continue
		}
		data := ConfigData[k]
		if data == nil || data == "" {
			return fmt.Errorf("missing config %s", k)
		}

		if c.configs[k] != data {
			log.Info("config", "setting config", OverrideMeta{k}, log.GenerateCoi(nil))
		}
		c.configs[k] = data
	}

	c.configs[RefreshConfigTimeoutInSeconds] = ConfigData[RefreshConfigTimeoutInSeconds]

	go c.refreshConfig(app, keys)
	return nil
}

func (c Configs) refreshConfig(app string, keys []string) {
	for {
		sleepSeconds := DefaultRefreshTimeoutInSeconds
		if c.GetConfigAsInt(RefreshConfigTimeoutInSeconds) != 0 {
			sleepSeconds = c.GetConfigAsInt(RefreshConfigTimeoutInSeconds)
		}
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
		if err := c.LoadConfig(app, keys); err != nil {
			log.Error("config", "error LoadConfig", model.NewMetaError(err), log.GenerateCoi(nil))
			break
		}

	}
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
