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
	Times []TimeStamp
}

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

	// Benchmark Allocation
	benchmark := make(map[string]Benchmark, 0)
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
		benchmark[hostID] = Benchmark{ID: ID, Times: make([]TimeStamp, 0)}
	}

	// Client Loop
	buffer := make([]byte, 128)
	startTime := time.Now()
	fmt.Println(time.Now().Sub(startTime))
	for time.Now().Sub(startTime).Seconds() < EXPTIME {
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
				currentTimeStamp.initialTime = time.Now()
				CheckError(err)

				// read buffer
				messageLength, err := connection.Read(buffer)
				currentTimeStamp.finalTime = time.Now()
				CheckError(err)

				// Add response to the response list
				responseList[hostID] = UIResponseStrip(string(buffer[:messageLength]))
				benchmark[hostID].Times = append(benchmark[hostID].Times, TimeStamp{initialTime: currentTimeStamp.initialTime, finalTime: currentTimeStamp.finalTime})

				// We don't use check error for this because we need to close the socket, then panic
				if err != nil {
					clientTrigger = true
				}
			}
			commandIndex += QuoromCheck(responseList, commandList[commandIndex])
		}
		commandIndex = 0
	}
	for hostID, connection := range connections {
		fmt.Printf("Closing connection for Host: %s\r\n", hostID)
		errc := connection.Close()
		CheckError(errc)
	}

	// Calculations
	var timeDifferences map[int][]float64
	var averages map[int]float64
	var medians map[int]float64
	var ZNines map[int]float64
	var ONines map[int]float64

	// Benchmark Calculations
	for _, host := range benchmark {
		timeDifferences[host.ID] = make([]float64, len(host.Times))
		for i := 0; i < len(timeDifferences[host.ID]); i++ {
			timeDifferences[host.ID][i] = float64(host.Times[i].finalTime.Sub(host.Times[i].initialTime).Milliseconds())
		}
		averages[host.ID] = Average(timeDifferences[host.ID])
		medians[host.ID] = Median(timeDifferences[host.ID])
		ZNines[host.ID] = ZeroNinePercentile(timeDifferences[host.ID])
		ONines[host.ID] = OneNinePercentile(timeDifferences[host.ID])
	}

	// Benchmark Printouts
	fmt.Print("Averages:\r\n\t")
	fmt.Println(averages)
	fmt.Print("Medians:\r\n\t")
	fmt.Println(medians)
	fmt.Print("Zero Nines:\r\n\t")
	fmt.Println(ZNines)
	fmt.Print("One Nines:\r\n\t")
	fmt.Println(ONines)

}
