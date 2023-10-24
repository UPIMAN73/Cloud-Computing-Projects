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

var serverTrigger bool = false
var connMap map[string]net.Conn

// Pong role as a socket server
func RunServerSocket(ID int, config Config) {
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

	// Starting threading (Goroutine) for client server connection handling
	connMap = make(map[string]net.Conn)
	go ClientHandling(&server)

	// Wait for client before fully starting server
	for len(connMap) == 0 {
		// Do Nothing
	}

	// Server Loop
	for serverTrigger = false; !serverTrigger; {
		// Client loops
		for hostID, connection := range connMap {
			// Read the message from client
			messageLength, err := connection.Read(buffer)

			// If error from connection from client
			if err != nil {
				fmt.Printf("Host: %s Connection Closed", hostID)
				errc := connection.Close()
				delete(connMap, hostID)
				fmt.Println(errc)
				continue
			}

			// If no connection from client
			if messageLength == 0 {
				errc := connection.Close()
				delete(connMap, hostID)
				fmt.Println(errc)
				continue
			}

			// Benchmark Condition
			if string(buffer[:messageLength]) == "benchmark()" {
				fmt.Println("Benchmark Setup")
				// Message Out
				messageOut := strconv.Itoa(int(ID))

				// Write
				_, err = connection.Write([]byte(messageOut))

				// We don't use check error for this because we need to close the socket, then panic if another error occurs
				if err != nil {
					fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
					fmt.Printf("Host: %s Connection Closed", hostID)
					defer CheckError(err)
					errc := connection.Close()
					delete(connMap, hostID)
					CheckError(errc)
					continue
				}
				fmt.Println("Benchmark Start")
				continue
			}

			// Prints out the message for ease of understanding
			fmt.Println(string(buffer[:messageLength]))
			cmdMessage := string(buffer[:messageLength])

			// DB Message processing
			dbCommand, dbArgs := UICMDStrip(cmdMessage)
			dbResponse := UICMDPass(dbCommand, dbArgs)

			// Check if DB Command responds in ack
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

			// We don't use check error for this because we need to close the socket, then panic if another error occurs
			if err != nil {
				fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
				fmt.Printf("Host: %s Connection Closed", hostID)
				defer CheckError(err)
				errc := connection.Close()
				delete(connMap, hostID)
				CheckError(errc)
				continue
			}
		}
	}
}

// Client handling for server using goroutines
func ClientHandling(server *net.Listener) {
	// Server Trigger
	for serverTrigger == false {
		// Setup a connection to ping
		connection, err := (*server).Accept()
		if err != nil {
			fmt.Println(err)
			connection.Close()
		}
		fmt.Println("Client Connected!")
		connMap[connection.RemoteAddr().String()] = connection
	}
}
