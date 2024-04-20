package viper_config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigFileType = "yaml"
	defaultPrefix         = ""
	defaultPrintDebugLogs = false
)

var defaultConfigFilePaths = []string{"./", "./configs"}

type Config struct {
}

func SetUpConfig(prefix string, configFileType string, printDebugLogs bool) *Config {
	//environmental variables
	if prefix != "" {
		viper.SetEnvPrefix(prefix)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//config file type
	validatedConfigFileType := validateConfigFileType(configFileType)
	viper.SetConfigType(validatedConfigFileType)

	//config file paths
	for _, path := range defaultConfigFilePaths {
		viper.AddConfigPath(path)
	}

	readConfigWithOptionalOverride("config", "config.override", validatedConfigFileType)

	config := &Config{}

	if printDebugLogs {
		config.PrintDebugLogs(prefix, validatedConfigFileType)
	}

	return config
}

func readConfigWithOptionalOverride(baseConfig string, overrideConfig string, configFileType string) {
	if configFileType == "" {
		configFileType = "yaml"
	}

	//check for presence of config files
	_, err := os.Stat(baseConfig + "." + configFileType)
	if err != nil {
		log.Printf("%s.%s not found. No config file loaded.", baseConfig, configFileType)
		return
	}

	// read the primary config file
	viper.SetConfigName(baseConfig)
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading base config file, %s. No config file loaded.", err)
	} else {
		log.Printf("Base config applied: %s\n", viper.ConfigFileUsed())
	}

	// check for presence of override config file
	_, err = os.Stat(overrideConfig + "." + configFileType)
	if err != nil {
		log.Printf("%s.%s not found. No override configuration loaded.", overrideConfig, configFileType)
		return
	}

	// merge if override config file if it exists
	viper.SetConfigName(overrideConfig)
	err = viper.MergeInConfig()
	if err != nil {
		log.Printf("Error reading override config file, %s. No override configuration loaded.", err)
	} else {
		log.Printf("Override config applied: %s\n", viper.ConfigFileUsed())
	}
}

func validateConfigFileType(configFileType string) string {
	if configFileType == "" {
		return defaultConfigFileType
	}

	//JSON, TOML, YAML, HCL, envfile
	if configFileType != "json" &&
		configFileType != "toml" &&
		configFileType != "yaml" &&
		configFileType != "hcl" &&
		configFileType != "envfile" {
		log.Printf("Invalid config file type: %s. Defaulting to yaml", configFileType)
		return defaultConfigFileType
	}

	return configFileType
}

// ---------------------------------------------------------

func (c *Config) PrintDebugLogs(prefix string, configFileType string) {
	log.Println("----- Config Setup -----")
	log.Printf("|| Prefix   %s", prefix)
	log.Printf("|| replacer %s", ". -> _")
	log.Printf("|| autoEnv  %s", "true")
	log.Printf("|| name     %s", "config OR config.override")
	log.Printf("|| paths    %s", defaultConfigFilePaths)
	log.Printf("|| fileType %s", configFileType)
	log.Println("------------------------")
}

// ---------------------------------------------------------

// Fallback configurations will be applied if there is no config file or environment variables
func (c *Config) SetFallbackConfigs(configs map[string]interface{}) {

	for k, v := range configs {
		if !viper.IsSet(k) {
			viper.Set(k, v)
		}
	}
}
