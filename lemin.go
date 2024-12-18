package main

import (
	"fmt"
	"os"
)

type room struct {
	Name      string
	Coords    [2]int
	Occupants []int    // ant names are integers
	Links     []string // neighbouring room names
	Role      string   // "start", "normal" or "end"
}

type ant struct {
	Name  int
	Route route
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
		start.Occupants = append(start.Occupants, a.Name)
	}
}

// printSolution prints the contents of the file and the moves taken
func printSolution(file string, moves []string) {
	fmt.Printf("%s\n\n", file)
	for _, v := range moves {
		fmt.Println(v)
	}
}

func main() {
	if len(os.Args) != 2 {
		handleError(fmt.Errorf("ERROR: provide the input file in one argument"))
	}
	in, err1 := os.ReadFile(os.Args[1])
	handleError(err1)

	// read and save the number of ants and information about the rooms
	nAnts, rooms := getStartValues(removeCarRet(string(in)))
	verifyRooms(rooms)

	// find all routes connecting "start" to "end" and all unique combinations of non-crossing routes
	var routes []route
	findRoutes(rooms[getStartInd(rooms)], route{}, &routes, &rooms)
	//fmt.Println("Routes", time.Now().Format("05.00"))
	sortRoutes(&routes)
	separateRoutes := getSepRoutes(routes)

	/*
		Two crossing routes work effectively as one single route because of the
		bottleneck, so we focus only on combinations of separate routes

		Optimal route combinations include:
		- A combination with the shortest route (always the best option for one ant)
		- A combination with the most routes (best option for a large amount of ants)
		- A combination with the lowest average route length (possibly for a medium amount of ants)
	*/

	optimals := reduceOptimals([][]route{shortCombo(separateRoutes, routes), longCombo(separateRoutes), bestScoreCombo(separateRoutes)})
	//fmt.Println(time.Now().Format("05.00"), "Optimals")
	setsOfAnts := makeAnts(optimals, nAnts)
	//fmt.Println(time.Now().Format("05.00"), "Ants made")
	assignRoutes(optimals, &setsOfAnts)
	//fmt.Println(time.Now().Format("05.00"), "Assigned")
	_, optI := bestSolution(optimals, setsOfAnts)
	//fmt.Println(time.Now().Format("05.00"), "Best")
	populateStart(&rooms, setsOfAnts[optI])
	//fmt.Println(time.Now().Format("05.00"), "Populated")

	// Move ants and save the moves
	turns := moveAnts(&rooms, setsOfAnts[optI])
	//fmt.Println(time.Now().Format("05.00"), "Moved")

	// Print out the file contents and the moves
	printSolution(string(in), turns)
	//fmt.Println()
	//fmt.Println(len(turns), "turns")
}
