package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	SeedUser SeedUserConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port   int
	Secret string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	Server     string
	DBUser     string
	DBPassword string
	DBName     string
}

// SeedUserConfigurations exported
type SeedUserConfigurations struct {
	Email           string
	InitialPassword string
}

// GetConfig returns the configuration
func GetConfig(path string) (*Configurations, error) {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath(path)
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("toml")
	var conf Configurations

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file, %v", err)
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "knowledge-keeper")

	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode configurations into struct, %v", err)
	}

	return &conf, nil
}
