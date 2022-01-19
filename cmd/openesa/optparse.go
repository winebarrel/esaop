package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/winebarrel/openesa"
)

var version string

const (
	DefaultConfig = "openesa.toml"
)

func parseArgs() *openesa.Config {
	var config string
	flag.StringVar(&config, "config", "", "config file")
	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		printVersionAndEixt()
	}

	if config == "" {
		exePath, err := os.Executable()

		if err != nil {
			log.Fatal(err)
		}

		config = path.Join(filepath.Dir(exePath), DefaultConfig)
	}

	return loadConfig(config)
}

func loadConfig(path string) *openesa.Config {
	rawCfg, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	cfg := &openesa.Config{}
	err = toml.Unmarshal(rawCfg, cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = cfg.Validate()

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func printVersionAndEixt() {
	fmt.Fprintln(os.Stderr, version)
	os.Exit(0)
}
