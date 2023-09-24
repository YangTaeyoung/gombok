package main

import (
	"log"
	"os"

	"github.com/YangTaeyoung/gombok/parser"

	"github.com/urfave/cli/v2"
)

func GombokAction(_ *cli.Context) error {
	parser.Run()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "gombok"
	app.Usage = "Gombok is Lombok Style Code Generator for Go"
	app.Version = "1.0.0"
	app.Action = GombokAction

	if err := app.Run(os.Args); err != nil {
		log.Panicf("gombok error: %v", err)
	}
}
