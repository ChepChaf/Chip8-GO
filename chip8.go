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
	PC     int16
	V      [16]byte
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

func (c *chip8) execNextOperation() int {
	switch c.Memory[c.PC] & 0xF0 {
	case 0x60:
		x := c.Memory[c.PC] & 0x0F
		c.V[x] = c.Memory[c.PC+1]

		c.PC += 2
	default:
		fmt.Printf("\nOperation not implemented: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

		return -1
	}

	return 1
}

func startChip8() *chip8 {
	c := chip8{}
	c.PC = 0x200

	return &c
}

func main() {
	fmt.Print("Starting Chip8 emulator, by Chip.\n")

	var chipEmulator = startChip8()
	chipEmulator.readRom("./c8games/PONG")

	for {
		// fmt.Print("Memory: ", chipEmulator.Memory)
		fmt.Printf("\nPC: %X", chipEmulator.PC)
		fmt.Print("\nV: ", chipEmulator.V)

		if chipEmulator.execNextOperation() < 0 {
			break
		}
	}
}
