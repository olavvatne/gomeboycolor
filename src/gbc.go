package main

import (
	"apu"
	"bufio"
	"cartridge"
	"cpu"
	"flag"
	"fmt"
	"gpu"
	"inputoutput"
	"log"
	"mmu"
	"os"
	"strings"
	"types"
	"utils"
)

var printState *bool = flag.Bool("dump", false, "Print state of machine after each cycle (WARNING - WILL RUN SLOW)")
var skipBoot *bool = flag.Bool("noboot", false, "Skip boot sequence")
var debug *bool = flag.Bool("d", false, "Enable debugger")
var breakOn *string = flag.String("b", "0x0000", "Break into debugger when PC equals a given value between 0x0000 and 0xFFFF")

const FRAME_CYCLES = 70224
const TITLE string = "gomeboycolor"
const VERSION float32 = 0.1

//A type
type GameboyColor struct {
	gpu          *gpu.GPU
	cpu          *cpu.Z80
	mmu          *mmu.GbcMMU
	io           *inputoutput.IO
	apu          *apu.APU
	debugOptions *DebugOptions
	cpuClockAcc  int
	frameCount   int
	stepCount    int
	inBootMode   bool
}

func NewGBC() *GameboyColor {
	gbc := new(GameboyColor)

	gbc.mmu = mmu.NewGbcMMU()
	gbc.cpu = cpu.NewCPU()
	gbc.cpu.LinkMMU(gbc.mmu)

	gbc.io = inputoutput.NewIO(inputoutput.DefaultControlScheme)
	gbc.gpu = gpu.NewGPU()
	gbc.apu = apu.NewAPU()

	gbc.gpu.LinkScreen(gbc.io.Display)
	//mmu will process interrupt requests from GPU (i.e. it will set appropriate flags)
	gbc.gpu.LinkIRQHandler(gbc.mmu)

	gbc.mmu.ConnectPeripheral(gbc.apu, 0xFF10, 0xFF30)
	gbc.mmu.ConnectPeripheral(gbc.gpu, 0x8000, 0x9FFF)
	gbc.mmu.ConnectPeripheral(gbc.gpu, 0xFE00, 0xFE9F)
	gbc.mmu.ConnectPeripheral(gbc.gpu, 0xFF40, 0xFF45)
	gbc.mmu.ConnectPeripheral(gbc.gpu, 0xFF47, 0xFF4B)
	gbc.mmu.ConnectPeripheral(gbc.gpu, 0xFF51, 0xFF70)
	gbc.mmu.ConnectPeripheral(gbc.io.KeyHandler, 0xFF00, 0xFF00)

	return gbc
}

func (gbc *GameboyColor) DoFrame() {
	for gbc.cpuClockAcc < FRAME_CYCLES {
		if gbc.debugOptions.debuggerOn && gbc.cpu.PC == gbc.debugOptions.breakWhen {
			gbc.Pause()
		}

		if *printState {
			fmt.Println(gbc.cpu)
		}
		gbc.Step()
	}
}

func (gbc *GameboyColor) Step() {
	cycles := gbc.cpu.Step()
	gbc.gpu.Step(cycles)
	gbc.cpuClockAcc += cycles
	gbc.stepCount++
	//value in FF50 means gameboy has finished booting
	if gbc.inBootMode {
		if gbc.mmu.ReadByte(0xFF50) != 0x00 {
			gbc.mmu.SetInBootMode(false)
			gbc.inBootMode = false
			gbc.cpu.PC = 0x0100
			log.Println("Finished GB boot program, launching game...")
		}
	}
}

func (gbc *GameboyColor) Run() {
	log.Println("Starting emulator")
	for {
		gbc.DoFrame()
		gbc.frameCount++
		gbc.cpuClockAcc = 0
	}
}

func main() {
	flag.Parse()
	log.Println(TITLE, VERSION)
	log.Println(strings.Repeat("*", 80))

	if flag.NArg() != 1 {
		log.Fatalf("Please specify the location of a ROM to boot")
		return
	}

	rom, err := RetrieveROM(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error retrieving ROM file: %v", err)
	}

	cart, err := cartridge.NewCartridge(rom)
	if err != nil {
		log.Println(err)
		return
	}

	var gbc *GameboyColor = NewGBC()
	b, er := gbc.mmu.LoadBIOS(BOOTROM)
	if !b {
		log.Println("Error loading bootrom:", er)
		return
	}

	gbc.mmu.LoadCartridge(cart)
	gbc.debugOptions = new(DebugOptions)
	gbc.debugOptions.Init()

	if *debug {
		log.Println("Emulator will start in debug mode")
		gbc.debugOptions.debuggerOn = true

		//set breakpoint if defined
		if b, err := utils.StringToWord(*breakOn); err != nil {
			log.Fatalln("Cannot parse breakpoint:", *breakOn, "\n\t", err)
		} else {
			gbc.debugOptions.breakWhen = types.Word(b)
			log.Println("Emulator will break into debugger when PC = ", gbc.debugOptions.breakWhen)
		}
	}

	screenInitErr := gbc.io.Init(TITLE, onClose)
	if screenInitErr != nil {
		log.Fatalf("%v", screenInitErr)
	}

	gbc.setupBoot()

	log.Println("Completed setup")
	log.Println(strings.Repeat("*", 80))

	gbc.Run()
}

