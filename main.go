package main

import (
	"context"
	"flag"
	"log"

	"github.com/OpenNebula/terraform-provider-opennebula/opennebula"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// Provider version - set via ldflags during build
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Create context for mux server
	ctx := context.Background()

	// Create mux server that combines both SDKv2 and Framework providers
	// The first provider (SDKv2) handles all current resources
	// The second provider (Framework) will handle migrated resources
	providers := []func() tfprotov5.ProviderServer{
		opennebula.Provider().GRPCProvider,                      // SDKv2 provider
		providerserver.NewProtocol5(opennebula.New(version)()), // Framework provider
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	// Serve the muxed provider
	err = tf5server.Serve(
		"registry.terraform.io/OpenNebula/opennebula",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
