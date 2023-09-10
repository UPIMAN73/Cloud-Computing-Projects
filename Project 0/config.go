/**
 * @file config.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

// Imports
import (
	"os"

	"gopkg.in/yaml.v3"
)

// Connection ports based on protocol
type Ports_t struct {
	RPC    string `yaml:"RPC"`
	Socket string `yaml:"Socket"`
}

// Role stores the ping and pong messages associated with the config
type Roles_t struct {
	Ping string/*[128]byte */ `yaml:"Ping"` // Ping or initator role sends the ping message and records data1
	Pong string/*[128]byte */ `yaml:"Pong"` // Pong role just sends the pong message when it receives ping
}

// Config used to determine the structure of the system
type Config struct {
	Roles Roles_t `yaml:"Roles"`
	Host  string  `yaml:"Host"`
	Ports Ports_t `yaml:"Ports"`
}

// Load items from the config file into an config object
func LoadConfig(FileName string, config *Config) {
	// Defining configuration
	confData, err := os.ReadFile(FileName)
	CheckError(err)

	// Yaml config parsing
	err = yaml.Unmarshal(confData, &config)
	CheckError(err)
}
