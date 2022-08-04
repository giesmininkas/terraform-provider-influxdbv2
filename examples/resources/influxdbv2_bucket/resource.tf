locals {
  org_id = "example_org_id"
}

resource "influxdbv2_bucket" "example_bucket" {
  name        = "example_bucket_1"
  org_id      = local.org_id
  description = "example description"
  retention_rules {
    every_seconds                = 3600
    shard_group_duration_seconds = 1800
  }
}
