package main

import (
	"log"
	"os"
	"strconv"

	"github.com/shoooooman/mg-rs/run"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("input length is too small")
	}

	arg := os.Args[1]
	if arg == "master" {
		run.Master()
		return
	}

	id, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatal("id error")
	}
	run.Run(id)
}
