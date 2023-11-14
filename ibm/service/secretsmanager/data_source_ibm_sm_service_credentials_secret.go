// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmServiceCredentialsSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmServiceCredentialsSecretRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				Description:  "The ID of the secret.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				RequiredWith: []string{"secret_group_name"},
				Description:  "The human-readable name of your secret.",
			},

			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"secret_group_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"name"},
				Description:  "The human-readable name of your secret group.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
			},
			"state_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A text representation of the secret state.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"versions_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of versions of the secret.",
			},
			"version_custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "The secret version metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ttl": &schema.Schema{
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: StringIsIntBetween(60, 7776000),
				Description:  "The time-to-live (TTL) or lease duration to assign to generated credentials.",
			},
			"rotation": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Determines whether Secrets Manager rotates your secrets automatically.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotate": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Determines whether Secrets Manager rotates your secret automatically.Default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined interval.",
						},
						"interval": &schema.Schema{
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							Description:      "The length of the secret rotation time interval.",
							DiffSuppressFunc: rotationAttributesDiffSuppress,
						},
						"unit": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							Description:      "The units for the secret rotation time interval.",
							DiffSuppressFunc: rotationAttributesDiffSuppress,
						},
					},
				},
			},
			"next_rotation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that the secret is scheduled for automatic rotation. The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The properties of the service credentials secret payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apikey": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The API key that is generated for this secret.",
						},
						"cos_hmac_keys": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Cloud Object Storage HMAC keys that are returned after you create a service credentials secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The access key ID for Cloud Object Storage HMAC credentials.",
									},
									"secret_access_key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The secret access key ID for Cloud Object Storage HMAC credentials.",
									},
								},
							},
						},
						"endpoints": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoints that are returned after you create a service credentials secret.",
						},
						"iam_apikey_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the generated IAM API key.",
						},
						"iam_apikey_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the generated IAM API key.",
						},
						"iam_role_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM role CRN that is returned after you create a service credentials secret.",
						},
						"iam_serviceid_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM serviceId CRN that is returned after you create a service credentials secret.",
						},
						"resource_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource instance CRN that is returned after you create a service credentials secret.",
						},
					},
				},
			},
			"source_service": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "The properties of the source service credentials secret payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "A CRN that uniquely identifies a service credentials target.",
									},
								},
							},
						},
						"role": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The role identifier for creating a service-id.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The role identifier for creating a service-id.",
									},
								},
							},
						},
						"iam": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"apikey": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the generated IAM API key.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The description of the generated IAM API key.",
												},
											},
										},
									},
									"role": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM role CRN that is returned after you create a service credentials secret.",
												},
											},
										},
									},
									"serviceid": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM serviceId CRN that is returned after you create a service credentials secret.",
												},
											},
										},
									},
								},
							},
						},
						"resource_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource key CRN that is returned after you create a service credentials secret.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource key name that is returned after you create a service credentials secret.",
									},
								},
							},
						},
						"parameters": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The collection of parameters for the service credentials target.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSmServiceCredentialsSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ServiceCredentialsSecretIntf, region, instanceId, diagError := getSecretByIdOrByName(context, d, meta, ServiceCredentialsSecretType)
	if diagError != nil {
		return diagError
	}

	ServiceCredentialsSecret := ServiceCredentialsSecretIntf.(*secretsmanagerv2.ServiceCredentialsSecret)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *ServiceCredentialsSecret.ID))

	var err error
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("created_by", ServiceCredentialsSecret.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("created_at", DateTimeToRFC3339(ServiceCredentialsSecret.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("crn", ServiceCredentialsSecret.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if ServiceCredentialsSecret.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(ServiceCredentialsSecret.CustomMetadata))
		for k, v := range ServiceCredentialsSecret.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting custom_metadata: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting custom_metadata %s", err))
		}
	}

	if err = d.Set("description", ServiceCredentialsSecret.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("downloaded", ServiceCredentialsSecret.Downloaded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting downloaded: %s", err))
	}

	if ServiceCredentialsSecret.Labels != nil {
		if err = d.Set("labels", ServiceCredentialsSecret.Labels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting labels: %s", err))
		}
	}

	if err = d.Set("locks_total", flex.IntValue(ServiceCredentialsSecret.LocksTotal)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting locks_total: %s", err))
	}

	if err = d.Set("name", ServiceCredentialsSecret.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("secret_group_id", ServiceCredentialsSecret.SecretGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_group_id: %s", err))
	}

	if err = d.Set("secret_type", ServiceCredentialsSecret.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}

	if err = d.Set("state", flex.IntValue(ServiceCredentialsSecret.State)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	if err = d.Set("state_description", ServiceCredentialsSecret.StateDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state_description: %s", err))
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(ServiceCredentialsSecret.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("versions_total", flex.IntValue(ServiceCredentialsSecret.VersionsTotal)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting versions_total: %s", err))
	}

	if err = d.Set("ttl", ServiceCredentialsSecret.TTL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ttl: %s", err))
	}

	rotation := []map[string]interface{}{}
	if ServiceCredentialsSecret.Rotation != nil {
		modelMap, err := dataSourceIbmSmServiceCredentialsSecretRotationPolicyToMap(ServiceCredentialsSecret.Rotation.(*secretsmanagerv2.CommonRotationPolicy))
		if err != nil {
			return diag.FromErr(err)
		}
		rotation = append(rotation, modelMap)
	}
	if err = d.Set("rotation", rotation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rotation %s", err))
	}

	if err = d.Set("next_rotation_date", DateTimeToRFC3339(ServiceCredentialsSecret.NextRotationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting next_rotation_date: %s", err))
	}

	return nil
}

func dataSourceIbmSmServiceCredentialsSecretRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	return modelMap, nil
}
