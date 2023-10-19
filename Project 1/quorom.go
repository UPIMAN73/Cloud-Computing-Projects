package main

var Quarom int

// Checks to make sure that the responses from the servers are accurate with a quorom
func QuoromCheck(responseList map[string]Response, cmd string) int {
	// Important definitions
	cmd, values := UICMDStrip(cmd)                  // Pulls out the command values
	cmdRunStatus := UICMDRunStatus(UICMDStrip(cmd)) // Easier to deal with than having
	incrementable := 0                              // used for quorom calculations

	// Quarom loop
	for _, response := range responseList {
		// Checks the run status from the response
		if response.DBRStats == cmdRunStatus {
			// Checks the opcodes to make sure they match
			if DBCMDtoDBOC(cmd) == response.DBRStats.OPCode {
				// Determines the validity based on DB Operations (DB-Op Codes)
				switch response.DBRStats.OPCode {

				// Determines if the last write was proper (in reference to the user)
				case Write:
					if stringSlicesEqual(response.Values, values) {
						incrementable++
					}

				// We need to fix this later on (since not all hosts will be sending the same value, we need to find a way to test this)
				case Read:
					incrementable++

				// Delete case determines if response actually went through
				case Delete:
					if len(response.Values) > 1 {
						if values[0] == response.Values[0] && response.Values[1] == "" {
							incrementable++
						}
					}
				}
			}
		}
	}

	// If we have reached a quarom verdict
	if incrementable >= Quarom {
		return 1
	} else {
		return 0
	}

}
