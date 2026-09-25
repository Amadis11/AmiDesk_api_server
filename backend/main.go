package main

import (
	"amidesk-api-server/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
