package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/server"
	"github.com/byxorna/flow/version"
	"github.com/sirupsen/logrus"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Version: %s Commit: %s Branch: %s Date: %s\n", version.Version, version.Commit, version.Branch, version.BuildDate)
		PrintDefaults()
	}

	servercfg := config.LoadServerConfigFromArgs(os.Args[1:])
	etcdcfg := config.LoadEtcdConfigFromArgs(os.Args[1:])
	s := server.New(servercfg, etcdcfg)
	logrus.Fatal(s.ListenAndServe())
}
