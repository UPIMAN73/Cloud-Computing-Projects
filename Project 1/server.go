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
)

// Pong role as a socket server
func RunServerSocket(config Config) {
	// Initilize the database
	// var db map[string]string
	// db = make(map[string]string)

	// Environment definitions
	buffer := make([]byte, 128) // Message buffer for the ping pong messages to be read
	address := ":" + config.Ports.Socket

	// Start pong server
	fmt.Println("Server role begin...")
	server, err := net.Listen("tcp", address)
	CheckError(err)
	defer server.Close()

	// Describe the server being associated with the pong role
	fmt.Printf("Server Role:\n\tPort: %s\n", address)

	// Server Loop
	for serverTrigger := false; !serverTrigger; {
		// Setup a connection to ping
		connection, err := server.Accept()
		CheckError(err)
		fmt.Println("Client Connected!")

		// BIG TODO
		// DB Implementation
		// Client isn't working (continuosly) because we need to have this hosting on different hosts (or we need to thread it)

		// Read the message from client
		messageLength, err := connection.Read(buffer)
		CheckError(err)

		fmt.Println(string(buffer[:messageLength]))

		// Sending message
		_, err = connection.Write(buffer[:messageLength])

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
