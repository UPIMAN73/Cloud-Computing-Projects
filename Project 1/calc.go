/**
 * @file calc.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import (
	"math"
	"sort"
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
	// Autosort if it has not already been sorted
	if !sort.Float64sAreSorted(list) {
		sort.Float64s(list)
	}

	// Definitions
	listLength := len(list)

	// Median calculation
	if listLength == 0 {
		return 0
	} else if listLength%2 == 0 {
		return (list[listLength/2-1] + list[listLength/2]) / 2
	} else {
		return list[listLength/2]
	}
}

// This function allows us to find the 99% percentile calculation used to determine performance
func ZeroNinePercentile(list []float64) float64 {
	// Autosort if it has not already been sorted
	if !sort.Float64sAreSorted(list) && len(list) > 0 {
		sort.Float64s(list)
	}

	// Find the 99% of the given list
	if len(list) > 0 {
		return list[int(math.Round(float64(len(list))*0.99))]
	} else {
		return 0.0
	}
}

// This function allows us to find the 99% percentile calculation used to determine performance
func OneNinePercentile(list []float64) float64 {
	// Autosort if it has not already been sorted
	if !sort.Float64sAreSorted(list) && len(list) > 0 {
		sort.Float64s(list)
	}

	// Find the 99% of the given list
	if len(list) > 0 {
		return list[int(math.Round(float64(len(list))*0.999))]
	} else {
		return 0.0
	}
}
