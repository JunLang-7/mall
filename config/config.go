package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

const ServerName = "mall"

var (
	localConfigPath string
)

type Config struct {
	Server Server `yaml:"Server"`
	MySQL  MySQL  `yaml:"MySQL"`
}

type Server struct {
	HttpPort int    `yaml:"HttpPort"`
	LogLevel string `yaml:"LogLevel"`
}
type MySQL struct {
	Dialect  string `yaml:"dialect"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	ShowSql  bool   `yaml:"show_sql"`
	MaxOpen  int    `yaml:"max_open"`
	MaxIdle  int    `yaml:"max_idle"`
}

func (m *MySQL) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yml", "default config path")
}

func InitConfig() *Config {
	tempConf, err := getFromLocal()
	if err != nil {
		panic(err)
	}
	return tempConf
}

func getFromLocal() (*Config, error) {
	tempConf := Config{}
	if _, err := os.Stat(localConfigPath); err == nil {
		content, err := os.ReadFile(localConfigPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(content, &tempConf)
		if err != nil {
			return nil, err
		}
		return &tempConf, nil
	}
	return nil, fmt.Errorf("local config file not found ,file_name: %s", localConfigPath)
}
