package main

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	sepMu sync.Mutex
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
func findRoom(rooms []room, name string) int {
	for i, r := range rooms {
		if r.Name == name {
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
		for _, link := range curRoom.Links {
			nextRoom := (*rooms)[findRoom(*rooms, link)]
			findRoutes(nextRoom, curRoute, routes, rooms)
		}
	}
}

// sortRoutes sorts a slice of routes from shortest to longest
func sortRoutes(rts *[]route) error {
	if len(*rts) < 1 {
		return fmt.Errorf("ERROR: invalid data format, no valid routes")
	}

	for i := 0; i < len(*rts)-1; i++ {
		for j := i + 1; j < len(*rts); j++ {
			if len((*rts)[i]) > len((*rts)[j]) {
				(*rts)[i], (*rts)[j] = (*rts)[j], (*rts)[i]
			}
		}
	}

	return nil
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
func findSeparates(routes, curCombo []route, combosOfSeparates *[][]route, ind int, wg *sync.WaitGroup) {

	// add this route to the current combination
	curCombo = append(curCombo, routes[ind])

	// only look at routes after this one to avoid duplicates in different order
	routes = routes[ind+1:]

	// filter out the ones that are no longer separate
	newRoutes := []route{}
	for _, potentialRoute := range routes {
		separate := true
		for _, foundRoute := range curCombo {
			if !areSeparate(&foundRoute, &potentialRoute) {
				separate = false
				break
			}
		}
		if separate {
			newRoutes = append(newRoutes, potentialRoute)
		}
	}

	// Grow the combo from each available route and add to all combinations
	for i := range newRoutes {
		wg.Add(1)
		go findSeparates(newRoutes, curCombo, combosOfSeparates, i, wg)
	}

	sepMu.Lock()
	*combosOfSeparates = append(*combosOfSeparates, curCombo)
	sepMu.Unlock()

	wg.Done()
}

// comboAvgLength calculates the average length of a slice of routes
func comboAvgLength(combo []route) float64 {
	lens := 0.0
	for _, route := range combo {
		lens += float64(len(route))
	}
	return lens / float64(len(combo))
}

// shortCombos returns all the longest combinations of routes that includes at least one of the shortest routes
func shortCombos(combosOfSeparates [][]route, routes []route) [][]route {

	shortestLength := len(routes[0])
	longestComboWithShortest := 0
	for _, combo := range combosOfSeparates {
		// First route in a combo is always the shortest
		if len(combo[0]) == shortestLength && len(combo) > longestComboWithShortest {
			longestComboWithShortest = len(combo)
		}
	}

	shorts := [][]route{}
	for _, combo := range combosOfSeparates {
		if len(combo[0]) == shortestLength && len(combo) == longestComboWithShortest {
			shorts = append(shorts, combo)
		}
	}

	return shorts
}

// lowAverages finds the lowest average length combinations for each number of routes
func lowAverages(combosOfSeparates [][]route) [][]route {

	combosByLength := make(map[int][][]route)
	bestCombosByLength := make(map[int][][]route)
	var longestCombo int
	lowAvgs := [][]route{}

	// Organize combinations by number of routes
	for _, combo := range combosOfSeparates {
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

	// Add the best longest ones since longest is always the best for a large amount of ants
	lowAvgs = append(lowAvgs, bestCombosByLength[longestCombo]...)

	// Add shorter combinations if average length is lower
	benchmark := comboAvgLength(bestCombosByLength[longestCombo][0])
	for i := longestCombo - 1; i > 0; i-- {
		for _, combo := range bestCombosByLength[i] {
			if comboAvgLength(combo) < benchmark {
				lowAvgs = append(lowAvgs, combo)
				benchmark = comboAvgLength(combo)
			}
		}
	}

	return lowAvgs
}

// isSubset checks if combo2 is functionally a subset of or similar to combo1.
func isSubset(combo1 []route, combo2 []route) bool {
	if len(combo1) < len(combo2) {
		return false
	}
	for i := range combo2 {
		if len(combo2[i]) != len(combo1[i]) {
			return false
		}
	}
	return true
}

// removeRedundant removes duplicates, subsets and functionally equal combinations
func removeRedundant(optimals [][]route) [][]route {
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

	// remove functionally similar and functionally subsets (same length routes)
	functionalUniques := [][]route{uniques[0]}
	for _, uniq := range uniques {
		found := true
		for _, fUniq := range functionalUniques {
			if isSubset(fUniq, uniq) {
				found = false
				break
			}
		}
		if found {
			functionalUniques = append(functionalUniques, uniq)
		}
	}

	return functionalUniques
}

func optimalsToRooms(optimals [][]route, rooms *[]room) [][][](*room) {

	roomCombos := [][][](*room){} // multiple combinations of routes

	for i, combo := range optimals {
		roomCombos = append(roomCombos, [][](*room){}) // combination of routes

		for j, route := range combo {
			roomCombos[i] = append(roomCombos[i], [](*room){}) // one route

			for _, roomName := range route {
				thisRoom := &(*rooms)[findRoom(*rooms, roomName)]
				roomCombos[i][j] = append(roomCombos[i][j], thisRoom)
			}
		}
	}

	return roomCombos
}

// bestSolution measures known optimal route combinations for the given number
// of ants and returns the shortest one and its index
func bestSolution(optimals [][]route, setsOfAnts [][]ant) int {
	if len(optimals) == 1 {
		return 0
	}

	longestRoutes := make([]int, len(optimals))
	for i, combo := range optimals {
		longest := 0
		for _, route := range combo {
			// count ants on this route
			ants := 0
			for _, ant := range setsOfAnts[i] {
				if reflect.DeepEqual(ant.Route, route) {
					ants++
				}
			}

			// turns to complete this route (only compare to longest if active)
			turns := len(route) - 1 + ants
			if ants > 0 && turns > longest {
				longest = turns
			}
		}
		longestRoutes[i] = longest
	}

	// find which optimal route is the quickest for these ants
	quickI := 0
	shortestLong := longestRoutes[0]
	for i, n := range longestRoutes {
		if n < shortestLong {
			shortestLong = n
			quickI = i
		}
	}

	return quickI
}
