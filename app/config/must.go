package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func must(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		s := fmt.Sprintf("you can set by example: -%s '...'", envName)
		f := flag.String(envName, "", s)
		flag.Parse()

		if f == nil || *f == "" {
			log.Fatalf("%s is undefined, use -help for examples",envName)
		}
		env = *f
	}

	return env
}
