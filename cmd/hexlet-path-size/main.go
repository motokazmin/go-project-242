package main

import (
	"context"
	"log"
	"os"

	cli "github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args().Slice()

			if len(args) == 0 {
				return cli.ShowAppHelp(cmd)
			}

			log.Printf("Начинаем анализ пути: %s\n", cmd.Args().First())
			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
