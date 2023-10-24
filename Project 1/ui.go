/**
 * @file ui.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"strconv"
	"strings"
)

// Response is used to determine if a
type Response struct {
	DBRStats DBRunStatus // Database response stats
	Values   []string    // Server response value
}

// User Interface to DB Command
func UICMDStrip(cmd string) (string, []string) {
	// filtering arguments & inputs
	var filterArgs []string
	filter := strings.Split(cmd, "(")

	// Proplery order the command into a set of values in a slice
	if len(filter) > 1 {
		filter[0] = strings.ToLower(filter[0])

		// Managing user input if a second argument exists
		if strings.Contains(filter[1], ", ") {
			filterArgs = strings.Split(filter[1], ", ") // easier to split
		} else if strings.Contains(filter[1], ",") {
			filterArgs = strings.Split(filter[1], ",") // easier to split
		} else {
			// fmt.Println("You did not properly write the db command.\r\n\tFormat:\t CMDTYPE(KEY,VALUE) or CMDTYPE(KEY, VALUE)")
			// return DBERR, ""
			// Do Nothing
			filterArgs = make([]string, 0)
			filterArgs = append(filterArgs, filter[1])
		}

		// Cleaning arguments
		if len(filterArgs) > 0 {
			filterArgs[len(filterArgs)-1] = filterArgs[len(filterArgs)-1][0 : len(filterArgs[len(filterArgs)-1])-1]
		} else {
			// "Passed Args had no values: " + (strings.Join(filter, "(")
		}
	}
	return filter[0], filterArgs
}

// User Interface Command Pass
func UICMDPass(cmd string, args []string) DBRunStatus {
	// Command Action
	switch cmd {
	case "benchmark":
		return DBRunStatus{DBACK, Noop}
	// Read
	case "get":
		return DBRunStatus{DBACK, Read}

	// Put
	case "put":
		DBEnqueue(DBCommand{Write, args[0], args[1]})
		return DBRunStatus{DBACK, Write}

		// Delete
	case "delete":
		DBEnqueue(DBCommand{Delete, args[0], ""})
		return DBRunStatus{DBACK, Delete}

	default:
		return DBRunStatus{DBERR, Noop}
	}
}

// User Interface Command Pass
func UICMDRunStatus(cmd string, args []string) DBRunStatus {
	// Command Action
	switch cmd {
	// Benchmark
	case "benchmark":
		return DBRunStatus{DBACK, Noop}

	// Read
	case "get":
		return DBRunStatus{DBACK, Read}

	// Put
	case "put":
		return DBRunStatus{DBACK, Write}

		// Delete
	case "delete":
		return DBRunStatus{DBACK, Delete}

	default:
		return DBRunStatus{DBERR, Noop}
	}
}

// UI Response strip to figure out what our response holds
func UIResponseStrip(response string) Response {
	// Split string
	values := strings.Split(response, ",")
	// fmt.Println(values)
	// fmt.Println()
	dbrc, err := strconv.Atoi(values[0])
	CheckError(err)
	dboc, err := strconv.Atoi(values[1])
	CheckError(err)
	var result Response
	result.DBRStats = DBRunStatus{DBRC(dbrc), DBOC(dboc)}
	result.Values = []string{values[2]}
	if len(values) > 3 {
		result.Values = append(result.Values, values[3])
	}
	return result
}
