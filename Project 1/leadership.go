/**
 * @file leadership.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import "fmt"

// Leadership Client Socket Function
func LeadershipClientSocket(config Config) {
	fmt.Println(config)
	// // TODO FIND HOST
	// // Establish a connection
	// connections := make(map[string]net.Conn, len(config.Hosts))

	// // Connect and check to see if there are any failed server connections
	// for i := 0; i < len(config.Hosts); i++ {
	// 	// Address Setup
	// 	address := config.Hosts[i] + ":" + config.Ports.Socket

	// 	// Connection setup
	// 	connection, err := net.Dial("tcp", address)
	// 	CheckError(err)

	// 	// Set connection item to connection type
	// 	for key, item := range connections {
	// 		fmt.Printf("Connection Already Established, please fix your config file. \n\tYou have a duplicate for...\n\t\tHOST: %s\n", key)
	// 		if config.Hosts[i] == key {
	// 			err := item.Close()
	// 			CheckError(err)
	// 		} else {
	// 			continue
	// 		}
	// 	}
	// 	connections[config.Hosts[i]] = connection
	// }

	// // Client Loop
	// buffer := make([]byte, 128)
	// for clientTrigger := false; !clientTrigger; {
	// 	for hostID, connection := range connections {
	// 		// send ping string
	// 		_, err := connection.Write([]byte("Put(A, Hello World!)"))
	// 		CheckError(err)

	// 		// read pong string
	// 		messageLength, err := connection.Read(buffer)
	// 		CheckError(err)

	// 		fmt.Printf("\tHost: %s\n\tMessage: %s\n", hostID, string(buffer[:messageLength]))

	// 		// We don't use check error for this because we need to close the socket, then panic
	// 		if err != nil {
	// 			fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
	// 			defer CheckError(err)
	// 			errc := connection.Close()
	// 			CheckError(errc)
	// 			clientTrigger = true
	// 		}
	// 		clientTrigger = true
	// 	}
	// }
}
