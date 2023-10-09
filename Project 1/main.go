/**
 * @file main.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"flag"
	"fmt"
)

// Main system function
func main() {
	// Flag declarations
	var configFile string    // <Config-File-Name>.yaml
	var displayHelp bool     // 1 = help, 0 = no help prompt
	var executionType string // ls = leadership, ll = leaderless
	var role string          // c = client, s = server

	// Assign flags to variable types
	flag.StringVar(&configFile, "f", "config.yaml", "Specifies the config file.     Options: \"config-file-name.yaml\".\n")
	flag.StringVar(&role, "r", "s", "Specifies the config file.     Options: \"config-file-name.yaml\".\n")
	flag.StringVar(&executionType, "t", "ll", "Specifies the data-replication type.     Options: leaderless (ll), leadership (ls).\n")
	flag.BoolVar(&displayHelp, "h", false, "Prints out the help screen.")

	// Parse command-line flags
	flag.Parse()

	// Flag control flow
	if displayHelp || configFile == "" || role == "" {
		// If 'help' flag is set or role/config file is not specified, display the usage information.
		fmt.Println(DefaultString())
		return
	} else {
		// Definitions
		var config Config
		// var db map[string]string

		// Loading config file
		LoadConfig(configFile, &config)

		// Initilize the database
		// db = make(map[string]string)

		// Run test
		if role == "s" {
			// DBTest()
			RunServerSocket(config)
		} else if role == "c" {
			switch executionType {
			case "ll":
				LeaderlessClientSocket(config)
			case "ls":
				LeadershipClientSocket(config)
				// RunClientSocket(config)
			default:
				PrintCommandHelp()
			}
		}
	}
}

// Usage of the program (default help string)
func DefaultString() string {
	output := "Examples:\n"
	output += "\tkvstore -t ll -r c \t\t Launches KVStore as a leaderless client"
	output += "\tkvstore -t leadless -r c \t\t Launches KVStore as a leaderless client"
	output += "\tkvstore -t ls -r c \t\t Launches KVStore as a leadership client"
	output += "\tkvstore -t leadership -r c \t\t Launches KVStore as a leadership client"
	output += "\tkvstore -r s \t\t Launches KVStore as a server"
	output += "For more information, kvstore -h\n"
	return output
}

func PrintCommandHelp() {
	flag.PrintDefaults()
	fmt.Println(DefaultString())
}
