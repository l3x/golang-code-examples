/*
filename:  custom-error-handling.go
author:    Lex Sheehan
copyright: Lex Sheehan LLC
license:   GPL
status:    published
comments:  http://l3x.github.io/golang-code-examples/2014/07/14/custom-error-handling.html
*/
package main

import (
	"fmt"
	"time"
)

type Err struct {
	errNo int
	when time.Time
	msg string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v [%d] %s", e.when, e.errNo, e.msg)
}
func (err Err) errorNumber() int {
	return err.errNo
}

type ErrWidget_A struct {
	Err  					// Err is an embedded struct - ErrWidget_A inherits it's data and behavior
}
// a behavior only available for the ErrWidget_A
func (e ErrWidget_A) Error() string {
	fmt.Println("do special ErrWidget_A thing...")
	return fmt.Sprintf("%s [%d] %s", e.when, e.errNo, e.msg)
}
// a behavior only available for the ErrWidget_A
func (e ErrWidget_A) optionalErrHandlingOperation() {
	fmt.Println("Email the admins...\n")
}

type ErrWidget_B struct {
	Err						// Err is an embedded struct - ErrWidget_B inherits it's data and behavior
}
// a behavior only available for the Widget_B
func (e ErrWidget_B) Error() string {
	fmt.Println("do special Widget_B thing...")
	return fmt.Sprintf("%s [%d] %s", e.when, e.errNo, e.msg)
}
// a behavior only available for the Widget_B
func (e ErrWidget_B) optionalErrHandlingOperation() {
	fmt.Println("SMS operations manager...\n")
}

func run() error {
	return &Err{
		8001,
		time.Now(),
		"generic error occurred\n",
	}
}

func run2() *ErrWidget_B {
	errB := new(ErrWidget_B)
	errB.errNo = 6001
	errB.when = time.Now()
	errB.msg = "Widget_B error occurred"
	return errB
}

func RunWidget(modelNo int) (string, error) {
	// Run valid widgets
	switch modelNo {
	case 1:
		return fmt.Sprintf("run widget model %d", modelNo), nil
	case 2:
		return fmt.Sprintf("run widget model %d", modelNo), nil
	default:
		// Error condition - unknown widget model number
		errA := new(ErrWidget_A)
		errA.errNo = 5002
		errA.when = time.Now()
		errA.msg = "Widget_A error occurred"
		return fmt.Sprintf("unable to run unknown model %d", modelNo), errA
	}
}

// Split multiple (variadic) return values into a slice of values
// in this case, where [0] = value and [1] = the error message
func split(args ...interface{}) []interface{} {
	return args
}

func main() {

	// Execute RunWidget function and handle error if necessary
	msg := ""
	// RunWidget(1) succeeds
	x := split(RunWidget(1))
	msg = "\n\n"; if x[1] != nil {msg = fmt.Sprintf(", err(%v)\n\n", x[1])}
	fmt.Printf("RunWidget(1) => result(%s)" + msg, x[0])

	// RunWidget(2) succeeds
	x = split(RunWidget(2))
	msg = "\n\n"; if x[1] != nil {msg = fmt.Sprintf(", err(%v)\n\n", x[1])}
	fmt.Printf("RunWidget(2) => result(%s)" + msg, x[0])

	// RunWidget(666) fails -
	x = split(RunWidget(666))
	msg = "\n\n"; if x[1] != nil {msg = fmt.Sprintf(", err(%v)\n\n", x[1])}
	fmt.Printf("RunWidget(666) => result(%s)" + msg, x[0])


	// Throw generic custom error type and handle it
	if err := run(); err != nil { fmt.Println(err) }

	// Throw ErrWidget_B error and handle it by printing and running optional custom behavior
	widget_B_error := run2(); if widget_B_error.errNo != 0 {
		fmt.Println(widget_B_error)
	}
	fmt.Println("")


	timeNow := time.Now()
	// Create and print ErrWidget_A, then call custom behavior
	a := ErrWidget_A {Err{5001, timeNow, "test"}}
	fmt.Println(a)  // fmt will execute Error() method that can have special behavior
	fmt.Println("A ErrWidget_A has this error number: ", a.errorNumber())
	a.optionalErrHandlingOperation()  // Widget_A emails admins

	// Create ErrWidget_B, then call custom behavior
	b := ErrWidget_B {Err{6001, timeNow, "test"}}
	fmt.Println("A ErrWidget_B has this error number: ", b.errorNumber())
	b.optionalErrHandlingOperation()  // Widget_B sends SMS message to managers
	// Since b was not printed by fmt, the special ErrWidget_B behavior is not triggered
}
