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
func assignRoutes(optimals [][]route, setsOfAnts *[][]ant) {
	//wgCombos := sync.WaitGroup{}
	for i, routeCombo := range optimals {
		//wgCombos.Add(1)
		//go func() {
		// how many ants on each route in this combo
		onRoutes := make([]int, len(routeCombo))
		//wgAnts := sync.WaitGroup{}
		// loop over the set of ants pertaining to the combo of routes
		for j := 0; j < len((*setsOfAnts)[i]); j++ {
			//wgAnts.Add(1)
			//go func() {
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
			onRoutes[shortest]++
			//wgAnts.Done()
			//}()
		}
		//wgAnts.Wait()
		//wgCombos.Done()
		//}()
	}
	//wgCombos.Wait()
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
	for i, rm := range *rooms {
		for _, occ := range rm.Occupants {
			if occ == a.Name {
				curr = &(*rooms)[i]
			}
		}
	}
	if curr.Role == "end" {
		return false, curr, next
	}

	next = nextRoom(rooms, *curr, a)

	// false if this link was already used on this turn
	for _, link := range usedLinks {
		if [2]string{curr.Name, next.Name} == link {
			return false, curr, next
		}
	}

	// true if next room is empty or the end
	return len(next.Occupants) < 1 || next.Role == "end", curr, next
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
	next := nextRoom(&rms, curr, a)
	return curr, *next, roomsOut
}

func removeAntFromRoom(thisRoom *room, a ant) {
	//before := time.Now()
	newOccupants := []int{}
	for _, oc := range thisRoom.Occupants {
		//fmt.Println(time.Since(before), "0")
		if oc != a.Name {
			newOccupants = append(newOccupants, oc)
		}
		//fmt.Println(time.Since(before), "1")
		//before = time.Now()
	}
	thisRoom.Occupants = newOccupants
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
			NextOk, currentRoom, nextRoom := nextIsOk(ants[i], rms, linksUsed)
			if NextOk {
				removeAntFromRoom(currentRoom, ants[i])
				linksUsed = append(linksUsed, [2]string{currentRoom.Name, nextRoom.Name}) // mark this link as used
				addToRoom(rms, nextRoom, ants[i])

				// add move to current turn
				moves += "L" + strconv.Itoa(ants[i].Name) + "-" + nextRoom.Name + " "
			}
		}

		if moves != "" {
			turns = append(turns, moves[:len(moves)-1]) // append this turn without the last space character
		}
	}

	return turns
}
