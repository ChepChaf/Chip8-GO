package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var renderer *sdl.Renderer
var surface *sdl.Surface

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type chip8 struct {
	Memory  [0x1000]uint8
	PC      uint16
	V       [16]uint8
	I       uint16
	SP      uint8
	stack   [16]uint16
	Display [0x440]uint8
	Running bool
}

func (c *chip8) readRom(romName string) {
	fmt.Printf("Reading: %s\n", romName)
	rom, err := os.Open(romName)
	check(err)

	b := make([]uint8, 1)

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
		step++
	}

	rom.Close()
}

func (c *chip8) execNextOperation() {
	currentCommand := c.Memory[c.PC] & 0xF0

	switch currentCommand {
	case 0x20:
		c.SP++
		c.stack[c.SP] = c.PC
		c.PC = uint16(c.Memory[c.PC]&0x0F)<<8 | uint16(c.Memory[c.PC+1])
	case 0x60:
		x := c.Memory[c.PC] & 0x0F
		c.V[x] = c.Memory[c.PC+1]

		c.PC += 2
	case 0xA0:
		c.I = uint16(c.Memory[c.PC]&0x0F)<<8 | uint16(c.Memory[c.PC+1])

		c.PC += 2
	case 0xD0:
		vX := c.V[c.Memory[c.PC]&0x0F]
		vY := c.V[c.Memory[c.PC+1]&0xF0>>4]
		n := c.Memory[c.PC+1] & 0x0F

		for i := 0; i < int(n); i++ {
			idx := uint8(i)
			c.Display[vX+vY+(8*idx)] = c.Display[vX+vY+(8*idx)] ^ c.Memory[c.I+uint16(idx)]
		}

		c.PC += 2
	case 0xF0:
		currentCommand = c.Memory[(c.PC+1)] & 0xFF

		switch currentCommand {
		case 0x33:
			index := (c.PC & 0x0F) >> 4
			c.Memory[c.I] = c.V[index] / 100
			c.Memory[c.I+1] = (c.V[index] % 100) / 10
			c.Memory[c.I+2] = c.V[index] % 10

			fmt.Printf("vX: %X", c.V[(c.PC)&0x0F>>4])
			fmt.Printf("c.I: %X", c.Memory[c.I])
			fmt.Printf("c.I: %X", c.Memory[c.I+1])
			fmt.Printf("c.I: %X", c.Memory[c.I+2])

			c.PC += 2
		case 0x65:
			x := c.Memory[c.PC] & 0x0F

			for i := 0; i < int(x+1); i++ {
				c.V[i] = c.Memory[c.I+uint16(i)]
			}

			c.PC += 2
		default:
			fmt.Printf("\nOperation not implemented: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

			c.Running = false
		}
	default:
		fmt.Printf("\nOperation not implemented: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

		c.Running = false
	}

	if c.Running {
		c.updateScreen()
		time.Sleep(1 * time.Second)
	}
}

func (c *chip8) updateScreen() {
	// fmt.Print("Memory: ", c.Memory)
	sdl.Log("\nPC: %X", c.PC)
	sdl.Log("\nV: ", c.V)
	sdl.Log("\nI: %X", c.I)
	sdl.Log("\nOPCODE: %X%X", c.Memory[c.PC], c.Memory[c.PC+1])

	sdl.Log("\nUpdating Window")

	//  Create Surface
	newSurface, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&c.Display), 64, 32, 1, 8, 0, 0, 0, 0xFF)
	check(err)

	newSurface = gfx.RotoZoomSurface(newSurface, 0, 10, 0)

	newSurface.Blit(nil, surface, nil)
	window.UpdateSurface()

	fmt.Print("\nWindow updated")
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
	fmt.Print("\nIniting Window")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	// defer sdl.Quit()

	fmt.Print("\nCreating Window.")
	var err error
	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 320, sdl.WINDOW_SHOWN)
	check(err)
	// defer window.Destroy()

	fmt.Print("\nCreating Renderer.")
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	check(err)

	fmt.Print("\nCreating Surface")
	surface, err = window.GetSurface()
	check(err)

	surface.FillRect(nil, 0)
	window.UpdateSurface()

	c.Running = true

	fmt.Print("\nWindow Inited")
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
