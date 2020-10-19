package dynatrace

import (
	"github.com/Kissy/go-dynatrace/dynatrace"
	"github.com/Kissy/go-dynatrace/dynatrace/client"
	"github.com/Kissy/go-dynatrace/dynatrace/client/notifications"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNotificationConfig() *schema.Resource {
    return &schema.Resource{
        Create: resourceNotificationConfigCreate,
        Read: resourceNotificationConfigRead,
        Update: resourceNotificationConfigUpdate,
        Delete: resourceNotificationConfigDelete,
        
        Schema: map[string]*schema.Schema{
            "name": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "alerting_profile": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "active": {
                Type: schema.TypeBool,
                Required: true,
                
            },
            "type": {
                Type: schema.TypeString,
                Required: true,
                
            },
        },
    }
}

func resourceNotificationConfigCreate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)

    params := notification_configs.NewCreateNotificationConfigParams().WithBody(buildNotificationConfigStruct(d))
    window, err := apiClient.NotificationConfigs.CreateNotificationConfig(params, nil)

    if err != nil {
        if v, ok := err.(*notification_configs.CreateNotificationConfigBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    d.SetId(*window.Payload.ID)

    return nil
}

func resourceNotificationConfigRead(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := notifications.NewGetNotificationParams().WithID(strfmt.UUID(d.Id()))
    notificationConfig, err := apiClient.NotificationConfigs.GetNotificationConfig(params, nil)

    if err != nil {
        return handleNotFoundError(err, d)
    }
        _ = d.Set("name", notificationConfig.Payload.Name)
        _ = d.Set("alerting_profile", notificationConfig.Payload.AlertingProfile)
        _ = d.Set("active", notificationConfig.Payload.Active)
        _ = d.Set("type", notificationConfig.Payload.Type)

    return nil
}

func resourceNotificationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := notifications.NewUpdateNotificationParams().WithID(strfmt.UUID(d.Id()))
    params = params.WithBody(buildNotificationConfigStruct(d))
    _, _, err := apiClient.NotificationConfigs.UpdateNotificationConfig(params, nil)
    if err != nil {
        if v, ok := err.(*notifications.UpdateNotificationConfigBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    return nil
}

func resourceNotificationConfigDelete(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := notifications.NewDeleteNotificationParams().WithID(strfmt.UUID(d.Id()))
    _, err := apiClient.NotificationConfigs.DeleteNotificationConfig(params, nil)
    if err != nil {
        return err
    }

    d.SetId("")

    return nil
}

func buildNotificationConfigStruct(d *schema.ResourceData) *dynatrace.NotificationConfig {
	notificationConfig := dynatrace.NotificationConfig()
        //Name: swag.String(d.Get("name").(string)),
        //AlertingProfile: swag.String(d.Get("alerting_profile").(string)),
        //Active: swag.Bool(d.Get("active").(bool)),
        //Type: swag.String(d.Get("type").(string)),
	}
	

    return notificationConfig
}

