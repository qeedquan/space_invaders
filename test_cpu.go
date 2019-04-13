// +build ignore
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var (
	dflag = flag.Bool("d", false, "show disassembled output")
)

func main() {
	flag.Parse()
	testrom("cpu_tests/TST8080.COM")
	testrom("cpu_tests/CPUTEST.COM")
	testrom("cpu_tests/8080PRE.COM")
	testrom("cpu_tests/8080EXM.COM")
}

func testrom(name string) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("%s", err)
	}

	t := time.Now()
	c := NewCPU()
	copy(c.M[0x100:], b)
	c.PC = 0x100

	c.Write8(5, 0xc9)

	for {
		curpc := c.PC

		if c.HLT {
			fmt.Printf("HLT at %04X\n", c.PC)
		}

		if c.PC == 5 {
			if c.R[C] == 9 {
				for i := c.de(); ; i++ {
					v := c.Read8(i)
					if v == 0x24 {
						break
					}
					fmt.Printf("%c", v)
				}
			}
			if c.R[C] == 2 {
				fmt.Printf("%c", c.R[E])
			}
		}

		if *dflag {
			fmt.Printf("%s", c.Disasm(c.PC))
		}
		c.Step()
		if c.PC == 0 {
			fmt.Printf("\nJumped to 0x0000 from 0x%04X\n", curpc)
			break
		}
	}
	fmt.Printf("\nTook %v\n\n", time.Since(t))
}
