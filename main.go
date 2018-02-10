package main

import (
	//	"flag"
	//"fmt"
	"os"

	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/server"
	//"github.com/byxorna/flow/version"
	//"github.com/cobratbq/flagtag"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "main"})
)

func main() {

	/*
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Version: %s Commit: %s Branch: %s Date: %s\n", version.Version, version.Commit, version.Branch, version.BuildDate)
			flag.PrintDefaults()
		}
	*/

	servercfg, err := config.LoadServerConfigFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	etcdcfg, err := config.LoadEtcdConfigFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	//TODO make `-h` output work with the flagtags
	//use flagset maybe?

	s, err := server.New(servercfg, etcdcfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.ListenAndServe())
}
