package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	cfg  *AppConfig
	once sync.Once
)

// AppConfig holds all configuration
type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Session  SessionConfig  `yaml:"session"`
	Security SecurityConfig `yaml:"security"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`

	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	TablePrefix  string `yaml:"table_prefix"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
	LogLevel     string `yaml:"log_level"`
}

// DSN returns the MySQL DSN string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
	PoolSize int    `yaml:"pool_size"`
}

// Addr returns the Redis address
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type SessionConfig struct {
	Expire       int    `yaml:"expire"`
	CookiePrefix string `yaml:"cookie_prefix"`
}

type SecurityConfig struct {
	MD5Key  string `yaml:"md5_key"`
	APISalt string `yaml:"api_salt"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type LogConfig struct {
	ErrorPath string `yaml:"error_path"`
	Level     string `yaml:"level"`
}

// Load reads config from yaml file
func Load(path string) (*AppConfig, error) {
	var err error
	once.Do(func() {
		var data []byte
		data, err = os.ReadFile(path)
		if err != nil {
			return
		}
		cfg = &AppConfig{}
		err = yaml.Unmarshal(data, cfg)
	})
	return cfg, err
}

// Get returns the loaded config (must call Load first)
func Get() *AppConfig {
	return cfg
}
