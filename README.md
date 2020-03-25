# downtime

Downtime allows you to check if a domain is down or does not return HTTP code 200.

It uses the github.com/ardanlabs/conf/ package, which allows you to set the configuration using:
 - Environment variables
 - Command line flags
 - Any other source
In this case, I've implemented a subpackage that reads the configuration from AWS Parameter Store.

The application can read the list of domains from AWS S3 or from command line. The Builder design pattern allows to implement any other source, for example a file.

The appplication notifies through text message via Twilio API and it can also send the output to screen. Notification via email can be implemented as well.

Cross compile
GOARCH=amd64 GOOS=linux go build
