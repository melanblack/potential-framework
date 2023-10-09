package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/melanblack/potential-framework/cfghive"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stderr)
	app := &cli.App{
		UseShortOptionHandling: true,
		Name:                   "cfghive cli",
		Usage:                  "A command line interface for cfghive",
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "Creates a new hive",
				ArgsUsage: "path",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "compress",
						Value:    false,
						Usage:    "Compress the hive",
						Aliases:  []string{"c"},
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					file, err := os.Create(c.Args().Get(0))
					if err != nil {
						return err
					}
					defer func(file *os.File) {
						err := file.Close()
						if err != nil {
							log.Fatal("failed to close file after creation")
						}
					}(file)

					reader := bufio.NewReader(file)
					writer := bufio.NewWriter(file)
					buf := bufio.NewReadWriter(reader, writer)

					hive := cfghive.NewBinHive(c.Bool("compress"), 9)
					hive.Stream = buf
					err = hive.Save()
					if err != nil {
						return err
					}
					err = writer.Flush()
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:      "import",
				Usage:     "Loads data from a json file into a hive",
				ArgsUsage: "data hive",
				Action: func(c *cli.Context) error {
					dataFile, err := os.Open(c.Args().Get(0))
					if err != nil {
						return err
					}
					defer dataFile.Close()

					hiveFile, err := os.OpenFile(c.Args().Get(1), os.O_RDWR, 0644)
					if err != nil {
						return err
					}
					defer hiveFile.Close()
					reader := bufio.NewReader(hiveFile)
					writer := bufio.NewWriter(hiveFile)
					buf := bufio.NewReadWriter(reader, writer)
					hive := cfghive.NewBinHive(false, 0)
					hive.Stream = buf
					err = hive.Load()
					_, err = hiveFile.Seek(0, 0)
					if err != nil {
						return err
					}
					if err != nil {
						return err
					}
					data := make(map[string]interface{})
					codec := json.NewDecoder(dataFile)
					err = codec.Decode(&data)
					if err != nil {
						return err
					}
					for k, v := range data {
						err = hive.Set(k, v)
						if err != nil {
							return err
						}
					}
					err = hive.Save()
					if err != nil {
						return err
					}
					err = writer.Flush()
					if err != nil {
						return err
					}
					fmt.Printf("new hive size: %d\n", cfghive.HiveSize(*hive.GetData()))
					return nil
				},
			},
			{
				Name:      "dump",
				ArgsUsage: "<hive file>",
				Action: func(context *cli.Context) error {
					file, err := os.Open(context.Args().Get(0))
					if err != nil {
						return err
					}
					defer file.Close()

					reader := bufio.NewReader(file)
					writer := bufio.NewWriter(file)
					buf := bufio.NewReadWriter(reader, writer)
					hive := cfghive.NewBinHive(false, 0)
					hive.Stream = buf
					err = hive.Load()
					if err != nil {
						log.Fatal(err)
					}

					data := hive.GetData()
					cfghive.HiveDump(data)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
