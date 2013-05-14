package cpu

import (
	"bufio"
	"cartridge"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"types"
	"utils"
)

var c *GbcCPU

func before() {
	c = NewCPU()
	c.LinkMMU(NewMockMMU())
}

func AssertTimings(c *GbcCPU, t *testing.T, instr byte, expectedTiming int, isCB bool) {
	tick := c.Step()

	if isCB {
		log.Println("0xCB "+utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "testing that instruction runs for", expectedTiming, "cycles")
	} else {
		log.Println(utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "testing that instruction runs for", expectedTiming, "cycles")
	}

	if tick != expectedTiming {
		if isCB {
			t.Log("-----> For instruction 0xCB", utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "Expected", expectedTiming, "but got", tick)
		} else {
			t.Log("-----> For instruction", utils.ByteToString(instr)+" ("+c.CurrentInstruction.Description+")", "Expected", expectedTiming, "but got", tick)

		}
		t.FailNow()
	}
}

func RunCBInstrAndAssertTimings(instr byte, flags []int, t *testing.T) {
	before()
	//	expectedTiming += 4 //+4 because CB adds another cycle
	c.PC = 0x0000
	c.WriteByte(c.PC, 0xCB)
	c.WriteByte(c.PC+1, instr)
	var expectedTiming int = getTimingFromInstructionServer(instr)
	AssertTimings(c, t, instr, expectedTiming, true)
}

func RunInstrAndAssertTimings(instr byte, flags []int, t *testing.T) {
	before()
	if flags != nil {
		for _, f := range flags {
			if f == Z {
				c.SetFlag(Z)
			}
			if f == C {
				c.SetFlag(C)
			}
			if f == N {
				c.SetFlag(N)
			}
			if f == H {
				c.SetFlag(H)
			}
		}
	} else {
		c.ResetFlag(Z)
		c.ResetFlag(H)
		c.ResetFlag(N)
		c.ResetFlag(C)
		c.R.F = 0x00
	}
	c.PC = 0x0000
	c.WriteByte(c.PC, instr)
	var expectedTiming int = getTimingFromInstructionServer(instr)
	AssertTimings(c, t, instr, expectedTiming, false)
}

func getTimingFromInstructionServer(instr byte) int {
	conn, err := net.Dial("tcp", "localhost:8012")
	if err != nil {
		fmt.Println("error!")
	}
	var writer *bufio.Writer = bufio.NewWriter(conn)
	var reader *bufio.Reader = bufio.NewReader(conn)
	instrString := utils.ByteToString(instr)
	writer.WriteString(instrString + "\n")
	writer.Flush()
	result, _ := reader.ReadString('\n')
	conn.Close()
	i, _ := strconv.Atoi(result)
	return int(i)
}

