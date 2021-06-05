package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/morsby/kbu"
	"github.com/morsby/kbu/db"
)

func main() {
	inputFlag := flag.String("input", "", "JSON file to parse")
	outputFlag := flag.String("output", "data.json", "path to a JSON file to write to")
	owFlag := flag.Bool("ow", false, "whether the output file should be overwritten if it exists")

	flag.Parse()

	if *inputFlag == "" {
		fmt.Println("needs an input file!")
		return
	}

	if *outputFlag == "" {
		fmt.Println("needs a an output path (defalt: data.json)!")
		return
	}

	if _, err := os.Stat(*outputFlag); !os.IsNotExist(err) && !*owFlag {
		fmt.Println("output file already exists, append -ow to overwrite")
		return
	}

	input, err := os.Open(*inputFlag)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	rounds, err := kbu.ParseData(input)
	if err != nil {
		panic(err)
	}
	dbConn, err := db.Connect()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(dbConn, &kbu.Round{}, &kbu.Selection{}, &kbu.Position{})
	db.InsertRounds(dbConn, &rounds)
}
