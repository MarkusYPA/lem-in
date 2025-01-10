package main

import (
	"os"
	"sync"
	"testing"
)

type testCaseGood struct {
	name     string
	input    *os.File
	expected int // number of turns
}

var testCasesGood = []testCaseGood{
	{
		name:     "example00",
		input:    getFile("testcases/example00.txt"),
		expected: 6,
	},
	{
		name:     "example01",
		input:    getFile("testcases/example01.txt"),
		expected: 8,
	},
	{
		name:     "example02",
		input:    getFile("testcases/example02.txt"),
		expected: 11,
	},
	{
		name:     "example03",
		input:    getFile("testcases/example03.txt"),
		expected: 6,
	},
	{
		name:     "example04",
		input:    getFile("testcases/example04.txt"),
		expected: 6,
	},
	{
		name:     "example05",
		input:    getFile("testcases/example05.txt"),
		expected: 8,
	},
	{
		name:     "example06",
		input:    getFile("testcases/example06.txt"),
		expected: 52,
	},
	{
		name:     "example07",
		input:    getFile("testcases/example07.txt"),
		expected: 502,
	},
}

func getFile(s string) *os.File {
	file, err1 := os.Open(s)
	handleError(err1)
	return file
}

func TestMoveAntsGood(t *testing.T) {
	for _, tc := range testCasesGood {
		t.Run(tc.name, func(t *testing.T) {

			nAnts, rooms, err := getStartValues(tc.input)
			handleError(err)
			err = verifyRooms(rooms)
			handleError(err)

			var routes []route
			startRoom := rooms[getStartInd(rooms)]
			findRoutes(startRoom, route{}, &routes, &rooms)
			sortRoutes(&routes)

			combosOfSeparates := [][]route{}
			wg := sync.WaitGroup{}
			for i := range routes {
				wg.Add(1)
				go findSeparates(routes, []route{}, &combosOfSeparates, i, &wg)
			}
			wg.Wait()

			optimals := shortCombos(combosOfSeparates, routes)
			optimals = append(optimals, lowAverages(combosOfSeparates)...)
			optimals = removeRedundant(optimals)
			optiRooms := optimalsToRooms(optimals, &rooms) // optimal routes as slices of rooms instead of slices of room names

			setsOfAnts := makeAnts(optimals, nAnts)
			assignRoutes(optimals, optiRooms, &setsOfAnts)
			optI := bestSolution(optiRooms, setsOfAnts)
			populateStart(&rooms, setsOfAnts[optI])
			turns := moveAnts(setsOfAnts[optI])

			result := len(turns)

			if tc.expected != result {
				t.Errorf("\n\"%s\"\nwant: %v\ngot:  %v", tc.name, tc.expected, result)
			}
		})
	}
}

type testCaseBad struct {
	name     string
	input    *os.File
	expected string // error message
}

var testCasesBad = []testCaseBad{
	{
		name:     "badexample00",
		input:    getFile("testcases/badexample00.txt"),
		expected: "ERROR: invalid data format, invalid number of Ants: 0",
	},
	{
		name:     "badexample01",
		input:    getFile("testcases/badexample01.txt"),
		expected: "ERROR: invalid data format, no valid routes",
	},
	{
		name:     "bad02",
		input:    getFile("testcases/bad02.txt"),
		expected: "ERROR: invalid data format, too many start rooms",
	},
	{
		name:     "bad03",
		input:    getFile("testcases/bad03.txt"),
		expected: "ERROR: invalid data format, too many end rooms",
	},
	{
		name:     "bad04",
		input:    getFile("testcases/bad04.txt"),
		expected: "ERROR: invalid data format, no end room",
	},
	{
		name:     "bad05",
		input:    getFile("testcases/bad05.txt"),
		expected: "ERROR: invalid data format, duplicate room name: 3",
	},
	{
		name:     "bad06",
		input:    getFile("testcases/bad06.txt"),
		expected: "ERROR: invalid data format, bad link: 2 > 5",
	},
	{
		name:     "example00",
		input:    getFile("testcases/example00.txt"),
		expected: "",
	},
	{
		name:     "example01",
		input:    getFile("testcases/example01.txt"),
		expected: "",
	},
	{
		name:     "example02",
		input:    getFile("testcases/example02.txt"),
		expected: "",
	},
	{
		name:     "example03",
		input:    getFile("testcases/example03.txt"),
		expected: "",
	},
	{
		name:     "example04",
		input:    getFile("testcases/example04.txt"),
		expected: "",
	},
	{
		name:     "example05",
		input:    getFile("testcases/example05.txt"),
		expected: "",
	},
	{
		name:     "example06",
		input:    getFile("testcases/example06.txt"),
		expected: "",
	},
	{
		name:     "example07",
		input:    getFile("testcases/example07.txt"),
		expected: "",
	},
}

func TestForErrors(t *testing.T) {
	for _, tc := range testCasesBad {
		t.Run(tc.name, func(t *testing.T) {

			// Errors for wrongly formatted data
			_, rooms, err := getStartValues(tc.input)

			// Errors for invalid data
			if err == nil {
				err = verifyRooms(rooms)
			}

			// Error for 0 routes
			if err == nil {
				var routes []route
				startRoom := rooms[getStartInd(rooms)]
				findRoutes(startRoom, route{}, &routes, &rooms)
				err = sortRoutes(&routes)
			}

			var result string
			if err != nil {
				result = err.Error()
			}

			if tc.expected != result {
				t.Errorf("\n\"%s\"\nwant: %v\ngot:  %v", tc.name, tc.expected, result)
			}
		})
	}
}
