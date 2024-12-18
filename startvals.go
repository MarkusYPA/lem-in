package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// removeCarRet removes all carriage returns from a string
func removeCarRet(s string) (out string) {
	for _, ru := range s {
		if ru != 13 {
			out += string(ru)
		}
	}
	return
}

// getStartValues reads the number of ants and the slice of rooms from a text file
func getStartValues(s string) (int, []room) {
	lines := strings.Split(s, "\n")
	ants := 0
	var err error
	rooms := []room{}

	for i, l := range lines {
		reRoom := regexp.MustCompile(`^\w+\s[-+]?\d+\s[-+]?\d+$`) // "string int int"
		reLink := regexp.MustCompile(`^\w+-\w+$`)                 // "string-string"

		if i == 0 {
			ants, err = strconv.Atoi(l)
			if err != nil || ants < 1 {
				handleError(fmt.Errorf("ERROR: invalid data format, invalid number of Ants: " + l))
			}
			continue
		}

		if len(l) > 0 && l[0] == '#' {
			continue
		}

		// If it's a room
		if reRoom.MatchString(l) {
			var thisRoom room

			if len(lines[i-1]) > 6 && lines[i-1][0:7] == "##start" {
				thisRoom.Role = "start"
			} else if len(lines[i-1]) > 4 && lines[i-1][0:5] == "##end" {
				thisRoom.Role = "end"
			} else {
				thisRoom.Role = "normal"
			}

			roomWds := strings.Fields(l)
			if len(roomWds) == 3 {
				thisRoom.Name = roomWds[0]
				var errC1 error
				thisRoom.Coords[0], errC1 = strconv.Atoi(roomWds[1])
				if errC1 != nil {
					handleError(fmt.Errorf("ERROR: invalid data format, " + errC1.Error()))
				}
				var errC2 error
				thisRoom.Coords[1], errC2 = strconv.Atoi(roomWds[2])
				if errC2 != nil {
					handleError(fmt.Errorf("ERROR: invalid data format, " + errC2.Error()))
				}
			} else {
				handleError(fmt.Errorf("ERROR: invalid data format, " + l))
			}

			thisRoom.Occupants = make(map[int]bool)
			rooms = append(rooms, thisRoom)
		}

		// If it's a link
		if reLink.MatchString(l) {
			pair := strings.Split(l, "-") // When regexp matches, we have two strings separated by a dash

			for j, ro := range rooms {
				if ro.Name == pair[0] {
					rooms[j].Links = append(rooms[j].Links, pair[1])
				}
				if ro.Name == pair[1] {
					rooms[j].Links = append(rooms[j].Links, pair[0])
				}
			}
		}
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
				handleError(fmt.Errorf("ERROR: invalid data format, duplicate room name: " + rooms[i].Name))
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
				handleError(fmt.Errorf("ERROR: invalid data format, bad link: " + rooms[i].Name + " > " + ln))
			}
		}
	}

	if starts != 1 {
		if starts == 0 {
			handleError(fmt.Errorf("ERROR: invalid data format, no start room found"))
		} else {
			handleError(fmt.Errorf("ERROR: invalid data format, too many start rooms"))
		}
	}

	if ends != 1 {
		if ends == 0 {
			handleError(fmt.Errorf("ERROR: invalid data format, no end room found"))
		} else {
			handleError(fmt.Errorf("ERROR: invalid data format, too many end rooms"))
		}
	}
}
