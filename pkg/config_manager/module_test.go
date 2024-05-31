package config_manager

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigNoPrefix(t *testing.T) {
	prefix := ""
	config := SetUpConfig(prefix, "yaml")

	// Assert that prefix is not set
	assert.Equal(t, "", viper.GetEnvPrefix())

	// Assert that config is not nil
	assert.NotNil(t, config)

	viper.Reset()
}

func TestNewConfigWithPrefix(t *testing.T) {
	prefix := "test"
	config := SetUpConfig(prefix, "yaml")

	// Assert that prefix is set
	assert.Equal(t, "test", viper.GetEnvPrefix())

	// Assert that config is not nil
	assert.NotNil(t, config)

	viper.Reset()
}

func TestReadConfigWithOverride(t *testing.T) {
	viper.Reset()

	baseConfigName := "config"
	overrideConfigName := "config.override"

	getConfigNameFromPath := func(configFilePath string) string {
		configFilePathArray := strings.Split(configFilePath, "/")
		configFileName := configFilePathArray[len(configFilePathArray)-1]

		return configFileName
	}

	t.Run("TestReadConfigWithNoConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		//set a test key
		viper.SetDefault("test_key", "test_value")

		// Read config with no config files
		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml")
		assert.Equal(t, "", viper.ConfigFileUsed())

		// Assert that the test key is still set
		assert.Equal(t, "test_value", viper.GetString("test_key"))
	})

	t.Run("TestReadConfigWithBaseConfigOnly", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()
		defer os.Remove(baseConfigName + ".yaml")

		//set a default test key
		viper.SetDefault("test_key", "test_value")

		// write base config file with new test value
		os.WriteFile(baseConfigName+".yaml", []byte("test_key: new_test_value"), 0644)

		viper.AddConfigPath(".")
		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml")

		usedConfigFileName := getConfigNameFromPath(viper.ConfigFileUsed())
		assert.Equal(t, baseConfigName+".yaml", usedConfigFileName, "Base config file should be used")

		// check if test key is overwritten by the new value in the base config file
		assert.Equal(t, "new_test_value", viper.GetString("test_key"))
	})

	t.Run("TestReadConfigWithOverrideConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()
		defer os.Remove("config.yaml")
		defer os.Remove("config.override.yaml")

		//set a default test key
		viper.SetDefault("test_key", "test_value")

		// write base and override config files with new test values
		os.WriteFile(baseConfigName+".yaml", []byte("test_key: new_test_value"), 0644)
		os.WriteFile(overrideConfigName+".yaml", []byte("test_key: override_test_value"), 0644)

		viper.AddConfigPath(".")
		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml")

		usedConfigFileName := getConfigNameFromPath(viper.ConfigFileUsed())

		assert.Equal(t, overrideConfigName+".yaml", usedConfigFileName, "Override config file should be used")

		// check if test key is overwritten by the new value in the override config file
		assert.Equal(t, "override_test_value", viper.GetString("test_key"))
	})
}
func TestValidateConfigFileType(t *testing.T) {
	t.Run("TestEmptyConfigFileType", func(t *testing.T) {
		result := validateConfigFileType("")
		assert.Equal(t, defaultConfigFileType, result)
	})

	t.Run("TestValidConfigFileType", func(t *testing.T) {
		validConfigFileTypes := []string{"json", "toml", "yaml", "hcl", "envfile"}
		for _, fileType := range validConfigFileTypes {
			result := validateConfigFileType(fileType)
			assert.Equal(t, fileType, result)
		}
	})

	t.Run("TestInvalidConfigFileType", func(t *testing.T) {
		invalidConfigFileTypes := []string{"txt", "xml", "csv"}
		for _, fileType := range invalidConfigFileTypes {
			result := validateConfigFileType(fileType)
			assert.Equal(t, defaultConfigFileType, result)
		}
	})
}
