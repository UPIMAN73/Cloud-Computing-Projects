package kvstorage

// Check for errors & stop program if one occurs
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
