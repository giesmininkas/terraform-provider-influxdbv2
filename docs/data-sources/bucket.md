---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "influxdbv2_bucket Data Source - terraform-provider-influxdbv2"
subcategory: ""
description: |-
  InfluxDB Bucket data source
---

# influxdbv2_bucket (Data Source)

InfluxDB Bucket data source

## Example Usage

```terraform
data "influxdbv2_bucket" "example_bucket" {
  id = "BUCKET_ID"
  // or
  name = "BUCKET_NAME"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) Bucket id.
- `name` (String) Bucket name.

### Read-Only

- `created_at` (String) Bucket creation date.
- `description` (String) Description of the bucket.
- `org_id` (String) ID of organization in which to create a bucket.
- `retention_rules` (Set of Object) Rules to expire or retain data. No rules means data never expires. (see [below for nested schema](#nestedatt--retention_rules))
- `type` (String) Bucket type.
- `updated_at` (String) Last bucket update date.

<a id="nestedatt--retention_rules"></a>
### Nested Schema for `retention_rules`

Read-Only:

- `every_seconds` (Number)
- `shard_group_duration_seconds` (Number)


