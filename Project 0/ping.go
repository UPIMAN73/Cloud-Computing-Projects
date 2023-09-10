/**
 * @file ping.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"time"
)

// Find the average of a float64 list
func Average(list []float64) float64 {
	var output float64
	for i := 0; i < len(list); i++ {
		output += list[i]
	}
	return output / float64(len(list))
}

// Find the median of a float64 list
func Median(list []float64) float64 {
	var median float64
	listLength := len(list)
	if listLength == 0 {
		return 0
	} else if listLength%2 == 0 {
		median = (list[listLength/2-1] + list[listLength/2]) / 2
	} else {
		median = list[listLength/2]
	}
	return median
}

// Ping role as a socket client
func RunPingSocket(config Config) {
	// Stats definition
	var rtt RTT
	rtt_list := make([]time.Duration, 0)
	rtt_int_list := make([]float64, RUN_EXECUTIONS)

	// Establish a connection
	address := config.Host + ":" + config.Ports.Socket
	connection, err := net.Dial("tcp", address)
	CheckError(err)

	// Ping Pong
	buffer := make([]byte, 128)
	defer connection.Close()
	// Send ping message 10,000 times
	for i := 0; i < RUN_EXECUTIONS; i++ {
		// send ping string
		_, err = connection.Write([]byte(config.Roles.Ping))
		rtt.Start = time.Now()
		CheckError(err)

		// read pong string
		messageLength, err := connection.Read(buffer)
		rtt.End = time.Now()
		CheckError(err)

		// check if pong string matches the one in the config file
		if string(buffer[:messageLength]) == config.Roles.Pong {
			rtt_list = append(rtt_list, CalculateRTT(rtt))
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

	// Sort the data
	for i := 0; i < RUN_EXECUTIONS; i++ {
		rtt_int_list[i] = float64(rtt_list[i].Microseconds())
	}
	sort.Float64s(rtt_int_list)

	// Write RTT list to a file
	fmt.Println("Writing RTT stats to file: 'rtt_socket_output_stats.txt'")
	f, err := os.Create("rtt_socket_output_stats.txt")
	CheckError(err)

	// Write Average
	f.WriteString("Average: ")
	f.WriteString(fmt.Sprintf("%f\n", Average(rtt_int_list)))

	// Write Median
	f.WriteString("Median: ")
	f.WriteString(fmt.Sprintf("%f\n", Median(rtt_int_list)))

	// Write 99th
	f.WriteString("99th: ")
	f.Close()

	// Write list to file
	fmt.Println("Writing RTT stats to file: 'rtt_socket_output_file.txt'")
	f, err = os.Create("rtt_socket_output_file.txt")
	CheckError(err)
	f.WriteString(ConvertDurationToBytes(rtt_list))
	f.Close()
}
