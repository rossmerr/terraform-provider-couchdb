package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rossmerr/terraform-provider-couchdb/couchdb"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: couchdb.Provider})
}
