package main

import (
	"gopkg.in/gcfg.v1"
	"log"
)

type Settings struct {
	Webservice struct {
		Bind string
	}

	Proxy struct {
		Server string
	}

	OAuth2 struct {
		TokenEndpoint string
		ClientId      string
		ClientSecret  string
	}
}

func loadSettings(path string) Settings {
	var s Settings
	err := gcfg.ReadFileInto(&s, path)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
