package influxdbv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func setAuthorizationData(data *schema.ResourceData, authorization *domain.Authorization) diag.Diagnostics {
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

	return nil
}
