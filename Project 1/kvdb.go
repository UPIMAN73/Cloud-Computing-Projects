/**
 * @file kvglobalDB.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"fmt"
	"strings"
)

// Database used for program
var globalDB map[string]string = make(map[string]string, 0)

// Making DB Queue
var dbQueue DBQueue = DBQueue{make([]DBCommand, 0), MAXQUEUESIZE}

// Declares operations used to DB
type DBOC int

// DB Op Codes
const (
	Noop DBOC = iota
	Read
	Write
	Delete
)

// Operation Response Codes
type DBRC int

// Operation Response Codes
const (
	DBNOP DBRC = 0 // DB No Run
	DBACK      = 1 // DB Ack
	DBERR      = 2 // DB Error
)

// Database Command used to determine
type DBCommand struct {
	Type  DBOC   // Declares the type of operation going to be performed
	Key   string // Key of item you want to operate on
	Value string // Value is used for the write commands
}

// Database queue of commands that are ready to be processed
const (
	MAXQUEUESIZE = 5
)

// DB Queue Type
type DBQueue struct {
	History []DBCommand // The slice of items in the queue for operating on when it is ready
	Size    int         // Maximum Size of the queue
}

// Defines whether or not the operation was successful
type DBRunStatus struct {
	ResponseCode DBRC // Response code
	OPCode       DBOC // Operation Type
}

/* QUEUE Functions */

// Enqueue or append an command to the end of the DBQueue (We only queue Write & Delete commands not Read)
func DBEnqueue(DBC DBCommand) {
	if len(dbQueue.History) == dbQueue.Size {
		fmt.Println("Overflow!!!")
		DBRun(DBDenqueue())
	}
	dbQueue.History = append(dbQueue.History, DBC)
}

// Denqueue or remove a command from the begining of the queue
func DBDenqueue() DBCommand {
	var result DBCommand
	if len(dbQueue.History) > 1 {
		result = dbQueue.History[0]
		dbQueue.History = dbQueue.History[1:len(dbQueue.History)]
	} else if len(dbQueue.History) == 1 {
		result = dbQueue.History[0]
		dbQueue.History = nil
	} else {
		// Do Nothing
		result = DBCommand{Noop, "", ""}
	}
	return result
}

// DB Queue Run select number of items to run
func DBQueueRun(Items int) []DBRunStatus {
	// Making sure we properly set the number of items
	if Items > len(dbQueue.History) {
		Items = len(dbQueue.History)
	}

	// Used to determine if the run state was proper
	results := make([]DBRunStatus, Items)

	// Iterativly run items
	for i := 0; i < Items; i++ {
		results[i] = DBRun(dbQueue.History[0])
		if results[i].ResponseCode == DBERR || results[i].ResponseCode == DBNOP {
			break
		} else {
			DBDenqueue()
		}
	}

	// Return the results of all of the sates of the run
	return results
}

// Run all items in queue
func DBQueueFlush() []DBRunStatus {
	return DBQueueRun(dbQueue.Size)
}

/* Run Commands */
func DBRun(DBC DBCommand) DBRunStatus {
	var result DBRC
	switch DBC.Type {
	// No Operation
	case Noop:
		return DBRunStatus{DBNOP, Noop}

	// Read
	case Read:
		DBGet(DBC.Key)
		return DBRunStatus{DBACK, Read}

	// Write
	case Write:
		DBSet(DBC.Key, DBC.Value)
		if globalDB[DBC.Key] == DBC.Value {
			result = DBACK
		} else {
			result = DBERR
		}
		return DBRunStatus{result, Write}

	// Delete
	case Delete:
		DBDelete(DBC.Key)
		if globalDB[DBC.Key] == "" {
			result = DBACK
		} else {
			result = DBERR
		}
		return DBRunStatus{result, Delete}
	}
	return DBRunStatus{DBNOP, Noop}
}

// DB Get item based on key
func DBGet(Key string) string {
	return Key + ", " + globalDB[Key]
}

// DB Set Item based on key and value
func DBSet(Key string, Value string) {
	// Adding item to globalDB (if it doesn't exist), otherwise set key to value
	globalDB[Key] = Value
}

// DB remove item
func DBDelete(Key string) {
	if globalDB[Key] != "" {
		delete(globalDB, Key)
	} else {
		return
	}
}

// User Interface to DB Command
func UICMDStrip(cmd string) DBRunStatus {
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
			return DBRunStatus{DBERR, Noop}
		}
	} else {
		return DBRunStatus{DBERR, Noop}
	}

	// Command Action
	if filter[0] == "get" {
		// Read
		return DBRunStatus{DBACK, Read}
	} else if filter[0] == "put" {
		// Put
		DBEnqueue(DBCommand{Write, filterArgs[0], filterArgs[1]})
		return DBRunStatus{DBACK, Write}
	} else if filter[0] == "delete" {
		// Delete
		DBEnqueue(DBCommand{Delete, filterArgs[0], ""})
		return DBRunStatus{DBACK, Delete}
	} else {
		// Do Nothing
		// fmt.Println("")
		// "You did not pass a proper command!"
		return DBRunStatus{DBERR, Noop}
	}
}

// Database test
func DBTest() {

	// Adding items to DB Queue
	fmt.Println(UICMDStrip("put(h5, Test)"))

	//  Run before a get
	DBQueueFlush()
	fmt.Println(UICMDStrip("get(h5)"))

	// Flush on a put
	fmt.Println(UICMDStrip("put(h5, Joshua)"))
	DBQueueFlush()

	// Flush before a get
	DBQueueFlush()
	fmt.Println(UICMDStrip("get(h5)"))

	// Flush after a delete
	fmt.Println(UICMDStrip("delete(h5)"))
	DBQueueFlush()

	// Flush before a get after a delete
	DBQueueFlush()
	fmt.Println(UICMDStrip("get(h5)"))
	fmt.Println(dbQueue)
	// DBEnqueue(DBCommand{Write, "h5", "Test"})
	// DBEnqueue(DBCommand{Delete, "h5", ""})
	// DBEnqueue(DBCommand{Write, "h5", "Joshua!"})
	// DBQueueRun(1)
	// fmt.Println(DBGet("h5"))
	// DBEnqueue(DBCommand{Delete, "h5", ""})
	// DBQueueRun(1)
	// fmt.Println(DBGet("h5"))

	// Commands works!
}
