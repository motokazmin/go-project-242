package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"code"
	cli "github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
				Value:   false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args().Slice()

			if len(args) == 0 {
				return cli.ShowAppHelp(cmd)
			}

			path := cmd.Args().First()
			human := cmd.Bool("human")

			// Вызываем функцию GetSize из библиотеки
			result, err := code.GetSize(path, false, human, false)
			if err != nil {
				return fmt.Errorf("ошибка: %w", err)
			}

			// Выводим результат в формате: <размер>\t<путь>
			fmt.Println(result)
			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
