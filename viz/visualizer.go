package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ant struct {
	Name     string   // a number
	Moves    []string // target room, max one per turn
	Sequence string
}

type PageData struct {
	Ants    []ant
	Turns   [][2]string
	Moves   [][]string
	Drawing string
	Start   string
	End     string
}

var ants []ant
var gvname string
var turnList []string
var startGlob string
var endGlob string

func checkErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func scanInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	input := []string{}
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	return input
}

func parseInput(input []string) (string, string, []string, []string, []string) {
	start := ""
	end := ""
	links := []string{}
	turns := []string{}
	rooms := []string{}
	var readingMoves bool
	for i, line := range input {
		if !readingMoves {
			if len(line) > 1 && line[0:2] != "##" && !strings.Contains(line, " ") && strings.Contains(line, "-") {
				twoRooms := strings.Split(line, "-")
				links = append(links, twoRooms[0]+" -> "+twoRooms[1]+";")
			}
			if len(line) > 1 && line[0:2] != "##" && strings.Contains(line, " ") && !strings.Contains(line, "-") {
				roomInfo := strings.Split(line, " ")
				rooms = append(rooms, roomInfo[0])
			}
			if line == "##start" {
				start = strings.Fields(input[i+1])[0]
			}
			if line == "##end" {
				end = strings.Fields(input[i+1])[0]
			}
		} else {
			if len(line) > 1 && line[0:2] != "##" && line[0] == 'L' && strings.Contains(line, "-") {
				turns = append(turns, line)
			}
		}
		if line == "" {
			readingMoves = true
		}
	}
	return start, end, links, turns, rooms
}

func makeAnts(amount int, turns []string) {
	for i := range amount {
		ants = append(ants, ant{Name: strconv.Itoa(i + 1), Moves: make([]string, len(turns))})
	}

	for h, ant := range ants {
		prev := startGlob
		for i := range ant.Moves {
			allMovesThisTurn := strings.Fields(turns[i])
			for _, move := range allMovesThisTurn {
				twoParts := strings.Split(move, "-")
				if twoParts[0][1:] == ant.Name {
					ant.Moves[i] = prev + "->" + twoParts[1]
					if ants[h].Sequence == "" {
						ants[h].Sequence += prev + " > " + twoParts[1]
					} else {
						ants[h].Sequence += " > " + twoParts[1]
					}
					prev = twoParts[1]
				}
			}
		}
	}
}

func createGVfile(start, end string, rooms, links []string) string {

	// Some links (and so svg paths) are back-to-front for my purposes
	// Write down format for link from ant moves, original and reverse
	differentMoves := [][2]string{}
	for _, ant := range ants {
		for _, m := range ant.Moves {
			if len(m) > 0 {
				parts := strings.Split(m, "->")
				newForm1 := parts[0] + " -> " + parts[1]
				newForm2 := parts[1] + " -> " + parts[0]
				differentMoves = append(differentMoves, [2]string{newForm1, newForm2})
			}
		}
	}

	// Change link direction to ant movement direction when necessary
	for _, mv := range differentMoves {
		for j, link := range links {
			if link == mv[1]+";" {
				links[j] = mv[0]
			}
		}
	}

	gvfile := "digraph G {\n    ratio=1;\n    pad=0.5;\n    edge [arrowhead=none];\n\n"
	gvfile += "    " + start + " [shape=box];\n" + "    " + end + " [shape=box];\n\n"
	for _, room := range rooms {
		gvfile += "    " + room + " " + "[fixedsize=true height= 0, width=0 color=\"transparent\"];\n"
	}
	gvfile += "\n"
	for _, ln := range links {
		gvfile += "    " + ln + "\n"
	}
	gvfile += "}"

	// Make a .gv file for graphviz
	nameGV := "graph.gv"
	nameSVG := "graph.svg"
	path := "static/"
	err := os.WriteFile(nameGV, []byte(gvfile), 0644)
	checkErr(err)

	// Run GraphViz to create the .svg map of the ant farm
	cmd := exec.Command("dot", "-Tsvg", nameGV, "-o", path+nameSVG)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to create static/graph.svg, use a compatible premade graph from premadegraphs/")
	}

	return nameSVG
}

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {

	antMoves := [][]string{}
	for _, ant := range ants {
		antMoves = append(antMoves, ant.Moves)
	}

	doubleTurns := [][2]string{}
	for i, t := range turnList {
		doubleTurns = append(doubleTurns, [2]string{strconv.Itoa(i + 1), t})
	}

	data := PageData{
		Ants:    ants,
		Turns:   doubleTurns,
		Moves:   antMoves,
		Drawing: gvname,
		Start:   startGlob,
		End:     endGlob,
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound) // error 404
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request", http.StatusNotFound)
		return
	}

	tpl.Execute(w, data)
}

func startServer() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", homeHandler)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	input := scanInput()
	start, end, links, turns, rooms := parseInput(input)
	turnList = turns
	startGlob = start
	endGlob = end

	antsAmount, err := strconv.Atoi(input[0])
	checkErr(err)

	makeAnts(antsAmount, turns)
	gvname = createGVfile(start, end, rooms, links)

	startServer()
}
