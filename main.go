// Copyright (c) HashiCorp, Inc.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/prochac/terraform-provider-casdoor/internal/provider"
)

// Generate documentation.
//go:generate go tool -modfile=go.tool.mod tfplugindocs generate --provider-dir . -provider-name casdoor

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/prochac/casdoor",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
