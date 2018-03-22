/*
Copyright 2017 Pilosa Corp.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.
*/

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/goterm"
	pilosa "github.com/pilosa/go-pilosa"
)

func main() {
	var err error

	if len(os.Args) != 4 {
		log.Fatal(fmt.Sprintf("Usage: %s PILOSA_ADDR PATH.csv|PATH.csv.gz BATCH_SIZE", os.Args[0]))
	}
	pilosaAddr := os.Args[1]
	path := os.Args[2]
	batchSize, err := strconv.Atoi(strings.Replace(os.Args[3], "_", "", -1))
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	client, err := pilosa.NewClient(pilosaAddr)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var bitIterator pilosa.BitIterator
	if strings.HasSuffix(path, ".gz") {
		if bitIterator, err = csvGZipIterator(f); err != nil {
			log.Fatal(err)
		}
	} else {
		bitIterator = csvIterator(f)
	}

	_, f1, err := ensureSchema(client)
	if err != nil {
		log.Fatal(err)
	}

	statusChan := make(chan pilosa.ImportStatusUpdate, 1000)

	fmt.Printf("Pilosa addr:         %s\n", pilosaAddr)
	fmt.Printf("Batch size:          %d\n", batchSize)
	fmt.Printf("===\n\n")

	go func() {
		err := client.ImportFrameWithStatus(f1, bitIterator, uint(batchSize), statusChan)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var status pilosa.ImportStatusUpdate
	totalImported := 0
	tic := time.Now()
	ok := true
	for ok {
		select {
		case status, ok = <-statusChan:
			if !ok {
				break
			}
			totalImported += status.ImportedCount
			goterm.MoveCursorUp(1)
			goterm.Print(fmt.Sprintf("Imported %d bits in %d s. Speed: %d bits/s.",
				totalImported, int(time.Since(tic).Seconds()), int(float64(status.ImportedCount)/status.Time.Seconds())))
			goterm.Flush()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Printf("Imported %d bits in %d milliseconds", totalImported, time.Since(tic).Nanoseconds()/1000000)
}

func csvGZipIterator(f *os.File) (*pilosa.CSVBitIterator, error) {
	reader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	return pilosa.NewCSVBitIterator(reader), nil
}

func csvIterator(f *os.File) *pilosa.CSVBitIterator {
	reader := bufio.NewReader(f)
	return pilosa.NewCSVBitIterator(reader)
}

func ensureSchema(client *pilosa.Client) (index *pilosa.Index, frame *pilosa.Frame, err error) {
	schema, err := client.Schema()
	if err != nil {
		return nil, nil, err
	}
	i1, err := schema.Index("i1")
	if err != nil {
		return nil, nil, err
	}
	f1, err := i1.Frame("f1")
	if err != nil {
		return nil, nil, err
	}
	err = client.SyncSchema(schema)
	if err != nil {
		return nil, nil, err
	}
	return i1, f1, nil
}