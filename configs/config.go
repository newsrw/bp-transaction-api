package configs

import (
	"strings"

	"github.com/spf13/viper"
)

var (
	// Conf configs model
	Conf = &Config{}
)

const (
	configName     = "config"
	configTypeYAML = "yml"
)

// Environment environment
type Environment struct {
	Config *Config `mapstructure:"CONFIG"`
}

// Config config
type Config struct {
	Server struct {
		Address string `mapstructure:"ADDRESS"`
	} `mapstructure:"SERVER"`
	Client struct {
		Port                 int `mapstructure:"PORT"`
		GracefulStopTimeout  int `mapstructure:"GRACEFULSTOPTIMEOUT"`
		MonitorDelayDuration int `mapstructure:"MONITORDELAYDURATION"`
	} `mapstructure:"CLIENT"`
}

// Read read config environment from path
func Read(path string) (*Environment, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(path)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType(configTypeYAML)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	env := &Environment{}
	err := bind(v, &env)
	if err != nil {
		return nil, err
	}

	Conf = env.Config

	return env, nil
}

func bind(v *viper.Viper, i interface{}) error {
	err := v.Unmarshal(i)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig get config
func GetConfig() *Config {
	return Conf
}
