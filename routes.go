package main

import (
	"fmt"
	"reflect"
)

// isOnRoute tells if a room is on a slice
func isOnRoute(route route, room room) bool {
	for _, r := range route {
		if room.Name == r {
			return true
		}
	}
	return false
}

// findRoom returns the index of a room on a slice by room name
func findRoom(rms []room, nm string) int {
	for i, r := range rms {
		if r.Name == nm {
			return i
		}
	}
	return -1
}

// findRoutes reads a slice of rooms and appends all routes from start to end to a slice of routes.
// No multiple visits in a room are allowed.
func findRoutes(curRoom room, curRoute route, routes *[]route, rooms *[]room) {

	// reached the end, add to routes
	if curRoom.Role == "end" {
		curRoute = append(curRoute, curRoom.Name)
		toSave := make(route, len(curRoute))
		copy(toSave, curRoute) // copy values to a new route to avoid pointer problems
		*routes = append(*routes, toSave)

		return
	}

	// add new room to current route and proceed
	if !isOnRoute(curRoute, curRoom) {
		curRoute = append(curRoute, curRoom.Name)
		for _, ln := range curRoom.Links {
			nextRoom := (*rooms)[findRoom(*rooms, ln)]
			findRoutes(nextRoom, curRoute, routes, rooms)
		}
	}
}

// sortRoutes sorts a slice of routes from shortest to longest
func sortRoutes(rts *[]route) {
	if len(*rts) < 1 {
		handleError(fmt.Errorf("ERROR: invalid data format, no valid routes"))
	}

	for i := 0; i < len(*rts)-1; i++ {
		for j := i + 1; j < len(*rts); j++ {
			if len((*rts)[i]) > len((*rts)[j]) {
				(*rts)[i], (*rts)[j] = (*rts)[j], (*rts)[i]
			}
		}
	}
}

// areSeparate tells if two routes share intermediary rooms
func areSeparate(rt1, rt2 *route) bool {
	// compare all rooms except start and end
	for _, room1 := range (*rt1)[1 : len(*rt1)-1] {
		for _, room2 := range (*rt2)[1 : len(*rt2)-1] {
			if room1 == room2 {
				return false
			}
		}
	}
	return true
}

// findSeparates recurs through available routes to create combinations of separate routes
func findSeparates(routes, combo []route, allCombos *[][]route, ind int) []route {

	// add this route to the combo
	combo = append(combo, routes[ind])

	// only look at routes after this one to avoid duplicates in different order
	routes = routes[ind+1:]

	// filter out the ones that are no longer separate
	nuRoutes := []route{}
	for _, potentialRoute := range routes {
		separate := true
		for _, foundRoute := range combo {
			if !areSeparate(&foundRoute, &potentialRoute) {
				separate = false
			}
		}
		if separate {
			nuRoutes = append(nuRoutes, potentialRoute)
		}
	}

	// Grow the combo from each available route and add to all combinations
	for i := range nuRoutes {
		*allCombos = append(*allCombos, findSeparates(nuRoutes, combo, allCombos, i))
	}
	return combo
}

func comboAvgLength(combo []route) float64 {
	lens := 0.0
	for _, route := range combo {
		lens += float64(len(route))
	}
	return lens / float64(len(combo))
}

// shortCombos returns all the longest combinations of routes that includes at least one of the shortest routes
func shortCombos(seps [][]route, routes []route) [][]route {

	shortestRoute := len(routes[0])
	longestComboWithShortest := 0
	//var lowestAvgLen float64
	for _, combo := range seps {
		// First route in a combo is always the shortest
		if len(combo[0]) == shortestRoute && len(combo) > longestComboWithShortest {
			longestComboWithShortest = len(combo)
			//lowestAvgLen = comboAvgLength(combo) // get a reasonable non-zero value for best average length
		}
	}

	shorts := [][]route{}
	for _, combo := range seps {
		if len(combo[0]) == shortestRoute && len(combo) == longestComboWithShortest {
			shorts = append(shorts, combo)
		}
	}

	return shorts
}

