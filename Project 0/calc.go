/**
 * @file calc.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 0 - Ping Pong Project
 * @date 2023-09-10
 */
package main

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
