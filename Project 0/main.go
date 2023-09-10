/**
 * @file main.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

// Imports
import (
	"flag"
	"fmt"
	"time"
)

// RTT
type RTT struct {
	Start time.Time // Initial or start point of time when a message has been sent
	End   time.Time // Final or End point of time when a message has been received
}

// Calculate RTT based on rtt types
func CalculateRTT(rtt RTT) time.Duration {
	return rtt.End.Sub(rtt.Start)
}

// Check for errors & stop program if one occurs
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// Convert a list of durations into a byte string
func ConvertDurationToBytes(list []time.Duration) string {
	if len(list) > 1 {
		var output string
		for i := 0; i < len(list)-1; i++ {
			output += list[i].String() + ","
		}
		output += list[len(list)-1].String()
		return output
	} else if len(list) == 1 {
		return list[0].String()
	} else {
		return ""
	}
}

const (
	RUN_EXECUTIONS = 10000
)

// Main function
func main() {
	// Define command-line flags
	var displayHelp bool
	var role string
	var connectionType string
	var configFile string

	// Assign flags to variable types
	flag.StringVar(&role, "r", "", "Specifies the role to perform. Options: \"ping\" or \"pong\".")
	flag.StringVar(&connectionType, "t", "socket", "Specifies the connection type. Options: \"rpc\", \"socket\", or \"both\".")
	flag.StringVar(&configFile, "f", "config.yaml", "Specifies the config file.     Options: \"config-file\".")
	flag.BoolVar(&displayHelp, "h", false, "Prints out the help screen.")

	// Parse command-line flags
	flag.Parse()

	// Flag control flow
	if displayHelp || role == "" || connectionType == "" || configFile == "" {
		// If 'help' flag is set or role/connectionType is not specified, display the usage information.
		fmt.Println(DefaultString())
		return
	} else {
		// Making sure the arguments correlate to the specified values
		if role != "ping" && role != "pong" {
			fmt.Println(DefaultString())
			return
		}

		// Making sure the arguments correlate to the specified values
		if connectionType != "rpc" && connectionType != "socket" && connectionType != "both" {
			fmt.Println(DefaultString())
			return
		}

		// Definitions
		var config Config

		// Loading config file
		LoadConfig(configFile, &config)

		// Deciding role
		switch role {
		case "ping":
			if connectionType == "socket" {
				RunPingSocket(config)
			} else if connectionType == "rpc" {
				RunPingSocket(config) // Change to RPC later on
			} else {
				RunPingSocket(config) // add RPC later on
			}
		case "pong":
			if connectionType == "socket" {
				RunPongSocket(config)
			} else if connectionType == "rpc" {
				RunPongSocket(config) // Change to RPC later on
			} else {
				RunPongSocket(config) // add RPC later on
			}
		}
	}
}

// Usage of the program (default help string)
func DefaultString() string {
	output := "Usage: pingpong [OPTIONS] [ARGUMENTS]\n"
	output += "Description:\n"
	output += "\tThis command performs a network transportation delay analysis using the ping-pong protocol designed by me.\n"
	output += "\tThe idea is ping connects to a pong server and obtains time metrics based on response times (RTT). However,\n"
	output += "\tthere is a config file that determines whether the ping-pong communication can proceed. This limites the amount\n"
	output += "\tof bandwidth that is associated with the pong server, and allows for a service managable ping that is more secure than\n"
	output += "\tother alternatives. The config for this program must be a yaml defined file.\n\n\n"
	output += "Options:\n"
	output += "\t-h, --help         Show this help message and exit.\n\n"
	output += "Arguments:\n"
	output += "\t-r role               Specifies the role to perform. Options: \"ping\" or \"pong\".\n"
	output += "\t-t connection-type    Specifies the connection type. Options: \"rpc\", \"socket\", or \"both\".\n"
	output += "\t-f config-file        Specifies the config file.     Options: \"config-file\".\n\n"
	output += "Examples:\n"
	output += "\tpingpong -r ping -t rpc -f \"~/config.yaml\"               Performs the \"ping\" role using \"RPC\" as the connection type.\n"
	output += "\tpingpong -r pong -t socket -f \"../config.yaml\"           Performs the \"pong\" role using \"socket\" as the connection type.\n"
	output += "\tpingpong -r ping -t both -f \"config.yaml\"                Performs the \"ping\" role using both \"RPC\" and \"socket\" as the connection type.\n\n"
	output += "For more information, pingpong -h\n"
	return output
}
