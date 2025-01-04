package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

const VERBOSE = false

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
	sol2(program)
}

func sol1(program []int) {
	network := startNetwork(program, 50)
	router := Router{network: network}
	router.Start()
}

func sol2(program []int) {
	network := startNetwork(program, 50)
	router := Router{network: network, nat: &NAT{}}
	router.Start()
}

func startNetwork(program []int, n int) Network {
	network := make(Network, n)
	for i := range n {
		network[i] = NewHost(program, byte(i))
		network[i].Start()
	}

	log("started %d hosts\n", len(network))
	return network
}

type Address = byte

type Packet struct {
	Target Address
	X, Y   int
}

type Network = []*Host

type Host struct {
	addr Address
	vm   *vm.VM

	sendInput vmio.SendInputFunc
	stopInput context.CancelFunc
	outbuf    *bytes.Buffer

	outbound chan []Packet
	inbound  chan Packet
}

func NewHost(program []int, addr Address) *Host {
	m, inbuf, outbuf := vm.NewVMPiped(program, vm.WithNumericID(int(addr)))
	sendInput, stop := vmio.StartInputSender(inbuf)

	// assign address to host
	addrAssign := fmt.Sprintf("%d\n", addr)
	sendInput([]byte(addrAssign))
	m.RunForInput()

	return &Host{
		addr: addr,
		vm:   m,

		sendInput: sendInput,
		stopInput: stop,
		outbuf:    outbuf,

		// FIXME: These channels could be unbuffered with better goroutine handling.
		// Currently a host could block the main routine (while other hosts stay alive, thus not "deadlocking" fully)
		// if the router is sending to host.Inbound while the host is writing to host.Outbound
		outbound: make(chan []Packet, 10),
		inbound:  make(chan Packet, 10),
	}
}

const HOST_TICKRATE = 1 * time.Millisecond

func (h *Host) Start() {
	go func() {
		ticker := time.NewTicker(HOST_TICKRATE)

		for {
			select {
			case p := <-h.inbound:
				log("%02d - inputting %+v\n", h.addr, p)

				data := fmt.Sprintf("%d\n", p.X)
				h.sendInput([]byte(data))
				h.vm.RunForInput()

				data = fmt.Sprintf("%d\n", p.Y)
				h.sendInput([]byte(data))
				h.vm.RunForInput()

				ticker.Reset(HOST_TICKRATE)

			case <-ticker.C:
				h.sendInput([]byte("-1\n"))
				h.vm.RunForInput()
			}

			output := vmio.GetOutput(h.outbuf)

			if len(output)%3 == 0 {
				packets := make([]Packet, 0)
				for i := 0; i < len(output)-1; i += 3 {
					p := Packet{
						Target: byte(output[i]),
						X:      output[i+1], Y: output[i+2],
					}
					packets = append(packets, p)
				}
				log("%02d - outputting %d packets: %+v\n", h.addr, len(packets), packets)
				h.outbound <- packets
			}
		}
	}()
}

func (h *Host) Inbound() chan<- Packet {
	return h.inbound
}

func (h *Host) Outbound() <-chan []Packet {
	return h.outbound
}

type Router struct {
	network Network
	nat     *NAT

	lastResumePacket Packet

	pending []Packet
}

func (r *Router) Start() {
	resumeTickrate := HOST_TICKRATE * 3
	tResume := time.NewTicker(resumeTickrate)

	for {
		// recieve messages
		for _, host := range r.network {
			select {
			case p := <-host.Outbound():
				r.pending = append(r.pending, p...)
			default:
			}
		}

		// send messages
		for _, p := range r.pending {
			if p.Target == 255 {
				if r.nat == nil {
					fmt.Printf("solution 1: %d\n", p.Y)
					return
				}

				r.nat.Set(p)
			}

			if int(p.Target) > len(r.network)-1 {
				log("Skipped packet: %+v\n", p)
				continue
			}
			host := r.network[p.Target]
			host.Inbound() <- p
		}

		// clear sent messages
		if len(r.pending) > 0 {
			tResume.Reset(resumeTickrate)
			r.pending = r.pending[:0]
		}

		if r.nat != nil {
			select {
			case <-tResume.C:
				last := r.nat.Get()

				log("restarting comms with: %+v\n", last)
				if last.Y == r.lastResumePacket.Y {
					fmt.Printf("solution 2: %d\n", last.Y)
					return
				}

				r.pending = append(r.pending, Packet{Target: 0, X: last.X, Y: last.Y})
				r.lastResumePacket = last
				tResume.Reset(resumeTickrate)
			default:
			}

		}
	}
}

type NAT struct {
	lastPacket Packet
}

func (n *NAT) Set(p Packet) {
	n.lastPacket = p
}

func (n *NAT) Get() Packet {
	return n.lastPacket
}

func log(msg string, args ...any) {
	if VERBOSE {
		fmt.Printf(msg, args...)
	}
}
