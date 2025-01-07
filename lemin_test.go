package main

import (
	"os"
	"sync"
	"testing"
)

type testCase struct {
	name     string
	input    *os.File
	expected int // number of turns
}

var testCases = []testCase{
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
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			nAnts, rooms := getStartValues(tc.input)
			verifyRooms(rooms)
			var routes []route
			findRoutes(rooms[getStartInd(rooms)], route{}, &routes, &rooms)
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

			setsOfAnts := makeAnts(optimals, nAnts)
			assignRoutes(optimals, &setsOfAnts)
			_, optI := bestSolution(optimals, setsOfAnts)
			populateStart(&rooms, setsOfAnts[optI])
			turns := moveAnts(&rooms, setsOfAnts[optI])

			result := len(turns)

			if tc.expected != result {
				t.Errorf("\n\"%s\"\nwant: %v\ngot:  %v", tc.name, tc.expected, result)
			}
		})
	}
}
