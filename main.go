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

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("id error")
	}

	run.Run(id)
}
