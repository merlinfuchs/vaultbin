package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var K = koanf.New(".")

const DefaultConfigName = "vaultbin.toml"
const envVarPrefix = "VBIN__"

var CfgFile string

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func InitConfig() {
	setupDefaults()

	if err := K.Load(file.Provider(DefaultConfigName), toml.Parser()); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error(fmt.Sprintf("Failed to load config file %s", DefaultConfigName))
			panic(nil)
		}
	}

	if err := K.Load(env.Provider(envVarPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envVarPrefix)), "_", ".", -1)
	}), nil); err != nil {
		slog.Error("Failed to load env vars")
		panic(nil)
	}
}
