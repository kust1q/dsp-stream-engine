package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

// HeaderSize is the offset in bytes where the data chunk starts in our custom WAV file.
const HeaderSize = 56

// printSamples reads and prints the first count samples from the given WAV file.
func printSamples(filename string, count int) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Seek(HeaderSize, io.SeekStart)

	fmt.Printf("Samples from %s:\n", filename)
	for i := 0; i < count; i++ {
		var val float32
		err := binary.Read(f, binary.LittleEndian, &val)
		if err != nil {
			break
		}
		fmt.Printf("  Sample %d: %f\n", i, val)
	}
	fmt.Println()
}

// main analyzes the test output.
func main() {
	printSamples("test.wav", 10)
	printSamples("output.wav", 10)
}
