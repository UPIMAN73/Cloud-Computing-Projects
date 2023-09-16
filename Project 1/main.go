package kvstorage

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var configFile string
	var displayHelp bool
	fmt.Println("Hello World!")
	fmt.Println(time.Now())

	// Assign flags to variable types
	flag.StringVar(&configFile, "f", "config.yaml", "Specifies the config file.     Options: \"config-file\".")
	flag.BoolVar(&displayHelp, "h", false, "Prints out the help screen.")

	// Parse command-line flags
	flag.Parse()

	// Flag control flow
	if displayHelp || configFile == "" {
		// If 'help' flag is set or role/config file is not specified, display the usage information.
		fmt.Println(DefaultString())
		return
	} else {
		// Definitions
		var config Config

		// Loading config file
		LoadConfig(configFile, &config)

		fmt.Println(config)
	}
}

// Usage of the program (default help string)
func DefaultString() string {
	output := "Usage: kvstorage [OPTIONS] [ARGUMENTS]\n"
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
	output += "\t-f config-file        Specifies the config file.     Options: \"config-file\".\n\n"
	output += "Examples:\n"
	output += "\tpingpong -r ping -f \"~/config.yaml\"               Performs the \"ping\" role using a different path for the config.yaml.\n"
	output += "\tpingpong -r pong -f \"../config.yaml\"           Performs the \"pong\" role using a different path for the config.yaml.\n"
	output += "\tpingpong -r ping -f \"config.yaml\"                Performs the \"ping\" role a local path for the config.yaml.\n\n"
	output += "For more information, pingpong -h\n"
	return output
}
