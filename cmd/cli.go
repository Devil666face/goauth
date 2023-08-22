package cmd

import (
	"flag"
)

const (
	START     int = 0
	MIGRATE   int = 1
	SUPERUSER int = 2
)

func SetCli() int {
	migrate := flag.Bool("migrate", false, "Set migrations in database")
	superuser := flag.Bool("superuser", false, "Create superuser set envs SUUSER=User and SUPASS=Password")
	flag.Parse()
	if *migrate {
		return MIGRATE
	}
	if *superuser {
		return SUPERUSER
	}
	return START
}
