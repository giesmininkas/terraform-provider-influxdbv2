package influxdbv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func resourceAuthorization() *schema.Resource {
	return &schema.Resource{
		Description:   "InfluxDB Bucket resource",
		CreateContext: resourceAuthorizationCreate,
		ReadContext:   resourceAuthorizationRead,
		UpdateContext: resourceAuthorizationUpdate,
		DeleteContext: resourceAuthorizationDelete,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Description: "ID of the organization that the authorization is scoped to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"permissions": {
				Description: "List of permissions for an authorization. An authorization must have at least one permission.",
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Enum: 'read'|'write'.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"resource": {
							Description: "Resource info.",
							Type:        schema.TypeSet,
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Type of resource.",
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
									},
									"id": {
										Description: "If ID is set that is a permission for a specific resource. if it is not set it is a permission for all resources of that resource type.",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
									"org_id": {
										Description: "If orgID is set that is a permission for all resources owned my that org. if it is not set it is a permission for all resources of that resource type.",
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Description: "A description of the token.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "Status of the token. If inactive, requests using the token will be rejected.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"user_id": {
				Description: "ID of the user that created and owns the token.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "Token used to authenticate API requests.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Authorization creation date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Last authorization update date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAuthorizationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	authClient := client.AuthorizationsAPI()

	orgId := data.Get("org_id").(string)
	userId, userOk := data.GetOk("user_id")
	description, descriptionOk := data.GetOk("description_id")
	active := data.Get("active").(bool)

	authorization := &domain.Authorization{
		OrgID:       &orgId,
		Permissions: &[]domain.Permission{},
	}

	if active {
		tmp := domain.AuthorizationUpdateRequestStatusActive
		authorization.AuthorizationUpdateRequest.Status = &tmp
	} else {
		tmp := domain.AuthorizationUpdateRequestStatusInactive
		authorization.AuthorizationUpdateRequest.Status = &tmp
	}

	if userOk {
		tmp := userId.(string)
		authorization.UserID = &tmp
	}

	if descriptionOk {
		tmp := description.(string)
		authorization.AuthorizationUpdateRequest.Description = &tmp
	}

	permissions := []domain.Permission{}
	for _, permissionData := range data.Get("permissions").(*schema.Set).List() {
		permissionDataMap := permissionData.(map[string]interface{})
		resourceDataMap := permissionDataMap["resource"].(*schema.Set).List()[0].(map[string]interface{})

		resourceId, resourceIdOk := resourceDataMap["id"]
		resourceOrgId, resourceOrgIdOk := resourceDataMap["org_id"]

		permission := domain.Permission{
			Action: domain.PermissionAction(permissionDataMap["action"].(string)),
			Resource: domain.Resource{
				Type: domain.ResourceType(resourceDataMap["type"].(string)),
			},
		}

		if resourceIdOk {
			tmp := resourceId.(string)
			permission.Resource.Id = &tmp
		}

		if resourceOrgIdOk {
			tmp := resourceOrgId.(string)
			permission.Resource.OrgID = &tmp
		}

		permissions = append(permissions, permission)
	}

	authorization.Permissions = &permissions

	authorization, err := authClient.CreateAuthorization(ctx, authorization)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(*authorization.Id)

	diags := resourceAuthorizationRead(ctx, data, meta)

	return diags
}

func resourceAuthorizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	authClient := client.AuthorizationsAPI()

	authorizations, err := authClient.GetAuthorizations(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	var authorization *domain.Authorization = nil
	for _, auth := range *authorizations {
		if *auth.Id == data.Id() {
			authorization = &auth
			break
		}
	}

	if authorization == nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No InfluxDB authorization with ID " + data.Id(),
			},
		}
	}

	data.Set("org_id", *authorization.OrgID)
	data.Set("description", authorization.Description)
	data.Set("user_id", authorization.UserID)
	data.Set("token", authorization.Token)
	data.Set("created_at", authorization.CreatedAt.String())
	data.Set("updated_at", authorization.UpdatedAt.String())

	switch *authorization.Status {
	case "active":
		data.Set("active", true)
		break
	case "inactive":
		data.Set("active", false)
		break
	}

	var permissions []map[string]interface{}
	//permissions := schema.Set{}
	for _, permission := range *authorization.Permissions {
		tmp := map[string]interface{}{
			"action": permission.Action,
			"resource": []map[string]interface{}{
				{
					"id":     permission.Resource.Id,
					"org_id": permission.Resource.OrgID,
					"type":   permission.Resource.Type,
				},
			},
		}

		permissions = append(permissions, tmp)
	}

	data.Set("permissions", permissions)

	//bucketsClient := client.BucketsAPI()
	//
	//bucket, err := bucketsClient.FindBucketByID(ctx, data.Id())
	//
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//_ = data.Set("org_id", *bucket.OrgID)
	//_ = data.Set("name", bucket.Name)
	//_ = data.Set("description", bucket.Description)
	//_ = data.Set("created_at", bucket.CreatedAt.String())
	//_ = data.Set("updated_at", bucket.UpdatedAt.String())
	//_ = data.Set("type", bucket.Type)
	//
	//var retentionRules []map[string]interface{}
	//for _, rule := range bucket.RetentionRules {
	//	mapped := map[string]interface{}{
	//		"every_seconds":                rule.EverySeconds,
	//		"shard_group_duration_seconds": *rule.ShardGroupDurationSeconds,
	//	}
	//	retentionRules = append(retentionRules, mapped)
	//}
	//
	//_ = data.Set("retention_rules", retentionRules)
	//
	return nil
}

func resourceAuthorizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := *(meta.(*influxdb2.Client))
	//bucketsClient := client.BucketsAPI()
	//
	//bucket, diags := mapToBucket(data)
	//
	//if diags.HasError() {
	//	return diags
	//}
	//
	//if len(bucket.RetentionRules) == 0 {
	//	rule := domain.RetentionRule{
	//		EverySeconds:              0,
	//		ShardGroupDurationSeconds: nil,
	//	}
	//	bucket.RetentionRules = append(bucket.RetentionRules, rule)
	//}
	//
	//bucketId := data.Id()
	//bucket.Id = &bucketId
	//
	//bucket, err := bucketsClient.UpdateBucket(ctx, bucket)
	//
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//

	diags := resourceAuthorizationRead(ctx, data, meta)

	return diags
}

func resourceAuthorizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	authClient := client.AuthorizationsAPI()

	err := authClient.DeleteAuthorizationWithID(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

//func mapToBucket(data *schema.ResourceData) (*domain.Bucket, diag.Diagnostics) {
//	orgId := data.Get("org_id").(string)
//
//	bucket := domain.Bucket{
//		Name:  data.Get("name").(string),
//		OrgID: &orgId,
//	}
//
//	description, ok := data.GetOk("description")
//	if ok {
//		bucket.Description = description.(*string)
//	}
//
//	retentionRules := domain.RetentionRules{}
//	for _, retentionRule := range data.Get("retention_rules").(*schema.Set).List() {
//		mapped, err := mapToRetentionRule(retentionRule.(map[string]interface{}))
//		if err.HasError() {
//			return nil, err
//		}
//		retentionRules = append(retentionRules, mapped)
//	}
//	bucket.RetentionRules = retentionRules
//
//	return &bucket, nil
//}

//func mapToRetentionRule(data map[string]interface{}) (domain.RetentionRule, diag.Diagnostics) {
//	rule := domain.RetentionRule{
//		EverySeconds: int64(data["every_seconds"].(int)),
//		Type:         domain.RetentionRuleTypeExpire,
//	}
//
//	shardGroupDuration, ok := data["shard_group_duration_seconds"]
//	if ok {
//		tmp := int64(shardGroupDuration.(int))
//		rule.ShardGroupDurationSeconds = &tmp
//	}
//
//	var diags diag.Diagnostics
//
//	if rule.ShardGroupDurationSeconds != nil && *rule.ShardGroupDurationSeconds > rule.EverySeconds {
//		diags = append(diags, diag.Diagnostic{
//			Severity: diag.Error,
//			Summary:  "Shard Group duration longer than Retention Period.",
//		})
//	}
//
//	return rule, diags
//}
