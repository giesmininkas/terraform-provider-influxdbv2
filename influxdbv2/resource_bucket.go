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
		UpdateContext: resourceBucketUpdate,
		DeleteContext: resourceBucketDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Bucket name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "ID of organization in which to create a bucket.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description of the bucket.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"retention_rules": {
				Description: "Rules to expire or retain data. No rules means data never expires.",
				Type:        schema.TypeSet,
				Optional:    true,
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

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBucketCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	bucket, diags := mapToBucket(data)

	if diags.HasError() {
		return diags
	}

	bucket, err := bucketsClient.CreateBucket(ctx, bucket)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(*bucket.Id)

	diags = append(diags, resourceBucketRead(ctx, data, meta)...)

	return diags
}

func resourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	bucket, err := bucketsClient.FindBucketByID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	_ = data.Set("org_id", *bucket.OrgID)
	_ = data.Set("name", bucket.Name)
	_ = data.Set("description", bucket.Description)
	_ = data.Set("created_at", bucket.CreatedAt.String())
	_ = data.Set("updated_at", bucket.UpdatedAt.String())
	_ = data.Set("type", bucket.Type)

	var retentionRules []map[string]interface{}
	for _, rule := range bucket.RetentionRules {
		mapped := map[string]interface{}{
			"every_seconds":                rule.EverySeconds,
			"shard_group_duration_seconds": *rule.ShardGroupDurationSeconds,
		}
		retentionRules = append(retentionRules, mapped)
	}

	_ = data.Set("retention_rules", retentionRules)

	return nil
}

func resourceBucketUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	bucketsClient := client.BucketsAPI()

	bucket, diags := mapToBucket(data)

	if diags.HasError() {
		return diags
	}

	if len(bucket.RetentionRules) == 0 {
		rule := domain.RetentionRule{
			EverySeconds:              0,
			ShardGroupDurationSeconds: nil,
		}
		bucket.RetentionRules = append(bucket.RetentionRules, rule)
	}

	bucketId := data.Id()
	bucket.Id = &bucketId

	bucket, err := bucketsClient.UpdateBucket(ctx, bucket)

	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, resourceBucketRead(ctx, data, meta)...)

	return diags
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

func mapToBucket(data *schema.ResourceData) (*domain.Bucket, diag.Diagnostics) {
	orgId := data.Get("org_id").(string)

	bucket := domain.Bucket{
		Name:  data.Get("name").(string),
		OrgID: &orgId,
	}

	description, ok := data.GetOk("description")
	if ok {
		tmp := description.(string)
		bucket.Description = &tmp
	}

	retentionRules := domain.RetentionRules{}
	for _, retentionRule := range data.Get("retention_rules").(*schema.Set).List() {
		mapped, err := mapToRetentionRule(retentionRule.(map[string]interface{}))
		if err.HasError() {
			return nil, err
		}
		retentionRules = append(retentionRules, mapped)
	}
	bucket.RetentionRules = retentionRules

	return &bucket, nil
}

func mapToRetentionRule(data map[string]interface{}) (domain.RetentionRule, diag.Diagnostics) {
	rule := domain.RetentionRule{
		EverySeconds: int64(data["every_seconds"].(int)),
		Type:         domain.RetentionRuleTypeExpire,
	}

	shardGroupDuration, ok := data["shard_group_duration_seconds"]
	if ok {
		tmp := int64(shardGroupDuration.(int))
		rule.ShardGroupDurationSeconds = &tmp
	}

	var diags diag.Diagnostics

	if rule.ShardGroupDurationSeconds != nil && *rule.ShardGroupDurationSeconds > rule.EverySeconds {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Shard Group duration longer than Retention Period.",
		})
	}

	return rule, diags
}
