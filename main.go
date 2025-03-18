package main

import (
	"go-extractor/cmd"
	"go-extractor/config"
)

func main() {
	config.Init()
	cmd.Execute()
}
