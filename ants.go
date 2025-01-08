package main

import (
	"strconv"
)

// makeAnts makes a group of ants for each combination of routes to be tested
func makeAnts(optimals [][]route, n int) [][]ant {
	setsOfAnts := [][]ant{}

	for i := range optimals {
		setsOfAnts = append(setsOfAnts, []ant{})
		for j := range n {
			setsOfAnts[i] = append(setsOfAnts[i], ant{Name: j + 1})
		}
	}
	return setsOfAnts
}

// assignRoutes gives each ant a route to follow
func assignRoutes(optimals [][]route, optiRooms [][][]*room, setsOfAnts *[][]ant, startRoom *room, rooms *[]room) {
	for i, routeCombo := range optimals {
		// how many ants on each route in this combo
		onRoutes := make([]int, len(routeCombo))

		// loop over the set of ants pertaining to the combo of routes
		for j := 0; j < len((*setsOfAnts)[i]); j++ {
			// find the shortest route for this ant (length = route length + ants already taking it)
			shortest := 0
			shortD := len(routeCombo[0]) + onRoutes[0]
			for k := 0; k < len(routeCombo); k++ {
				if len(routeCombo[k])+onRoutes[k] < shortD {
					shortest = k
					shortD = len(routeCombo[k]) + onRoutes[k]
				}
			}
			(*setsOfAnts)[i][j].Route = routeCombo[shortest]
			(*setsOfAnts)[i][j].Route2 = optiRooms[i][shortest]

			onRoutes[shortest]++
		}
	}
}

// nextRoom finds the next room on an ant's route
func nextRoom(rms *[]room, curr room, a ant) *room {
	var next *room
	for i, roomName := range a.Route {
		if roomName == curr.Name && i < len(a.Route)-1 {
			next = &(*rms)[findRoom((*rms), a.Route[i+1])]
		}
	}
	return next
}

// nextIsOk returns true if the next room has space and
// the route to it hasn't been used on this turn already
func nextIsOk(a ant, rooms *[]room, usedLinks [][2]string) (bool, *room, *room) {
	var curr *room
	var next *room

	curr = a.Route2[a.routeIndex]

	if curr.Role == "end" {
		return false, curr, next
	}

	next = a.Route2[a.routeIndex+1]

	// false if this link was already used on this turn
	for _, link := range usedLinks {
		if [2]string{curr.Name, next.Name} == link {
			return false, curr, next
		}
	}

	// true if next room is empty or the end
	nextIsEmpty := len(next.Occupants) < 1
	return nextIsEmpty || next.Role == "end", curr, next
}

// moveAnts moves the ants across the farm and returns the commands to do so
func moveAnts(rms *[]room, ants []ant) []string {
	turns := []string{}
	antsAtEnd := 0

	// move ants until all are in the last room
	for antsAtEnd < len(ants) {
		moves := ""
		linksUsed := [][2]string{}

		// try to move each ant
		for i := 0; i < len(ants); i++ {
			if !ants[i].atEnd {
				NextOk, currentRoom, nextRoom := nextIsOk(ants[i], rms, linksUsed)
				if NextOk {
					delete(currentRoom.Occupants, ants[i].Name)
					linksUsed = append(linksUsed, [2]string{currentRoom.Name, nextRoom.Name}) // mark this link as used
					nextRoom.Occupants[ants[i].Name] = true
					ants[i].routeIndex++
					if nextRoom.Role == "end" {
						ants[i].atEnd = true
						antsAtEnd++
					}
					// add move to current turn
					moves += "L" + strconv.Itoa(ants[i].Name) + "-" + nextRoom.Name + " "
				}
			}
		}

		turns = append(turns, moves[:len(moves)-1]) // append this turn without the last space character
	}

	return turns
}
