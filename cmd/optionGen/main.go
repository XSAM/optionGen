package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(optionGen + ": ")
	
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get working directory: %v", err)
	}
	
	parseDir(wd)
}
