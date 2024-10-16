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

// getSepRoutes creates slices of route combinations where the routes don't share rooms on the way from start to end
func getSepRoutes(rts []route) [][]route {
	if len(rts) == 1 {
		return [][]route{rts}
	}

	// For each route, make a slice that includes that one and any of the next routes that don't
	// share any intermediary rooms. Many will be 1 long.
	sepRts := make([][]route, len(rts)-1)
	for i := 0; i < len(rts)-1; i++ {
		sepRts[i] = append(sepRts[i], rts[i])
		for j := i + 1; j < len(rts); j++ {
			separate := true
			for _, foundRt := range sepRts[i] {
				if !areSeparate(foundRt, rts[j]) {
					separate = false
				}
			}
			if separate {
				sepRts[i] = append(sepRts[i], rts[j])
			}
		}
	}
	return sepRts
}

// areSeparate tells if two routes share intermediary rooms
func areSeparate(rt1, rt2 route) bool {
	// compare all rooms except start and end
	for _, room1 := range rt1[1 : len(rt1)-1] {
		for _, room2 := range rt2[1 : len(rt2)-1] {
			if room1 == room2 {
				return false
			}
		}
	}
	return true
}

// shortCombo returns the longest combination of routes that includes at least one of the shortest routes
func shortCombo(seps [][]route, routes []route) []route {
	shortestOfAll := len(routes[0])
	shortComb := seps[0]
	for _, combo := range seps {
		// Looking for a combo with the same length shortest and more routes
		if len(combo[0]) == shortestOfAll && len(combo) > len(seps[0]) {
			shortComb = combo
		}
	}
	return shortComb
}

// longCombo returns the longest combination of routes. In case of tie,
// it returns the one with the fewest links
func longCombo(seps [][]route) []route {
	longest := seps[0]
	for _, comb := range seps {
		if len(comb) > len(longest) {
			longest = comb
		}
		if len(comb) == len(longest) { // Which one has fewer links (steps)
			var links1, links2 int
			for _, rt := range comb {
				links1 += len(rt) - 1
			}
			for _, rt := range longest {
				links2 += len(rt) - 1
			}
			if links1 < links2 {
				longest = comb
			}
		}
	}
	return longest
}

// bestScoreCombo returns the combination of routes with the shortest average length
func bestScoreCombo(seps [][]route) []route {
	best := seps[0]
	bestScore := linksToRoutesScore(best)
	for _, combo := range seps {
		if linksToRoutesScore(combo) < bestScore {
			best = combo
		}
	}
	return best
}

// linksToRoutesScore calculates the links to routes ratio
func linksToRoutesScore(comb []route) int {
	var links int
	for _, rt := range comb {
		links += len(rt) - 1
	}
	return links / len(comb)
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
