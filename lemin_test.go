package main

import (
	"os"
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected int // number of turns
}

var testCases = []testCase{
	{
		name:     "example00",
		input:    fileToString("testcases/example00.txt"),
		expected: 6,
	},
	{
		name:     "example01",
		input:    fileToString("testcases/example01.txt"),
		expected: 8,
	},
	{
		name:     "example02",
		input:    fileToString("testcases/example02.txt"),
		expected: 11,
	},
	{
		name:     "example03",
		input:    fileToString("testcases/example03.txt"),
		expected: 6,
	},
	{
		name:     "example04",
		input:    fileToString("testcases/example04.txt"),
		expected: 6,
	},
	{
		name:     "example05",
		input:    fileToString("testcases/example05.txt"),
		expected: 8,
	},
	{
		name:     "example06",
		input:    fileToString("testcases/example06.txt"),
		expected: 52,
	},
	{
		name:     "example07",
		input:    fileToString("testcases/example07.txt"),
		expected: 502,
	},
}

func fileToString(s string) string {
	file, _ := os.ReadFile(s)
	return removeCarRet(string(file))
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
			for i := range routes {
				combosOfSeparates = append(combosOfSeparates, findSeparates(routes, []route{}, &combosOfSeparates, i))
			}

			optimals := shortCombos(combosOfSeparates, routes)
			//optimals = append(optimals, longCombos(combosOfSeparates)...)
			optimals = append(optimals, lowAverages(combosOfSeparates)...)
			optimals = reduceOptimals(optimals)

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
