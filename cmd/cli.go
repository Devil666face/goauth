package cmd

import (
	"flag"
	"fmt"
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
		fmt.Println("Start migrations")
		return MIGRATE
	}
	if *superuser {
		fmt.Println("Create superuser")
		return SUPERUSER
	}
	fmt.Println("Start server")
	return START
}
