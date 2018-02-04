package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"io"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// ICMPServer holds a few methods on the management of our ICMP server
type ICMPServer struct {
	Seeds chan *int64
	numRe *regexp.Regexp
}

// NewICMPServer creates a new ICMPServer struct and initiates the Seeds channel
// Warning: you need to start a consumer for the channel before starting to
// listen for ICMP packets or you will have a deadlock
func NewICMPServer() *ICMPServer {
	s := ICMPServer{}
	s.Seeds = make(chan *int64)
	s.numRe = regexp.MustCompile("[0-9]+")
	return &s
}

// Listen listens on the ICMP socket of your choice and calculates a seed based in
// the current time and the sender's IP.
// The icmp package is only supported in Mac and Linux.
// Be careful, this does not close the Seeds channel.
func (s *ICMPServer) Listen(network, address string) error {
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	default:
		return fmt.Errorf("listening for icmp packets is not supported on %s", runtime.GOOS)
	}

	c, err := icmp.ListenPacket(network, address)
	if err != nil {
		return fmt.Errorf("ListenPacket returned: %s", err)
	}
	defer s.Close(c)

	for {
		// We don't care about the message itself so no buffer allocation is needed
		_, peer, err := c.ReadFrom(nil)
		if err != nil {
			return fmt.Errorf("ReadFrom returned: %s", err)
		}
		seed := int64(time.Now().Nanosecond())
		for _, nb := range s.numRe.FindAllString(peer.String(), -1) {
			if n, err := strconv.ParseInt(nb, 10, 32); err == nil {
				seed = seed - n
			}
		}
		s.Seeds <- &seed
	}
}

// Close is used to handle errors when closing resources
func (s *ICMPServer) Close(resource io.Closer) {
	if err := resource.Close(); err != nil {
		log.Fatal(err)
	}
}

// ListenIPv4 listens on the ICMPv4 raw socket for all ICMPv4 packages and calls
// wg.Done() when exiting.
// Be careful, this does not close the Seeds channel.
func (s *ICMPServer) ListenIPv4(wg *sync.WaitGroup) {
	if err := s.Listen("ip4:1", "0.0.0.0"); err != nil {
		log.Printf("Fatal error reported by the IPv6 server: %s", err)
	}
	wg.Done()
}

// ListenIPv6 listens on the ICMPv6 raw socket for all ICMPv4 packages and calls
// wg.Done() when exiting.
// Be careful, this does not close the Seeds channel.
func (s *ICMPServer) ListenIPv6(wg *sync.WaitGroup) {
	if err := s.Listen("ip6:58", "::"); err != nil {
		log.Printf("Fatal error reported by the IPv6 server: %s", err)
	}
	wg.Done()
}

// ListenIPv4AndIPv6 starts an ICMP server that listens on ICMPv4 and ICMPv6
// raw sockets so it needs root privileges.
// It also takes care of closing the Seeds channel.
// Note: You don't need to adjust net.ipv4.ping_group_range as we are NOT using
// unprivileged ping.
func (s *ICMPServer) ListenIPv4AndIPv6() {
	var wg sync.WaitGroup
	wg.Add(2)
	go s.ListenIPv4(&wg)
	go s.ListenIPv6(&wg)
	wg.Wait()
	close(s.Seeds)
}
