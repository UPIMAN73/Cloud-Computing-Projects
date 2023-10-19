/**
 * @file config.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
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
	Socket string `yaml:"Socket"`
}

// Config used to determine the structure of the system
type Config struct {
	Hosts []string `yaml:"Hosts"`
	Ports Ports_t  `yaml:"Ports"`
}

// Config Global
var config Config

// Load items from the config file into an config object
func LoadConfig(FileName string, cfg *Config) {
	// Defining configuration
	confData, err := os.ReadFile(FileName)
	CheckError(err)

	// Yaml config parsing
	err = yaml.Unmarshal(confData, &cfg)
	CheckError(err)
}
