/**
 * @file kvglobalDB.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"fmt"
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
	DBNOP DBRC = iota // DB No Run
	DBACK      = 1    // DB Ack
	DBERR      = 2    // DB Error
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
	return Key + "," + globalDB[Key]
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

// Converts a string into a DB Operation Command Value (Identifiable value)
func DBCMDtoDBOC(cmd string) DBOC {
	switch cmd {
	// Read
	case "get":
		return Read

	// Put
	case "put":
		return Write

		// Delete
	case "delete":
		return Delete

	default:
		return Noop
	}
}
