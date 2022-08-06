package influxdbv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func setBucketData(data *schema.ResourceData, bucket *domain.Bucket) diag.Diagnostics {
	data.Set("org_id", *bucket.OrgID)
	data.Set("name", bucket.Name)
	data.Set("description", bucket.Description)
	data.Set("created_at", bucket.CreatedAt.String())
	data.Set("updated_at", bucket.UpdatedAt.String())
	data.Set("type", bucket.Type)

	var retentionRules []map[string]interface{}
	for _, rule := range bucket.RetentionRules {
		mapped := map[string]interface{}{
			"every_seconds":                rule.EverySeconds,
			"shard_group_duration_seconds": rule.ShardGroupDurationSeconds,
		}
		retentionRules = append(retentionRules, mapped)
	}

	data.Set("retention_rules", retentionRules)

	return nil
}
