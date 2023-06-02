package main

import "passKeeper/internal/cmd"

func main() {

	root := cmd.NewRootCommand()
	cmd.Execute(root)

}