func TestTest(t *testing.T) {
	/*
		RunInstrAndAssertTimings(0x00, nil, t)
		RunInstrAndAssertTimings(0x01,  nil, t)
		RunInstrAndAssertTimings(0x02,  nil, t)
		RunInstrAndAssertTimings(0x03,  nil, t)
		RunInstrAndAssertTimings(0x04,  nil, t)
		RunInstrAndAssertTimings(0x05,  nil, t)
		RunInstrAndAssertTimings(0x06,  nil, t)
		RunInstrAndAssertTimings(0x07,  nil, t)
		RunInstrAndAssertTimings(0x08,  nil, t)
		RunInstrAndAssertTimings(0x09,  nil, t)
		RunInstrAndAssertTimings(0x0A,  nil, t)
		RunInstrAndAssertTimings(0x0B,  nil, t)
		RunInstrAndAssertTimings(0x0C,  nil, t)
		RunInstrAndAssertTimings(0x0D,  nil, t)
		RunInstrAndAssertTimings(0x0E,  nil, t)
		RunInstrAndAssertTimings(0x0F,  nil, t)

		RunInstrAndAssertTimings(0x10,  nil, t)
		RunInstrAndAssertTimings(0x11,  nil, t)
		RunInstrAndAssertTimings(0x12, nil, t)
		RunInstrAndAssertTimings(0x13, nil, t)
		RunInstrAndAssertTimings(0x14, nil, t)
		RunInstrAndAssertTimings(0x15, nil, t)
		RunInstrAndAssertTimings(0x16, nil, t)
		RunInstrAndAssertTimings(0x17, nil, t)
		RunInstrAndAssertTimings(0x18,  nil, t)
		RunInstrAndAssertTimings(0x19, nil, t)
		RunInstrAndAssertTimings(0x1A, nil, t)
		RunInstrAndAssertTimings(0x1B, nil, t)
		RunInstrAndAssertTimings(0x1C, nil, t)
		RunInstrAndAssertTimings(0x1D, nil, t)
		RunInstrAndAssertTimings(0x1E, nil, t)
		RunInstrAndAssertTimings(0x1F, nil, t)

		RunInstrAndAssertTimings(0x21, nil, t)
		RunInstrAndAssertTimings(0x22, nil, t)
		RunInstrAndAssertTimings(0x23,  nil, t)
		RunInstrAndAssertTimings(0x24,  nil, t)
		RunInstrAndAssertTimings(0x25,  nil, t)
		RunInstrAndAssertTimings(0x26,  nil, t)
		RunInstrAndAssertTimings(0x27,  nil, t)
		RunInstrAndAssertTimings(0x29, nil, t)
		RunInstrAndAssertTimings(0x2A,  nil, t)
		RunInstrAndAssertTimings(0x2B,  nil, t)
		RunInstrAndAssertTimings(0x2C,  nil, t)
		RunInstrAndAssertTimings(0x2D,  nil, t)
		RunInstrAndAssertTimings(0x2E,  nil, t)
		RunInstrAndAssertTimings(0x2F,  nil, t)

		RunInstrAndAssertTimings(0x31,  nil, t)
		RunInstrAndAssertTimings(0x32,  nil, t)
		RunInstrAndAssertTimings(0x33,  nil, t)
		RunInstrAndAssertTimings(0x34,  nil, t)
		RunInstrAndAssertTimings(0x35,  nil, t)
		RunInstrAndAssertTimings(0x36,  nil, t)
		RunInstrAndAssertTimings(0x37, nil, t)
		RunInstrAndAssertTimings(0x39, nil, t)
		RunInstrAndAssertTimings(0x3A, nil, t)
		RunInstrAndAssertTimings(0x3B, nil, t)
		RunInstrAndAssertTimings(0x3C, nil, t)
		RunInstrAndAssertTimings(0x3D, nil, t)
		RunInstrAndAssertTimings(0x3E, nil, t)
		RunInstrAndAssertTimings(0x3F, nil, t)

		RunInstrAndAssertTimings(0x40, nil, t)
		RunInstrAndAssertTimings(0x41,  nil, t)
		RunInstrAndAssertTimings(0x42,  nil, t)
		RunInstrAndAssertTimings(0x43,  nil, t)
		RunInstrAndAssertTimings(0x44,  nil, t)
		RunInstrAndAssertTimings(0x45,  nil, t)
		RunInstrAndAssertTimings(0x46,  nil, t)
		RunInstrAndAssertTimings(0x47,  nil, t)
		RunInstrAndAssertTimings(0x48,  nil, t)
		RunInstrAndAssertTimings(0x49,  nil, t)
		RunInstrAndAssertTimings(0x4A,  nil, t)
		RunInstrAndAssertTimings(0x4B,  nil, t)
		RunInstrAndAssertTimings(0x4C,  nil, t)
		RunInstrAndAssertTimings(0x4D,  nil, t)
		RunInstrAndAssertTimings(0x4E,  nil, t)
		RunInstrAndAssertTimings(0x4F,  nil, t)

		RunInstrAndAssertTimings(0x50, nil, t)
		RunInstrAndAssertTimings(0x51,  nil, t)
		RunInstrAndAssertTimings(0x52,  nil, t)
		RunInstrAndAssertTimings(0x53,  nil, t)
		RunInstrAndAssertTimings(0x54,  nil, t)
		RunInstrAndAssertTimings(0x55,  nil, t)
		RunInstrAndAssertTimings(0x56,  nil, t)
		RunInstrAndAssertTimings(0x57,  nil, t)
		RunInstrAndAssertTimings(0x58,  nil, t)
		RunInstrAndAssertTimings(0x59,  nil, t)
		RunInstrAndAssertTimings(0x5A,  nil, t)
		RunInstrAndAssertTimings(0x5B,  nil, t)
		RunInstrAndAssertTimings(0x5C,  nil, t)
		RunInstrAndAssertTimings(0x5D,  nil, t)
		RunInstrAndAssertTimings(0x5E,  nil, t)
		RunInstrAndAssertTimings(0x5F,  nil, t)

		//0x6x
		RunInstrAndAssertTimings(0x60,  nil, t)
		RunInstrAndAssertTimings(0x61,  nil, t)
		RunInstrAndAssertTimings(0x62,  nil, t)
		RunInstrAndAssertTimings(0x63,  nil, t)
		RunInstrAndAssertTimings(0x64,  nil, t)
		RunInstrAndAssertTimings(0x65,  nil, t)
		RunInstrAndAssertTimings(0x66,  nil, t)
		RunInstrAndAssertTimings(0x67, nil, t)
		RunInstrAndAssertTimings(0x68,  nil, t)
		RunInstrAndAssertTimings(0x69,  nil, t)
		RunInstrAndAssertTimings(0x6A,  nil, t)
		RunInstrAndAssertTimings(0x6B,  nil, t)
		RunInstrAndAssertTimings(0x6C,  nil, t)
		RunInstrAndAssertTimings(0x6D,  nil, t)
		RunInstrAndAssertTimings(0x6E,  nil, t)
		RunInstrAndAssertTimings(0x6F,  nil, t)

		//0x7x
		RunInstrAndAssertTimings(0x70,  nil, t)
		RunInstrAndAssertTimings(0x71,  nil, t)
		RunInstrAndAssertTimings(0x72,  nil, t)
		RunInstrAndAssertTimings(0x73,  nil, t)
		RunInstrAndAssertTimings(0x74,  nil, t)
		RunInstrAndAssertTimings(0x75,  nil, t)
		RunInstrAndAssertTimings(0x76,  nil, t)
		RunInstrAndAssertTimings(0x77,  nil, t)
		RunInstrAndAssertTimings(0x78,  nil, t)
		RunInstrAndAssertTimings(0x79,  nil, t)
		RunInstrAndAssertTimings(0x7A,  nil, t)
		RunInstrAndAssertTimings(0x7B,  nil, t)
		RunInstrAndAssertTimings(0x7C,  nil, t)
		RunInstrAndAssertTimings(0x7D,  nil, t)
		RunInstrAndAssertTimings(0x7E,  nil, t)
		RunInstrAndAssertTimings(0x7F,  nil, t)

			//0x8x
		RunInstrAndAssertTimings(0x80, nil, t)
		RunInstrAndAssertTimings(0x81,  nil, t)
		RunInstrAndAssertTimings(0x82,  nil, t)
		RunInstrAndAssertTimings(0x83,  nil, t)
		RunInstrAndAssertTimings(0x84,  nil, t)
		RunInstrAndAssertTimings(0x85,  nil, t)
		RunInstrAndAssertTimings(0x86,  nil, t)
		RunInstrAndAssertTimings(0x87,  nil, t)
		RunInstrAndAssertTimings(0x88,  nil, t)
		RunInstrAndAssertTimings(0x89,  nil, t)
		RunInstrAndAssertTimings(0x8A,  nil, t)
		RunInstrAndAssertTimings(0x8B,  nil, t)
		RunInstrAndAssertTimings(0x8C,  nil, t)
		RunInstrAndAssertTimings(0x8D, nil, t)
		RunInstrAndAssertTimings(0x8E,  nil, t)
		RunInstrAndAssertTimings(0x8F,  nil, t)

		//0x9x
		RunInstrAndAssertTimings(0x90,  nil, t)
		RunInstrAndAssertTimings(0x91,  nil, t)
		RunInstrAndAssertTimings(0x92,  nil, t)
		RunInstrAndAssertTimings(0x93,  nil, t)
		RunInstrAndAssertTimings(0x94, nil, t)
		RunInstrAndAssertTimings(0x95,  nil, t)
		RunInstrAndAssertTimings(0x96,  nil, t)
		RunInstrAndAssertTimings(0x97,  nil, t)
		RunInstrAndAssertTimings(0x98,  nil, t)
		RunInstrAndAssertTimings(0x99,  nil, t)
		RunInstrAndAssertTimings(0x9A, nil, t)
		RunInstrAndAssertTimings(0x9B,  nil, t)
		RunInstrAndAssertTimings(0x9C,  nil, t)
		RunInstrAndAssertTimings(0x9D,  nil, t)
		RunInstrAndAssertTimings(0x9E,  nil, t)
		RunInstrAndAssertTimings(0x9F,  nil, t)

		//0xAx
		RunInstrAndAssertTimings(0xA0,  nil, t)
		RunInstrAndAssertTimings(0xA1,  nil, t)
		RunInstrAndAssertTimings(0xA2,  nil, t)
		RunInstrAndAssertTimings(0xA3,  nil, t)
		RunInstrAndAssertTimings(0xA4,  nil, t)
		RunInstrAndAssertTimings(0xA5,  nil, t)
		RunInstrAndAssertTimings(0xA6,  nil, t)
		RunInstrAndAssertTimings(0xA7,  nil, t)
		RunInstrAndAssertTimings(0xA8,  nil, t)
		RunInstrAndAssertTimings(0xA9,  nil, t)
		RunInstrAndAssertTimings(0xAA,  nil, t)
		RunInstrAndAssertTimings(0xAB,  nil, t)
		RunInstrAndAssertTimings(0xAC,  nil, t)
		RunInstrAndAssertTimings(0xAD,  nil, t)
		RunInstrAndAssertTimings(0xAE,  nil, t)
		RunInstrAndAssertTimings(0xAF,  nil, t)

		RunInstrAndAssertTimings(0xB0,  nil, t)
		RunInstrAndAssertTimings(0xB1,  nil, t)
		RunInstrAndAssertTimings(0xB2,  nil, t)
		RunInstrAndAssertTimings(0xB3,  nil, t)
		RunInstrAndAssertTimings(0xB4,  nil, t)
		RunInstrAndAssertTimings(0xB5,  nil, t)
		RunInstrAndAssertTimings(0xB6,  nil, t)
		RunInstrAndAssertTimings(0xB7,  nil, t)
		RunInstrAndAssertTimings(0xB8,  nil, t)
		RunInstrAndAssertTimings(0xB9,  nil, t)
		RunInstrAndAssertTimings(0xBA,  nil, t)
		RunInstrAndAssertTimings(0xBB,  nil, t)
		RunInstrAndAssertTimings(0xBC,  nil, t)
		RunInstrAndAssertTimings(0xBD,  nil, t)
		RunInstrAndAssertTimings(0xBE,  nil, t)
		RunInstrAndAssertTimings(0xBF,  nil, t)

		//0xCx
		RunInstrAndAssertTimings(0xC1,  nil, t)
		RunInstrAndAssertTimings(0xC3,  nil, t)
		RunInstrAndAssertTimings(0xC5,  nil, t)
		RunInstrAndAssertTimings(0xC6, nil, t)
		RunInstrAndAssertTimings(0xC7,  nil, t)
		RunInstrAndAssertTimings(0xC9,  nil, t)
		RunInstrAndAssertTimings(0xCD,  nil, t)
		RunInstrAndAssertTimings(0xCE, nil, t)
		RunInstrAndAssertTimings(0xCF,  nil, t)

		RunInstrAndAssertTimings(0xD1,  nil, t)
		RunInstrAndAssertTimings(0xD5,  nil, t)
		RunInstrAndAssertTimings(0xD6, nil, t)
		RunInstrAndAssertTimings(0xD7,  nil, t)
		RunInstrAndAssertTimings(0xD9,  nil, t)
		RunInstrAndAssertTimings(0xDE, nil, t)
		RunInstrAndAssertTimings(0xDF,  nil, t)

		//0xEx
		RunInstrAndAssertTimings(0xE0,  nil, t)
		RunInstrAndAssertTimings(0xE1,  nil, t)
		RunInstrAndAssertTimings(0xE2, nil, t)
		RunInstrAndAssertTimings(0xE5,  nil, t)
		RunInstrAndAssertTimings(0xE6, nil, t)
		RunInstrAndAssertTimings(0xE7,  nil, t)
		RunInstrAndAssertTimings(0xE8,  nil, t)
		RunInstrAndAssertTimings(0xE9, nil, t)
		RunInstrAndAssertTimings(0xEA,  nil, t)
		RunInstrAndAssertTimings(0xEE, nil, t)
		RunInstrAndAssertTimings(0xEF,  nil, t)

		//0xFx
		RunInstrAndAssertTimings(0xF0,  nil, t)
		RunInstrAndAssertTimings(0xF1,  nil, t)
		RunInstrAndAssertTimings(0xF2,  nil, t)
		RunInstrAndAssertTimings(0xF3, nil, t)
		RunInstrAndAssertTimings(0xF5,  nil, t)
		RunInstrAndAssertTimings(0xF6,  nil, t)
		RunInstrAndAssertTimings(0xF7,  nil, t)
		RunInstrAndAssertTimings(0xF8,  nil, t)
		RunInstrAndAssertTimings(0xF9, nil, t)
		RunInstrAndAssertTimings(0xFA,  nil, t)
		RunInstrAndAssertTimings(0xFB,  nil, t)
		RunInstrAndAssertTimings(0xFE,  nil, t)
		RunInstrAndAssertTimings(0xFF,  nil, t)


		//CONDITIONALS
		RunInstrAndAssertTimings(0x20, nil, t)
		RunInstrAndAssertTimings(0x28, []int{Z}, t)
		RunInstrAndAssertTimings(0x30,  nil, t)
		RunInstrAndAssertTimings(0x38,  []int{C}, t)
		RunInstrAndAssertTimings(0xC0,  nil, t)
		RunInstrAndAssertTimings(0xC2,  nil, t)
		RunInstrAndAssertTimings(0xC4,  nil, t)
		RunInstrAndAssertTimings(0xC8,  []int{Z}, t)
		RunInstrAndAssertTimings(0xCA,  []int{Z}, t)
		RunInstrAndAssertTimings(0xCC,  []int{Z}, t)
		RunInstrAndAssertTimings(0xD0,  nil, t)
		RunInstrAndAssertTimings(0xD2,  nil, t)
		RunInstrAndAssertTimings(0xD4,  nil, t)
		RunInstrAndAssertTimings(0xD8,  []int{C}, t)
		RunInstrAndAssertTimings(0xDA,  []int{C}, t)
		RunInstrAndAssertTimings(0xDC,  []int{C}, t)


		RunInstrAndAssertTimings(0x20, []int{Z}, t)
		RunInstrAndAssertTimings(0x28,  nil, t)
		RunInstrAndAssertTimings(0x30, []int{C}, t)
		RunInstrAndAssertTimings(0x38,  nil, t)
		RunInstrAndAssertTimings(0xC0,  []int{Z}, t)
		RunInstrAndAssertTimings(0xC2,  []int{Z}, t)
		RunInstrAndAssertTimings(0xC4,  []int{Z}, t)
		RunInstrAndAssertTimings(0xC8,  nil, t)
		RunInstrAndAssertTimings(0xCA,  nil, t)
		RunInstrAndAssertTimings(0xCC,  nil, t)
		RunInstrAndAssertTimings(0xD0,  []int{C}, t)
		RunInstrAndAssertTimings(0xD2,  []int{C}, t)
		RunInstrAndAssertTimings(0xD4,  []int{C}, t)
		RunInstrAndAssertTimings(0xD8,  nil, t)
		RunInstrAndAssertTimings(0xDA,  nil, t)
		RunInstrAndAssertTimings(0xDC,  nil, t)
	*/

	RunCBInstrAndAssertTimings(0x0, nil, t)
	RunCBInstrAndAssertTimings(0x1, nil, t)
	RunCBInstrAndAssertTimings(0x2, nil, t)
	RunCBInstrAndAssertTimings(0x3, nil, t)
	RunCBInstrAndAssertTimings(0x4, nil, t)
	RunCBInstrAndAssertTimings(0x5, nil, t)
	RunCBInstrAndAssertTimings(0x6, nil, t)
	RunCBInstrAndAssertTimings(0x7, nil, t)
	RunCBInstrAndAssertTimings(0x8, nil, t)
	RunCBInstrAndAssertTimings(0x9, nil, t)
	RunCBInstrAndAssertTimings(0xA, nil, t)
	RunCBInstrAndAssertTimings(0xB, nil, t)
	RunCBInstrAndAssertTimings(0xC, nil, t)
	RunCBInstrAndAssertTimings(0xD, nil, t)
	RunCBInstrAndAssertTimings(0xE, nil, t)
	RunCBInstrAndAssertTimings(0xF, nil, t)
	RunCBInstrAndAssertTimings(0x10, nil, t)
	RunCBInstrAndAssertTimings(0x11, nil, t)
	RunCBInstrAndAssertTimings(0x12, nil, t)
	RunCBInstrAndAssertTimings(0x13, nil, t)
	RunCBInstrAndAssertTimings(0x14, nil, t)
	RunCBInstrAndAssertTimings(0x15, nil, t)
	RunCBInstrAndAssertTimings(0x16, nil, t)
	RunCBInstrAndAssertTimings(0x17, nil, t)
	RunCBInstrAndAssertTimings(0x18, nil, t)
	RunCBInstrAndAssertTimings(0x19, nil, t)
	RunCBInstrAndAssertTimings(0x1A, nil, t)
	RunCBInstrAndAssertTimings(0x1B, nil, t)
	RunCBInstrAndAssertTimings(0x1C, nil, t)
	RunCBInstrAndAssertTimings(0x1D, nil, t)
	RunCBInstrAndAssertTimings(0x1E, nil, t)
	RunCBInstrAndAssertTimings(0x1F, nil, t)
	RunCBInstrAndAssertTimings(0x20, nil, t)
	RunCBInstrAndAssertTimings(0x21, nil, t)
	RunCBInstrAndAssertTimings(0x22, nil, t)
	RunCBInstrAndAssertTimings(0x23, nil, t)
	RunCBInstrAndAssertTimings(0x24, nil, t)
	RunCBInstrAndAssertTimings(0x25, nil, t)
	RunCBInstrAndAssertTimings(0x26, nil, t)
	RunCBInstrAndAssertTimings(0x27, nil, t)
	RunCBInstrAndAssertTimings(0x28, nil, t)
	RunCBInstrAndAssertTimings(0x29, nil, t)
	RunCBInstrAndAssertTimings(0x2A, nil, t)
	RunCBInstrAndAssertTimings(0x2B, nil, t)
	RunCBInstrAndAssertTimings(0x2C, nil, t)
	RunCBInstrAndAssertTimings(0x2D, nil, t)
	RunCBInstrAndAssertTimings(0x2E, nil, t)
	RunCBInstrAndAssertTimings(0x2F, nil, t)
	RunCBInstrAndAssertTimings(0x30, nil, t)
	RunCBInstrAndAssertTimings(0x31, nil, t)
	RunCBInstrAndAssertTimings(0x32, nil, t)
	RunCBInstrAndAssertTimings(0x33, nil, t)
	RunCBInstrAndAssertTimings(0x34, nil, t)
	RunCBInstrAndAssertTimings(0x35, nil, t)
	RunCBInstrAndAssertTimings(0x36, nil, t)
	RunCBInstrAndAssertTimings(0x37, nil, t)
	RunCBInstrAndAssertTimings(0x38, nil, t)
	RunCBInstrAndAssertTimings(0x39, nil, t)
	RunCBInstrAndAssertTimings(0x3A, nil, t)
	RunCBInstrAndAssertTimings(0x3B, nil, t)
	RunCBInstrAndAssertTimings(0x3C, nil, t)
	RunCBInstrAndAssertTimings(0x3D, nil, t)
	RunCBInstrAndAssertTimings(0x3E, nil, t)
	RunCBInstrAndAssertTimings(0x3F, nil, t)
	RunCBInstrAndAssertTimings(0x40, nil, t)
	RunCBInstrAndAssertTimings(0x41, nil, t)
	RunCBInstrAndAssertTimings(0x42, nil, t)
	RunCBInstrAndAssertTimings(0x43, nil, t)
	RunCBInstrAndAssertTimings(0x44, nil, t)
	RunCBInstrAndAssertTimings(0x45, nil, t)
	RunCBInstrAndAssertTimings(0x46, nil, t)
	RunCBInstrAndAssertTimings(0x47, nil, t)
	RunCBInstrAndAssertTimings(0x48, nil, t)
	RunCBInstrAndAssertTimings(0x49, nil, t)
	RunCBInstrAndAssertTimings(0x4A, nil, t)
	RunCBInstrAndAssertTimings(0x4B, nil, t)
	RunCBInstrAndAssertTimings(0x4C, nil, t)
	RunCBInstrAndAssertTimings(0x4D, nil, t)
	RunCBInstrAndAssertTimings(0x4E, nil, t)
	RunCBInstrAndAssertTimings(0x4F, nil, t)
	RunCBInstrAndAssertTimings(0x50, nil, t)
	RunCBInstrAndAssertTimings(0x51, nil, t)
	RunCBInstrAndAssertTimings(0x52, nil, t)
	RunCBInstrAndAssertTimings(0x53, nil, t)
	RunCBInstrAndAssertTimings(0x54, nil, t)
	RunCBInstrAndAssertTimings(0x55, nil, t)
	RunCBInstrAndAssertTimings(0x56, nil, t)
	RunCBInstrAndAssertTimings(0x57, nil, t)
	RunCBInstrAndAssertTimings(0x58, nil, t)
	RunCBInstrAndAssertTimings(0x59, nil, t)
	RunCBInstrAndAssertTimings(0x5A, nil, t)
	RunCBInstrAndAssertTimings(0x5B, nil, t)
	RunCBInstrAndAssertTimings(0x5C, nil, t)
	RunCBInstrAndAssertTimings(0x5D, nil, t)
	RunCBInstrAndAssertTimings(0x5E, nil, t)
	RunCBInstrAndAssertTimings(0x5F, nil, t)
	RunCBInstrAndAssertTimings(0x60, nil, t)
	RunCBInstrAndAssertTimings(0x61, nil, t)
	RunCBInstrAndAssertTimings(0x62, nil, t)
	RunCBInstrAndAssertTimings(0x63, nil, t)
	RunCBInstrAndAssertTimings(0x64, nil, t)
	RunCBInstrAndAssertTimings(0x65, nil, t)
	RunCBInstrAndAssertTimings(0x66, nil, t)
	RunCBInstrAndAssertTimings(0x67, nil, t)
	RunCBInstrAndAssertTimings(0x68, nil, t)
	RunCBInstrAndAssertTimings(0x69, nil, t)
	RunCBInstrAndAssertTimings(0x6A, nil, t)
	RunCBInstrAndAssertTimings(0x6B, nil, t)
	RunCBInstrAndAssertTimings(0x6C, nil, t)
	RunCBInstrAndAssertTimings(0x6D, nil, t)
	RunCBInstrAndAssertTimings(0x6E, nil, t)
	RunCBInstrAndAssertTimings(0x6F, nil, t)
	RunCBInstrAndAssertTimings(0x70, nil, t)
	RunCBInstrAndAssertTimings(0x71, nil, t)
	RunCBInstrAndAssertTimings(0x72, nil, t)
	RunCBInstrAndAssertTimings(0x73, nil, t)
	RunCBInstrAndAssertTimings(0x74, nil, t)
	RunCBInstrAndAssertTimings(0x75, nil, t)
	RunCBInstrAndAssertTimings(0x76, nil, t)
	RunCBInstrAndAssertTimings(0x77, nil, t)
	RunCBInstrAndAssertTimings(0x78, nil, t)
	RunCBInstrAndAssertTimings(0x79, nil, t)
	RunCBInstrAndAssertTimings(0x7A, nil, t)
	RunCBInstrAndAssertTimings(0x7B, nil, t)
	RunCBInstrAndAssertTimings(0x7C, nil, t)
	RunCBInstrAndAssertTimings(0x7D, nil, t)
	RunCBInstrAndAssertTimings(0x7E, nil, t)
	RunCBInstrAndAssertTimings(0x7F, nil, t)
	RunCBInstrAndAssertTimings(0x80, nil, t)
	RunCBInstrAndAssertTimings(0x81, nil, t)
	RunCBInstrAndAssertTimings(0x82, nil, t)
	RunCBInstrAndAssertTimings(0x83, nil, t)
	RunCBInstrAndAssertTimings(0x84, nil, t)
	RunCBInstrAndAssertTimings(0x85, nil, t)
	RunCBInstrAndAssertTimings(0x86, nil, t)
	RunCBInstrAndAssertTimings(0x87, nil, t)
	RunCBInstrAndAssertTimings(0x88, nil, t)
	RunCBInstrAndAssertTimings(0x89, nil, t)
	RunCBInstrAndAssertTimings(0x8A, nil, t)
	RunCBInstrAndAssertTimings(0x8B, nil, t)
	RunCBInstrAndAssertTimings(0x8C, nil, t)
	RunCBInstrAndAssertTimings(0x8D, nil, t)
	RunCBInstrAndAssertTimings(0x8E, nil, t)
	RunCBInstrAndAssertTimings(0x8F, nil, t)
	RunCBInstrAndAssertTimings(0x90, nil, t)
	RunCBInstrAndAssertTimings(0x91, nil, t)
	RunCBInstrAndAssertTimings(0x92, nil, t)
	RunCBInstrAndAssertTimings(0x93, nil, t)
	RunCBInstrAndAssertTimings(0x94, nil, t)
	RunCBInstrAndAssertTimings(0x95, nil, t)
	RunCBInstrAndAssertTimings(0x96, nil, t)
	RunCBInstrAndAssertTimings(0x97, nil, t)
	RunCBInstrAndAssertTimings(0x98, nil, t)
	RunCBInstrAndAssertTimings(0x99, nil, t)
	RunCBInstrAndAssertTimings(0x9A, nil, t)
	RunCBInstrAndAssertTimings(0x9B, nil, t)
	RunCBInstrAndAssertTimings(0x9C, nil, t)
	RunCBInstrAndAssertTimings(0x9D, nil, t)
	RunCBInstrAndAssertTimings(0x9E, nil, t)
	RunCBInstrAndAssertTimings(0x9F, nil, t)
	RunCBInstrAndAssertTimings(0xA0, nil, t)
	RunCBInstrAndAssertTimings(0xA1, nil, t)
	RunCBInstrAndAssertTimings(0xA2, nil, t)
	RunCBInstrAndAssertTimings(0xA3, nil, t)
	RunCBInstrAndAssertTimings(0xA4, nil, t)
	RunCBInstrAndAssertTimings(0xA5, nil, t)
	RunCBInstrAndAssertTimings(0xA6, nil, t)
	RunCBInstrAndAssertTimings(0xA7, nil, t)
	RunCBInstrAndAssertTimings(0xA8, nil, t)
	RunCBInstrAndAssertTimings(0xA9, nil, t)
	RunCBInstrAndAssertTimings(0xAA, nil, t)
	RunCBInstrAndAssertTimings(0xAB, nil, t)
	RunCBInstrAndAssertTimings(0xAC, nil, t)
	RunCBInstrAndAssertTimings(0xAD, nil, t)
	RunCBInstrAndAssertTimings(0xAE, nil, t)
	RunCBInstrAndAssertTimings(0xAF, nil, t)
	RunCBInstrAndAssertTimings(0xB0, nil, t)
	RunCBInstrAndAssertTimings(0xB1, nil, t)
	RunCBInstrAndAssertTimings(0xB2, nil, t)
	RunCBInstrAndAssertTimings(0xB3, nil, t)
	RunCBInstrAndAssertTimings(0xB4, nil, t)
	RunCBInstrAndAssertTimings(0xB5, nil, t)
	RunCBInstrAndAssertTimings(0xB6, nil, t)
	RunCBInstrAndAssertTimings(0xB7, nil, t)
	RunCBInstrAndAssertTimings(0xB8, nil, t)
	RunCBInstrAndAssertTimings(0xB9, nil, t)
	RunCBInstrAndAssertTimings(0xBA, nil, t)
	RunCBInstrAndAssertTimings(0xBB, nil, t)
	RunCBInstrAndAssertTimings(0xBC, nil, t)
	RunCBInstrAndAssertTimings(0xBD, nil, t)
	RunCBInstrAndAssertTimings(0xBE, nil, t)
	RunCBInstrAndAssertTimings(0xBF, nil, t)
	RunCBInstrAndAssertTimings(0xC0, nil, t)
	RunCBInstrAndAssertTimings(0xC1, nil, t)
	RunCBInstrAndAssertTimings(0xC2, nil, t)
	RunCBInstrAndAssertTimings(0xC3, nil, t)
	RunCBInstrAndAssertTimings(0xC4, nil, t)
	RunCBInstrAndAssertTimings(0xC5, nil, t)
	RunCBInstrAndAssertTimings(0xC6, nil, t)
	RunCBInstrAndAssertTimings(0xC7, nil, t)
	RunCBInstrAndAssertTimings(0xC8, nil, t)
	RunCBInstrAndAssertTimings(0xC9, nil, t)
	RunCBInstrAndAssertTimings(0xCA, nil, t)
	RunCBInstrAndAssertTimings(0xCB, nil, t)
	RunCBInstrAndAssertTimings(0xCC, nil, t)
	RunCBInstrAndAssertTimings(0xCD, nil, t)
	RunCBInstrAndAssertTimings(0xCE, nil, t)
	RunCBInstrAndAssertTimings(0xCF, nil, t)
	RunCBInstrAndAssertTimings(0xD0, nil, t)
	RunCBInstrAndAssertTimings(0xD1, nil, t)
	RunCBInstrAndAssertTimings(0xD2, nil, t)
	RunCBInstrAndAssertTimings(0xD3, nil, t)
	RunCBInstrAndAssertTimings(0xD4, nil, t)
	RunCBInstrAndAssertTimings(0xD5, nil, t)
	RunCBInstrAndAssertTimings(0xD6, nil, t)
	RunCBInstrAndAssertTimings(0xD7, nil, t)
	RunCBInstrAndAssertTimings(0xD8, nil, t)
	RunCBInstrAndAssertTimings(0xD9, nil, t)
	RunCBInstrAndAssertTimings(0xDA, nil, t)
	RunCBInstrAndAssertTimings(0xDB, nil, t)
	RunCBInstrAndAssertTimings(0xDC, nil, t)
	RunCBInstrAndAssertTimings(0xDD, nil, t)
	RunCBInstrAndAssertTimings(0xDE, nil, t)
	RunCBInstrAndAssertTimings(0xDF, nil, t)
	RunCBInstrAndAssertTimings(0xE0, nil, t)
	RunCBInstrAndAssertTimings(0xE1, nil, t)
	RunCBInstrAndAssertTimings(0xE2, nil, t)
	RunCBInstrAndAssertTimings(0xE3, nil, t)
	RunCBInstrAndAssertTimings(0xE4, nil, t)
	RunCBInstrAndAssertTimings(0xE5, nil, t)
	RunCBInstrAndAssertTimings(0xE6, nil, t)
	RunCBInstrAndAssertTimings(0xE7, nil, t)
	RunCBInstrAndAssertTimings(0xE8, nil, t)
	RunCBInstrAndAssertTimings(0xE9, nil, t)
	RunCBInstrAndAssertTimings(0xEA, nil, t)
	RunCBInstrAndAssertTimings(0xEB, nil, t)
	RunCBInstrAndAssertTimings(0xEC, nil, t)
	RunCBInstrAndAssertTimings(0xED, nil, t)
	RunCBInstrAndAssertTimings(0xEE, nil, t)
	RunCBInstrAndAssertTimings(0xEF, nil, t)
	RunCBInstrAndAssertTimings(0xF0, nil, t)
	RunCBInstrAndAssertTimings(0xF1, nil, t)
	RunCBInstrAndAssertTimings(0xF2, nil, t)
	RunCBInstrAndAssertTimings(0xF3, nil, t)
	RunCBInstrAndAssertTimings(0xF4, nil, t)
	RunCBInstrAndAssertTimings(0xF5, nil, t)
	RunCBInstrAndAssertTimings(0xF6, nil, t)
	RunCBInstrAndAssertTimings(0xF7, nil, t)
	RunCBInstrAndAssertTimings(0xF8, nil, t)
	RunCBInstrAndAssertTimings(0xF9, nil, t)
	RunCBInstrAndAssertTimings(0xFA, nil, t)
	RunCBInstrAndAssertTimings(0xFB, nil, t)
	RunCBInstrAndAssertTimings(0xFC, nil, t)
	RunCBInstrAndAssertTimings(0xFD, nil, t)
	RunCBInstrAndAssertTimings(0xFE, nil, t)
	RunCBInstrAndAssertTimings(0xFF, nil, t)
}

type MockMMU struct {
	memory map[types.Word]byte
}

func NewMockMMU() *MockMMU {
	var m *MockMMU = new(MockMMU)
	m.memory = make(map[types.Word]byte)
	return m
}

func (m *MockMMU) WriteByte(address types.Word, value byte) {
	m.memory[address] = value
}

func (m *MockMMU) WriteWord(address types.Word, value types.Word) {
	m.memory[address] = byte(value >> 8)
	m.memory[address+1] = byte(value & 0x00FF)
}

func (m *MockMMU) ReadByte(address types.Word) byte {
	return m.memory[address]
}

func (m *MockMMU) ReadWord(address types.Word) types.Word {
	a, b := m.memory[address], m.memory[address+1]
	return (types.Word(a) << 8) ^ types.Word(b)
}

func (m *MockMMU) SetInBootMode(mode bool) {
}

func (m *MockMMU) Reset() {
	m.memory = make(map[types.Word]byte)
}

func (m *MockMMU) LoadBIOS(data []byte) (bool, error) {
	return true, nil
}

func (m *MockMMU) LoadCartridge(cart *cartridge.Cartridge) {
}
