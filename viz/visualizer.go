package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ant struct {
	name  string   // a number
	moves []string // target room, max one per turn
}

func checkErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	output := []string{}
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	antsAmount, err := strconv.Atoi(output[0])
	checkErr(err)

	links := []string{}
	turns := []string{}
	var readingMoves bool
	for _, line := range output {
		if !readingMoves {
			if len(line) > 1 && line[0:2] != "##" && !strings.Contains(line, " ") && strings.Contains(line, "-") {
				twoRooms := strings.Split(line, "-")
				links = append(links, twoRooms[0]+" --> "+twoRooms[1]+";")
			}
		} else {
			if line[0:2] != "##" && line[0] == 'L' && strings.Contains(line, "-") {
				turns = append(turns, line)
			}
		}
		if line == "" {
			readingMoves = true
		}
	}

	ants := []ant{}
	for i := range antsAmount {
		ants = append(ants, ant{name: strconv.Itoa(i + 1), moves: make([]string, len(turns))})
	}

	for _, ant := range ants {
		for i := range ant.moves {
			allMovesThisTurn := strings.Fields(turns[i])
			for _, move := range allMovesThisTurn {
				twoParts := strings.Split(move, "-")
				if twoParts[0][1:] == ant.name {
					ant.moves[i] = twoParts[1]
				}
			}
		}
	}

	fmt.Println(links)
	for _, ant := range ants {
		fmt.Println(ant.name, ant.moves)
	}

}
