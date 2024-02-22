package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/suyashkumar/dicom"
	"github.com/urfave/cli/v2"
	"github.com/wyll-io/dicomizer/internal/check"
	"github.com/wyll-io/dicomizer/internal/database"
	"github.com/wyll-io/dicomizer/internal/scheduler"
	"github.com/wyll-io/dicomizer/internal/storage/s3"
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
	Version: "1.0.0",
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
					Value:       AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
				&cli.StringFlag{
					Name:     "pacs",
					Category: "DICOM",
					Usage:    "PACS server",
					Required: true,
					EnvVars:  []string{"PACS_SERVER"},
				},
				&cli.StringFlag{
					Name:     "aet",
					Category: "DICOM",
					Usage:    "AET",
					Required: true,
					EnvVars:  []string{"AET"},
				},
				&cli.StringFlag{
					Name:     "aec",
					Category: "DICOM",
					Usage:    "AEC",
					Required: true,
					EnvVars:  []string{"AEC"},
				},
				&cli.StringFlag{
					Name:     "aem",
					Category: "DICOM",
					Usage:    "AEM",
					Required: true,
					EnvVars:  []string{"AEM"},
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

				if strings.Split(ctx.Args().First(), ":")[1] == "" {
					return fmt.Errorf("missing port in HTTP address")
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

				dbClient := database.New(awsCfg, ctx.String("dynamodb-table"))
				s3Client := s3.NewClient(awsCfg)

				s, err := scheduler.Create(
					ctx.Context,
					crontab,
					s3Client,
					dbClient,
					ctx.String("pacs"),
					ctx.String("aet"),
					ctx.String("aec"),
					ctx.String("aem"),
				)
				if err != nil {
					return err
				}
				defer s.Shutdown()

				fmt.Printf("starting scheduler with \"%s\"...\n", crontab)
				s.Start()

				fmt.Printf("web server started at http://%s\n", addr)
				return http.ListenAndServe(addr, web.RegisterHandlers(awsCfg, dbClient))
			},
		},
		{
			Name:        "check-patient",
			Description: "Check a patient for new DICOM files",
			ArgsUsage:   "<pk>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dynamodb-table",
					Category:    "AWS",
					Usage:       "DynamoDB table name",
					DefaultText: AWS_DYNAMODB_TABLE,
					Value:       AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
				&cli.StringFlag{
					Name:     "pacs",
					Category: "DICOM",
					Usage:    "PACS server",
					Required: true,
					EnvVars:  []string{"PACS_SERVER"},
				},
				&cli.StringFlag{
					Name:     "aet",
					Category: "DICOM",
					Usage:    "AET",
					Required: true,
					EnvVars:  []string{"AET"},
				},
				&cli.StringFlag{
					Name:     "aec",
					Category: "DICOM",
					Usage:    "AEC",
					Required: true,
					EnvVars:  []string{"AEC"},
				},
				&cli.StringFlag{
					Name:     "aem",
					Category: "DICOM",
					Usage:    "AEM",
					Required: true,
					EnvVars:  []string{"AEM"},
				},
			},
			Before: func(ctx *cli.Context) error {
				if ctx.Args().First() == "" {
					return fmt.Errorf("missing pk")
				}

				return nil
			},
			Action: func(ctx *cli.Context) error {
				dbClient := database.New(awsCfg, ctx.String("dynamodb-table"))
				s3Client := s3.NewClient(awsCfg)

				fmt.Println("Fetching patient info...")
				pInfo, err := dbClient.GetPatientInfo(ctx.Context, ctx.Args().First())
				if err != nil {
					return err
				}

				return check.CheckPatientDCM(
					ctx.Context,
					s3Client,
					dbClient,
					ctx.String("pacs"),
					ctx.String("aet"),
					ctx.String("aec"),
					ctx.String("aem"),
					pInfo,
				)
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
