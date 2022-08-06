package influxdbv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func dataSourceAuthorization() *schema.Resource {
	return &schema.Resource{
		Description: "InfluxDB Bucket resource",
		ReadContext: dataSourceAuthorizationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Authorization ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"org_id": {
				Description: "ID of the organization that the authorization is scoped to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permissions": {
				Description: "List of permissions for an authorization. An authorization must have at least one permission.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Enum: 'read'|'write'.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"resource": {
							Description: "Resource info.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Type of resource.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"id": {
										Description: "If ID is set, that is a permission for a specific resource. If it is not set, it is a permission for all resources of that resource type.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"org_id": {
										Description: "If orgID is set, that is a permission for all resources owned by that org. If it is not set, it is a permission for all resources of that resource type.",
										Type:        schema.TypeString,
										Computed:    true,
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
				Computed:    true,
			},
			"active": {
				Description: "Status of the token. If inactive, requests using the token will be rejected.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"user_id": {
				Description: "ID of the user that created and owns the token.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"token": {
				Description: "Token used to authenticate API requests.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
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
	}
}

func dataSourceAuthorizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := *(meta.(*influxdb2.Client))
	authClient := client.AuthorizationsAPI()

	authorizations, err := authClient.GetAuthorizations(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	var authorization *domain.Authorization = nil
	id := data.Get("id").(string)
	for _, auth := range *authorizations {
		if *auth.Id == id {
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

	diags := setAuthorizationData(data, authorization)
	data.Set("id", authorization.Id)
	data.SetId(*authorization.Id)

	return diags
}
