package dynatrace

import (
	"github.com/Kissy/go-dynatrace/dynatrace"
	"github.com/Kissy/go-dynatrace/dynatrace/client"
	"github.com/Kissy/go-dynatrace/dynatrace/client/alerting_profiles"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlertingProfile() *schema.Resource {
    return &schema.Resource{
        Create: resourceAlertingProfileCreate,
        Read: resourceAlertingProfileRead,
        Update: resourceAlertingProfileUpdate,
        Delete: resourceAlertingProfileDelete,
        
        Schema: map[string]*schema.Schema{
            "display_name": {
                Type: schema.TypeString,
                Required: true,
                
            },
            "rules": {
                Type: schema.TypeList,
                MinItems: 0,
                MaxItems: 20,
                Optional: true,
                Elem: &schema.Resource{ 
                    Schema: map[string]*schema.Schema{
                        "severity_level": {
                            Type: schema.TypeString,
                            Required: true,
                            
                        },
                        "tag_filter": {
                            Type: schema.TypeList,
                            MinItems: 1,
                            MaxItems: 1,
                            Required: true,
                            Elem: &schema.Resource{ 
                                Schema: map[string]*schema.Schema{
                                    "include_mode": {
                                        Type: schema.TypeString,
                                        Required: true,
                                        
                                    },
                                    "tag_filters": {
                                        Type: schema.TypeList,
                                        Optional: true,
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
                                },
                            },
                        },
                        "delay_in_minutes": {
                            Type: schema.TypeInt,
                            Required: true,
                            
                        },
                    },
                },
            },
            "management_zone_id": {
                Type: schema.TypeInt,
                Optional: true,
                
            },
            "event_type_filters": {
                Type: schema.TypeList,
                MinItems: 0,
                MaxItems: 20,
                Optional: true,
                Elem: &schema.Resource{ 
                    Schema: map[string]*schema.Schema{
                        "predefined_event_filter": {
                            Type: schema.TypeList,
                            MinItems: 1,
                            MaxItems: 1,
                            Optional: true,
                            Elem: &schema.Resource{ 
                                Schema: map[string]*schema.Schema{
                                    "event_type": {
                                        Type: schema.TypeString,
                                        Required: true,
                                        
                                    },
                                    "negate": {
                                        Type: schema.TypeBool,
                                        Required: true,
                                        
                                    },
                                },
                            },
                        },
                        "custom_event_filter": {
                            Type: schema.TypeList,
                            MinItems: 1,
                            MaxItems: 1,
                            Optional: true,
                            Elem: &schema.Resource{ 
                                Schema: map[string]*schema.Schema{
                                    "custom_title_filter": {
                                        Type: schema.TypeList,
                                        MinItems: 1,
                                        MaxItems: 1,
                                        Optional: true,
                                        Elem: &schema.Resource{ 
                                            Schema: map[string]*schema.Schema{
                                                "enabled": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                                "value": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "operator": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "negate": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                                "case_insensitive": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                            },
                                        },
                                    },
                                    "custom_description_filter": {
                                        Type: schema.TypeList,
                                        MinItems: 1,
                                        MaxItems: 1,
                                        Optional: true,
                                        Elem: &schema.Resource{ 
                                            Schema: map[string]*schema.Schema{
                                                "enabled": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                                "value": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "operator": {
                                                    Type: schema.TypeString,
                                                    Required: true,
                                                    
                                                },
                                                "negate": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                                "case_insensitive": {
                                                    Type: schema.TypeBool,
                                                    Required: true,
                                                    
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}

func resourceAlertingProfileCreate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)

    params := alerting_profiles.NewCreateAlertingProfileParams().WithBody(buildAlertingProfileStruct(d))
    window, err := apiClient.AlertingProfiles.CreateAlertingProfile(params, nil)

    if err != nil {
        if v, ok := err.(*alerting_profiles.CreateAlertingProfileBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    d.SetId(*window.Payload.ID)

    return nil
}

func resourceAlertingProfileRead(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := alerting_profiles.NewGetAlertingProfileParams().WithID(strfmt.UUID(d.Id()))
    alertingProfile, err := apiClient.AlertingProfiles.GetAlertingProfile(params, nil)

    if err != nil {
        return handleNotFoundError(err, d)
    }
        _ = d.Set("display_name", alertingProfile.Payload.DisplayName)
        _ = d.Set("management_zone_id", alertingProfile.Payload.ManagementZoneID)

    return nil
}

func resourceAlertingProfileUpdate(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := alerting_profiles.NewUpdateAlertingProfileParams().WithID(strfmt.UUID(d.Id()))
    params = params.WithBody(buildAlertingProfileStruct(d))
    _, _, err := apiClient.AlertingProfiles.UpdateAlertingProfile(params, nil)
    if err != nil {
        if v, ok := err.(*alerting_profiles.UpdateAlertingProfileBadRequest); ok {
            return handleBadRequestPayload(v.Payload)
        }
        return err
    }

    return nil
}

func resourceAlertingProfileDelete(d *schema.ResourceData, meta interface{}) error {
    apiClient := meta.(*client.GoDynatrace)
    params := alerting_profiles.NewDeleteAlertingProfileParams().WithID(strfmt.UUID(d.Id()))
    _, err := apiClient.AlertingProfiles.DeleteAlertingProfile(params, nil)
    if err != nil {
        return err
    }

    d.SetId("")

    return nil
}

func buildAlertingProfileStruct(d *schema.ResourceData) *dynatrace.AlertingProfile {
	alertingProfile := &dynatrace.AlertingProfile{
        DisplayName: swag.String(d.Get("display_name").(string)),
	}
	
    if attr, ok := d.GetOk("rules"); ok {
        alertingProfile.Rules = expandAlertingProfileSeverityRule(attr)
    }
    if attr, ok := d.GetOk("management_zone_id"); ok {
        alertingProfile.ManagementZoneID = int64(attr.(int))
    }
    if attr, ok := d.GetOk("event_type_filters"); ok {
        alertingProfile.EventTypeFilters = expandAlertingEventTypeFilter(attr)
    }

    return alertingProfile
}

func expandTagFilter(tagFilterData interface{}) []*dynatrace.TagFilter {
    var tagFilters []*dynatrace.TagFilter

	for _, v := range tagFilterData.([]interface{}) {
	    d := v.(map[string]interface{})
		tagFilter := &dynatrace.TagFilter{
                        Context: swag.String(d["context"].(string)),
                        Key: swag.String(d["key"].(string)),
    	}
        if attr, ok := d["value"]; ok {
                tagFilter.Value = attr.(string)
        }

		tagFilters = append(tagFilters, tagFilter)
	}

	return tagFilters
}

func expandAlertingProfileTagFilter(alertingProfileTagFilterData interface{}) *dynatrace.AlertingProfileTagFilter {
	d := alertingProfileTagFilterData.([]interface{})[0].(map[string]interface{})
    alertingProfileTagFilter := &dynatrace.AlertingProfileTagFilter{
                    IncludeMode: swag.String(d["include_mode"].(string)),
	}
    if attr, ok := d["tag_filters"]; ok {
            if len(attr.([]interface{})) > 0 {
                alertingProfileTagFilter.TagFilters = expandTagFilter(attr)
            }
    }

	return alertingProfileTagFilter
}

func expandAlertingProfileSeverityRule(alertingProfileSeverityRuleData interface{}) []*dynatrace.AlertingProfileSeverityRule {
    var alertingProfileSeverityRules []*dynatrace.AlertingProfileSeverityRule

	for _, v := range alertingProfileSeverityRuleData.([]interface{}) {
	    d := v.(map[string]interface{})
		alertingProfileSeverityRule := &dynatrace.AlertingProfileSeverityRule{
                        SeverityLevel: swag.String(d["severity_level"].(string)),
                    TagFilter: expandAlertingProfileTagFilter(d["tag_filter"]),
                    DelayInMinutes: swag.Int32(int32(d["delay_in_minutes"].(int))),
    	}

		alertingProfileSeverityRules = append(alertingProfileSeverityRules, alertingProfileSeverityRule)
	}

	return alertingProfileSeverityRules
}

func expandAlertingPredefinedEventFilter(alertingPredefinedEventFilterData interface{}) *dynatrace.AlertingPredefinedEventFilter {
	d := alertingPredefinedEventFilterData.([]interface{})[0].(map[string]interface{})
    alertingPredefinedEventFilter := &dynatrace.AlertingPredefinedEventFilter{
                    EventType: swag.String(d["event_type"].(string)),
                Negate: swag.Bool(d["negate"].(bool)),
	}

	return alertingPredefinedEventFilter
}

func expandAlertingCustomTextFilter(alertingCustomTextFilterData interface{}) *dynatrace.AlertingCustomTextFilter {
	d := alertingCustomTextFilterData.([]interface{})[0].(map[string]interface{})
    alertingCustomTextFilter := &dynatrace.AlertingCustomTextFilter{
                Enabled: swag.Bool(d["enabled"].(bool)),
                    Value: swag.String(d["value"].(string)),
                    Operator: swag.String(d["operator"].(string)),
                Negate: swag.Bool(d["negate"].(bool)),
                CaseInsensitive: swag.Bool(d["case_insensitive"].(bool)),
	}

	return alertingCustomTextFilter
}

func expandAlertingCustomEventFilter(alertingCustomEventFilterData interface{}) *dynatrace.AlertingCustomEventFilter {
	d := alertingCustomEventFilterData.([]interface{})[0].(map[string]interface{})
    alertingCustomEventFilter := &dynatrace.AlertingCustomEventFilter{
	}
    if attr, ok := d["custom_title_filter"]; ok {
            if len(attr.([]interface{})) > 0 {
                alertingCustomEventFilter.CustomTitleFilter = expandAlertingCustomTextFilter(attr)
            }
    }
    if attr, ok := d["custom_description_filter"]; ok {
            if len(attr.([]interface{})) > 0 {
                alertingCustomEventFilter.CustomDescriptionFilter = expandAlertingCustomTextFilter(attr)
            }
    }

	return alertingCustomEventFilter
}

func expandAlertingEventTypeFilter(alertingEventTypeFilterData interface{}) []*dynatrace.AlertingEventTypeFilter {
    var alertingEventTypeFilters []*dynatrace.AlertingEventTypeFilter

	for _, v := range alertingEventTypeFilterData.([]interface{}) {
	    d := v.(map[string]interface{})
		alertingEventTypeFilter := &dynatrace.AlertingEventTypeFilter{
    	}
        if attr, ok := d["predefined_event_filter"]; ok {
                if len(attr.([]interface{})) > 0 {
                    alertingEventTypeFilter.PredefinedEventFilter = expandAlertingPredefinedEventFilter(attr)
                }
        }
        if attr, ok := d["custom_event_filter"]; ok {
                if len(attr.([]interface{})) > 0 {
                    alertingEventTypeFilter.CustomEventFilter = expandAlertingCustomEventFilter(attr)
                }
        }

		alertingEventTypeFilters = append(alertingEventTypeFilters, alertingEventTypeFilter)
	}

	return alertingEventTypeFilters
}


func flattenTagFilter(tagFilterDatas []*dynatrace.TagFilter) interface{} {
    var tagFilters []interface{}

	for _, tagFilterData := range tagFilterDatas {
		tagFilter := map[string]interface{}{
                    "context": tagFilterData.Context,
                    "key": tagFilterData.Key,
    	}
        if tagFilterData.Value != "" {
            tagFilter["value"] = tagFilterData.Value
        }

		tagFilters = append(tagFilters, tagFilter)
	}

	return tagFilters
}

func flattenAlertingProfileTagFilter(alertingProfileTagFilterData *dynatrace.AlertingProfileTagFilter) interface{} {
    alertingProfileTagFilter := map[string]interface{}{
                "include_mode": alertingProfileTagFilterData.IncludeMode,
	}
    if alertingProfileTagFilterData.TagFilters != nil {
        alertingProfileTagFilter["tag_filters"] = flattenTagFilter(alertingProfileTagFilterData.TagFilters)
    }

    return []interface{}{ alertingProfileTagFilter }
    
}

func flattenAlertingProfileSeverityRule(alertingProfileSeverityRuleDatas []*dynatrace.AlertingProfileSeverityRule) interface{} {
    var alertingProfileSeverityRules []interface{}

	for _, alertingProfileSeverityRuleData := range alertingProfileSeverityRuleDatas {
		alertingProfileSeverityRule := map[string]interface{}{
                    "severity_level": alertingProfileSeverityRuleData.SeverityLevel,
                    "tag_filter": flattenAlertingProfileTagFilter(alertingProfileSeverityRuleData.TagFilter),
                    "delay_in_minutes": alertingProfileSeverityRuleData.DelayInMinutes,
    	}

		alertingProfileSeverityRules = append(alertingProfileSeverityRules, alertingProfileSeverityRule)
	}

	return alertingProfileSeverityRules
}

func flattenAlertingPredefinedEventFilter(alertingPredefinedEventFilterData *dynatrace.AlertingPredefinedEventFilter) interface{} {
    alertingPredefinedEventFilter := map[string]interface{}{
                "event_type": alertingPredefinedEventFilterData.EventType,
                "negate": alertingPredefinedEventFilterData.Negate,
	}

    return []interface{}{ alertingPredefinedEventFilter }
    
}

func flattenAlertingCustomTextFilter(alertingCustomTextFilterData *dynatrace.AlertingCustomTextFilter) interface{} {
    alertingCustomTextFilter := map[string]interface{}{
                "enabled": alertingCustomTextFilterData.Enabled,
                "value": alertingCustomTextFilterData.Value,
                "operator": alertingCustomTextFilterData.Operator,
                "negate": alertingCustomTextFilterData.Negate,
                "case_insensitive": alertingCustomTextFilterData.CaseInsensitive,
	}

    return []interface{}{ alertingCustomTextFilter }
    
}

func flattenAlertingCustomEventFilter(alertingCustomEventFilterData *dynatrace.AlertingCustomEventFilter) interface{} {
    alertingCustomEventFilter := map[string]interface{}{
	}
    if alertingCustomEventFilterData.CustomTitleFilter != nil {
        alertingCustomEventFilter["custom_title_filter"] = flattenAlertingCustomTextFilter(alertingCustomEventFilterData.CustomTitleFilter)
    }
    if alertingCustomEventFilterData.CustomDescriptionFilter != nil {
        alertingCustomEventFilter["custom_description_filter"] = flattenAlertingCustomTextFilter(alertingCustomEventFilterData.CustomDescriptionFilter)
    }

    return []interface{}{ alertingCustomEventFilter }
    
}

func flattenAlertingEventTypeFilter(alertingEventTypeFilterDatas []*dynatrace.AlertingEventTypeFilter) interface{} {
    var alertingEventTypeFilters []interface{}

	for _, alertingEventTypeFilterData := range alertingEventTypeFilterDatas {
		alertingEventTypeFilter := map[string]interface{}{
    	}
        if alertingEventTypeFilterData.PredefinedEventFilter != nil {
            alertingEventTypeFilter["predefined_event_filter"] = flattenAlertingPredefinedEventFilter(alertingEventTypeFilterData.PredefinedEventFilter)
        }
        if alertingEventTypeFilterData.CustomEventFilter != nil {
            alertingEventTypeFilter["custom_event_filter"] = flattenAlertingCustomEventFilter(alertingEventTypeFilterData.CustomEventFilter)
        }

		alertingEventTypeFilters = append(alertingEventTypeFilters, alertingEventTypeFilter)
	}

	return alertingEventTypeFilters
}
