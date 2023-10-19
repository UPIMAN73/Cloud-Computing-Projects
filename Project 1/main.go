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
	"math"
)

// Main system function
func main() {
	// Flag declarations
	var commandFile string   // <DB-Command-File-Name>.txt
	var configFile string    // <Config-File-Name>.yaml
	var displayHelp bool     // 1 = help, 0 = no help prompt
	var executionType string // ls = leadership, ll = leaderless
	var role string          // c = client, s = server
	var id int               // id = # (only needed for server)

	// Assign flags to variable types
	flag.StringVar(&commandFile, "d", "dbcmd.txt", "Specifies the Database Command List file.     Options: \"dbcmd.txt\" (NO .CSV please, it doesn't work well).\r\n")
	flag.StringVar(&configFile, "f", "config.yaml", "Specifies the config file.     Options: \"config-file-name.yaml\".\r\n")
	flag.StringVar(&role, "r", "s", "Specifies the config file.     Options: \"config-file-name.yaml\".\r\n")
	flag.StringVar(&executionType, "t", "ll", "Specifies the data-replication type.     Options: leaderless (ll), leadership (ls).\r\n")
	flag.IntVar(&id, "id", 0, "Specifies the server ID needed for when a client is connecting, it knows which server it is connecting too. Example: '1'\r\n")
	flag.BoolVar(&displayHelp, "h", false, "Prints out the help screen.\r\n")

	// Parse command-line flags
	flag.Parse()

	// Flag control flow
	if displayHelp || configFile == "" || role == "" {
		// If 'help' flag is set or role/config file is not specified, display the usage information.
		flag.PrintDefaults()
		fmt.Println(DefaultString())
		return
	} else {
		// var db map[string]string

		// Loading config file
		LoadConfig(configFile, &config)

		// Define Quarom
		Quarom = int(math.Floor(float64(len(config.Hosts)/2))) + 1

		// Initilize the database
		dbCmds := ReadDBCmdFile(commandFile)
		if len(dbCmds) == 0 && role == "c" {
			fmt.Println("We need commands to be populated in the db-commands file (example name: 'dbcmds.txt')")
			fmt.Println("The organized way that we have setup for this is defined as such:\r\nFormat:")
			fmt.Println("\tput(A, Hello World!)\r\n\tget(A)")
			fmt.Println("You don't necessarily need tabs but it is okay if you have them.")
			fmt.Println("We highly recommend that you don't use tabs or spaces before a command.")
			return
		}

		// Run
		if role == "s" {
			if id > 0 {
				RunServerSocket(id, config)
			} else {
				fmt.Println("ID Was not defined, please define it with \"-id #\" or \"-id=#\".")
			}
		} else if role == "c" {
			switch executionType {
			case "ll":
				LeaderlessClientSocket(config, dbCmds)
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
