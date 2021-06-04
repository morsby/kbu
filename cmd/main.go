package main

import (
	"encoding/json"
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
	selections, err := kbu.ParseRawJSON(input)
	flattened := kbu.FlattenSelections(selections)
	if err != nil {
		panic(err)
	}

	output, err := os.Create(*outputFlag)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	enc := json.NewEncoder(output)
	enc.Encode(flattened)

	fmt.Printf("Found %d selections in file and wrote to %s\n", len(selections), *outputFlag)

	database := db.Connect()
	db.CreateTables(database)

	seeds := db.Seeds{
		Regions:      []kbu.Region{kbu.RegionH, kbu.RegionMidt, kbu.RegionNord, kbu.RegionSj, kbu.RegionSyd},
		Universities: []kbu.University{kbu.UniversityAAU, kbu.UniversityAU, kbu.UniversityKU, kbu.UniversitySDU, kbu.UniversityNA},
	}
	err = db.Seed(database, seeds)
	if err != nil {
		panic(err)
	}
}
