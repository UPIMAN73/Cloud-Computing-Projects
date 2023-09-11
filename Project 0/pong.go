/**
 * @file pong.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

import (
	"fmt"
	"net"
)

// Pong role as a gRPC server
func RunPongRPC(config Config) {
	//
}

// Pong role as a socket server
func RunPongSocket(config Config) {
	// Address
	address := config.Host + ":" + config.Ports.Socket

	// Start pong server
	fmt.Println("Pong role begin...")
	server, err := net.Listen("tcp", address)
	CheckError(err)
	defer server.Close()

	// Describe the server being associated with the pong role
	fmt.Println("Pong Role:\n\t" + address)

	// Setup a connection to ping
	connection, err := server.Accept()
	CheckError(err)
	fmt.Println("Ping Connected!")

	// Ping-Pong start
	fmt.Println("Running Ping-Pong")
	buffer := make([]byte, 128) // Message buffer for the ping pong messages to be read

	// Run this loop about 10,000 times
	for i := 0; i < RUN_EXECUTIONS; i++ {
		// Read the message from ping
		messageLength, err := connection.Read(buffer)
		CheckError(err)

		// Filter out message to make sure it is associated with the one in the config file
		if string(buffer[:messageLength]) == config.Roles.Ping {
			// Send Response
			_, err = connection.Write([]byte(config.Roles.Pong))
			CheckError(err)
		} else {
			// Kill the connection and run a response to fix the config
			fmt.Println("Not the right response, please try again when you have the correct configuration.")
			err := connection.Close()
			CheckError(err)
			return
		}
	}

	// Update user with status
	fmt.Println("Ping-Pong status complete")
}
