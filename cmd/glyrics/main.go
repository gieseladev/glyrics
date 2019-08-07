package main

import (
	"os"
)

func main() {
	app := GetApp()
	if err := app.Run(os.Args); err != nil {
		errlog.Fatal(err)
	}
}
