/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	tjconfig "github.com/crossplane/terrajet/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-jet-grafana/config/alertnotification"
	"github.com/crossplane-contrib/provider-jet-grafana/config/apikey"
	"github.com/crossplane-contrib/provider-jet-grafana/config/datasource"
	datasourcepermission "github.com/crossplane-contrib/provider-jet-grafana/config/datasource/permission"
	"github.com/crossplane-contrib/provider-jet-grafana/config/folder"
	folderpermission "github.com/crossplane-contrib/provider-jet-grafana/config/folder/permission"
	"github.com/crossplane-contrib/provider-jet-grafana/config/organization"
	"github.com/crossplane-contrib/provider-jet-grafana/config/team"
)

const (
	resourcePrefix = "grafana"
	modulePath     = "github.com/crossplane-contrib/provider-jet-grafana"
)

// GetProvider returns provider configuration
func GetProvider(resourceMap map[string]*schema.Resource) *tjconfig.Provider {
	defaultResourceFn := func(name string, terraformResource *schema.Resource, opts ...tjconfig.ResourceOption) *tjconfig.Resource {
		r := tjconfig.DefaultResource(name, terraformResource)
		// Add any provider-specific defaulting here. For example:
		r.ExternalName = tjconfig.IdentifierFromProvider
		return r
	}

	pc := tjconfig.NewProvider(resourceMap, resourcePrefix, modulePath,
		tjconfig.WithDefaultResourceFn(defaultResourceFn),
		tjconfig.WithIncludeList([]string{
			"grafana_organization$",
			"grafana_folder",
			// "grafana_data_source$", TODO: https://github.com/crossplane/terrajet/issues/208
			"grafana_folder_permission$",
			// "grafana_data_source_permission$", TODO: https://github.com/crossplane/terrajet/issues/208
			"grafana_api_key$",
			// "grafana_alert_notification$", TODO: https://github.com/crossplane/terrajet/issues/208
			"grafana_team$",
		}))

	for _, configure := range []func(provider *tjconfig.Provider){
		// add custom config functions
		organization.Configure,
		folder.Configure,
		folderpermission.Configure,
		team.Configure,
		datasource.Configure,
		datasourcepermission.Configure,
		alertnotification.Configure,
		apikey.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
