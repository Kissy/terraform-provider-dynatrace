package dynatrace

import (
	"github.com/Kissy/go-dynatrace/dynatrace"
	"github.com/Kissy/go-dynatrace/dynatrace/client"
	"github.com/Kissy/go-dynatrace/dynatrace/client/maintenance_windows"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMaintenanceWindow() *schema.Resource {
    return &schema.Resource{
        Create: resourceMaintenanceWindowCreate,
        Read: resourceMaintenanceWindowRead,
        Update: resourceMaintenanceWindowUpdate,
        Delete: resourceMaintenanceWindowDelete,
        
        Schema: map[string]*schema.Schema{
            "name": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "description": {
                Type: schema.TypeString,
                Optional: true,
                Default: "Managed by terraform",
                
            },
            "type": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "suppression": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "scope": {
                Type: schema.TypeList,
                MinItems: 1,
                MaxItems: 1,
                Optional: true,
                Elem: &schema.Resource{ 
                    Schema: map[string]*schema.Schema{
                        "entities": {
                            Type: schema.TypeList,
                            Required: true,
                            Elem: &schema.Schema{
                                    Type: schema.TypeString,
                                    },
                        },
                        "matches": {
                            Type: schema.TypeList,
                            Required: true,
                            Elem: &schema.Resource{ 
                                Schema: map[string]*schema.Schema{
                                    "type": {
                                        Type: schema.TypeString,
                                        Optional: true,
                                        
                                    },
                                    "management_zone_id": {
                                        Type: schema.TypeInt,
                                        Optional: true,
                                        
                                    },
                                    "tags": {
                                        Type: schema.TypeList,
                                        Required: true,
                                        Elem: &schema.Resource{ 
                                            Schema: map[string]*schema.Schema{
                                                "context": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "key": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "value": {
                                                    Type: schema.TypeString,
                                                    Optional: true,
                                                    
                                                },
                                            },
                                        },
                                    },
                                    "tag_combination": {
                                        Type: schema.TypeString,
                                        Optional: true,
                                        
                                    },
                                },
                            },
                        },
                    },
                },
            },
            "schedule": {
                Type: schema.TypeList,
                MinItems: 1,
                MaxItems: 1,
                Required: true,
                Elem: &schema.Resource{ 
                    Schema: map[string]*schema.Schema{
                        "recurrence_type": {
                            Type: schema.TypeString,
                            Required: true,
                            
                        },
                        "recurrence": {
                            Type: schema.TypeList,
                            MinItems: 1,
                            MaxItems: 1,
                            Optional: true,
                            Elem: &schema.Resource{ 
                                Schema: map[string]*schema.Schema{
                                    "day_of_week": {
                                        Type: schema.TypeString,
                                        Optional: true,
                                        
                                    },
                                    "day_of_month": {
                                        Type: schema.TypeInt,
                                        Optional: true,
                                        
                                    },
                                    "start_time": {
                                        Type: schema.TypeString,
                                        Required: true,
                                        
                                    },
                                    "duration_minutes": {
                                        Type: schema.TypeInt,
                                        Required: true,
                                        
                                    },
                                },
                            },
                        },
                        "start": {
                            Type: schema.TypeString,
                            Required: true,
                            
                        },
                        "end": {
                            Type: schema.TypeString,
                            Required: true,
                            
                        },
                        "zone_id": {
                            Type: schema.TypeString,
                            Required: true,
                            
                        },
                    },
                },
            },
        },
    }
}

func resourceMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)

    params := maintenance_windows.NewCreateMaintenanceWindowParams().WithBody(buildMaintenanceWindowStruct(d))
    window, err := apiClient.MaintenanceWindows.CreateMaintenanceWindow(params, nil)

    if err != nil {
        if v, ok := err.(*maintenance_windows.CreateMaintenanceWindowBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    d.SetId(*window.Payload.ID)

    return nil
}

func resourceMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := maintenance_windows.NewGetMaintenanceWindowParams().WithID(strfmt.UUID(d.Id()))
    maintenanceWindow, err := apiClient.MaintenanceWindows.GetMaintenanceWindow(params, nil)

    if err != nil {
        return handleNotFoundError(err, d)
    }
        _ = d.Set("name", maintenanceWindow.Payload.Name)
    if maintenanceWindow.Payload.Description == nil {
        _ = d.Set("description", "null")
    } else {
        _ = d.Set("description", maintenanceWindow.Payload.Description)
    }
        _ = d.Set("type", maintenanceWindow.Payload.Type)
        _ = d.Set("suppression", maintenanceWindow.Payload.Suppression)
    if maintenanceWindow.Payload.Scope != nil {
		if err := d.Set("scope", flattenScope(maintenanceWindow.Payload.Scope)); err != nil {
			return err
		}
	}
    if maintenanceWindow.Payload.Schedule != nil {
		if err := d.Set("schedule", flattenSchedule(maintenanceWindow.Payload.Schedule)); err != nil {
			return err
		}
	}

    return nil
}

func resourceMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := maintenance_windows.NewUpdateMaintenanceWindowParams().WithID(strfmt.UUID(d.Id()))
    params = params.WithBody(buildMaintenanceWindowStruct(d))
    _, _, err := apiClient.MaintenanceWindows.UpdateMaintenanceWindow(params, nil)
    if err != nil {
        if v, ok := err.(*maintenance_windows.UpdateMaintenanceWindowBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    return nil
}

func resourceMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := maintenance_windows.NewDeleteMaintenanceWindowParams().WithID(strfmt.UUID(d.Id()))
    _, err := apiClient.MaintenanceWindows.DeleteMaintenanceWindow(params, nil)
    if err != nil {
        return err
    }

    d.SetId("")

    return nil
}

func buildMaintenanceWindowStruct(d *schema.ResourceData) *dynatrace.MaintenanceWindow {
	maintenanceWindow := &dynatrace.MaintenanceWindow{
        Name: swag.String(d.Get("name").(string)),
        Type: swag.String(d.Get("type").(string)),
        Suppression: swag.String(d.Get("suppression").(string)),
        Schedule: expandSchedule(d.Get("schedule")),
	}
	
    if attr, ok := d.GetOk("description"); ok {
        maintenanceWindow.Description = swag.String(attr.(string))
    }
    if attr, ok := d.GetOk("scope"); ok {
        maintenanceWindow.Scope = expandScope(attr)
    }

    return maintenanceWindow
}

func expandTagInfo(tagInfoData interface{}) []*dynatrace.TagInfo {
    var tagInfos []*dynatrace.TagInfo

	for _, v := range tagInfoData.([]interface{}) {
	    d := v.(map[string]interface{})
		tagInfo := &dynatrace.TagInfo{
                        Context: swag.String(d["context"].(string)),
                        Key: swag.String(d["key"].(string)),
    	}
        if attr, ok := d["value"]; ok {
                tagInfo.Value = swag.String(attr.(string))
        }

		tagInfos = append(tagInfos, tagInfo)
	}

	return tagInfos
}

func expandMonitoredEntityFilter(monitoredEntityFilterData interface{}) []*dynatrace.MonitoredEntityFilter {
    var monitoredEntityFilters []*dynatrace.MonitoredEntityFilter

	for _, v := range monitoredEntityFilterData.([]interface{}) {
	    d := v.(map[string]interface{})
		monitoredEntityFilter := &dynatrace.MonitoredEntityFilter{
                    Tags: expandTagInfo(d["tags"]),
    	}
        if attr, ok := d["type"]; ok {
                monitoredEntityFilter.Type = attr.(string)
        }
        if attr, ok := d["management_zone_id"]; ok {
                monitoredEntityFilter.ManagementZoneID = int64(attr.(int))
        }
        if attr, ok := d["tag_combination"]; ok {
                monitoredEntityFilter.TagCombination = attr.(string)
        }

		monitoredEntityFilters = append(monitoredEntityFilters, monitoredEntityFilter)
	}

	return monitoredEntityFilters
}

func expandScope(scopeData interface{}) *dynatrace.Scope {
	d := scopeData.([]interface{})[0].(map[string]interface{})
    scope := &dynatrace.Scope{
                Matches: expandMonitoredEntityFilter(d["matches"]),
	}
    if attr, ok := d["entities"]; ok {
        var entities []string

        for _, v := range attr.([]interface{}) {
            entities = append(entities, v.(string))
        }

        scope.Entities = entities
    }

	return scope
}

