package main

import (
	"encoding/binary"
	"math"
	"os"
)

// main generates a test IEEE Float WAV file with a 440Hz sine wave.
func main() {
	f, _ := os.Create("test.wav")
	defer f.Close()

	sampleRate := 44100
	numSamples := sampleRate * 2

	f.WriteString("RIFF")
	binary.Write(f, binary.LittleEndian, uint32(50+numSamples*4))
	f.WriteString("WAVE")

	f.WriteString("fmt ")
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(3))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint32(sampleRate))
	binary.Write(f, binary.LittleEndian, uint32(sampleRate*4))
	binary.Write(f, binary.LittleEndian, uint16(4))
	binary.Write(f, binary.LittleEndian, uint16(32))

	f.WriteString("fact")
	binary.Write(f, binary.LittleEndian, uint32(4))
	binary.Write(f, binary.LittleEndian, uint32(numSamples))

	f.WriteString("data")
	binary.Write(f, binary.LittleEndian, uint32(numSamples*4))

	for i := 0; i < numSamples; i++ {
		val := float32(math.Sin(2.0 * math.Pi * 440.0 * float64(i) / float64(sampleRate)))
		binary.Write(f, binary.LittleEndian, val)
	}
}
