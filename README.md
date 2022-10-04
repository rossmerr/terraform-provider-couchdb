Terraform Provider CouchDB
==========================

A Terraform provider for CouchDB.
* Support for the CouchDB API v3
* Lazy initialization of the CouchDB connection
* Retry logic for provisioning a database

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/rossmerr/terraform-provider-couchdb`

```sh
$ mkdir -p $GOPATH/src/github.com/rossmerr; cd $GOPATH/src/github.com/rossmerr
$ git clone git@github.com/rossmerr/terraform-provider-couchdb
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/rossmerr/terraform-provider-couchdb
$ go install
```
## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin)

`$HOME/.terraform.d/plugins`

After placing it into your plugins directory, run `terraform init` to initialize it.

Or run the `make install` and you can reference the provider locally using 

```
terraform {
  required_providers {
    couchdb = {
      source = "github.com/rossmerr/couchdb"
    }
  }
}
``` 

### Configuring the provider

```
provider "couchdb" {
    endpoint = "localhost:5984"
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
    views = <<EOF
    {
        "people" : {
            "map" : "function(doc) { if (doc.type == 'person') { emit(doc); } }",
            "reduce": ""           
        }
    }
    EOF
}
```
