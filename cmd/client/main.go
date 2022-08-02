package main

import (
	"flag"
	"github.com/lordvidex/wordle-wf/cmd/client/local"
	"github.com/lordvidex/wordle-wf/cmd/client/remote"
)

func main() {
	tp := flag.String("type", "local", "the type of client to run, either local or remote")
	flag.Parse()
	if *tp == "remote" {
		remote.Start()
	} else {
		local.Start()
	}
}
