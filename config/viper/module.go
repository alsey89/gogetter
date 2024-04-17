package viper_config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
}

func NewConfig(prefix string, printDebugLogs bool) *Config {

	//SETUP -----------------------------
	if prefix != "" {
		viper.SetEnvPrefix(prefix)
	} else {
		log.Printf("No prefix set for config")
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//CONFIG FILES -------------------------------------
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")

	readConfigWithOptionalOverride("config", "config.override")

	config := &Config{}

	if printDebugLogs {
		config.PrintDebugLogs(prefix)
	}

	return config
}

func (c *Config) SetConfigs(configs map[string]interface{}) {

	for k, v := range configs {
		if !viper.IsSet(k) {
			viper.Set(k, v)
		}
	}
}

func (c *Config) PrintDebugLogs(prefix string) {
	log.Println("----- Config Setup -----")
	log.Printf("|| Prefix   %s", prefix)
	log.Printf("|| replacer %s", ". -> _")
	log.Printf("|| autoEnv  %s", "true")
	log.Printf("|| name     %s", "config")
	log.Printf("|| path     %s", "./ OR ./configs")
	log.Println("------------------------")
}

// readConfigs reads and applies the configuration files using Viper.
// It takes two parameters: baseConfig (the name of the primary config file) and overrideConfig (the name of the override config file).
// If the primary config file exists, it is read and applied. If the override config file exists, it is merged with the primary config file.
func readConfigWithOptionalOverride(baseConfig string, overrideConfig string) {
	// read the primary config file
	viper.SetConfigName(baseConfig)
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Printf("Error reading config file, %s", err)
		} else {
			log.Println("No base config")
		}
	} else {
		log.Printf("Base config applied: %s\n", viper.ConfigFileUsed())
	}

	// merge if override config file if it exists
	viper.SetConfigName(overrideConfig)
	err = viper.MergeInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Printf("Error reading override config file, %s", err)
		} else {
			log.Println("No override config")
		}
	} else {
		log.Printf("Override config applied: %s\n", viper.ConfigFileUsed())
	}
}
