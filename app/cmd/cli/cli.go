package cli

import (
	"flag"
	"fmt"
)

var (
	ReturnStart     int = 0
	ReturnMigrate   int = 1
	ReturnSuperuser int = 2
)

func SetCli() int {
	migrate := flag.Bool("migrate", false, "Set migrations in database")
	superuser := flag.Bool("superuser", false, "Create superuser set envs SUUSER=User and SUPASS=Password")
	flag.Parse()
	if *migrate {
		fmt.Println("Start migrations")
		return ReturnMigrate
	}
	if *superuser {
		fmt.Println("Create superuser")
		return ReturnSuperuser
	}
	fmt.Println("Start server")
	return ReturnStart
}
