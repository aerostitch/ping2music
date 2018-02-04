package main

import "fmt"

func main() {
	srv := NewICMPServer()
	go func() {
		for c := range srv.Seeds {
			fmt.Printf("%d\n", *c)
		}
	}()
	srv.ListenIPv4AndIPv6()
}
