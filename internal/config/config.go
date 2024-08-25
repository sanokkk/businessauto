package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type Config struct {
	EnvConfig
	DbConfig      `yaml:"db_config"`
	JwtConfig     `yaml:"jwt_config"`
	ApiConfig     `yaml:"api_config"`
	ContentConfig `yaml:"content_storage_config"`
}

type DbConfig struct {
	DbConnectionString string `yaml:"connection_string"`
	MigrationsPath     string `yaml:"migrations_path"`
	MigrationsTable    string `yaml:"migrations_table"`
}

type JwtConfig struct {
	Secret             string        `yaml:"secret"`
	ExpireAfter        time.Duration `yaml:"expire_after"`
	RefreshExpireAfter time.Duration `yaml:"refresh_expire_after"`
}

type ApiConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type EnvConfig struct {
	Env        string `env:"ENV" env-default:"dev"`
	ConfigPath string `env:"CONFIG" env-default:"config.yml"`
}

type ContentConfig struct {
	UseContentStorage bool `yaml:"use_content_storage"`
}

var once sync.Once
var config Config

func MustLoadConfig() *Config {
	once.Do(func() {
		var envCfg EnvConfig
		if err := cleanenv.ReadEnv(&envCfg); err != nil {
			panic("error while reading config: " + err.Error())
		}

		parseConfig(envCfg.ConfigPath)

		config.Env = envCfg.Env
	})

	return &config
}

func GetConfigByPath(path string) *Config {
	parseConfig(path)

	return &config
}

func parseConfig(path string) {
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("error while reading config: " + err.Error())
	}
}
