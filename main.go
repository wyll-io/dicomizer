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
	dicomFlags = []cli.Flag{
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
	}
)

var app = &cli.App{
	Name:    "dicomizer",
	Version: "v1.0.5",
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
- CRONTAB`,
			ArgsUsage: "[CRONTAB]",
			Flags: append([]cli.Flag{
				&cli.StringFlag{
					Name:        "dynamodb-table",
					Category:    "AWS",
					Usage:       "DynamoDB table name",
					DefaultText: AWS_DYNAMODB_TABLE,
					Value:       AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
				&cli.StringFlag{
					Name:     "center",
					Category: "AWS",
					Usage:    "Laboratory name used as root folder in AWS S3",
					Required: true,
					EnvVars:  []string{"LABORATORY"},
				},
				&cli.StringFlag{
					Name:     "bind",
					Category: "HTTP",
					Usage:    "Address and port to listen and bind",
					EnvVars:  []string{"HTTP_BIND"},
					Value:    "localhost:80",
				},
				&cli.StringFlag{
					Name:     "crontab",
					Category: "CRON",
					Usage:    "crontab to execute pacs check. E.g: \"* * * * *\"",
					EnvVars:  []string{"CRONTAB"},
					Required: true,
				},
			}, dicomFlags...),
			Before: func(ctx *cli.Context) error {
				return checkMandatoryStringFlags(
					[]string{"pacs", "aet", "aem", "aec", "center"},
					ctx,
				)
			},
			Action: func(ctx *cli.Context) error {
				dbClient := database.New(awsCfg, ctx.String("dynamodb-table"))
				s3Client := s3.NewClient(awsCfg, "dicomizer")

				s, err := scheduler.Create(
					ctx.Context,
					ctx.String("crontab"),
					s3Client,
					dbClient,
					ctx.String("pacs"),
					ctx.String("aet"),
					ctx.String("aec"),
					ctx.String("aem"),
					ctx.String("center"),
				)
				if err != nil {
					return err
				}
				defer s.Shutdown()

				fmt.Printf("starting scheduler with \"%s\"...\n", ctx.String("crontab"))
				s.Start()

				fmt.Printf("web server started at http://%s\n", ctx.String("bind"))
				return http.ListenAndServe(
					ctx.String("bind"),
					web.RegisterHandlers(s3Client, dbClient, ctx.String("center")),
				)
			},
		},
		{
			Name:        "check-patient",
			Description: "Check a patient for new DICOM files",
			ArgsUsage:   "<pk>",
			Flags: append([]cli.Flag{
				&cli.StringFlag{
					Name:        "dynamodb-table",
					Category:    "AWS",
					Usage:       "DynamoDB table name",
					DefaultText: AWS_DYNAMODB_TABLE,
					Value:       AWS_DYNAMODB_TABLE,
					EnvVars:     []string{"DYNAMODB_TABLE"},
				},
				&cli.StringFlag{
					Name:     "center",
					Category: "AWS",
					Usage:    "LABORATORY",
					Required: true,
					EnvVars:  []string{"LABORATORY"},
				},
			}, dicomFlags...),
			Before: func(ctx *cli.Context) error {
				if ctx.Args().First() == "" {
					return fmt.Errorf("missing pk")
				}

				return checkMandatoryStringFlags(
					[]string{"pacs", "aet", "aem", "aec", "center"},
					ctx,
				)
			},
			Action: func(ctx *cli.Context) error {
				dbClient := database.New(awsCfg, ctx.String("dynamodb-table"))
				s3Client := s3.NewClient(awsCfg, "dicomizer")
				pk := ctx.Args().First()
				if strings.Contains(pk, "PATIENT#") {
					fmt.Println(
						"prefix \"PATIENT#\" found in pk. It is not necessary to include it.",
					)
					pk = strings.Replace(pk, "PATIENT#", "", 1)
				}

				fmt.Println("Fetching patient info...")
				pInfo, err := dbClient.GetPatientInfo(ctx.Context, pk)
				if err != nil {
					return err
				}
				if pInfo == nil {
					return fmt.Errorf("patient \"%s\" not found", pk)
				}

				return check.CheckPatientDCM(
					ctx.Context,
					s3Client,
					dbClient,
					ctx.String("pacs"),
					ctx.String("aet"),
					ctx.String("aec"),
					ctx.String("aem"),
					ctx.String("center"),
					*pInfo,
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
		fmt.Println("error while loading .env file")
		os.Exit(1)
	}

	for _, env := range mandatoryEnvVars {
		if os.Getenv(env) == "" {
			fmt.Printf("%s is not set", env)
			os.Exit(1)
		}
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func checkMandatoryStringFlags(flags []string, ctx *cli.Context) error {
	for _, f := range flags {
		if v := ctx.String(f); v == "" {
			return fmt.Errorf("%s is mandatory", f)
		}
	}

	return nil
}
