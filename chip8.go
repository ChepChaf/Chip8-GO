package main

import (
	"fmt"
	"io"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type chip8 struct {
	Memory  [0x1000]byte
	PC      uint16
	V       [16]byte
	I       uint16
	Display [64][32]byte
	Running bool
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

func (c *chip8) execNextOperation() {
	switch c.Memory[c.PC] & 0xF0 {
	case 0x60:
		x := c.Memory[c.PC] & 0x0F
		c.V[x] = c.Memory[c.PC+1]

		c.PC += 2
	case 0xA0:
		c.I = uint16(c.Memory[c.PC]&0x0F)<<8 | uint16(c.Memory[c.PC+1])

		c.PC += 2
	default:
		fmt.Printf("\nOperation not implemented: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

		c.Running = false
	}

	if c.Running {
		c.updateScreen()
	}
}

func (c *chip8) updateScreen() {
	// fmt.Print("Memory: ", c.Memory)
	fmt.Printf("\nPC: %X", c.PC)
	fmt.Print("\nV: ", c.V)
	fmt.Printf("\nI: %X", c.I)
	fmt.Printf("\nOPCODE: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			c.Running = false
			break
		}
	}
}

func (c *chip8) initWindow() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	window.UpdateSurface()

	c.Running = true
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

	chipEmulator.initWindow()

	for chipEmulator.Running {
		chipEmulator.execNextOperation()
	}
}
