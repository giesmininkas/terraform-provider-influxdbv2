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
		ReadContext:   resourceBucketRead,
		DeleteContext: resourceBucketDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Bucket name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"org_id": {
				Description: "Organization ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBucketCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*influxdb2.Client)
	bucketsClient := (*client).BucketsAPI()

	orgId := data.Get("org_id").(string)
	name := data.Get("name").(string)

	bucket, err := bucketsClient.CreateBucketWithNameWithID(ctx, orgId, name)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(*bucket.Id)

	return nil
}

func resourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*influxdb2.Client)
	bucketsClient := (*client).BucketsAPI()

	bucket, err := bucketsClient.FindBucketByID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("org_id", *bucket.OrgID)
	data.Set("name", bucket.Name)

	return nil
}

func resourceBucketDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*influxdb2.Client)
	bucketsClient := (*client).BucketsAPI()

	err := bucketsClient.DeleteBucketWithID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
