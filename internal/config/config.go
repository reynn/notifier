package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"

	"github.com/reynn/notifier/internal/constants"
	"github.com/reynn/notifier/internal/paths"
)

var (
	configFilePath string
	port           int
)

type (
	Settings struct {
		HTTPPort       int    `env:"HTTP_PORT" envDefault:"8080"`
		ConfigFilePath string `env:"CONFIG_FILE_PATH"`
		Telemetry      Telemetry
	}

	Telemetry struct {
		Enabled bool `env:"TELEMETRY_ENABLED" envDefault:"false"`
	}
)

func Load() *Settings {
	s := &Settings{}
	if err := env.Parse(s); err != nil {
		panic(err)
	}
	return s
}

func init() {
	flag.StringVar(&configFilePath, "config-file", "", "Path to the configuration file")
	flag.Set("config-file", paths.ConfigFile(constants.DefaultConfigFileExt))

	flag.IntVar(&port, "port", 8080, "Port to listen on")

	flag.Usage = func() {
		fmt.Printf("Usage: %s.%s:%s [options]\n", constants.AppName, constants.AppModule, constants.AppVersion)
		flag.PrintDefaults()
		os.Exit(0) // Exit with usage message
	}
	flag.Parse()
}
