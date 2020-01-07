package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	dconf "github.com/roloum/downtime/internal/conf"
)

var appName = "downtime"

func main() {
	if err := run(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func run() error {

	//Load configuration
	var cfg struct {
		Input  string `conf:"default:inline"`
		Output struct {
			Screen bool `conf:"default:true"`
			Twilio bool `conf:"default:false"`
			Email  bool `conf:"default:false"`
		}
		Twilio struct {
			Sid   string `conf:"default:ID,noprint"`
			Token string `conf:"default:TOKEN,noprint"`
			From  string `conf:"default:+12247013610"`
			To    string `conf:"default:+14154079869"`
		}
		S3 struct {
			Bucket string `conf:"default:bucket"`
			Key    string `conf:"default:key"`
		}
	}

	log := log.New(os.Stdout, "Downtime: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	sources := []conf.Sourcer{}

	if os.Getenv("DOWNTIME_AWSPS") != "" {
		log.Println("main: Loading configuration from AWS Parameter Store")
		ps, err := dconf.NewPs(appName)
		if err != nil {
			return errors.Wrap(err, "loading parameters from AWS Parameter Store")
		}

		sources = append(sources, ps)
	}

	if err := conf.Parse(os.Args[1:], appName, &cfg, sources...); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage(appName, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	log.Println("main: Configuration loaded")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "Generating config for output")
	}
	log.Printf("main : Config : \n%v\n", out)

	return nil
}
