terraform {
  required_providers {
    docker = {
      source = "terraform-providers/docker"
      version = "2.7.2"
    }
    couchdb = {
      source = "github.com/RossMerr/couchdb"
    }
  }
}

locals {
  couchdb = {
    endpoint = "localhost"
    port = 8901
    username = "admin"
    password = "password"
  }
}

provider "docker" {
}

resource "docker_image" "couchdb" {
  name         = "couchdb:3.1.0"
  keep_locally = false
}

resource "docker_container" "couchdb" {
  image = docker_image.couchdb.latest
  name  = "couchdb_terraform"
  ports {
    internal = 5984
    external = local.couchdb.port
  }
  healthcheck {
    test = ["CMD", "curl", "-f", "http://localhost:5984/"]
  }
  restart = "unless-stopped"

  env = ["COUCHDB_USER=${local.couchdb.username}", "COUCHDB_PASSWORD=${local.couchdb.password}"]
}


provider "couchdb" {
  endpoint = "${local.couchdb.endpoint}:${local.couchdb.port}"
  username = local.couchdb.username
  password = local.couchdb.password
}

resource "couchdb_database" "db1" {
  depends_on = [
    docker_container.couchdb,
  ]
  name = "example"
}


resource "couchdb_database_design_document" "test" {
  database = couchdb_database.db1.name
  name = "test"
  views = <<EOF
  {
    "people" : {
        "map": "function(doc) { emit(doc._id, doc); }"
    }
  }
EOF
}

resource "couchdb_database" "user" {
  depends_on = [
    docker_container.couchdb,
  ]
  name = "_users"
}


resource "couchdb_user" "jenny" {
  depends_on = [
    couchdb_database.user,
  ]
  name = "jenny"
  password = "secret"
}

resource "couchdb_document" "spaghetti" {
    database = couchdb_database.db1.name
	doc = <<EOF
  {
		"description": "An Italian-American dish that usually consists of spaghetti, tomato sauce and meatballs.",
		"ingredients": [
			"spaghetti",
			"tomato sauce",
			"meatballs"
		],
		"name": "Spaghetti with meatballs"
	}
EOF
}

resource "couchdb_bulk_documents" "recipes" {
  database = couchdb_database.db1.name
  docs = <<EOF
  [
    {
          "_id": "9391913b56c655881fa57d60830008ac",
          "description": "An Italian-American dish that usually consists of spaghetti, tomato sauce and meatballs.",
          "ingredients": [
              "spaghetti",
              "tomato sauce",
              "meatballs"
          ],
          "name": "Spaghetti"
      }
  ]
EOF
}
