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
		UsageText: "hexlet-path-size [options] <path>\n\npath - path to a file or directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
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
			all := cmd.Bool("all")
			recursive := cmd.Bool("recursive")

			// Вызываем функцию GetPathSize из библиотеки
			sizeStr, err := code.GetPathSize(path, recursive, human, all)
			if err != nil {
				return fmt.Errorf("ошибка: %w", err)
			}

			// Выводим результат в формате: <размер>\t<путь>
			fmt.Printf("%s\t%s\n", sizeStr, path)
			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
