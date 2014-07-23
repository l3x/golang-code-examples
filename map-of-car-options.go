/*
filename:  map-of-car-options.go
author:    Lex Sheehan
copyright: Lex Sheehan LLC
license:   GPL
status:    published
comments:  http://l3x.github.io/golang-code-examples/2014/07/22/map-of-car-options.html
*/
package main

import (
	"fmt"
	"strings"
)

type Car struct {
	Make  string
	Model  string
	Options []string
}

func main() {

	dashes := strings.Repeat("-", 50)

	is250 := &Car{"Lexus", "IS250", []string{"GPS", "Alloy Wheels", "Roof Rack", "Power Outlets", "Heated Seats"}}
	accord := &Car{"Honda", "Accord", []string{"Alloy Wheels", "Roof Rack"}}
	blazer := &Car{"Chevy", "Blazer", []string{"GPS", "Roof Rack", "Power Outlets"}}

	cars := []*Car{is250, accord, blazer}
	fmt.Printf("Cars:\n%v\n\n", cars)  // cars is a slice of pointers to our three cars

	// Create a map to associate options with each car
	car_options := make(map[string][]*Car)

	fmt.Printf("CARS:\n%s\n", dashes)
	for _, car := range cars {
		fmt.Printf("%v\n", car)
		for _, option := range car.Options {
			// Associate this car with each of it's options
			car_options[option] = append(car_options[option], car)
			fmt.Printf("car_options[option]: %s\n", option)
		}
		fmt.Println(dashes)
	}
	fmt.Println(dashes)

	// Print a list of cars with the "GPS" option
	for _, p := range car_options["GPS"] {
		fmt.Println(p.Make, "has GPS.")
	}

	fmt.Println("")
	fmt.Println(len(car_options["Alloy Wheels"]), "has Alloy Wheels.")
}

