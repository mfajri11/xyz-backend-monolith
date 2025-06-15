package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg AppConfig
)

type AppConfig struct {
	Server    Server    `yaml:"server"`
	Database  Database  `yaml:"database"`
	KYCClient KYCClient `yaml:"kyc-client"`
	path      string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {

	env := getEnv("ENV", "dev")
	configFile := "config." + env + ".yaml"
	// little trick to change to working directory, really useful for testing
	// see https://brandur.org/fragments/testing-go-project-root
	cfg = AppConfig{}
	_, fname, _, _ := runtime.Caller(0)
	dir := filepath.Join(fname, "..")
	cfg.path = dir
	configPath := filepath.Join(dir, configFile)
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}

}

func Get() *AppConfig {
	return &cfg
}

type Server struct {
	Host            string        `yaml:"host" env-required:"true"`
	Port            int           `yaml:"port" env-default:"9000" env-layout:"int"`
	ReadTimeout     time.Duration `yaml:"read-timeout" env-default:"10s" env-layout:"time.Duration"`
	WriteTimeout    time.Duration `yaml:"write-timeout" env-default:"10s" env-layout:"time.Duration"`
	ShutDownTimeout time.Duration `yaml:"shutdown-timeout" env-default:"30s" env-layout:"time.Duration"`
}

type Database struct {
	Host                  string        `yaml:"host" env-required:"true"`
	Port                  int           `yaml:"port" env-required:"true"`
	Username              string        `env:"DB_USER"`
	Password              string        `env:"DB_PASSWORD" env-layout:"string"`
	DataBaseName          string        `env:"DB_NAME" env-layout:"string"`
	OpenConnection        int           `yaml:"max-open-connection" env-layout:"int"`
	IdleConnection        int           `yaml:"max-idle" env-layout:"int"`
	ConnectionMaxLifeTime time.Duration `yaml:"max-lifetime" env-layout:"time.Duration"`
}

type KYCClient struct {
	BaseURL string `yaml:"base-url" env-required:"true"`
	APIKey  string `yaml:"api-key" env-required:"true"`
	APPID   string `yaml:"app-id" env-required:"true" env-layout:"string" env-default:"xyz"`
}
