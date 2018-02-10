package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/server"
	"github.com/byxorna/flow/version"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "main"})
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Version: %s Commit: %s Branch: %s Date: %s\n", version.Version, version.Commit, version.Branch, version.BuildDate)
		flag.PrintDefaults()
	}

	var cfg config.Config
	arg.MustParse(&cfg)
	err := cfg.ValidateAndSetDefaults()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%v\n", cfg)
	s, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.ListenAndServe())
}
