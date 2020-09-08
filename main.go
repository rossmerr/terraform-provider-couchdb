package main

import (
	"github.com/RossMerr/terraform-provider-couchdb/couchdb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: couchdb.Provider})
}
