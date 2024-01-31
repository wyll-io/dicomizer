package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/suyashkumar/dicom"
	"github.com/urfave/cli/v2"
	"github.com/wyll-io/dicomizer/internal/scheduler"
	"github.com/wyll-io/dicomizer/internal/storage/glacier"
	"github.com/wyll-io/dicomizer/internal/web"
	"github.com/wyll-io/dicomizer/pkg/anonymize"
)

const (
	AWS_DYNAMODB_TABLE = "dicomizer"
)

var (
	awsCfg           aws.Config
	mandatoryEnvVars = []string{
		"JWT_SECRET",
		"ADMIN_PASSWORD",
	}
)

var app = &cli.App{
	Name:    "dicomizer",
	Version: "0.1.0",
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return nil
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
			Action: func(ctx *cli.Context) error {
				if !ctx.Args().Present() {
					return fmt.Errorf("missing input file")
				}

				dataset, err := dicom.ParseFile(ctx.Args().First(), nil)
				if err != nil {
					return err
				}

				if err := anonymize.AnonymizeDataset(&dataset); err != nil {
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
		{
			Name:        "upload",
			ArgsUsage:   "<input_file>...",
			Description: "Upload DICOM files to AWS Glacier",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dynamodb-table",
					Category:    "AWS",
					Usage:       "DynamoDB table name",
					DefaultText: AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
			},
			Action: func(ctx *cli.Context) error {
				if !ctx.Args().Present() {
					return fmt.Errorf("missing input file")
				}

				client := glacier.NewClient(awsCfg)
				for _, arg := range ctx.Args().Slice() {
					client.UploadFile(ctx.Context, arg, glacier.Options{
						VaultName: "dicomizer",
					})
				}

				return nil
			},
		},
		{
			Name: "start",
			Description: `Start the DICOMizer server. 
If no arguments are provided, the server will use the environment variables:
- HTTP_HOST
- HTTP_PORT
- CRONTAB`,
			ArgsUsage: "[HOST:PORT CRONTAB]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dynamodb-table",
					Category:    "AWS",
					Usage:       "DynamoDB table name",
					DefaultText: AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
			},
			Before: func(ctx *cli.Context) error {
				if ctx.Args().Len() < 2 {
					if os.Getenv("HTTP_HOST") == "" &&
						os.Getenv("HTTP_PORT") == "" &&
						os.Getenv("CRONTAB") == "" {
						return fmt.Errorf("missing HTTP host or HTTP port or crontab")
					}
				}

				return nil
			},
			Action: func(ctx *cli.Context) error {
				addr := ctx.Args().First()
				if addr == "" {
					addr = os.Getenv("HTTP_HOST") + ":" + os.Getenv("HTTP_PORT")
				}

				crontab := ctx.Args().Get(1)
				if crontab == "" {
					crontab = os.Getenv("CRONTAB")
				}

				s, err := scheduler.Create(crontab)
				if err != nil {
					return err
				}
				defer s.Shutdown()

				fmt.Printf("starting scheduler with \"%s\"...\n", crontab)
				s.Start()

				fmt.Printf("web server started at http://%s\n", addr)
				return http.ListenAndServe(addr, web.RegisterHandlers(awsCfg, ctx.String("dynamodb-table")))
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

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic("error while loading .env file")
	}
	for _, env := range mandatoryEnvVars {
		if os.Getenv(env) == "" {
			panic(env + " is not set")
		}
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
