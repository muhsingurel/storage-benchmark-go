package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Usage of storage-benchmark-go:
//   -filePath string
//         filePath for the read test. A file larger than 100GB is recommended. (default "256GB.dummy")
//   -offsetOption int
//         An optional offset number as number of chunks between 0-10000. Otherwise it is random (default -1)
//   -parallelOption
//         Optianal parallel read flag.

func main() {
	// Initialize random, give a seed with current time
	rand.Seed(time.Now().UnixNano())
	var inputFilePath = flag.String("filePath", "256GB.dummy", "filePath for the read test. A file larger than 100GB is recommended.")
	var parallelRead = flag.Bool("parallelOption", false, "Optianal parallel read flag.")
	var offsetOption = flag.Int("offsetOption", -1, "An optional offset number as number of chunks between 0-10000. Otherwise it is random")
	flag.Parse()
	fmt.Println(*inputFilePath)
	f, err := os.Open(*inputFilePath)
	check(err)
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The file is %d bytes long", fi.Size())
	var checksum uint64 = 0
	var singleChunkBufferLenght = 10000
	var maxOffsetAsChunks int64 = 10000
	var numberOfChunksToRead = (fi.Size() / (int64(singleChunkBufferLenght) * maxOffsetAsChunks))
	OffsetNumber := 0

	// If user inputs a number use it, otherwise get a random number for offset
	if *offsetOption < 0 || *offsetOption >= 10000 {
		OffsetNumber = rand.Intn(10000)
	} else {
		OffsetNumber = *offsetOption
	}

	// If user does not provide parallel argument, execute the code in sequential.
	if !*parallelRead {
		// Sequential read test
		start := time.Now()
		for chunkIndex := int64(0); chunkIndex < numberOfChunksToRead; chunkIndex++ {
			singleChunkBuffer := make([]byte, singleChunkBufferLenght)
			var position = (int64(maxOffsetAsChunks)*int64(chunkIndex)*int64(len(singleChunkBuffer)) + (int64(len(singleChunkBuffer)) * int64(OffsetNumber))) % fi.Size()
			_, err1 := f.Seek(position, 0)
			check(err1)
			_, err2 := f.Read(singleChunkBuffer)
			check(err2)

			for _, e := range singleChunkBuffer {
				checksum += uint64(e)
			}

			fmt.Printf("Offset value: %d  Chunk Number: %d  Position: %d  Filelength: %d \n", OffsetNumber, chunkIndex, position, fi.Size())
		}
		duration := time.Since(start)
		fmt.Printf("Checksum: %d Sequential read operation took %d ms \n", checksum, duration.Milliseconds())
	} else {
		// Parallel read test
		start := time.Now()

		// We need an array so that each tread can access to a different index without blocking other threads.
		// If we use the global checksum variable, it may cause race condition.
		// The best way is to store each trace result in a different checksumBuffer index
		var checksumBuffer = make([]uint64, numberOfChunksToRead)

		var wg sync.WaitGroup
		wg.Add(int(numberOfChunksToRead))
		// Start reading the chunks and calculate chunk sum as a checksum. Later we can verify if the file is actually read with other tools.
		// Each itteration creates a different file handler. When the itteration is complete the resource is disposed. Each itteration calculates the chunk position.
		for chunkIndex := int64(0); chunkIndex < numberOfChunksToRead; chunkIndex++ {
			go func(i int64) {
				defer wg.Done()
				singleChunkBuffer := make([]byte, singleChunkBufferLenght)
				var position = (int64(maxOffsetAsChunks)*i*int64(len(singleChunkBuffer)) + (int64(len(singleChunkBuffer)) * int64(OffsetNumber))) % fi.Size()
				parallelFile, err := os.Open(*inputFilePath)
				check(err)
				_, err2 := parallelFile.ReadAt(singleChunkBuffer, position)
				check(err2)

				for _, e := range singleChunkBuffer {
					checksumBuffer[i] += uint64(e)
				}

				fmt.Printf("Offset value: %d  Chunk Number: %d  Position: %d  Filelength: %d \n", OffsetNumber, i, position, fi.Size())
			}(chunkIndex)
		}

		//wait all functions in separate threads to complate.
		wg.Wait()
		for _, e := range checksumBuffer {
			checksum += e
			//fmt.Printf("index: %d valur %d \n", indexx, e)
		}
		duration := time.Since(start)
		fmt.Printf("Checksum: %d Sequential read operation took %d ms \n", checksum, duration.Milliseconds())
	}
}
