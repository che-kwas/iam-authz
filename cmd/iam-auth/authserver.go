// iam-auth is responsible for serving the ladon authorization request.
package main

import (
	"math/rand"
	"time"

	"github.com/che-kwas/iam-kit/config"
	"github.com/spf13/pflag"

	"iam-auth/internal/authserver"
)

var (
	name = "iam-auth"
	cfg  = pflag.StringP("config", "c", "", "config file")
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
		panic(err)
	}

	authserver.NewServer(name).Run()
}
