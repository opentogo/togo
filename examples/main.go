package main

import (
	"fmt"
	"os"

	"github.com/opentogo/togo"
	"github.com/opentogo/togo/examples/service"
)

func main() {
	config := togo.LoadJSONFile("config.json")

	t := togo.Init("togo-example-service", config)
	t.Register(service.NewService())

	if err := t.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
