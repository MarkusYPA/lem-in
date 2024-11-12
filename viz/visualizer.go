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
	Name  string   // a number
	Moves []string // target room, max one per turn
}

type PageData struct {
	Ants    []ant
	Turns   []string
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
			if line[0:2] != "##" && line[0] == 'L' && strings.Contains(line, "-") {
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

	for _, ant := range ants {
		prev := startGlob
		for i := range ant.Moves {
			allMovesThisTurn := strings.Fields(turns[i])
			for _, move := range allMovesThisTurn {
				twoParts := strings.Split(move, "-")
				if twoParts[0][1:] == ant.Name {
					ant.Moves[i] = prev + "->" + twoParts[1]
					prev = twoParts[1]
				}
			}
		}
	}
}

func createGVfile(start, end string, rooms, links []string) string {
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

	cmd := exec.Command("dot", "-Tsvg", nameGV, "-o", path+nameSVG)
	_, err = cmd.Output()
	checkErr(err)

	return nameSVG
}

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {

	antMoves := [][]string{}
	for _, ant := range ants {
		antMoves = append(antMoves, ant.Moves)
	}

	data := PageData{
		Ants:    ants,
		Turns:   turnList,
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

	for i := range turns {
		for _, a := range ants {
			fmt.Println(a.Name, a.Moves[i])
		}
		fmt.Println()
	}

	fmt.Println()

	startServer()
}
