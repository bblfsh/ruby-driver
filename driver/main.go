package main

import (
	_ "github.com/bblfsh/ruby-driver/driver/impl"
	"github.com/bblfsh/ruby-driver/driver/normalizer"

	"gopkg.in/bblfsh/sdk.v2/driver/server"
)

func main() {
	server.Run(normalizer.Transforms)
}
