package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
}

var SETUP = []string{
	"south",
	"take fuel cell",
	"south",
	"take manifold",
	"north",
	"north",
	"north",
	"take candy cane",
	"south",
	"west",
	"take mutex",
	"south",
	"south",
	"take coin",
	"west",
	"take dehydrated water",
	"south",
	"take prime numbers",
	"north",
	"east",
	"north",
	"east",
	"take cake",
	"north",
	"west",
	"south", // security check
}

var ITEMS = []string{
	"fuel cell",
	"manifold",
	"candy cane",
	"mutex",
	"coin",
	"dehydrated water",
	"prime numbers",
	"cake",
}

func sol1(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program, vm.WithIOMode(vm.IOAscii))
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	outchan := make(chan string, 40)

	go func() {
		for {
			if m.Halted() {
				return
			}

			m.RunForInput()

			out := outbuf.String()
			outbuf.Reset()

			if out == "" {
				continue
			}

			outchan <- out
		}
	}()

	// setup
	for _, inp := range SETUP {
		cmd := fmt.Sprintf("%s\n", inp)
		sendInput([]byte(cmd))
		expect(outchan, "Command")
	}

	// bruteforce
	p := Player{
		sendInput: sendInput,
		outchan:   outchan,
		carrying:  make([]int, 0),
	}
	p.dropAll()

	for _, items := range powerset(len(ITEMS)) {
		p.takeThese(items)

		sendInput([]byte("west\n"))
		msg, _ := expect(p.outchan, "keypad")

		if m.Halted() {
			formatted := make([]string, 0)
			for _, ix := range items {
				formatted = append(formatted, ITEMS[ix])
			}
			fmt.Printf("Items: %+v\n", formatted)
			fmt.Printf("RES: %s\n", msg)
			return
		}
		p.dropCarry()
	}

}

type Player struct {
	sendInput vmio.SendInputFunc
	outchan   chan string

	carrying []int
}

// it's dangerous to go alone
func (p *Player) takeThese(items []int) {
	for _, ix := range items {
		p.take(ix)
	}
}

func (p *Player) dropAll() {
	for ix := range ITEMS {
		p.drop(ix)
	}
	p.carrying = p.carrying[:0]
}

func (p *Player) dropCarry() {
	for _, ix := range p.carrying {
		p.drop(ix)
	}
	p.carrying = p.carrying[:0]
}

func (p *Player) drop(ix int) {
	cmd := fmt.Sprintf("drop %s\n", ITEMS[ix])
	p.sendInput([]byte(cmd))
	expect(p.outchan, "Command")
}

func (p *Player) take(ix int) {
	cmd := fmt.Sprintf("take %s\n", ITEMS[ix])
	p.sendInput([]byte(cmd))
	expect(p.outchan, "Command")
	p.carrying = append(p.carrying, ix)
}

func expect(msgs chan string, exp string) (string, bool) {
	t := time.NewTimer(10 * time.Millisecond)
	for {
		select {
		case msg := <-msgs:
			if strings.Contains(msg, exp) {
				return msg, true
			}
		case <-t.C:
			return "", false
		}
	}
}

func powerset(length int) [][]int {
	out := make([][]int, 0)

	// for a set of l=2 (i.e. {1, 2}) outs are:
	//
	// 	{} -> 00
	// 	{1} -> 01
	// 	{2} -> 10
	// 	{1,2} -> 11
	totalPositions := 1 << length

	for positions := range totalPositions {
		set := make([]int, 0)
		for ix := range length {
			if (positions>>ix)&1 == 1 {
				set = append(set, ix)
			}
		}
		out = append(out, set)
	}

	return out
}
