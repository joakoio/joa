package shell

import (
	"fmt"
	"os"
	"log"

	"github.com/urfave/cli"
	"github.com/joakoio/joa/server"
)

type flags struct {
	language string
}

var Flags = flags{}

func Run() {
	app := startCli()
	setFlags(app)
	setMainAction(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func startCli() *cli.App {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "server"
	app.Usage = "This will start the server!"
	return app
}

func setFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &Flags.language,
		},
	}
}

func setMainAction(app *cli.App) {
	app.Action = func(c *cli.Context) error {
		name := "someone"
		if c.NArg() > 0 {
			name = c.Args()[0]
		}
		if Flags.language == "spanish" {
			fmt.Println("Hola", name)

			fmt.Println("Starting server!")
			server.Start()
		} else {
			fmt.Println("Hello", name)
		}



		return nil
	}
}
