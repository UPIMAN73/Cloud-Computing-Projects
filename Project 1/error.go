/**
 * @file error.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

// Check for errors & stop program if one occurs
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
