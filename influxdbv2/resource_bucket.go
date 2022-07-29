package influxdbv2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Description: "InfluxDB Bucket resource",

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Bucket name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "Organization ID",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
