package main

import (
	"fmt"
	"os"

	"github.com/suyashkumar/dicom"
	"github.com/urfave/cli/v2"
	anonimize "github.com/wyll-io/dicomizer/pkg/anonymize"
)

var app = &cli.App{
	Name:    "dicomizer",
	Version: "0.1.0",
	Commands: []*cli.Command{
		{
			Name:      "anonymize",
			ArgsUsage: "<input_file> [output_dicom]",
			Flags:     []cli.Flag{},
			Action: func(ctx *cli.Context) error {
				if !ctx.Args().Present() {
					return fmt.Errorf("missing input file")
				}

				dataset, err := dicom.ParseFile(ctx.Args().First(), nil)
				if err != nil {
					return err
				}

				if err := anonimize.Anonymize(&dataset); err != nil {
					return err
				}

				output := "anonymized.dcm"
				if ctx.Args().Len() > 1 {
					output = ctx.Args().Get(2)
				}

				f, err := os.Create(output)
				if err != nil {
					return err
				}

				// ! Disable VR verification for PixelData in case it is OB instead of OW and
				// ! not a little endian.
				return dicom.Write(f, dataset, dicom.SkipVRVerification())
			},
		},
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			DefaultText: "enable verbose mode",
		},
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
