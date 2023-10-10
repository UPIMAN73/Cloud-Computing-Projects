package main

import (
	"os"
	"strings"
)

// Read DB Command File (list of commands) and provide a command list
func ReadDBCmdFile(FileName string) []string {
	// Defining configuration
	var dbcmds []string = make([]string, 0)
	var dbcmdsS string

	// Read file
	dbcmdsB, err := os.ReadFile(FileName)
	CheckError(err)

	// convert byte array to string
	dbcmdsS = string(dbcmdsB[:])

	// Proper splitting (Remember to not use comma seperated values)
	if strings.Contains(dbcmdsS, "\r\n") {
		dbcmds = strings.Split(dbcmdsS, "\r\n")
	} else if strings.Contains(dbcmdsS, "\n") {
		dbcmds = strings.Split(dbcmdsS, "\n")
	}

	for i := 0; i < len(dbcmds); i++ {
		// Trim tabs
		if strings.Contains(dbcmds[i], "    ") {
			dbcmds[i] = strings.Trim(dbcmds[i], "    ")
		} else if strings.Contains(dbcmds[i], "\t") {
			dbcmds[i] = strings.Trim(dbcmds[i], "\t")
		}
	}

	// Return the goods
	return dbcmds
}
