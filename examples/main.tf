terraform {
  required_providers {
    couchdb = {
      source = "hashicorp.com/RossMerr/couchdb"
    }
  }
}

provider "couchdb" {
  endpoint = "http://localhost:5900"
  username = "admin"
  password = "secret"
}

resource "couchdb_database" "db1" {
  name = "example"
}


resource "couchdb_database_design_document" "test" {
  database = "example"
  name = "types"
  view {
    name = "people"
    map = "function(doc) { if (doc.type == 'person') { emit(doc); } }"
    reduce = ""
  }
}

resource "couchdb_database" "user" {
  name = "_users"
}


resource "couchdb_user" "jenny" {
  name = "jenny"
  password = "secret"
}
