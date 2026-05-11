package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

type MCPOptions struct {
	GetClient func() (*api.Client, error)
}

func NewCmdMCP(opts *MCPOptions) *cobra.Command {
	return &cobra.Command{
		Use:    "mcp-serve",
		Short:  "Start MCP server (JSON-RPC via stdio)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return serveMCP(opts)
		},
	}
}

func serveMCP(opts *MCPOptions) error {
	s := server.NewMCPServer("simplemdm-cli", "0.1.0",
		server.WithToolCapabilities(true),
	)

	// === ACCOUNT ===
	s.AddTool(mcp.NewTool("account-get",
		mcp.WithDescription("Get SimpleMDM account details"),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/account", nil
	}))

	s.AddTool(mcp.NewTool("account-update",
		mcp.WithDescription("Update SimpleMDM account"),
		mcp.WithString("name", mcp.Description("Account name")),
		mcp.WithString("apple_store_country_code", mcp.Description("Apple Store country code")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		values := extractFormValues(r, "name", "apple_store_country_code")
		return "/account", values
	}))

	// === APPS ===
	s.AddTool(mcp.NewTool("app-list",
		mcp.WithDescription("List all apps"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("include_shared", mcp.Description("Include shared apps: true or false")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps" + buildMCPQuery(r, "limit", "starting_after", "direction", "include_shared"), nil
	}))

	s.AddTool(mcp.NewTool("app-get",
		mcp.WithDescription("Get app details"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id"), nil
	}))

	s.AddTool(mcp.NewTool("app-create",
		mcp.WithDescription("Create an app"),
		mcp.WithString("name", mcp.Description("App name")),
		mcp.WithString("app_store_id", mcp.Description("App Store ID")),
		mcp.WithString("bundle_id", mcp.Description("Bundle ID")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/apps", extractFormValues(r, "name", "app_store_id", "bundle_id")
	}))

	s.AddTool(mcp.NewTool("app-update",
		mcp.WithDescription("Update an app"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
		mcp.WithString("name", mcp.Description("New name")),
		mcp.WithString("deploy_to", mcp.Description("Deploy to: none, outdated, or all")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/apps/" + getString(r, "app_id"), extractFormValues(r, "name", "deploy_to")
	}))

	s.AddTool(mcp.NewTool("app-delete",
		mcp.WithDescription("Delete an app"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id"), nil
	}))

	s.AddTool(mcp.NewTool("app-installs",
		mcp.WithDescription("List app installs"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id") + "/installs" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("app-managed-configs",
		mcp.WithDescription("List managed app configs"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id") + "/managed_configs", nil
	}))

	s.AddTool(mcp.NewTool("app-managed-config-create",
		mcp.WithDescription("Create a managed config"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Config key")),
		mcp.WithString("value", mcp.Description("Config value")),
		mcp.WithString("value_type", mcp.Description("Value type: boolean, date, float, float array, integer, integer array, string, string array")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/apps/" + getString(r, "app_id") + "/managed_configs", extractFormValues(r, "key", "value", "value_type")
	}))

	s.AddTool(mcp.NewTool("app-managed-configs-push",
		mcp.WithDescription("Push managed configs to devices"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id") + "/managed_configs/push", nil
	}))

	s.AddTool(mcp.NewTool("app-managed-config-delete",
		mcp.WithDescription("Delete a managed config"),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
		mcp.WithString("config_id", mcp.Required(), mcp.Description("Config ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/apps/" + getString(r, "app_id") + "/managed_configs/" + getString(r, "config_id"), nil
	}))

	// === ASSIGNMENT GROUPS ===
	s.AddTool(mcp.NewTool("assignment-group-list",
		mcp.WithDescription("List all assignment groups"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-get",
		mcp.WithDescription("Get assignment group details"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-create",
		mcp.WithDescription("Create an assignment group"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Group name")),
		mcp.WithBoolean("auto_deploy", mcp.Description("Auto deploy apps")),
		mcp.WithString("priority", mcp.Description("Group priority")),
		mcp.WithString("type", mcp.Description("Group type")),
		mcp.WithString("install_type", mcp.Description("Install type")),
		mcp.WithString("app_track_location", mcp.Description("App track location")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/assignment_groups", extractFormValues(r, "name", "auto_deploy", "priority", "type", "install_type", "app_track_location")
	}))

	s.AddTool(mcp.NewTool("assignment-group-update",
		mcp.WithDescription("Update an assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("name", mcp.Description("Group name")),
		mcp.WithBoolean("auto_deploy", mcp.Description("Auto deploy apps")),
		mcp.WithString("priority", mcp.Description("Group priority")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/assignment_groups/" + getString(r, "group_id"), extractFormValues(r, "name", "auto_deploy", "priority")
	}))

	s.AddTool(mcp.NewTool("assignment-group-delete",
		mcp.WithDescription("Delete an assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-assign-app",
		mcp.WithDescription("Assign app to assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
		mcp.WithString("deployment_type", mcp.Description("Deployment type: standard or munki")),
		mcp.WithString("install_type", mcp.Description("Install type: managed, self_serve, default_installs, managed_updates")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/assignment_groups/" + getString(r, "group_id") + "/apps/" + getString(r, "app_id"), extractFormValues(r, "deployment_type", "install_type")
	}))

	s.AddTool(mcp.NewTool("assignment-group-unassign-app",
		mcp.WithDescription("Unassign app from assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("app_id", mcp.Required(), mcp.Description("App ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/apps/" + getString(r, "app_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-assign-device",
		mcp.WithDescription("Assign device to assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithBoolean("remove_others", mcp.Description("Remove device from other assignment groups")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/assignment_groups/" + getString(r, "group_id") + "/devices/" + getString(r, "device_id"), extractFormValues(r, "remove_others")
	}))

	s.AddTool(mcp.NewTool("assignment-group-unassign-device",
		mcp.WithDescription("Unassign device from assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-assign-device-group",
		mcp.WithDescription("Assign device group to assignment group (deprecated)"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-unassign-device-group",
		mcp.WithDescription("Unassign device group from assignment group (deprecated)"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-assign-profile",
		mcp.WithDescription("Assign profile to assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/profiles/" + getString(r, "profile_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-unassign-profile",
		mcp.WithDescription("Unassign profile from assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/profiles/" + getString(r, "profile_id"), nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-custom-attributes",
		mcp.WithDescription("Get custom attribute values for assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/custom_attribute_values", nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-set-custom-attribute",
		mcp.WithDescription("Set custom attribute value for assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("attribute_name", mcp.Required(), mcp.Description("Attribute name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Attribute value")),
	), makeFormHandler(opts, "PUT", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/assignment_groups/" + getString(r, "group_id") + "/custom_attribute_values/" + getString(r, "attribute_name"), extractFormValues(r, "value")
	}))

	s.AddTool(mcp.NewTool("assignment-group-push-apps",
		mcp.WithDescription("Push apps to all devices in assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/push_apps", nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-update-apps",
		mcp.WithDescription("Update apps on all devices in assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/update_apps", nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-sync-profiles",
		mcp.WithDescription("Sync profiles for assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/sync_profiles", nil
	}))

	s.AddTool(mcp.NewTool("assignment-group-clone",
		mcp.WithDescription("Clone an assignment group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/assignment_groups/" + getString(r, "group_id") + "/clone", nil
	}))

	// === CUSTOM ATTRIBUTES ===
	s.AddTool(mcp.NewTool("custom-attribute-list",
		mcp.WithDescription("List all custom attributes"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_attributes" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("custom-attribute-get",
		mcp.WithDescription("Get custom attribute details"),
		mcp.WithString("attribute_id", mcp.Required(), mcp.Description("Attribute ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_attributes/" + getString(r, "attribute_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-attribute-create",
		mcp.WithDescription("Create a custom attribute"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Attribute name")),
		mcp.WithString("default_value", mcp.Description("Default value")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/custom_attributes", extractFormValues(r, "name", "default_value")
	}))

	s.AddTool(mcp.NewTool("custom-attribute-update",
		mcp.WithDescription("Update a custom attribute"),
		mcp.WithString("attribute_id", mcp.Required(), mcp.Description("Attribute ID")),
		mcp.WithString("default_value", mcp.Required(), mcp.Description("Default value")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/custom_attributes/" + getString(r, "attribute_id"), extractFormValues(r, "default_value")
	}))

	s.AddTool(mcp.NewTool("custom-attribute-delete",
		mcp.WithDescription("Delete a custom attribute"),
		mcp.WithString("attribute_id", mcp.Required(), mcp.Description("Attribute ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_attributes/" + getString(r, "attribute_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-attribute-set-value",
		mcp.WithDescription("Set custom attribute value for multiple devices (JSON)"),
		mcp.WithString("attribute_name", mcp.Required(), mcp.Description("Attribute name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("JSON array: [{\"device_id\":1,\"value\":\"x\"},...]")),
	), makeJSONHandler(opts, "PUT", func(r mcp.CallToolRequest) (string, string) {
		return "/custom_attribute_values/" + getString(r, "attribute_name"), getString(r, "data")
	}))

	// === CUSTOM CONFIGURATION PROFILES ===
	s.AddTool(mcp.NewTool("custom-configuration-profile-list",
		mcp.WithDescription("List all custom configuration profiles"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("search", mcp.Description("Search by name")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles" + buildMCPQuery(r, "limit", "starting_after", "direction", "search"), nil
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-update",
		mcp.WithDescription("Update a custom configuration profile"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("name", mcp.Description("Profile name")),
		mcp.WithString("user_scope", mcp.Description("User scope")),
		mcp.WithString("attribute_support", mcp.Description("Enable attribute support: true or false")),
		mcp.WithString("escape_attributes", mcp.Description("Escape attributes: true or false")),
		mcp.WithString("reinstall_after_os_update", mcp.Description("Reinstall after OS update: true or false")),
		mcp.WithString("declarative", mcp.Description("Declarative profile: true or false")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id"), extractFormValues(r, "name", "user_scope", "attribute_support", "escape_attributes", "reinstall_after_os_update", "declarative")
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-delete",
		mcp.WithDescription("Delete a custom configuration profile"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-push-device",
		mcp.WithDescription("Push custom configuration profile to device"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-remove-device",
		mcp.WithDescription("Remove custom configuration profile from device"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-assign-device-group",
		mcp.WithDescription("Assign custom configuration profile to device group (deprecated)"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-configuration-profile-unassign-device-group",
		mcp.WithDescription("Unassign custom configuration profile from device group (deprecated)"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_configuration_profiles/" + getString(r, "profile_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	// === CUSTOM DECLARATIONS ===
	s.AddTool(mcp.NewTool("custom-declaration-list",
		mcp.WithDescription("List all custom declarations"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("search", mcp.Description("Search by name")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_declarations" + buildMCPQuery(r, "limit", "starting_after", "direction", "search"), nil
	}))

	s.AddTool(mcp.NewTool("custom-declaration-update",
		mcp.WithDescription("Update a custom declaration"),
		mcp.WithString("declaration_id", mcp.Required(), mcp.Description("Declaration ID")),
		mcp.WithString("name", mcp.Description("Declaration name")),
		mcp.WithString("declaration_type", mcp.Description("Declaration type")),
		mcp.WithString("user_scope", mcp.Description("User scope")),
		mcp.WithString("attribute_support", mcp.Description("Enable attribute support: true or false")),
		mcp.WithString("escape_attributes", mcp.Description("Escape attributes: true or false")),
		mcp.WithString("activation_predicate", mcp.Description("Activation predicate")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/custom_declarations/" + getString(r, "declaration_id"), extractFormValues(r, "name", "declaration_type", "user_scope", "attribute_support", "escape_attributes", "activation_predicate")
	}))

	s.AddTool(mcp.NewTool("custom-declaration-delete",
		mcp.WithDescription("Delete a custom declaration"),
		mcp.WithString("declaration_id", mcp.Required(), mcp.Description("Declaration ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_declarations/" + getString(r, "declaration_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-declaration-push-device",
		mcp.WithDescription("Push custom declaration to device"),
		mcp.WithString("declaration_id", mcp.Required(), mcp.Description("Declaration ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_declarations/" + getString(r, "declaration_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("custom-declaration-remove-device",
		mcp.WithDescription("Remove custom declaration from device"),
		mcp.WithString("declaration_id", mcp.Required(), mcp.Description("Declaration ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/custom_declarations/" + getString(r, "declaration_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	// === DEP SERVERS ===
	s.AddTool(mcp.NewTool("dep-server-list",
		mcp.WithDescription("List all DEP servers"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/dep_servers" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("dep-server-get",
		mcp.WithDescription("Get DEP server details"),
		mcp.WithString("server_id", mcp.Required(), mcp.Description("Server ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/dep_servers/" + getString(r, "server_id"), nil
	}))

	s.AddTool(mcp.NewTool("dep-server-devices",
		mcp.WithDescription("List DEP devices for server"),
		mcp.WithString("server_id", mcp.Required(), mcp.Description("Server ID")),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/dep_servers/" + getString(r, "server_id") + "/dep_devices" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("dep-server-device-get",
		mcp.WithDescription("Get DEP device details"),
		mcp.WithString("server_id", mcp.Required(), mcp.Description("Server ID")),
		mcp.WithString("dep_device_id", mcp.Required(), mcp.Description("DEP Device ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/dep_servers/" + getString(r, "server_id") + "/dep_devices/" + getString(r, "dep_device_id"), nil
	}))

	s.AddTool(mcp.NewTool("dep-server-sync",
		mcp.WithDescription("Sync DEP server"),
		mcp.WithString("server_id", mcp.Required(), mcp.Description("Server ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/dep_servers/" + getString(r, "server_id") + "/sync", nil
	}))

	// === DEVICE GROUPS ===
	s.AddTool(mcp.NewTool("device-group-list",
		mcp.WithDescription("List all device groups"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/device_groups" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("device-group-get",
		mcp.WithDescription("Get device group details"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/device_groups/" + getString(r, "group_id"), nil
	}))

	s.AddTool(mcp.NewTool("device-group-assign-device",
		mcp.WithDescription("Assign device to device group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/device_groups/" + getString(r, "group_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("device-group-clone",
		mcp.WithDescription("Clone a device group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/device_groups/" + getString(r, "group_id") + "/clone", nil
	}))

	s.AddTool(mcp.NewTool("device-group-custom-attributes",
		mcp.WithDescription("Get custom attribute values for device group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/device_groups/" + getString(r, "group_id") + "/custom_attribute_values", nil
	}))

	s.AddTool(mcp.NewTool("device-group-set-custom-attribute",
		mcp.WithDescription("Set custom attribute value for device group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("attribute_name", mcp.Required(), mcp.Description("Attribute name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Attribute value")),
	), makeFormHandler(opts, "PUT", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/device_groups/" + getString(r, "group_id") + "/custom_attribute_values/" + getString(r, "attribute_name"), extractFormValues(r, "value")
	}))

	// === DEVICES ===
	s.AddTool(mcp.NewTool("device-list",
		mcp.WithDescription("List all devices"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("search", mcp.Description("Search query")),
		mcp.WithString("include_awaiting_enrollment", mcp.Description("Include awaiting enrollment: true or false")),
		mcp.WithString("include_secret_custom_attributes", mcp.Description("Include secret custom attributes: true or false")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices" + buildMCPQuery(r, "limit", "starting_after", "direction", "search", "include_awaiting_enrollment", "include_secret_custom_attributes"), nil
	}))

	s.AddTool(mcp.NewTool("device-get",
		mcp.WithDescription("Get device details"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("include_secret_custom_attributes", mcp.Description("Include secret custom attributes: true or false")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + buildMCPQuery(r, "include_secret_custom_attributes"), nil
	}))

	s.AddTool(mcp.NewTool("device-create",
		mcp.WithDescription("Create a device"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Device name")),
		mcp.WithString("group_id", mcp.Description("Group ID (deprecated)")),
		mcp.WithString("static_group_ids", mcp.Description("Comma-separated static group IDs")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices", extractFormValues(r, "name", "group_id", "static_group_ids")
	}))

	s.AddTool(mcp.NewTool("device-update",
		mcp.WithDescription("Update a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("name", mcp.Description("New SimpleMDM device name")),
		mcp.WithString("device_name", mcp.Description("New device name (requires supervision)")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id"), extractFormValues(r, "name", "device_name")
	}))

	s.AddTool(mcp.NewTool("device-delete",
		mcp.WithDescription("Delete a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("device-installed-apps",
		mcp.WithDescription("List installed apps on device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/installed_apps" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("device-push-apps",
		mcp.WithDescription("Push assigned apps to device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/push_apps", nil
	}))

	s.AddTool(mcp.NewTool("device-refresh",
		mcp.WithDescription("Refresh device information"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/refresh", nil
	}))

	s.AddTool(mcp.NewTool("device-lock",
		mcp.WithDescription("Lock a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("message", mcp.Description("Lock message")),
		mcp.WithString("phone_number", mcp.Description("Phone number")),
		mcp.WithString("pin", mcp.Description("PIN code")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/lock", extractFormValues(r, "message", "phone_number", "pin")
	}))

	s.AddTool(mcp.NewTool("device-wipe",
		mcp.WithDescription("Wipe a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("pin", mcp.Description("PIN code")),
		mcp.WithBoolean("clear_custom_attributes", mcp.Description("Clear custom attributes")),
		mcp.WithBoolean("disable_activation_lock", mcp.Description("Disable activation lock")),
		mcp.WithBoolean("preserve_data_plan", mcp.Description("Preserve data plan (iOS)")),
		mcp.WithBoolean("disallow_proximity_setup", mcp.Description("Disallow proximity setup (iOS)")),
		mcp.WithBoolean("return_to_service", mcp.Description("Return to service (iOS 17+/tvOS 18+)")),
		mcp.WithString("wifi_network_id", mcp.Description("WiFi network ID for return to service")),
		mcp.WithString("obliteration_behavior", mcp.Description("Obliteration behavior (macOS 12+)")),
		mcp.WithBoolean("unassign_direct_profiles", mcp.Description("Unassign direct profiles")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/wipe", extractFormValues(r, "pin", "clear_custom_attributes", "disable_activation_lock", "preserve_data_plan", "disallow_proximity_setup", "return_to_service", "wifi_network_id", "obliteration_behavior", "unassign_direct_profiles")
	}))

	s.AddTool(mcp.NewTool("device-restart",
		mcp.WithDescription("Restart a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithBoolean("rebuild_kernel_cache", mcp.Description("Rebuild kernel cache (macOS 11+)")),
		mcp.WithBoolean("notify_user", mcp.Description("Notify user (macOS 11.3+)")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/restart", extractFormValues(r, "rebuild_kernel_cache", "notify_user")
	}))

	s.AddTool(mcp.NewTool("device-shutdown",
		mcp.WithDescription("Shut down a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/shutdown", nil
	}))

	s.AddTool(mcp.NewTool("device-clear-passcode",
		mcp.WithDescription("Clear device passcode"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/clear_passcode", nil
	}))

	s.AddTool(mcp.NewTool("device-clear-firmware-password",
		mcp.WithDescription("Clear firmware password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/clear_firmware_password", nil
	}))

	s.AddTool(mcp.NewTool("device-rotate-firmware-password",
		mcp.WithDescription("Rotate firmware password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/rotate_firmware_password", nil
	}))

	s.AddTool(mcp.NewTool("device-clear-recovery-lock",
		mcp.WithDescription("Clear recovery lock password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/clear_recovery_lock_password", nil
	}))

	s.AddTool(mcp.NewTool("device-rotate-recovery-lock",
		mcp.WithDescription("Rotate recovery lock password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/rotate_recovery_lock_password", nil
	}))

	s.AddTool(mcp.NewTool("device-rotate-filevault-key",
		mcp.WithDescription("Rotate FileVault recovery key"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/rotate_filevault_key", nil
	}))

	s.AddTool(mcp.NewTool("device-set-admin-password",
		mcp.WithDescription("Set admin password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("new_password", mcp.Required(), mcp.Description("New admin password")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/set_admin_password", extractFormValues(r, "new_password")
	}))

	s.AddTool(mcp.NewTool("device-rotate-admin-password",
		mcp.WithDescription("Rotate admin password"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/rotate_admin_password", nil
	}))

	s.AddTool(mcp.NewTool("device-clear-restrictions-password",
		mcp.WithDescription("Clear restrictions password (iOS/iPad)"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/clear_restrictions_password", nil
	}))

	s.AddTool(mcp.NewTool("device-bluetooth-enable",
		mcp.WithDescription("Enable Bluetooth"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/bluetooth", nil
	}))

	s.AddTool(mcp.NewTool("device-bluetooth-disable",
		mcp.WithDescription("Disable Bluetooth"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/bluetooth", nil
	}))

	s.AddTool(mcp.NewTool("device-remote-desktop-enable",
		mcp.WithDescription("Enable Remote Desktop"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/remote_desktop", nil
	}))

	s.AddTool(mcp.NewTool("device-remote-desktop-disable",
		mcp.WithDescription("Disable Remote Desktop"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/remote_desktop", nil
	}))

	s.AddTool(mcp.NewTool("device-set-timezone",
		mcp.WithDescription("Set device time zone"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("time_zone", mcp.Required(), mcp.Description("Time zone (TZ identifier)")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/set_time_zone", extractFormValues(r, "time_zone")
	}))

	s.AddTool(mcp.NewTool("device-delete-user",
		mcp.WithDescription("Delete a user from device (macOS)"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/users/" + getString(r, "user_id"), nil
	}))

	s.AddTool(mcp.NewTool("device-update-os",
		mcp.WithDescription("Update device OS"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("os_update_mode", mcp.Description("Update mode: smart_update, download_only, notify_only, install_asap, force_update")),
		mcp.WithString("version_type", mcp.Description("Version type: latest_minor_version or latest_major_version")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/update_os", extractFormValues(r, "os_update_mode", "version_type")
	}))

	s.AddTool(mcp.NewTool("device-profiles",
		mcp.WithDescription("List device profiles"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/profiles" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("device-users",
		mcp.WithDescription("List device users (macOS only)"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/users" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("device-unenroll",
		mcp.WithDescription("Unenroll a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithBoolean("unassign_direct_profiles", mcp.Description("Unassign direct profiles")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/unenroll", extractFormValues(r, "unassign_direct_profiles")
	}))

	s.AddTool(mcp.NewTool("device-custom-attributes",
		mcp.WithDescription("Get device custom attribute values"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/custom_attribute_values", nil
	}))

	s.AddTool(mcp.NewTool("device-set-custom-attribute",
		mcp.WithDescription("Set a device custom attribute"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("attribute_name", mcp.Required(), mcp.Description("Attribute name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Attribute value")),
	), makeFormHandler(opts, "PUT", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/custom_attribute_values/" + getString(r, "attribute_name"), extractFormValues(r, "value")
	}))

	s.AddTool(mcp.NewTool("device-set-custom-attributes",
		mcp.WithDescription("Set multiple custom attribute values for a device (JSON)"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("data", mcp.Required(), mcp.Description("JSON array: [{\"name\":\"attr\",\"value\":\"x\"},...]")),
	), makeJSONHandler(opts, "PUT", func(r mcp.CallToolRequest) (string, string) {
		return "/devices/" + getString(r, "device_id") + "/custom_attribute_values", getString(r, "data")
	}))

	// Lost Mode
	s.AddTool(mcp.NewTool("device-lost-mode-enable",
		mcp.WithDescription("Enable Lost Mode on device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
		mcp.WithString("message", mcp.Description("Message to display")),
		mcp.WithString("phone_number", mcp.Description("Phone number")),
		mcp.WithString("footnote", mcp.Description("Footnote")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/devices/" + getString(r, "device_id") + "/lost_mode", extractFormValues(r, "message", "phone_number", "footnote")
	}))

	s.AddTool(mcp.NewTool("device-lost-mode-disable",
		mcp.WithDescription("Disable Lost Mode on device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/lost_mode", nil
	}))

	s.AddTool(mcp.NewTool("device-lost-mode-play-sound",
		mcp.WithDescription("Play Lost Mode sound"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/lost_mode/play_sound", nil
	}))

	s.AddTool(mcp.NewTool("device-lost-mode-update-location",
		mcp.WithDescription("Request location update"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/devices/" + getString(r, "device_id") + "/lost_mode/update_location", nil
	}))

	// === ENROLLMENTS ===
	s.AddTool(mcp.NewTool("enrollment-list",
		mcp.WithDescription("List all enrollments"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/enrollments" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("enrollment-get",
		mcp.WithDescription("Get enrollment details"),
		mcp.WithString("enrollment_id", mcp.Required(), mcp.Description("Enrollment ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/enrollments/" + getString(r, "enrollment_id"), nil
	}))

	s.AddTool(mcp.NewTool("enrollment-delete",
		mcp.WithDescription("Delete an enrollment"),
		mcp.WithString("enrollment_id", mcp.Required(), mcp.Description("Enrollment ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/enrollments/" + getString(r, "enrollment_id"), nil
	}))

	s.AddTool(mcp.NewTool("enrollment-send-invitation",
		mcp.WithDescription("Send enrollment invitation"),
		mcp.WithString("enrollment_id", mcp.Required(), mcp.Description("Enrollment ID")),
		mcp.WithString("contact", mcp.Required(), mcp.Description("Email or phone number")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/enrollments/" + getString(r, "enrollment_id") + "/invitations", extractFormValues(r, "contact")
	}))

	// === INSTALLED APPS ===
	s.AddTool(mcp.NewTool("installed-app-get",
		mcp.WithDescription("Get installed app details"),
		mcp.WithString("installed_app_id", mcp.Required(), mcp.Description("Installed App ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/installed_apps/" + getString(r, "installed_app_id"), nil
	}))

	s.AddTool(mcp.NewTool("installed-app-delete",
		mcp.WithDescription("Delete an installed app"),
		mcp.WithString("installed_app_id", mcp.Required(), mcp.Description("Installed App ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/installed_apps/" + getString(r, "installed_app_id"), nil
	}))

	s.AddTool(mcp.NewTool("installed-app-update",
		mcp.WithDescription("Update an installed app"),
		mcp.WithString("installed_app_id", mcp.Required(), mcp.Description("Installed App ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/installed_apps/" + getString(r, "installed_app_id") + "/update", nil
	}))

	s.AddTool(mcp.NewTool("installed-app-request-management",
		mcp.WithDescription("Request management of an installed app"),
		mcp.WithString("installed_app_id", mcp.Required(), mcp.Description("Installed App ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/installed_apps/" + getString(r, "installed_app_id") + "/request_management", nil
	}))

	// === LOGS ===
	s.AddTool(mcp.NewTool("log-list",
		mcp.WithDescription("List all logs"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("serial_number", mcp.Description("Filter by device serial number")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/logs" + buildMCPQuery(r, "limit", "starting_after", "direction", "serial_number"), nil
	}))

	s.AddTool(mcp.NewTool("log-get",
		mcp.WithDescription("Get log details"),
		mcp.WithString("log_id", mcp.Required(), mcp.Description("Log ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/logs/" + getString(r, "log_id"), nil
	}))

	// === PROFILES ===
	s.AddTool(mcp.NewTool("profile-list",
		mcp.WithDescription("List all profiles"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
		mcp.WithString("search", mcp.Description("Search by name or type")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles" + buildMCPQuery(r, "limit", "starting_after", "direction", "search"), nil
	}))

	s.AddTool(mcp.NewTool("profile-get",
		mcp.WithDescription("Get profile details"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles/" + getString(r, "profile_id"), nil
	}))

	s.AddTool(mcp.NewTool("profile-assign-device",
		mcp.WithDescription("Assign profile to device"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles/" + getString(r, "profile_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("profile-unassign-device",
		mcp.WithDescription("Unassign profile from device"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles/" + getString(r, "profile_id") + "/devices/" + getString(r, "device_id"), nil
	}))

	s.AddTool(mcp.NewTool("profile-assign-device-group",
		mcp.WithDescription("Assign profile to device group (deprecated)"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "POST", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles/" + getString(r, "profile_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	s.AddTool(mcp.NewTool("profile-unassign-device-group",
		mcp.WithDescription("Unassign profile from device group (deprecated)"),
		mcp.WithString("profile_id", mcp.Required(), mcp.Description("Profile ID")),
		mcp.WithString("device_group_id", mcp.Required(), mcp.Description("Device Group ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/profiles/" + getString(r, "profile_id") + "/device_groups/" + getString(r, "device_group_id"), nil
	}))

	// === PUSH CERTIFICATE ===
	s.AddTool(mcp.NewTool("push-certificate-get",
		mcp.WithDescription("Get push certificate details"),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/push_certificate", nil
	}))

	s.AddTool(mcp.NewTool("push-certificate-scsr",
		mcp.WithDescription("Get signed CSR for push certificate"),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/push_certificate/scsr", nil
	}))

	// === SCRIPTS ===
	s.AddTool(mcp.NewTool("script-list",
		mcp.WithDescription("List all scripts"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/scripts" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("script-get",
		mcp.WithDescription("Get script details"),
		mcp.WithString("script_id", mcp.Required(), mcp.Description("Script ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/scripts/" + getString(r, "script_id"), nil
	}))

	s.AddTool(mcp.NewTool("script-update",
		mcp.WithDescription("Update a script"),
		mcp.WithString("script_id", mcp.Required(), mcp.Description("Script ID")),
		mcp.WithString("name", mcp.Description("Script name")),
		mcp.WithString("variable_support", mcp.Description("Enable variable support: true or false")),
	), makeFormHandler(opts, "PATCH", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/scripts/" + getString(r, "script_id"), extractFormValues(r, "name", "variable_support")
	}))

	s.AddTool(mcp.NewTool("script-delete",
		mcp.WithDescription("Delete a script"),
		mcp.WithString("script_id", mcp.Required(), mcp.Description("Script ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/scripts/" + getString(r, "script_id"), nil
	}))

	// === SCRIPT JOBS ===
	s.AddTool(mcp.NewTool("script-job-list",
		mcp.WithDescription("List all script jobs"),
		mcp.WithNumber("limit", mcp.Description("Limit results")),
		mcp.WithString("starting_after", mcp.Description("Pagination cursor")),
		mcp.WithString("direction", mcp.Description("Sort direction: asc or desc")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/script_jobs" + buildMCPQuery(r, "limit", "starting_after", "direction"), nil
	}))

	s.AddTool(mcp.NewTool("script-job-get",
		mcp.WithDescription("Get script job details"),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("Job ID")),
	), makeHandler(opts, "GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/script_jobs/" + getString(r, "job_id"), nil
	}))

	s.AddTool(mcp.NewTool("script-job-create",
		mcp.WithDescription("Create a script job"),
		mcp.WithString("script_id", mcp.Required(), mcp.Description("Script ID")),
		mcp.WithString("device_ids", mcp.Description("Comma-separated device IDs")),
		mcp.WithString("assignment_group_ids", mcp.Description("Comma-separated assignment group IDs")),
		mcp.WithString("group_ids", mcp.Description("Comma-separated device group IDs (deprecated)")),
		mcp.WithString("custom_attribute", mcp.Description("Custom attribute filter")),
		mcp.WithString("custom_attribute_regex", mcp.Description("Custom attribute regex filter")),
	), makeFormHandler(opts, "POST", func(r mcp.CallToolRequest) (string, map[string]string) {
		return "/script_jobs", extractFormValues(r, "script_id", "device_ids", "assignment_group_ids", "group_ids", "custom_attribute", "custom_attribute_regex")
	}))

	s.AddTool(mcp.NewTool("script-job-cancel",
		mcp.WithDescription("Cancel a script job"),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("Job ID")),
	), makeHandler(opts, "DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/script_jobs/" + getString(r, "job_id"), nil
	}))

	return server.ServeStdio(s)
}

// === HELPERS ===

func makeHandler(opts *MCPOptions, method string, pathFn func(mcp.CallToolRequest) (string, interface{})) server.ToolHandlerFunc {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := opts.GetClient()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		path, _ := pathFn(r)
		var data []byte

		switch method {
		case "GET":
			data, err = client.Get(path)
		case "POST":
			data, err = client.Post(path, nil)
		case "DELETE":
			data, err = client.Delete(path)
		default:
			return mcp.NewToolResultError("unsupported method: " + method), nil
		}

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func makeFormHandler(opts *MCPOptions, method string, pathFn func(mcp.CallToolRequest) (string, map[string]string)) server.ToolHandlerFunc {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := opts.GetClient()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		path, values := pathFn(r)
		var data []byte

		switch method {
		case "POST":
			data, err = client.DoForm(method, path, values)
		case "PATCH":
			data, err = client.DoForm(method, path, values)
		case "PUT":
			data, err = client.DoForm(method, path, values)
		default:
			return mcp.NewToolResultError("unsupported method: " + method), nil
		}

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func getString(r mcp.CallToolRequest, key string) string {
	return r.GetString(key, "")
}

func extractFormValues(r mcp.CallToolRequest, keys ...string) map[string]string {
	values := make(map[string]string)
	args := r.GetArguments()
	for _, k := range keys {
		if v, ok := args[k]; ok {
			s := fmt.Sprintf("%v", v)
			if s != "" && s != "<nil>" {
				values[k] = s
			}
		}
	}
	return values
}

func makeJSONHandler(opts *MCPOptions, method string, pathFn func(mcp.CallToolRequest) (string, string)) server.ToolHandlerFunc {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := opts.GetClient()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		path, body := pathFn(r)
		data, err := client.Do(method, path, strings.NewReader(body))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func buildMCPQuery(r mcp.CallToolRequest, params ...string) string {
	var parts []string
	args := r.GetArguments()
	for _, p := range params {
		if v, ok := args[p]; ok {
			s := fmt.Sprintf("%v", v)
			if s != "" && s != "0" && s != "<nil>" {
				parts = append(parts, p+"="+s)
			}
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "?" + strings.Join(parts, "&")
}
