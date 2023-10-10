/**
 * @file leaderless.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"fmt"
	"net"
)

// Leaderless Client Socket Function
func LeaderlessClientSocket(config Config, commandList []string) {
	// Establish a connection
	connections := make(map[string]net.Conn, len(config.Hosts))

	// Connect and check to see if there are any failed server connections
	for i := 0; i < len(config.Hosts); i++ {
		// Address Setup
		address := config.Hosts[i] + ":" + config.Ports.Socket

		// Connection setup
		connection, err := net.Dial("tcp", address)
		CheckError(err)

		// Set connection item to connection type
		for key, item := range connections {
			fmt.Printf("Connection Already Established, please fix your config file. \n\tYou have a duplicate for...\n\t\tHOST: %s\n", key)
			if config.Hosts[i] == key {
				err := item.Close()
				CheckError(err)
			} else {
				continue
			}
		}

		// Establish connection to conenction map
		connections[config.Hosts[i]] = connection
	}

	// Command List
	commandIndex := 0 // Used for inferencing the commands associated with a file list

	// Quarom definitions
	responseList := make(map[string]Response, 0)
	clientTrigger := false
	var err error

	// Client Loop
	buffer := make([]byte, 128)
	for commandIndex < len(commandList) {
		for hostID, connection := range connections {
			if clientTrigger {
				fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
				defer CheckError(err)
				errc := connection.Close()
				CheckError(errc)
				clientTrigger = false
			}
			// send value
			_, err := connection.Write([]byte(commandList[commandIndex]))
			CheckError(err)

			// read buffer
			messageLength, err := connection.Read(buffer)
			CheckError(err)

			// Add response to the response list
			responseList[hostID] = UIResponseStrip(string(buffer[:messageLength]))
			// fmt.Printf("\tHost: %s\n\tMessage: %s\n", hostID, responseList[hostID])

			// We don't use check error for this because we need to close the socket, then panic
			if err != nil {
				clientTrigger = true
			}
		}
		commandIndex += QuoromCheck(responseList, commandList[commandIndex])
	}
	for hostID, connection := range connections {
		fmt.Printf("Closing connection for Host: %s\r\n", hostID)
		errc := connection.Close()
		CheckError(errc)
	}
}
