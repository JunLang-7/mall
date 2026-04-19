package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const ServerName = "mall"

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", ServerName)
	etcdAddr        string
	localConfigPath string
	GlobalConfig    Config
)

type Config struct {
	Server  Server            `yaml:"server" mapstructure:"server"`
	MySQL   MySQL             `yaml:"mysql" mapstructure:"mysql"`
	Redis   Redis             `yaml:"redis" mapstructure:"redis"`
	AppConf map[int32]AppConf `yaml:"app_conf" mapstructure:"app_conf"`
	BizConf BizConf           `yaml:"biz_conf" mapstructure:"biz_conf"`
}

type Server struct {
	HttpPort    int    `yaml:"http_port" mapstructure:"http_port"`
	Env         string `yaml:"env" mapstructure:"env"`
	EnablePprof bool   `yaml:"enable_pprof" mapstructure:"enable_pprof"`
	LogLevel    string `yaml:"log_level" mapstructure:"log_level"`
}
type MySQL struct {
	Dialect  string `yaml:"dialect" mapstructure:"dialect"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Database string `yaml:"database" mapstructure:"database"`
	Charset  string `yaml:"charset" mapstructure:"charset"`
	ShowSql  bool   `yaml:"show_sql" mapstructure:"show_sql"`
	MaxOpen  int    `yaml:"max_open" mapstructure:"max_open"`
	MaxIdle  int    `yaml:"max_idle" mapstructure:"max_idle"`
}

func (m *MySQL) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

type Redis struct {
	Addr    string `yaml:"addr" mapstructure:"addr"`
	PWD     string `yaml:"password" mapstructure:"password"`
	DBIndex int    `yaml:"db_index" mapstructure:"db_index"`
	MaxIdle int    `yaml:"max_idle" mapstructure:"max_idle"`
	MaxOpen int    `yaml:"max_open" mapstructure:"max_open"`
}

type AppConf struct {
	AppType   string `yaml:"app_type" mapstructure:"app_type"`
	AppName   string `yaml:"app_name" mapstructure:"app_name"`
	AppID     string `yaml:"app_id" mapstructure:"app_id"`
	AppSecret string `yaml:"app_secret" mapstructure:"app_secret"`
}

type BizConf struct {
	LarkGroupID string `yaml:"lark_group_id" mapstructure:"lark_group_id"`
}

func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default etcd address")
}

func InitConfig() *Config {
	var (
		err      error
		tempConf *Config
		vipConf  = viper.New()
	)
	vipConf.SetConfigType("yaml")
	// 优先使用ectd配置
	if etcdAddr != "" {
		tempConf, err = getFromRemoteAndWatchUpdate(vipConf)
		if err != nil {
			panic(err)
		}
	} else {
		// 本地配置
		tempConf, err = getFromLocal()
		if err != nil {
			panic(err)
		}
	}
	return tempConf
}

// getFromRemoteAndWatchUpdate 从远程配置中心获取配置，并监听配置更新
func getFromRemoteAndWatchUpdate(v *viper.Viper) (*Config, error) {
	tempConf := Config{}
	if err := v.AddRemoteProvider("etcd3", etcdAddr, etcdKey); err != nil {
		return nil, err
	}
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			if err := v.WatchRemoteConfig(); err == nil {
				_ = v.Unmarshal(&GlobalConfig)
				fmt.Println(">>> etcd config hot-reloaded: ", gconv.String(GlobalConfig))
			}
		}
	}()
	return &tempConf, nil
}

// getFromLocal 从本地文件获取配置
func getFromLocal() (*Config, error) {
	tempConf := Config{}
	if _, err := os.Stat(localConfigPath); err == nil {
		content, err := os.ReadFile(localConfigPath)
		if err != nil {
			return nil, err
		}
		// 展开环境变量 ${VAR} -> os.Getenv("VAR")
		content = []byte(os.ExpandEnv(string(content)))
		err = yaml.Unmarshal(content, &tempConf)
		if err != nil {
			return nil, err
		}
		return &tempConf, nil
	}
	return nil, fmt.Errorf("local config file not found ,file_name: %s", localConfigPath)
}
