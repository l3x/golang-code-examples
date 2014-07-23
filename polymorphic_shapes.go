/*
filename:  polymorphic_shapes.go
author:    Lex Sheehan
copyright: Lex Sheehan LLC
license:   GPL
status:    published
comments:  http://l3x.github.io/golang-code-examples/2014/07/14/custom-error-handling.html
*/
package main	// Executable commands must always use package main.

import (
	"fmt"		// fmt.Println formats output to console
	"math"  	// provides math.Sqrt function
)

// ----------------------
//    Shape interface
// ----------------------
// Shape interface defines a method set (consisting of the area method)
type Shape interface {
	area() float64  		// any type that implements an area method is considered a Shape
}
// Calculate total area of all shapes via polymorphism (all shapes implement the area method)
func totalArea(shapes ...Shape) float64 {	// Use interface type as as function argument
	var area float64						// "..." makes shapes "variadic" (can send one or more)
	for _, s := range shapes {
		area += s.area()  	// the current Shape implements/receives the area method
	}						// go passes the pointer to the shape to the area method
	return area
}

// ----------------------
//    Drawer interface
// ----------------------
type Drawer interface {
	draw()					// does not return a type
}
func drawShape(d Drawer) {	// associate this method with the Drawer interface
	d.draw()
}

// ----------------------
//      Circle Type
// ----------------------
type Circle struct {		// Since "Circle" is capitalized, it is visible outside this package
	x, y, r float64			// a Circle struct is a collection of fields: x, y, r
}
// Circle implements Shape interface b/c it has an area method
// area is a method, which is special type of function that is associated with the Circle struct
// The Circle struct becomes the "receiver" of this method, so we can use the "." operator
func (c *Circle) area() float64 {  	// dereference Circle type (data pointed to by c)
	return math.Pi * c.r * c.r		// Pi is a constant in the math package
}
func (c Circle) draw() {
	fmt.Println("Circle drawing with radius: ", c.r)	// encapsulated draw implementation for Circle type
}
// ----------------------
//     Rectangle Type
// ----------------------
type Rectangle struct {		// a struct contains named fields of data
	x1, y1, x2, y2 float64	// define multiple fields with same data type on one line
}
func distance(x1, y1, x2, y2 float64) float64 {  		// lowercase functin name visible only in this package
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a * a + b * b)
}
// Rectangle implements Shape interface b/c it has an area method
func (r *Rectangle) area() float64 {
	l := distance(r.x1, r.y1, r.x1, r.y2)  // define and assign local variable "l"
	w := distance(r.x1, r.y1, r.x2, r.y1)  // l and w only available within scope of area function
	return l * w
}
func (r Rectangle) draw() {
	fmt.Printf("Rectangle drawing with point1: (%f, %f) and point2: (%f, %f)\n", r.x1, r.y1, r.x2, r.y2)
}
// ----------------------
//    MultiShape Type
// ----------------------
type MultiShape struct {
	shapes []Shape  // shapes field is a slice of interfaces
}
//
func (m *MultiShape) area() float64 {
	var area float64
	for _, shape := range m.shapes { 	// iterate through shapes ("_" indicates that index is not used)
		area += shape.area()			// execute polymorphic area method for this shape
	}
	return area
}

func main() {
	c := Circle{0, 0, 5}					        				// initialize new instance of Circle type by field order "struct literal"
																	// The new function allocates memory for all  fields, sets each to their zero value and returns a pointer
	c2 := new(Circle)												// c2 is a pointer to the instantiated Circle type
	c2.x = 0; c2.y = 0; c2.r = 10									// initialize data with multiple statements on one line
	fmt.Println("Circle Area:", totalArea(&c))						// pass address of circle (c)
	fmt.Println("Circle2 Area:", totalArea(c2))						// c2 was defined using built-in new function
	r := Rectangle{x1: 0, y1: 0, x2: 5, y2: 5}						// "struct literal" rectangle (r) initialized by field name
	fmt.Println("Rectangle Area:", totalArea(&r))   				// pass address of rectangle (r)
	fmt.Println("Rectangle + Circle Area:", totalArea(&c, c2, &r))  // can pass multiple shapes
	m := MultiShape{[]Shape{&r, &c, c2}}							// pass slice of shapes
	fmt.Println("Multishape Area:", totalArea(&m))					// calculate total area of all shapes
	fmt.Println("Area Totals:", totalArea(&c, c2, &r))  			// c2 is a pointer to a circle, &c and &r are addresses of shapes
	fmt.Println("2 X Area Totals:", totalArea(&c, c2, &r, &m))		// twice the size of all areas
	drawShape(c)													// execute polymorphic method call
	drawShape(c2)
	drawShape(r)
}
