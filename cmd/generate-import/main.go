package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Generator interface {
	nextBit() (int64, int64, bool)
}

type RandomGenerator struct {
	rnd      *rand.Rand
	maxRowID int64
	maxColID int64
	remBits  int
}

func NewRandomGenerator(args []string) (*RandomGenerator, error) {
	if len(args) != 4 {
		return nil, errors.New("Required params: RANDOM_SEED MAX_ROW_ID MAX_COL_ID BIT_COUNT")
	}
	seed := int64(parseInt(args[0]))
	maxRowID := int64(parseInt(args[1]))
	maxColID := int64(parseInt(args[2]))
	bitCount := parseInt(args[3])
	return &RandomGenerator{
		rnd:      rand.New(rand.NewSource(seed)),
		maxRowID: maxRowID,
		maxColID: maxColID,
		remBits:  bitCount,
	}, nil
}

func (g *RandomGenerator) nextBit() (int64, int64, bool) {
	if g.remBits <= 0 {
		return 0, 0, false
	}
	g.remBits -= 1
	return g.rnd.Int63n(g.maxRowID), g.rnd.Int63n(g.maxColID), true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s GEN_STRATEGY [STRATEGY_PARAMS]\n", os.Args[0])
		os.Exit(1)
	}

	var gen Generator
	var err error
	strategy := os.Args[1]

	switch strategy {
	case "random":
		gen, err = NewRandomGenerator(os.Args[2:])
	default:
		err = fmt.Errorf("Unknown strategy: %s", strategy)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	outputFile := "-"

	var f *os.File

	if outputFile == "-" {
		f = os.Stdout
	} else {
		f, err = os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	// w := bufio.NewWriter(f)
	w := f

	for {
		rid, cid, ok := gen.nextBit()
		if !ok {
			break
		}
		_, err := w.WriteString(fmt.Sprintf("%d,%d\n", rid, cid))
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
