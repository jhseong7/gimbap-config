package config

import (
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap"
	"github.com/spf13/viper"
)

type (
	ConfigService struct {
		option ConfigOption
		logger ecl.Logger
	}

	ConfigOption struct {
		// File path list of the configuration files to use
		ConfigFilePathList []string

		// Raw configuration data
		ConfigData map[string]interface{}

		// Toggle the automatic watching of the configuration file changes
		WatchConfigChange bool
	}
)

func (c *ConfigService) handleOnConfigChange() {
	c.logger.Info("Watching the configuration file changes")

	viper.OnConfigChange(func(e fsnotify.Event) {
		c.logger.Infof("Config file changed: %s", e.Name)

		// Reload the configurations
		c.loadConfig()
	})
	viper.WatchConfig()
}

func (c *ConfigService) handleConfigFilePathList() error {
	// Load configuration from the file path list
	for _, path := range c.option.ConfigFilePathList {
		// Check if the file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.logger.Warnf("Config file not found: %s", path)
			continue
		}

		// Get the extension of the file path
		split := strings.Split(path, ".")
		ext := split[len(split)-1]

		// Check if the extension is supported --> explicitly set the configuration type
		if ext == "env" {
			viper.SetConfigType("env")
		} else if ext == "yaml" || ext == "yml" {
			viper.SetConfigType("yaml")
		} else {
			c.logger.Warnf("Unsupported config file extension: %s. assuming as env", ext)
			viper.SetConfigType("env")
		}

		// Load configuration from the file path
		viper.SetConfigFile(path)

		if err := viper.MergeInConfig(); err != nil {
			c.logger.Errorf("Failed to load config file: %s. %v", path, err)
			continue
		}
	}

	return nil
}

func (c *ConfigService) handleRawConfigData() {
	for key, val := range c.option.ConfigData {
		// Show a warning if the key already exists
		if viper.IsSet(key) {
			c.logger.Warnf("Key already exists: %s. Overwriting", key)
		}

		viper.Set(key, val)
	}
}

func (c *ConfigService) loadConfig() {
	if err := c.handleConfigFilePathList(); err != nil {
		c.logger.Errorf("Failed to handle config file path list: %v", err)
	}

	c.handleRawConfigData()
}

func (c *ConfigService) Get(key string) interface{} {
	return viper.Get(key)
}

func (c *ConfigService) GetString(key string) string {
	return viper.GetString(key)
}

func (c *ConfigService) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c *ConfigService) GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func (c *ConfigService) GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func (c *ConfigService) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (c *ConfigService) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *ConfigService) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (c *ConfigService) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (c *ConfigService) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func NewConfigService(option ConfigOption) *ConfigService {
	c := &ConfigService{
		option: option,
		logger: ecl.NewLogger(ecl.LoggerOption{
			Name: "ConfigService",
		}),
	}

	c.loadConfig()

	if option.WatchConfigChange {
		c.handleOnConfigChange()
	}

	return c
}

var ConfigServiceProvider = gimbap.DefineProvider(gimbap.ProviderOption{
	Name:         "ConfigServiceProvider",
	Instantiator: NewConfigService,
})