func expandRecurrence(recurrenceData interface{}) *dynatrace.Recurrence {
	d := recurrenceData.([]interface{})[0].(map[string]interface{})
    recurrence := &dynatrace.Recurrence{
                    StartTime: swag.String(d["start_time"].(string)),
                DurationMinutes: swag.Int32(int32(d["duration_minutes"].(int))),
	}
    if attr, ok := d["day_of_week"]; ok {
            recurrence.DayOfWeek = attr.(string)
    }
    if attr, ok := d["day_of_month"]; ok {
            recurrence.DayOfMonth = int32(attr.(int))
    }

	return recurrence
}

func expandSchedule(scheduleData interface{}) *dynatrace.Schedule {
	d := scheduleData.([]interface{})[0].(map[string]interface{})
    schedule := &dynatrace.Schedule{
                    RecurrenceType: swag.String(d["recurrence_type"].(string)),
                    Start: swag.String(d["start"].(string)),
                    End: swag.String(d["end"].(string)),
                    ZoneID: swag.String(d["zone_id"].(string)),
	}
    if attr, ok := d["recurrence"]; ok {
            if len(attr.([]interface{})) > 0 {
                schedule.Recurrence = expandRecurrence(attr)
            }
    }

	return schedule
}


func flattenTagInfo(tagInfoDatas []*dynatrace.TagInfo) interface{} {
    var tagInfos []interface{}

	for _, tagInfoData := range tagInfoDatas {
		tagInfo := map[string]interface{}{
                    "context": tagInfoData.Context,
                    "key": tagInfoData.Key,
    	}
        if tagInfoData.Value != nil {
            tagInfo["value"] = tagInfoData.Value
        }

		tagInfos = append(tagInfos, tagInfo)
	}

	return tagInfos
}

func flattenMonitoredEntityFilter(monitoredEntityFilterDatas []*dynatrace.MonitoredEntityFilter) interface{} {
    var monitoredEntityFilters []interface{}

	for _, monitoredEntityFilterData := range monitoredEntityFilterDatas {
		monitoredEntityFilter := map[string]interface{}{
                    "tags": flattenTagInfo(monitoredEntityFilterData.Tags),
    	}
        if monitoredEntityFilterData.Type != "" {
            monitoredEntityFilter["type"] = monitoredEntityFilterData.Type
        }
        monitoredEntityFilter["management_zone_id"] = monitoredEntityFilterData.ManagementZoneID
        if monitoredEntityFilterData.TagCombination != "" {
            monitoredEntityFilter["tag_combination"] = monitoredEntityFilterData.TagCombination
        }

		monitoredEntityFilters = append(monitoredEntityFilters, monitoredEntityFilter)
	}

	return monitoredEntityFilters
}

func flattenScope(scopeData *dynatrace.Scope) interface{} {
    scope := map[string]interface{}{
                "matches": flattenMonitoredEntityFilter(scopeData.Matches),
	}
    if scopeData.Entities != nil {
        var entities []string

        for _, entitiesData := range scopeData.Entities {
            entities = append(entities, entitiesData)
        }

        scope["entities"] = entities
    }

    return []interface{}{ scope }
    
}

func flattenRecurrence(recurrenceData *dynatrace.Recurrence) interface{} {
    recurrence := map[string]interface{}{
                "start_time": recurrenceData.StartTime,
                "duration_minutes": recurrenceData.DurationMinutes,
	}
    if recurrenceData.DayOfWeek != "" {
        recurrence["day_of_week"] = recurrenceData.DayOfWeek
    }
    recurrence["day_of_month"] = recurrenceData.DayOfMonth

    return []interface{}{ recurrence }
    
}

func flattenSchedule(scheduleData *dynatrace.Schedule) interface{} {
    schedule := map[string]interface{}{
                "recurrence_type": scheduleData.RecurrenceType,
                "start": scheduleData.Start,
                "end": scheduleData.End,
                "zone_id": scheduleData.ZoneID,
	}
    if scheduleData.Recurrence != nil {
        schedule["recurrence"] = flattenRecurrence(scheduleData.Recurrence)
    }

    return []interface{}{ schedule }
    
}
