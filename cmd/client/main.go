package main

import (
	"passKeeper/internal/cmd"

	"github.com/charmbracelet/log"
)

func main() {

	root := cmd.NewRootCommand()
	err := cmd.Execute(root)
	if err != nil {
		log.Error(err)
	}

}
