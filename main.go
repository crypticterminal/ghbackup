package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"qvl.io/ghbackup/ghbackup"
)

const (
	// Printed for -help, -h or with wrong number of arguments
	usage = `Usage: %s [flags] directory

  directory  path to save the repositories to


  At least one of -account or -secret must be specified.

`
	accountUsage = `Github user or organization name to get repositories from.
	If not specified, all repositories the authenticated user has access to will be loaded.`
	secretUsage = `Authentication secret for Github API.
	Can use the users password or a personal access token (https://github.com/settings/tokens).
	Authentication increases rate limiting (https://developer.github.com/v3/#rate-limiting) and enables backup of private repositories.`
)

// Get command line arguments and start updating repositories
func main() {
	// Flags
	account := flag.String("account", "", accountUsage)
	secret := flag.String("secret", "", secretUsage)
	verbose := flag.Bool("verbose", false, "print progress information")

	// Parse args
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 || (*account == "" && *secret == "") {
		flag.Usage()
		os.Exit(1)
	}

	// Log updates
	updates := make(chan ghbackup.Update)
	go func() {
		for u := range updates {
			switch u.Type {
			case ghbackup.UErr:
				log.Println(u.Message)
			case ghbackup.UInfo:
				if *verbose {
					log.Println(u.Message)
				}
			}
		}
	}()

	err := ghbackup.Run(ghbackup.Config{
		Account: *account,
		Dir:     args[0],
		Secret:  *secret,
		Updates: updates,
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
