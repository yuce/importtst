package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s RANDOM_SEED MAX_ROW_ID MAX_COLUMN_ID BIT_COUNT PATH.csv\n", os.Args[0])
		fmt.Println("Pass - as the PATH to output to stdout")
		os.Exit(1)
	}

	randomSeed := int64(parseInt(os.Args[1]))
	maxRowID := int64(parseInt(os.Args[2]))
	maxColID := int64(parseInt(os.Args[3]))
	bitCount := parseInt(os.Args[4])
	outputFile := os.Args[5]

	var f *os.File
	var err error

	if outputFile == "-" {
		f = os.Stdout
	} else {
		f, err = os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	rand.Seed(randomSeed)

	for i := 0; i < bitCount; i++ {
		rid := rand.Int63n(maxRowID)
		cid := rand.Int63n(maxColID)
		_, err := f.WriteString(fmt.Sprintf("%d,%d\n", rid, cid))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(strings.Replace(s, "_", "", -1))
	if err != nil {
		log.Fatalf("Invalid int: %s", s)
	}
	return i
}
