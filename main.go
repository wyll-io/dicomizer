package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/suyashkumar/dicom"
	"github.com/urfave/cli/v2"
	"github.com/wyll-io/dicomizer/pkg/anonymize"
)

var (
	awsCfg aws.Config
)

var app = &cli.App{
	Name:    "dicomizer",
	Version: "0.1.0",
	Before: func(ctx *cli.Context) error {
		for _, arg := range ctx.Args().Slice() {
			if arg == "--help" || arg == "-h" {
				return nil
			}
		}

		var err error
		awsCfg, err = config.LoadDefaultConfig(ctx.Context)
		if err != nil {
			return err
		}

		return nil
	},
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

				if err := anonymize.Anonymize(&dataset); err != nil {
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
				defer f.Close()

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
