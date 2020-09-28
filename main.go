package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/jmespath/go-jmespath"
	"github.com/olivere/ndjson"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	Version   string
	GitCommit string
	GoVersion string
	BuildDate string
)

func main() {
	homeDir, _ := os.UserHomeDir()
	timestampFile := path.Join(homeDir, ".akamai-datastream-cli")

	app := &cli.App{
		Name:      "akamai-datastream-cli",
		Version:   "v1.0.0",
		Copyright: "(c) 2020 Coralogix Inc.",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Coralogix Inc.",
				Email: "info@coralogix.com",
			},
		},
		Usage:    "get Akamai DataStream flow",
		HelpName: "akamai-datastream-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Required: true,
				Usage:    "Akamai API endpoint `HOST`",
				EnvVars:  []string{"AKAMAI_HOST"},
			},
			&cli.StringFlag{
				Name:     "client-token",
				Required: true,
				Usage:    "Akamai API client token `CLIENT_TOKEN`",
				EnvVars:  []string{"AKAMAI_CLIENT_TOKEN"},
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Required: true,
				Usage:    "Akamai API client secret `CLIENT_SECRET`",
				EnvVars:  []string{"AKAMAI_CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:     "access-token",
				Required: true,
				Usage:    "Akamai API access token `ACCESS_TOKEN`",
				EnvVars:  []string{"AKAMAI_ACCESS_TOKEN"},
			},
			&cli.UintFlag{
				Name:     "stream-id",
				Required: true,
				Usage:    "Akamai DataStream ID `STREAM_ID`",
				EnvVars:  []string{"AKAMAI_STREAM_ID"},
			},
			&cli.StringFlag{
				Name:    "logs-type",
				Value:   "raw-logs",
				Usage:   "Akamai DataStream logs type `LOGS_TYPE`",
				EnvVars: []string{"AKAMAI_LOGS_TYPE"},
			},
			&cli.TimestampFlag{
				Name:        "start",
				Value:       cli.NewTimestamp(time.Now().UTC().Add(time.Duration(-15) * time.Minute)),
				Usage:       "Akamai DataStream start date `START_TIMESTAMP`",
				FilePath:    timestampFile,
				EnvVars:     []string{"START_TIMESTAMP"},
				Layout:      time.RFC3339,
				DefaultText: "15 minutes ago",
			},
			&cli.TimestampFlag{
				Name:        "end",
				Value:       cli.NewTimestamp(time.Now().UTC().Add(time.Duration(-1) * time.Minute)),
				Usage:       "Akamai DataStream end date `END_TIMESTAMP`",
				EnvVars:     []string{"END_TIMESTAMP"},
				Layout:      time.RFC3339,
				DefaultText: "current",
			},
			&cli.UintFlag{
				Name:        "max-records-limit",
				Value:       2000,
				Usage:       "maximal count of records in response `MAX_RECORDS_LIMIT`",
				EnvVars:     []string{"MAX_RECORDS_LIMIT"},
				DefaultText: "2000",
			},
			&cli.IntFlag{
				Name:        "max-body-size",
				Value:       128000,
				Usage:       "maximal body size of API response `MAX_BODY_SIZE`",
				EnvVars:     []string{"MAX_BODY_SIZE"},
				DefaultText: "128000",
			},
			&cli.StringFlag{
				Name:        "query",
				Aliases:     []string{"q"},
				Value:       "data",
				Usage:       "JMESPath compatible query `QUERY`",
				EnvVars:     []string{"QUERY"},
				DefaultText: "*",
			},
			&cli.BoolFlag{
				Name:    "flatten",
				Aliases: []string{"f"},
				Value:   true,
				Usage:   "flatten output",
				EnvVars: []string{"FLATTEN"},
			},
			&cli.BoolFlag{
				Name:    "keep-last-position",
				Aliases: []string{"p"},
				Value:   true,
				Usage:   "keep last queried timestamp",
				EnvVars: []string{"KEEP_LAST_POSITION"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "debug mode",
				EnvVars: []string{"DEBUG"},
			},
		},
		Action: func(ctx *cli.Context) error {
			var response_bytes []byte
			var response_json interface{}
			var query_result interface{}

			config := edgegrid.Config{
				Host:         ctx.String("host"),
				ClientToken:  ctx.String("client-token"),
				ClientSecret: ctx.String("client-secret"),
				AccessToken:  ctx.String("access-token"),
				MaxBody:      ctx.Int("max-body-size"),
				Debug:        ctx.Bool("debug"),
			}

			if ctx.String("logs-type") != "raw-logs" && ctx.String("logs-type") != "aggregate-logs" {
				return cli.Exit("Incorrect logs type! Allowed values: raw-logs, aggregate-logs", 2)
			}

			request, _ := client.NewRequest(config, "GET", fmt.Sprintf("/datastream-pull-api/v1/streams/%d/%s", ctx.Uint("stream-id"), ctx.String("logs-type")), nil)
			query := request.URL.Query()
			query.Add("start", ctx.Timestamp("start").Format(time.RFC3339))
			query.Add("end", ctx.Timestamp("end").Format(time.RFC3339))
			query.Add("size", strconv.FormatUint(uint64(ctx.Uint("max-records-limit")), 10))
			request.URL.RawQuery = query.Encode()
			response, err := client.Do(config, request)
			defer response.Body.Close()
			if err != nil {
				log.Println(err)
				return cli.Exit("Cannot execute request to Akamai API!", 2)
			}

			response_bytes, _ = ioutil.ReadAll(response.Body)
			if response.StatusCode == 204 {
				return cli.Exit("", 0)
			} else if response.StatusCode != 200 && ctx.Bool("debug") == true {
				log.Println("Status code: ", response.StatusCode)
				log.Println(string(response_bytes))
				return cli.Exit("Not success response code!", 2)
			}

			err = json.Unmarshal(response_bytes, &response_json)
			if err != nil {
				log.Println(err)
				return cli.Exit("Cannot parse response from Akamai API!", 2)
			}

			if ctx.String("query") != "" {
				query_result, err = jmespath.Search(ctx.String("query"), response_json)
				if err != nil {
					log.Println(err)
					return cli.Exit("Cannot execute query to Akamai API response!", 2)
				}
			} else {
				query_result = response_json
			}

			if ctx.Bool("flatten") == true {
				switch result := query_result.(type) {
				case []interface{}:
					writer := ndjson.NewWriter(os.Stdout)
					for _, record := range result {
						if err := writer.Encode(record); err != nil {
							log.Println(err)
							return cli.Exit("Cannot convert result to JSON string!", 2)
						}
					}
				default:
					json_result, err := json.Marshal(result)
					if err != nil {
						log.Println(err)
						return cli.Exit("Cannot convert query result to JSON string!", 2)
					}
					fmt.Println(string(json_result))
				}
			} else {
				json_result, err := json.Marshal(query_result)
				if err != nil {
					log.Println(err)
					return cli.Exit("Cannot convert query result to JSON string!", 2)
				}
				fmt.Println(string(json_result))
			}

			if ctx.Bool("keep-last-position") == true {
				track_file, err := os.Create(timestampFile)
				if err != nil {
					log.Println(err)
					return cli.Exit("Cannot save last check timestamp to track file!", 2)
				}
				track_file.WriteString(ctx.Timestamp("end").Format(time.RFC3339))
				track_file.Close()
			}

			return nil
		},
	}

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("%s: \"%s\", ", ctx.App.Name, ctx.App.Version)
		fmt.Printf("GitCommit: \"%s\", ", GitCommit)
		fmt.Printf("GoVersion: \"%s\", ", GoVersion)
		fmt.Printf("BuildDate: \"%s\"\n", BuildDate)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
