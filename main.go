package main

import "github.com/linuxsuren/atest-mcp-server/cmd"

func main() {
	c := cmd.NewRootCmd()
	err := c.Execute()
	if err != nil {
		panic(err)
	}
}
