terraform {
  required_providers {
    influxdbv2 = {
      source = "matrasas.dev/dev/influxdbv2"
      version = "0.1"
    }
  }
}

locals {
#  org_id = "293ca298c3861fa3" # laptop
  org_id = "9a4706ee542a75c5" # laptop
}

provider "influxdbv2" {
  host = "http://localhost:8086"
#  token = "J0HGVO8RGAq-gNfcppdPBXkqxNkTiSR9k4Ph3ilZYuC4mhqneFSBiTcjl3VoQx8gHJp81FdGSZcvx-9A_QTTkg==" # desktop
  token = "IROsuGp99nYeiori-EW5eogsAUVN6YLeyXQZMoL7d5QoRG0KUy2ZSCLZ7eiNDn1bpDKrtlmiuYrLLyE9WshYDg==" # laptop
}

resource "influxdbv2_bucket" "test_bucket" {
  name = "aaaaaaa"
  org_id = local.org_id
  description = "jkbkjkf"
#  retention_rules {
#    every_seconds = 3600
##    shard_group_duration_seconds = 1800
#  }
}

resource "influxdbv2_bucket" "test2_bucket" {
  name = "test2"
  org_id = local.org_id
}

resource "influxdbv2_authorization" "test_auth" {
  org_id = local.org_id
  permissions {
    action = "read"
    resource {
      id = "0f170a84fd13372b"
      org_id = local.org_id
      type = "buckets"
    }
  }
}
