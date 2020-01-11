package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/roloum/downtime/cmd/downtime/internal/notifier"
	"github.com/roloum/downtime/cmd/downtime/internal/reader"
	dconf "github.com/roloum/downtime/internal/conf"
	"github.com/roloum/downtime/internal/url"
)

var appName = "downtime"

var wg sync.WaitGroup

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
			AwsRegion string `conf:"default:region"`
			Bucket    string `conf:"default:bucket"`
			Key       string `conf:"default:key"`
		}
		Domain bool `conf:"default:true"`
	}

	log := log.New(os.Stdout, "Downtime: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	sources := []conf.Sourcer{}

	if os.Getenv("DOWNTIME_AWSPS") != "" {
		log.Println("main: Loading configuration from AWS Parameter Store")
		ps, err := dconf.NewPs(os.Getenv("AWS_REGION"), appName)
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

	//Read domains, using corresponding input based on configuration
	var input reader.Reader
	//Read from S3 bucket
	if cfg.Input == reader.S3 {
		input = &reader.InputS3Bucket{AwsRegion: cfg.S3.AwsRegion,
			Bucket: cfg.S3.Bucket, Key: cfg.S3.Key}
		//Read from command line
	} else {
		input = &reader.InputInline{}
	}
	i := &reader.Input{}
	uris, err := i.Read(input, log)
	if err != nil {
		return errors.Wrap(err, "Reading domain list")
	}

	errs := checkDownTime(uris, cfg.Domain)
	//No errors, return and finish
	if len(errs) == 0 {
		return nil
	}

	//Create message body
	body := "Errors checking domains:\n"
	for err := range errs {
		body += fmt.Sprintf("%v\n", err)
	}

	//Deliver notification
	var notifiers = []notifier.Notifier{}
	if cfg.Output.Screen {
		notifiers = append(notifiers, &notifier.Screen{})
	}

	for _, n := range notifiers {
		n.Notify(body, log)
	}

	return nil
}

func checkDownTime(uris []string, domain bool) chan error {

	count := len(uris)
	var ch = make(chan error, count)

	wg.Add(count)

	for _, uri := range uris {
		go checkurl(uri, domain, ch)
	}

	wg.Wait()
	close(ch)

	return ch
}

func checkurl(uri string, domain bool, ch chan error) {

	defer wg.Done()

	if err := url.Check(uri, domain); err != nil {
		ch <- errors.Wrap(err, uri)
	}
}
