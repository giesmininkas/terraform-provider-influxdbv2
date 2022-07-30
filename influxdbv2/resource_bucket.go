package influxdbv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Description:   "InfluxDB Bucket resource",
		CreateContext: resourceBucketCreate,

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

func resourceBucketCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(influxdb2.Client)
	bucketsClient := client.BucketsAPI()

	orgId := data.Get("orgId").(string)
	name := data.Get("name").(string)

	bucketsClient.CreateBucketWithNameWithID(ctx, orgId, name)

	return nil
}

func resourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(influxdb2.Client)
	bucketsClient := client.BucketsAPI()

	bucketsClient.FindBucketByID(ctx, data.Id())

	return nil
}
