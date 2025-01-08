package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type room struct {
	Name      string
	Coords    [2]int
	Occupants map[int]bool
	Links     []string // neighbouring room names
	Role      string   // "start", "normal" or "end"
}

type ant struct {
	Name       int
	Route      route
	Route2     [](*room)
	routeIndex int
	atEnd      bool
}

type route []string

// handleError is the default error handling
func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

// getStartInd returns the index of the "start" room
func getStartInd(rs []room) int {
	for i, r := range rs {
		if r.Role == "start" {
			return i
		}
	}
	return -1
}

// populateStart puts all the ants in the start room
func populateStart(rooms *[]room, ants []ant) {
	start := &(*rooms)[getStartInd(*rooms)]
	for _, a := range ants {
		start.Occupants[a.Name] = true
	}
}

// printSolution prints the contents of the file and the moves taken
func printSolution(file *os.File, turns []string) {
	// Reset file pointer to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("ERROR: seeking file failed:", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println()

	for _, v := range turns {
		fmt.Println(v)
	}
}

func main() {

	if len(os.Args) != 2 {
		handleError(fmt.Errorf("ERROR: provide the input file in one argument"))
	}
	file, err1 := os.Open(os.Args[1])
	handleError(err1)
	defer file.Close()

	// Read and save the number of ants and information about the rooms
	nAnts, rooms, err := getStartValues(file)
	handleError(err)
	err = verifyRooms(rooms)
	handleError(err)

	// Find all routes connecting "start" to "end"
	var routes []route
	startRoom := rooms[getStartInd(rooms)]
	findRoutes(startRoom, route{}, &routes, &rooms)
	err = sortRoutes(&routes)
	handleError(err)

	// Find all combinations of non-crossing routes
	combosOfSeparates := [][]route{}
	wg := sync.WaitGroup{}
	for i := range routes {
		wg.Add(1)
		go findSeparates(routes, []route{}, &combosOfSeparates, i, &wg)
	}
	wg.Wait()

	/*
		Two crossing routes work effectively as one single route because of the
		bottleneck, so we focus only on combinations of separate routes

		Optimal route combinations:
		- Combinations with the shortest route (best option for one ant)
		- Combinations with the lowest average length for each number of routes (sometimes best for an average amount of ants)
		- Combinations with the most routes (best for a large amount of ants)
	*/

	optimals := shortCombos(combosOfSeparates, routes)
	optimals = append(optimals, lowAverages(combosOfSeparates)...) // lowAverages() also finds the longest combinations
	optimals = removeRedundant(optimals)
	optiRooms := optimalsToRooms(optimals, &rooms) // optimal routes as slices of rooms instead of slices of room names

	setsOfAnts := makeAnts(optimals, nAnts)
	assignRoutes(optimals, optiRooms, &setsOfAnts, &startRoom, &rooms)
	optI := bestSolution(optimals, setsOfAnts)
	populateStart(&rooms, setsOfAnts[optI])

	// Move ants and save the moves
	turns := moveAnts(&rooms, setsOfAnts[optI])

	// Print out the file contents and the moves
	//printSolution(file, turns)
	fmt.Println("Turns taken:", len(turns)) // for testing
}