func lowAverages(seps [][]route) [][]route {

	combosByLength := make(map[int][][]route)
	bestCombosByLength := make(map[int][][]route)
	var longestCombo int
	lowAvgs := [][]route{}

	// Organize combinations by number of routes
	for _, combo := range seps {
		combosByLength[len(combo)] = append(combosByLength[len(combo)], combo)
		if len(combo) > longestCombo {
			longestCombo = len(combo)
		}
	}

	// Same organization but keep only the ones with the lowest average length
	for key, category := range combosByLength {
		bestAvgLen := comboAvgLength(category[0])
		for _, combo := range category {
			if comboAvgLength(combo) < bestAvgLen {
				bestAvgLen = comboAvgLength(combo)
			}
		}
		for _, combo := range category { // There may be several tied for shortest, therefore two loops
			if comboAvgLength(combo) == bestAvgLen {
				bestCombosByLength[key] = append(bestCombosByLength[key], combo)
			}
		}
	}

	// Add all longest ones since longest is always the best for a large amount of ants
	lowAvgs = append(lowAvgs, bestCombosByLength[longestCombo]...)

	// Remove worse solutions: if number of routes is lower and avg length is equal or greater
	benchmark := comboAvgLength(bestCombosByLength[longestCombo][0])
	for i := longestCombo - 1; i > 0; i-- {
		for _, combo := range bestCombosByLength[i] {
			if comboAvgLength(combo) < benchmark {
				lowAvgs = append(lowAvgs, combo)
				benchmark = comboAvgLength(combo)
			}
		}
	}

	/* 	fmt.Println("Low averages:")
	   	for _, combo := range lowAvgs {
	   		fmt.Println(combo)
	   	} */

	return lowAvgs
}

// longCombos returns all the longest combinations of routes with the lowest average lenght
func longCombos(seps [][]route) [][]route {
	longestCombo := 0
	var avgLen float64
	for _, combo := range seps {
		if len(combo) > longestCombo {
			longestCombo = len(combo)
			avgLen = comboAvgLength(combo)
		}
	}

	// Find the lowest average length in the longest combinations
	for _, combo := range seps {
		if len(combo) == longestCombo && comboAvgLength(combo) <= avgLen {
			avgLen = comboAvgLength(combo)
		}
	}

	longs := [][]route{}
	for _, combo := range seps {
		if len(combo) == longestCombo && comboAvgLength(combo) == avgLen {
			longs = append(longs, combo)
		}
	}

	return longs

}

// isSubset checks if combo2 is functionally a subset of or similar to combo1. combo2 must not be longer than combo1.
func isSubset(combo1 []route, combo2 []route) bool {
	for i := range combo2 {
		if len(combo2[i]) != len(combo1[i]) {
			return false
		}
	}
	return true
}

// reduceOptimals removes duplicates
func reduceOptimals(optimals [][]route) [][]route {
	uniques := [][]route{}
	for i := 0; i < len(optimals); i++ {
		found := false
		for j := 0; j < len(uniques); j++ {
			if reflect.DeepEqual(optimals[i], uniques[j]) {
				found = true
			}
		}
		if !found {
			uniques = append(uniques, optimals[i])
		}
	}

	// sort uniques by descending length
	for i := 0; i < len(uniques)-1; i++ {
		for j := i + 1; j < len(uniques); j++ {
			if len(uniques[i]) < len(uniques[j]) {
				uniques[i], uniques[j] = uniques[j], uniques[i]
			}
		}
	}

	// remove functionally similar and functional subsets (same length routes)
	functionalUniques := [][]route{uniques[0]}
	for _, uniq := range uniques {
		allowed := true
		for _, fUniq := range functionalUniques {
			if isSubset(fUniq, uniq) {
				allowed = false
				break
			}
		}
		if allowed {
			functionalUniques = append(functionalUniques, uniq)
		}
	}

	return uniques
}

// bestSolution measures known optimal route combinations for the given number
// of ants and returns the shortest one and its index
func bestSolution(opts [][]route, sAnts [][]ant) ([]route, int) {
	if len(opts) == 1 {
		return opts[0], 0
	}

	longestRoutes := make([]int, len(opts))
	for i, routes := range opts {
		longest := 0
		for _, rt := range routes {
			// count ants on this route
			ants := 0
			for _, ant := range sAnts[i] {
				if reflect.DeepEqual(ant.Route, rt) {
					ants++
				}
			}

			// turns to complete this route (only compare to longest if active)
			turns := len(rt) - 1 + ants
			if ants > 0 && turns > longest {
				longest = turns
			}
		}
		longestRoutes[i] = longest
	}

	// find which optimal route is the quickest for these ants
	quickest := opts[0]
	quickI := 0
	shortestLong := longestRoutes[0]
	for i, n := range longestRoutes {
		if n < shortestLong {
			shortestLong = n
			quickest = opts[i]
			quickI = i
		}
	}

	return quickest, quickI
}
