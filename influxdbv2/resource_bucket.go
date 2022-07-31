package influxdbv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
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
			"retention_seconds": {
				Description: "Time to keep the data in seconds. 0 == infinite",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"created_at": {
				Description: "Bucket creation date",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Last bucket update date",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Bucket type",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceBucketCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	orgId := data.Get("org_id").(string)
	name := data.Get("name").(string)
	retentionSeconds, ok := data.GetOk("retention_seconds")

	retentionRules := domain.RetentionRule{
		EverySeconds: 0,
	}

	if ok {
		retentionRules.EverySeconds = retentionSeconds.(int64)
	}

	bucket, err := bucketsClient.CreateBucketWithNameWithID(ctx, orgId, name, retentionRules)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(*bucket.Id)

	return nil
}

func resourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	bucket, err := bucketsClient.FindBucketByID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("org_id", *bucket.OrgID)
	data.Set("name", bucket.Name)
	data.Set("retention_seconds", bucket.RetentionRules[0].EverySeconds)
	data.Set("created_at", bucket.CreatedAt.String())
	data.Set("updated_at", bucket.UpdatedAt.String())
	data.Set("type", bucket.Type)

	return nil
}

func resourceBucketDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	err := bucketsClient.DeleteBucketWithID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
