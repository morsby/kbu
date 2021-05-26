package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/morsby/kbu"
)

func main() {
	fileFlag := flag.String("file", "", "JSON file to parse")
	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("needs a file!")
		return
	}

	f, err := os.Open(*fileFlag)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	selections, err := kbu.ParseRawJSON(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d selections in file\n", len(selections))

}
