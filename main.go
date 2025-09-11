package main

import (
	"github.com/linuxsuren/atest-mcp-server/cmd"
	"os"
)

func main() {
	c := cmd.NewRootCmd()
	err := c.Execute()
	if err != nil {
		os.Exit(1)
	}
}
