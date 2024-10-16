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
func assignRoutes(optimals [][]route, sAnts *[][]ant) {
	for iRs, routes := range optimals {
		// how many ants on each route in this combo
		onRoutes := make([]int, len(routes))

		// loop over the set of ants pertaining to the combo of routes
		for i := 0; i < len((*sAnts)[iRs]); i++ {

			// find the shortest route for this ant (length = route length + ants already taking it)
			shortest := 0
			shortD := len(routes[0]) + onRoutes[0]
			for j := 0; j < len(routes); j++ {
				if len(routes[j])+onRoutes[j] < shortD {
					shortest = j
					shortD = len(routes[j]) + onRoutes[j]
				}
			}
			(*sAnts)[iRs][i].Route = routes[shortest]
			onRoutes[shortest]++
		}
	}
}

// getEndInd returns the index of the "end" room
func getEndInd(rs []room) int {
	for i, r := range rs {
		if r.Role == "end" {
			return i
		}
	}
	return -1
}

// nextRoom finds the next room on an ant's route
func nextRoom(rms []room, curr room, a ant) room {
	var next room
	for i, roomName := range a.Route {
		if roomName == curr.Name && i < len(a.Route)-1 {
			next = rms[findRoom(rms, a.Route[i+1])]
		}
	}
	return next
}

// nextIsOk returns true if the next room has space and
// the route to it hasn't been used on this turn already
func nextIsOk(a ant, rooms []room, usedLinks [][2]string) bool {
	var curr room
	for _, rm := range rooms {
		for _, occ := range rm.Occupants {
			if occ == a.Name {
				curr = rm
			}
		}
	}
	if curr.Role == "end" {
		return false
	}

	next := nextRoom(rooms, curr, a)

	// false if this link was already used on this turn
	for _, link := range usedLinks {
		if [2]string{curr.Name, next.Name} == link {
			return false
		}
	}

	// true if next room is empty or the end
	return len(next.Occupants) < 1 || next.Role == "end"
}

// removeFromRoom removes ant a from the lists of occupants of all rooms and returns the current and next rooms
func removeFromRoom(rms []room, a ant) (room, room, []room) {
	roomsOut := []room{}
	var curr room
	for _, rm := range rms {
		nuRm := rm
		nuRm.Occupants = nil
		for _, occ := range rm.Occupants {
			if occ != a.Name {
				nuRm.Occupants = append(nuRm.Occupants, occ)
			} else {
				curr = rms[findRoom(rms, rm.Name)]
			}
		}
		roomsOut = append(roomsOut, nuRm)
	}
	next := nextRoom(rms, curr, a)
	return curr, next, roomsOut
}

// addToRoom adds an ant to the list of occupants in a room
func addToRoom(rms *[]room, nxt *room, a ant) {
	for i := 0; i < len(*rms); i++ {
		if (*rms)[i].Name == nxt.Name {
			(*rms)[i].Occupants = append((*rms)[i].Occupants, a.Name)
		}
	}
}

// moveAnts moves the ants across the farm and returns the commands to do so
func moveAnts(rms *[]room, ants []ant) []string {
	turns := []string{}

	// move ants until all are in the last room
	for len((*rms)[getEndInd(*rms)].Occupants) < len(ants) {
		moves := ""
		linksUsed := [][2]string{}

		// try to move each ant
		for i := 0; i < len(ants); i++ {

			if nextIsOk(ants[i], *rms, linksUsed) {
				// remove ant from any and all rooms, and mark this link as used
				var curr, next room
				curr, next, *rms = removeFromRoom(*rms, ants[i])
				linksUsed = append(linksUsed, [2]string{curr.Name, next.Name})

				// put ant to next room
				addToRoom(rms, &next, ants[i])

				// add move to current turn
				moves += "L" + strconv.Itoa(ants[i].Name) + "-" + next.Name + " "
			}
		}

		if moves != "" {
			turns = append(turns, moves[:len(moves)-1]) // append this turn without the last space character
		}
	}

	return turns
}
