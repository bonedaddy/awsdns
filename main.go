package main

import (
	"log"
	"os"

	"github.com/mitchellh/goamz/aws"
	route53 "github.com/segmentio/go-route53"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-dns-updater"
	app.Usage = "easily update route53 dns records"
	app.Action = func(c *cli.Context) error {
		auth, err := aws.EnvAuth()
		if err != nil {
			log.Println("failed to retrieve authentication information: ", err)
			return err
		}
		dns := route53.New(auth, aws.USWest2)
		switch c.String("operation") {
		case "update":
			_, err = dns.Zone(c.String("zone")).Upsert(c.String("record.type"), c.String("record.name"), c.String("record.value"))
			if err != nil {
				log.Println("failed to remove record: ", err)
				return err
			}
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "zone",
			Usage: "id of the zone to update",
		},
		cli.StringFlag{
			Name:  "operation",
			Value: "update",
			Usage: "type of operation to perform, values are: update, add, delete",
		},
		cli.StringFlag{
			Name:  "record.type",
			Value: "TXT",
			Usage: "the type of dns record, such as A, TXT, etc...",
		},
		cli.StringFlag{
			Name:  "record.name",
			Value: "foo.test.io",
			Usage: "the name of the record",
		},
		cli.StringFlag{
			Name:  "record.value",
			Value: "dnslink=/someipfshash",
			Usage: "the contents of the record",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
