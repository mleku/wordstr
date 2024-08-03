package main

import (
	"fmt"
	"os"

	"wordstr.mleku.dev/wordstr"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: wordstr from/to <nsec>/<word key>")
		os.Exit(1)
	}
	// calculate the place multipliers and assemble them in descending order.
	var err error
	places := wordstr.GetPlaces()
	var hsec, nsec, words string
	switch {
	case os.Args[1] == "from":
		if words, err = wordstr.FromNsec(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println(words)
	case os.Args[1] == "to":
		if hsec, nsec, err = wordstr.ToNsec(places, os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("HSEC=\"%s\"\nNSEC=\"%s\"\n", hsec, nsec)
	default:
		fmt.Fprintf(os.Stderr, "%s [from|to] <nsec>/<word key>\n", os.Args[0])
	}
}
