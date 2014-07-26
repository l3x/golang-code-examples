package main

import (
	"fmt"
	"time"
	"log"
	"github.com/l3x/jsoncfgo"
)

func main() {

	cfg := jsoncfgo.Load("/Users/lex/dev/go/data/jsoncfgo/simple-config.json")

	host := cfg.String("host")
	fmt.Printf("host: %v\n", host)
	bogusHost := cfg.String("bogusHost", "default_host_name")
	fmt.Printf("host: %v\n\n", bogusHost)

	port := cfg.Int("port")
	fmt.Printf("port: %v\n", port)
	bogusPort := cfg.Int("bogusPort", 9000)
	fmt.Printf("bogusPort: %v\n\n", bogusPort)

	bigNumber := cfg.Int64("bignNumber")
	fmt.Printf("bigNumber: %v\n", bigNumber)
	bogusBigNumber := cfg.Int64("bogusBigNumber", 9000000000000000000)
	fmt.Printf("bogusBigNumber: %v\n\n", bogusBigNumber)

	active := cfg.Bool("active")
	fmt.Printf("active: %v\n", active)
	bogusFalseActive := cfg.Bool("bogusFalseActive", false)
	fmt.Printf("bogusFalseActive: %v\n", bogusFalseActive)
	bogusTrueActive := cfg.Bool("bogusTrueActive", true)
	fmt.Printf("bogusTrueActive: %v\n\n", bogusTrueActive)

	appList := cfg.List("appList")
	fmt.Printf("appList: %v\n", appList)
	bogusAppList := cfg.List("bogusAppList", []string{"app1", "app2", "app3"})
	fmt.Printf("bogusAppList: %v\n\n", bogusAppList)

	numbers := cfg.IntList("numbers")
	fmt.Printf("numbers: %v\n", numbers)
	bogusSettings := cfg.IntList("bogusSettings", []int64{1, 2, 3})
	fmt.Printf("bogusAppList: %v\n\n", bogusSettings)

	if err := cfg.Validate(); err != nil {
		time.Sleep(100 * time.Millisecond)
		defer log.Fatalf("ERROR - Invalid config file...\n%v", err)
		return
	}
}

