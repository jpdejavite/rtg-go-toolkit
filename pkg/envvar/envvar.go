package envvar

import (
	"fmt"
	"os"
)

var envvars map[string]string

// LoadAll load all environment variable
func LoadAll(keys []string) {
	envvars = make(map[string]string)
	for _, k := range keys {
		envvars[k] = os.Getenv(k)
		if envvars[k] == "" {
			panic(fmt.Sprintf("missing env var %s", k))
		}
	}
}

// GetEnvVar get env var loaded
func GetEnvVar(key string) string {
	return envvars[key]
}
