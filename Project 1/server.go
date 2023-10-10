/**
 * @file server.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

import (
	"fmt"
	"net"
	"strconv"
)

// Pong role as a socket server
func RunServerSocket(config Config) {
	// Environment definitions
	var dbOutput string
	buffer := make([]byte, 128) // Message buffer
	address := ":" + config.Ports.Socket

	// Start Server
	fmt.Println("Server role begin...")
	server, err := net.Listen("tcp", address)
	CheckError(err)
	defer server.Close()

	// Server Host Parameters
	fmt.Printf("Server Role:\n\tPort: %s\n", address)

	// Server Loop
	for serverTrigger := false; !serverTrigger; {
		// Setup a connection to ping
		connection, err := server.Accept()
		CheckError(err)
		fmt.Println("Client Connected!")

		// BIG TODO
		// DB Implementation

		// Read the message from client
		messageLength, err := connection.Read(buffer)
		CheckError(err)

		fmt.Println(string(buffer[:messageLength]))
		cmdMessage := string(buffer[:messageLength])

		// Sending message
		dbCommand, dbArgs := UICMDStrip(cmdMessage)
		dbResponse := UICMDPass(dbCommand, dbArgs)
		if dbResponse.ResponseCode == DBACK {
			// Do Nothing
			if dbResponse.OPCode != Noop {
				DBQueueFlush()
				dbOutput = DBGet(dbArgs[0])
			} else {
				fmt.Println("No Operation")
			}
		}

		// Message Out
		messageOut := strconv.Itoa(int(dbResponse.ResponseCode)) + "," + strconv.Itoa(int(dbResponse.OPCode)) + "," + dbOutput
		_, err = connection.Write([]byte(messageOut))

		// We don't use check error for this because we need to close the socket, then panic
		if err != nil {
			fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
			defer CheckError(err)
			errc := connection.Close()
			CheckError(errc)
			serverTrigger = true
		}
	}
}
