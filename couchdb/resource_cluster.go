package couchdb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rossmerr/couchdb_go/client/server"
	"github.com/rossmerr/couchdb_go/models"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCluster,
		UpdateContext: updateCluster,
		ReadContext:   readCluster,
		DeleteContext: deleteCluster,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "enable_single_node, enable_cluster, add_node, finish_cluster",
			},
			// "bind_address": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The IP address to which to bind the current node. The special value 0.0.0.0 may be specified to bind to all interfaces on the host. (enable_cluster and enable_single_node only)",
			// },
			// "username": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The username of the server-level administrator to create. (enable_cluster and enable_single_node only), or the remote server’s administrator username (add_node)",
			// },
			// "password": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The password for the server-level administrator to create. (enable_cluster and enable_single_node only), or the remote server’s administrator username (add_node)",
			// },
			// "port": {
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	Description: "The TCP port to which to bind this node (enable_cluster and enable_single_node only) or the TCP port to which to bind a remote node (add_node only).",
			// },
			// "node_count": {
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	Description: "The total number of nodes to be joined into the cluster, including this one. Used to determine the value of the cluster’s n, up to a maximum of 3. (enable_cluster only)",
			// },
			// "remote_node": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The IP address of the remote node to setup as part of this cluster’s list of nodes. (enable_cluster only)",
			// },
			// "remote_current_user": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The username of the server-level administrator authorized on the remote node. (enable_cluster only)",
			// },
			// "remote_current_password": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The password of the server-level administrator authorized on the remote node. (enable_cluster only)",
			// },
			// "host": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The remote node IP of the node to add to the cluster. (add_node only)",
			// },
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the current node or cluster state",
			},
		},
	}
}

func createCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))

	if dd != nil {
		return append(diags, *dd)
	}
	cluster := &models.Body{
		Action: d.Get("action").(string),
	}

	params := server.NewClusterSetupPostParams().WithBody(cluster)
	_, _, err := client.Server.ClusterSetupPost(params)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to setup cluster")
	}

	d.SetId("cluster")

	return readCluster(ctx, d, meta)
}

func updateCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	return readCluster(ctx, d, meta)
}

func readCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, dd := connectToCouchDB(ctx, meta.(*CouchDBConfiguration))
	if dd != nil {
		return append(diags, *dd)
	}

	params := server.NewClusterSetupGetParams()
	response, err := client.Server.ClusterSetupGet(params)

	if err != nil {
		return AppendDiagnostic(diags, err, "Unable to read cluster states")
	}

	if response != nil && response.Payload != nil {
		d.Set("state", response.Payload.State)
	}

	return diags
}

func deleteCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	d.SetId("")
	return diags
}
