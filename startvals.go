package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// getStartValues reads the number of ants and the slice of rooms from a text file
func getStartValues(file *os.File) (int, []room) {
	ants := 0
	var err error
	rooms := []room{}

	prev := ""
	scanner := bufio.NewScanner(file)
	i := -1
	for scanner.Scan() {
		i++
		line := scanner.Text()

		reRoom := regexp.MustCompile(`^\w+\s[-+]?\d+\s[-+]?\d+$`) // "string int int"
		reLink := regexp.MustCompile(`^\w+-\w+$`)                 // "string-string"

		if i == 0 {
			ants, err = strconv.Atoi(line)
			if err != nil || ants < 1 {
				handleError(errors.New("ERROR: invalid data format, invalid number of Ants: " + line))
			}
			prev = line
			continue
		}

		if len(line) > 0 && line[0] == '#' {
			prev = line
			continue
		}

		// If it's a room
		if reRoom.MatchString(line) {
			var thisRoom room

			if len(prev) > 6 && prev[0:7] == "##start" {
				thisRoom.Role = "start"
			} else if len(prev) > 4 && prev[0:5] == "##end" {
				thisRoom.Role = "end"
			} else {
				thisRoom.Role = "normal"
			}

			roomWds := strings.Fields(line)
			if len(roomWds) == 3 {
				thisRoom.Name = roomWds[0]
				var errC1 error
				thisRoom.Coords[0], errC1 = strconv.Atoi(roomWds[1])
				if errC1 != nil {
					handleError(errors.New("ERROR: invalid data format, " + errC1.Error()))
				}
				var errC2 error
				thisRoom.Coords[1], errC2 = strconv.Atoi(roomWds[2])
				if errC2 != nil {
					handleError(errors.New("ERROR: invalid data format, " + errC2.Error()))
				}
			} else {
				handleError(errors.New("ERROR: invalid data format, " + line))
			}

			thisRoom.Occupants = make(map[int]bool)
			rooms = append(rooms, thisRoom)
		}

		// If it's a link
		if reLink.MatchString(line) {
			pair := strings.Split(line, "-") // When regexp matches, we have two strings separated by a dash

			for j, ro := range rooms {
				if ro.Name == pair[0] {
					rooms[j].Links = append(rooms[j].Links, pair[1])
				}
				if ro.Name == pair[1] {
					rooms[j].Links = append(rooms[j].Links, pair[0])
				}
			}
		}
		prev = line
	}

	return ants, rooms
}

// verifyRooms makes sure there is one start and one end and no duplicate room names
func verifyRooms(rooms []room) {
	starts := 0
	ends := 0
	for i := 0; i < len(rooms); i++ {
		if rooms[i].Role == "start" {
			starts++
		}
		if rooms[i].Role == "end" {
			ends++
		}
		for j := i + 1; j < len(rooms); j++ {
			if rooms[i].Name == rooms[j].Name {
				handleError(errors.New("ERROR: invalid data format, duplicate room name: " + rooms[i].Name))
			}
		}

		for _, ln := range rooms[i].Links {
			found := false
			for _, rm := range rooms {
				if ln == rm.Name {
					found = true
				}
			}
			if !found {
				handleError(errors.New("ERROR: invalid data format, bad link: " + rooms[i].Name + " > " + ln))
			}
		}
	}

	if starts != 1 {
		if starts == 0 {
			handleError(errors.New("ERROR: invalid data format, no start room found"))
		} else {
			handleError(errors.New("ERROR: invalid data format, too many start rooms"))
		}
	}

	if ends != 1 {
		if ends == 0 {
			handleError(errors.New("ERROR: invalid data format, no end room found"))
		} else {
			handleError(errors.New("ERROR: invalid data format, too many end rooms"))
		}
	}
}
