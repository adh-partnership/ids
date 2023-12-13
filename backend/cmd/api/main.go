package main

import (
	"fmt"
	"os"

	"github.com/adh-partnership/ids/backend/cmd/api/app"
)

func main() {
	app := app.NewRootCommand()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