func (gbc *GameboyColor) setupBoot() {
	if *skipBoot {
		log.Println("Boot sequence disabled")
		gbc.setupWithoutBoot()
	} else {
		log.Println("Boot sequence enabled")
		gbc.setupWithBoot()
	}
}

func (gbc *GameboyColor) setupWithBoot() {
	gbc.inBootMode = true
	gbc.mmu.WriteByte(0xFF50, 0x00)
}

func (gbc *GameboyColor) setupWithoutBoot() {
	gbc.inBootMode = false
	gbc.mmu.SetInBootMode(false)
	gbc.cpu.PC = 0x100
	gbc.cpu.R.A = 0x11
	gbc.cpu.R.F = 0xB0
	gbc.cpu.R.B = 0x00
	gbc.cpu.R.C = 0x13
	gbc.cpu.R.D = 0x00
	gbc.cpu.R.E = 0xD8
	gbc.cpu.R.H = 0x01
	gbc.cpu.R.L = 0x4D
	gbc.cpu.SP = 0xFFFE
	gbc.mmu.WriteByte(0xFF05, 0x00)
	gbc.mmu.WriteByte(0xFF06, 0x00)
	gbc.mmu.WriteByte(0xFF07, 0x00)
	gbc.mmu.WriteByte(0xFF10, 0x80)
	gbc.mmu.WriteByte(0xFF11, 0xBF)
	gbc.mmu.WriteByte(0xFF12, 0xF3)
	gbc.mmu.WriteByte(0xFF14, 0xBF)
	gbc.mmu.WriteByte(0xFF16, 0x3F)
	gbc.mmu.WriteByte(0xFF17, 0x00)
	gbc.mmu.WriteByte(0xFF19, 0xBF)
	gbc.mmu.WriteByte(0xFF1A, 0x7F)
	gbc.mmu.WriteByte(0xFF1B, 0xFF)
	gbc.mmu.WriteByte(0xFF1C, 0x9F)
	gbc.mmu.WriteByte(0xFF1E, 0xBF)
	gbc.mmu.WriteByte(0xFF20, 0xFF)
	gbc.mmu.WriteByte(0xFF21, 0x00)
	gbc.mmu.WriteByte(0xFF22, 0x00)
	gbc.mmu.WriteByte(0xFF23, 0xBF)
	gbc.mmu.WriteByte(0xFF24, 0x77)
	gbc.mmu.WriteByte(0xFF25, 0xF3)
	gbc.mmu.WriteByte(0xFF26, 0xF1)
	gbc.mmu.WriteByte(0xFF40, 0x91)
	gbc.mmu.WriteByte(0xFF42, 0x00)
	gbc.mmu.WriteByte(0xFF43, 0x00)
	gbc.mmu.WriteByte(0xFF45, 0x00)
	gbc.mmu.WriteByte(0xFF47, 0xFC)
	gbc.mmu.WriteByte(0xFF48, 0xFF)
	gbc.mmu.WriteByte(0xFF49, 0xFF)
	gbc.mmu.WriteByte(0xFF4A, 0x00)
	gbc.mmu.WriteByte(0xFF4B, 0x00)
	gbc.mmu.WriteByte(0xFF50, 0x00)
	gbc.mmu.WriteByte(0xFFFF, 0x00)
}

func onClose() {
	log.Println("Closing")
	os.Exit(0)
}

func (gbc *GameboyColor) Pause() {

	log.Println("DEBUGGER: Breaking because PC ==", gbc.debugOptions.breakWhen)
	b := bufio.NewWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)

	fmt.Fprintln(b, "Debug mode, type ? for help")
	for gbc.debugOptions.debuggerOn {
		var instruction string
		b.Flush()
		fmt.Fprint(b, "> ")
		b.Flush()
		instruction, _ = r.ReadString('\n')
		b.Flush()
		var instructions []string = strings.Split(strings.Replace(instruction, "\n", "", -1), " ")
		b.Flush()

		command := instructions[0]

		if command == "c" {
			break
		}

		//dispatch
		if v, ok := gbc.debugOptions.debugFuncMap[command]; ok {
			v(gbc, instructions[1:]...)
		} else {
			fmt.Fprintln(b, "Unknown command:", command)
			fmt.Fprintln(b, "Debug mode, type ? for help")
		}
	}
}

func (gbc *GameboyColor) Reset() {
	log.Println("Resetting system")
	gbc.cpu.Reset()
	gbc.gpu.Reset()
	gbc.mmu.Reset()
	gbc.apu.Reset()
	gbc.io.KeyHandler.Reset()
	gbc.io.Display.DrawFrame(&([144][160]types.RGB{}))
	gbc.setupBoot()
}

func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

var BOOTROM []byte = []byte{
	0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
	0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
	0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
	0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
	0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
	0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
	0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
	0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
	0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
	0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x4C,
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50}