package main

import (
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type chip8 struct {
	Memory [0x1000]byte
}

func (c *chip8) readRom(romName string) {
	fmt.Printf("Reading: %s\n", romName)
	rom, err := os.Open(romName)
	check(err)

	b := make([]byte, 2)

	initialMemory := 0x200
	step := 0
	for {
		opCode, err := rom.Read(b)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		c.Memory[initialMemory+step] = b[:opCode][0]
		c.Memory[initialMemory+step+1] = b[:opCode][1]

		step += 2
	}

	rom.Close()
}

func main() {
	fmt.Print("Starting Chip8 emulator, by Chip.\n")

	var chipEmulator = chip8{}
	chipEmulator.readRom("./c8games/PONG")

	fmt.Print("Memory: ", chipEmulator.Memory)
}
