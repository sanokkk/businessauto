package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
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
	DbConnectionString string `yaml:"connection_string" env:"DB_CONN"`
	MigrationsPath     string `yaml:"migrations_path"`
	MigrationsTable    string `yaml:"migrations_table"`
	SslMode            string `yaml:"ssl_mode"`
}

type JwtConfig struct {
	Secret             string        `yaml:"secret"`
	ExpireAfter        time.Duration `yaml:"expire_after"`
	RefreshExpireAfter time.Duration `yaml:"refresh_expire_after"`
}

type ApiConfig struct {
	Port            int    `yaml:"port"`
	Host            string `yaml:"host"`
	EnableAnyOrigin bool   `yaml:"enable_any_origin"`
}

type EnvConfig struct {
	Env        string `env:"ENV" env-default:"dev"`
	ConfigPath string `env:"CONFIG" env-default:"config.yml"`
}

type ContentConfig struct {
	UseContentStorage bool   `yaml:"use_content_storage"`
	Host              string `yaml:"host" env:"MINIO_HOST"`
	Port              string `yaml:"port" env:"MINIO_PORT"`
	User              string `yaml:"user" env:"MINIO_USER"`
	Secret            string `yaml:"secret" env:"MINIO_SECRET"`
	UseSsl            bool   `yaml:"use_ssl" env:"USE_SSL_MINIO"`
}

var once sync.Once
var config Config

func MustLoadConfig() *Config {
	once.Do(func() {
		readFile(&config)
		readEnv(&config)
	})

	return &config
}

func readFile(cfg *Config) {
	if err := cleanenv.ReadConfig("config.yml", cfg); err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(fmt.Sprintf("error while parsing config: %s", err.Error()))
	os.Exit(2)
}
