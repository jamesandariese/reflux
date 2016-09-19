package main

import (
	"flag"
	"log"
	"github.com/jamesandariese/reflux"
)

func main() {
	reflux.PrepareFlags("reflux")
	flag.Parse()
	
	if flag.NArg() % 2 != 1 {
		log.Fatalln("Command must have format: <stat_name> <field1-name> <field1-value> ... <fieldN-name> <fieldN-value>")
	}
	fields := make(map[string]interface{})
	for i := 1; i < flag.NArg(); i += 2 {
		fields[flag.Arg(i)] = flag.Arg(i+1)
	}
	if err := reflux.SendPointUsingFlags(flag.Arg(0), fields); err != nil {
		log.Fatalln(err)
	}
}
