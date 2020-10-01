Terraform Provider CouchDB
==========================

A Terraform provider for CouchDB, can be use as a near drop in replacement of [nicolai86](https://github.com/nicolai86/terraform-provider-couchdb) provider.
* Support for the CouchDB API v3
* Lazy initialization of the CouchDB connection
* Retry logic for provisioning a database

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/RossMerr/terraform-provider-couchdb`

```sh
$ mkdir -p $GOPATH/src/github.com/RossMerr; cd $GOPATH/src/github.com/RossMerr
$ git clone git@github.com/RossMerr/terraform-provider-couchdb
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/RossMerr/terraform-provider-couchdb
$ go install
```
## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin)
After placing it into your plugins directory, run `terraform init` to initialize it.

Or run the `make install` and you can reference the provider locally using 

```
terraform {
  required_providers {
    couchdb = {
      source = "github.com/RossMerr/couchdb"
    }
  }
}
``` 

### Configuring the provider

```
provider "couchdb" {
    endpoint = "http://localhost:5984"
    name = "jenny"
    password = "secret" 
}
 
resource "couchdb_database" "db1" {
    name = "example"
}

resource "couchdb_database_replication" "db2db" {
    name = "example"
    source = couchdb_database.db1.name
    target = "example-clone"
    create_target = true
    continuous = true
}

resource "couchdb_database_design_document" "test" {
    database = couchdb_database.db1.name
    name = "types"
    view = <<EOF
    {
        "name" : "people",
        "map" : "function(doc) { if (doc.type == 'person') { emit(doc); } }",
        "reduce": ""
    }
    EOF
}
```