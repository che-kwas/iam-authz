// iam-auth is responsible for serving the ladon authorization request.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/che-kwas/iam-kit/config"
	"github.com/spf13/pflag"
)

var (
	name = "iam-auth"
	cfg  = pflag.StringP("config", "c", "./iam-auth.yaml", "config file")
	help = pflag.BoolP("help", "h", false, "show help message")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// parse flag
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	if err := config.LoadConfig(*cfg, name); err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	fmt.Printf("Hello %s\n", name)
}
