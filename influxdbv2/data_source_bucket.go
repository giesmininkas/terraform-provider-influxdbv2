package influxdbv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func dataSourceBucket() *schema.Resource {
	return &schema.Resource{
		Description: "InfluxDB Bucket data source",
		ReadContext: dataSourceBucketRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description:   "Bucket id.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Description:   "Bucket name.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"org_id": {
				Description: "ID of organization in which to create a bucket.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the bucket.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"retention_rules": {
				Description: "Rules to expire or retain data. No rules means data never expires.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"every_seconds": {
							Description: "Duration in seconds for how long data will be kept in the database. 0 means infinite.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"shard_group_duration_seconds": {
							Description: "Shard duration measured in seconds.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"created_at": {
				Description: "Bucket creation date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Last bucket update date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Bucket type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	id, idOk := data.GetOk("id")
	name, nameOk := data.GetOk("name")

	if !idOk && !nameOk {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Must set either id or name",
			},
		}
	}

	var bucket *domain.Bucket
	var err error

	if idOk {
		bucket, err = bucketsClient.FindBucketByID(ctx, id.(string))
	} else if nameOk {
		bucket, err = bucketsClient.FindBucketByName(ctx, name.(string))
	}

	if err != nil {
		return diag.FromErr(err)
	}

	diags := setBucketData(data, bucket)
	data.Set("id", bucket.Id)
	data.SetId(*bucket.Id)

	return diags
}
