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
	"os"
	"strconv"
	"time"
)

const (
	EXPTIME = 60
)

// Time Stamps
type TimeStamp struct {
	initialTime time.Time
	finalTime   time.Time
}

// Benchmark
type Benchmark struct {
	ID    int
	Times []float64
}

// Leaderless Client Socket Function
func LeaderlessClientSocket(config Config, commandList []string) {
	// Establish a connection
	connections := make(map[string]net.Conn, len(config.Hosts))

	// Connect and check to see if there are any failed server connections
	for i := 0; i < len(config.Hosts); i++ {
		// Check for double connections
		for key, item := range connections {
			if config.Hosts[i] == key {
				fmt.Printf("Connection Already Established, please fix your config file. \n\tYou have a duplicate for...\n\t\tHOST: %s\n", key)
				err := item.Close()
				CheckError(err)
			} else {
				continue
			}
		}
		// Address Setup
		address := config.Hosts[i] + ":" + config.Ports.Socket

		// Connection setup
		connection, err := net.Dial("tcp", address)
		CheckError(err)

		// Establish connection to conenction map
		connections[config.Hosts[i]] = connection
	}

	// Command List
	commandIndex := 0 // Used for inferencing the commands associated with a file list

	// Quarom definitions
	responseList := make(map[string]Response, 0)
	clientTrigger := false
	var err error

	// Benchmark Allocation
	benchmark := make(map[string]*Benchmark, 0)
	currentTimeStamp := TimeStamp{initialTime: time.Now(), finalTime: time.Now()}
	for hostID, connection := range connections {
		_, err = connection.Write([]byte("benchmark()"))
		if err != nil {
			fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
			defer CheckError(err)
			errc := connection.Close()
			CheckError(errc)
		}

		IDBuffer := make([]byte, 128)
		messageLength, err := connection.Read(IDBuffer)
		if err != nil {
			fmt.Println("A read error occured to the socket stream, please check to make sure something did not happen to the client.")
			defer CheckError(err)
			errc := connection.Close()
			CheckError(errc)
		}
		ID, err := strconv.Atoi(string(IDBuffer[:messageLength]))
		if err != nil {
			fmt.Println("A read error occured to the socket stream, please check to make sure something did not happen to the client.")
			defer CheckError(err)
			errc := connection.Close()
			CheckError(errc)
		}
		benchmark[hostID] = &Benchmark{ID: ID, Times: make([]float64, 1)}
	}

	// Client Variables
	buffer := make([]byte, 128)
	startTime := time.Now()

	// Client Loop (Timed)
	for time.Since(startTime).Seconds() < EXPTIME {
		// Command index check
		for commandIndex < len(commandList) {
			// Connection Processing
			for hostID, connection := range connections {
				// If client close of connection is triggered, close socket stream
				if clientTrigger {
					fmt.Println("A write error occured to the socket stream, please check to make sure something did not happen to the client.")
					defer CheckError(err)
					errc := connection.Close()
					CheckError(errc)
					clientTrigger = false
				}

				// send value
				_, err := connection.Write([]byte(commandList[commandIndex]))
				currentTimeStamp.initialTime = time.Now()
				CheckError(err)

				// read buffer
				messageLength, err := connection.Read(buffer)
				currentTimeStamp.finalTime = time.Now()
				CheckError(err)

				// Add response to the response list
				responseList[hostID] = UIResponseStrip(string(buffer[:messageLength]))
				benchmark[hostID].Times = append(benchmark[hostID].Times, float64(currentTimeStamp.finalTime.Sub(currentTimeStamp.initialTime).Milliseconds()))

				// We don't use check error for this because we need to close the socket, then panic
				if err != nil {
					clientTrigger = true
				}
			}

			// Increment command index when the quorom is finalized
			commandIndex += QuoromCheck(responseList, commandList[commandIndex])
		}
		// Repeat the process until the time interrupt is met
		commandIndex = 0
	}

	// Close all connections to hosts
	for hostID, connection := range connections {
		fmt.Printf("Closing connection for Host: %s\r\n", hostID)
		errc := connection.Close()
		CheckError(errc)
	}

	// Calculations
	averages := make(map[int]float64, 0)
	medians := make(map[int]float64, 0)
	ZNines := make(map[int]float64, 0)
	ONines := make(map[int]float64, 0)

	// Benchmark Calculations
	for _, host := range benchmark {
		averages[host.ID] = Average(host.Times)
		medians[host.ID] = Median(host.Times)
		ZNines[host.ID] = ZeroNinePercentile(host.Times)
		ONines[host.ID] = OneNinePercentile(host.Times)
	}

	// Benchmark Printouts
	fmt.Println("Writing Benchmark stats to file: 'kvdb_stats.txt'")
	f, err := os.Create("kvdb_stats.txt")
	CheckError(err)

	// Write Average
	f.WriteString("Average: ")
	f.WriteString(fmt.Sprintf("%f\n", averages))

	// Write Median
	f.WriteString("Median: ")
	f.WriteString(fmt.Sprintf("%f\n", medians))

	// Write 99%
	f.WriteString("99th-Percentile: ")
	f.WriteString(fmt.Sprintf("%f\n", ZNines))

	// Write 99.9%
	f.WriteString("99.9th-Percentile: ")
	f.WriteString(fmt.Sprintf("%f\n", ONines))

	// Close file stream
	f.Close()

}
