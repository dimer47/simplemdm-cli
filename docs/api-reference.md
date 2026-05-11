# SimpleMDM API Reference

## Introduction

SimpleMDM provides a REST API for managing Apple devices. All API requests are made to the base URL:

```
https://a.simplemdm.com/api/v1
```

**Authentication:** HTTP Basic Auth with your SimpleMDM API key as the username and an empty password.

```bash
curl https://a.simplemdm.com/api/v1/account \
  -u "API_KEY:"
```

**Pagination:** List endpoints support `limit` and `starting_after` query parameters for cursor-based pagination.

---

## Account

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/account` | Get account details | `account-get` |
| PATCH | `/account` | Update account (name, country code) | `account-update` |

---

## Apps

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/apps` | List all apps | `app-list` |
| GET | `/apps/{app_id}` | Get app details | `app-get` |
| POST | `/apps` | Create an app | `app-create` |
| PATCH | `/apps/{app_id}` | Update an app | `app-update` |
| DELETE | `/apps/{app_id}` | Delete an app | `app-delete` |
| GET | `/apps/{app_id}/installs` | List app installs | `app-installs` |
| GET | `/apps/{app_id}/managed_configs` | List managed app configs | `app-managed-configs` |
| POST | `/apps/{app_id}/managed_configs` | Create a managed config | `app-managed-config-create` |
| POST | `/apps/{app_id}/managed_configs/push` | Push managed configs to devices | `app-managed-configs-push` |
| DELETE | `/apps/{app_id}/managed_configs/{config_id}` | Delete a managed config | `app-managed-config-delete` |

---

## Assignment Groups

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/assignment_groups` | List all assignment groups | `assignment-group-list` |
| GET | `/assignment_groups/{group_id}` | Get assignment group details | `assignment-group-get` |
| POST | `/assignment_groups` | Create an assignment group | `assignment-group-create` |
| PATCH | `/assignment_groups/{group_id}` | Update an assignment group | `assignment-group-update` |
| DELETE | `/assignment_groups/{group_id}` | Delete an assignment group | `assignment-group-delete` |
| POST | `/assignment_groups/{group_id}/apps/{app_id}` | Assign app to group | `assignment-group-assign-app` |
| DELETE | `/assignment_groups/{group_id}/apps/{app_id}` | Unassign app from group | `assignment-group-unassign-app` |
| POST | `/assignment_groups/{group_id}/devices/{device_id}` | Assign device to group | `assignment-group-assign-device` |
| DELETE | `/assignment_groups/{group_id}/devices/{device_id}` | Unassign device from group | `assignment-group-unassign-device` |
| POST | `/assignment_groups/{group_id}/device_groups/{device_group_id}` | Assign device group (deprecated) | `assignment-group-assign-device-group` |
| DELETE | `/assignment_groups/{group_id}/device_groups/{device_group_id}` | Unassign device group (deprecated) | `assignment-group-unassign-device-group` |
| POST | `/assignment_groups/{group_id}/profiles/{profile_id}` | Assign profile to group | `assignment-group-assign-profile` |
| DELETE | `/assignment_groups/{group_id}/profiles/{profile_id}` | Unassign profile from group | `assignment-group-unassign-profile` |
| GET | `/assignment_groups/{group_id}/custom_attribute_values` | Get custom attribute values | `assignment-group-custom-attributes` |
| PUT | `/assignment_groups/{group_id}/custom_attribute_values/{attribute_name}` | Set custom attribute value | `assignment-group-set-custom-attribute` |
| POST | `/assignment_groups/{group_id}/push_apps` | Push apps to all devices in group | `assignment-group-push-apps` |
| POST | `/assignment_groups/{group_id}/update_apps` | Update apps on all devices in group | `assignment-group-update-apps` |
| POST | `/assignment_groups/{group_id}/sync_profiles` | Sync profiles for group | `assignment-group-sync-profiles` |
| POST | `/assignment_groups/{group_id}/clone` | Clone an assignment group | `assignment-group-clone` |

---

## Custom Attributes

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/custom_attributes` | List all custom attributes | `custom-attribute-list` |
| GET | `/custom_attributes/{attribute_id}` | Get custom attribute details | `custom-attribute-get` |
| POST | `/custom_attributes` | Create a custom attribute | `custom-attribute-create` |
| PATCH | `/custom_attributes/{attribute_id}` | Update a custom attribute | `custom-attribute-update` |
| DELETE | `/custom_attributes/{attribute_id}` | Delete a custom attribute | `custom-attribute-delete` |

---

## Custom Configuration Profiles

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/custom_configuration_profiles` | List all custom configuration profiles | `custom-configuration-profile-list` |
| PATCH | `/custom_configuration_profiles/{profile_id}` | Update a custom configuration profile | `custom-configuration-profile-update` |
| DELETE | `/custom_configuration_profiles/{profile_id}` | Delete a custom configuration profile | `custom-configuration-profile-delete` |
| POST | `/custom_configuration_profiles/{profile_id}/devices/{device_id}` | Push profile to device | `custom-configuration-profile-push-device` |
| DELETE | `/custom_configuration_profiles/{profile_id}/devices/{device_id}` | Remove profile from device | `custom-configuration-profile-remove-device` |

---

## Custom Declarations

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/custom_declarations` | List all custom declarations | `custom-declaration-list` |
| PATCH | `/custom_declarations/{declaration_id}` | Update a custom declaration | `custom-declaration-update` |
| DELETE | `/custom_declarations/{declaration_id}` | Delete a custom declaration | `custom-declaration-delete` |
| POST | `/custom_declarations/{declaration_id}/devices/{device_id}` | Push declaration to device | `custom-declaration-push-device` |
| DELETE | `/custom_declarations/{declaration_id}/devices/{device_id}` | Remove declaration from device | `custom-declaration-remove-device` |

---

## DEP Servers

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/dep_servers` | List all DEP servers | `dep-server-list` |
| GET | `/dep_servers/{server_id}` | Get DEP server details | `dep-server-get` |
| GET | `/dep_servers/{server_id}/dep_devices` | List DEP devices for server | `dep-server-devices` |
| GET | `/dep_servers/{server_id}/dep_devices/{dep_device_id}` | Get DEP device details | `dep-server-device-get` |
| POST | `/dep_servers/{server_id}/sync` | Sync DEP server | `dep-server-sync` |

---

## Device Groups

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/device_groups` | List all device groups | `device-group-list` |
| GET | `/device_groups/{group_id}` | Get device group details | `device-group-get` |
| POST | `/device_groups/{group_id}/devices/{device_id}` | Assign device to group | `device-group-assign-device` |
| POST | `/device_groups/{group_id}/clone` | Clone a device group | `device-group-clone` |

---

## Devices

### CRUD Operations

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/devices` | List all devices | `device-list` |
| GET | `/devices/{device_id}` | Get device details | `device-get` |
| POST | `/devices` | Create a device | `device-create` |
| PATCH | `/devices/{device_id}` | Update a device | `device-update` |
| DELETE | `/devices/{device_id}` | Delete a device | `device-delete` |

### Device Information

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/devices/{device_id}/installed_apps` | List installed apps | `device-installed-apps` |
| GET | `/devices/{device_id}/profiles` | List device profiles | `device-profiles` |
| GET | `/devices/{device_id}/users` | List device users | `device-users` |
| GET | `/devices/{device_id}/custom_attribute_values` | Get custom attribute values | `device-custom-attributes` |
| PUT | `/devices/{device_id}/custom_attribute_values/{attribute_name}` | Set a custom attribute | `device-set-custom-attribute` |

### Device Actions

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| POST | `/devices/{device_id}/push_apps` | Push assigned apps | `device-push-apps` |
| POST | `/devices/{device_id}/refresh` | Refresh device information | `device-refresh` |
| POST | `/devices/{device_id}/lock` | Lock device | `device-lock` |
| POST | `/devices/{device_id}/wipe` | Wipe device | `device-wipe` |
| POST | `/devices/{device_id}/restart` | Restart device | `device-restart` |
| POST | `/devices/{device_id}/shutdown` | Shut down device | `device-shutdown` |
| POST | `/devices/{device_id}/unenroll` | Unenroll device | `device-unenroll` |
| POST | `/devices/{device_id}/update_os` | Update device OS | `device-update-os` |
| POST | `/devices/{device_id}/set_time_zone` | Set device time zone | `device-set-timezone` |
| DELETE | `/devices/{device_id}/users/{user_id}` | Delete a user from device (macOS) | `device-delete-user` |

### Security & Passwords

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| POST | `/devices/{device_id}/clear_passcode` | Clear device passcode | `device-clear-passcode` |
| POST | `/devices/{device_id}/clear_firmware_password` | Clear firmware password | `device-clear-firmware-password` |
| POST | `/devices/{device_id}/rotate_firmware_password` | Rotate firmware password | `device-rotate-firmware-password` |
| POST | `/devices/{device_id}/clear_recovery_lock_password` | Clear recovery lock password | `device-clear-recovery-lock` |
| POST | `/devices/{device_id}/rotate_recovery_lock_password` | Rotate recovery lock password | `device-rotate-recovery-lock` |
| POST | `/devices/{device_id}/rotate_filevault_key` | Rotate FileVault recovery key | `device-rotate-filevault-key` |
| POST | `/devices/{device_id}/set_admin_password` | Set admin password | `device-set-admin-password` |
| POST | `/devices/{device_id}/rotate_admin_password` | Rotate admin password | `device-rotate-admin-password` |
| POST | `/devices/{device_id}/clear_restrictions_password` | Clear restrictions password (iOS/iPad) | `device-clear-restrictions-password` |

### Bluetooth & Remote Desktop

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| POST | `/devices/{device_id}/bluetooth` | Enable Bluetooth | `device-bluetooth-enable` |
| DELETE | `/devices/{device_id}/bluetooth` | Disable Bluetooth | `device-bluetooth-disable` |
| POST | `/devices/{device_id}/remote_desktop` | Enable Remote Desktop | `device-remote-desktop-enable` |
| DELETE | `/devices/{device_id}/remote_desktop` | Disable Remote Desktop | `device-remote-desktop-disable` |

### Lost Mode

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| POST | `/devices/{device_id}/lost_mode` | Enable Lost Mode | `device-lost-mode-enable` |
| DELETE | `/devices/{device_id}/lost_mode` | Disable Lost Mode | `device-lost-mode-disable` |
| POST | `/devices/{device_id}/lost_mode/play_sound` | Play Lost Mode sound | `device-lost-mode-play-sound` |
| POST | `/devices/{device_id}/lost_mode/update_location` | Request location update | `device-lost-mode-update-location` |

---

## Enrollments

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/enrollments` | List all enrollments | `enrollment-list` |
| GET | `/enrollments/{enrollment_id}` | Get enrollment details | `enrollment-get` |
| DELETE | `/enrollments/{enrollment_id}` | Delete an enrollment | `enrollment-delete` |
| POST | `/enrollments/{enrollment_id}/invitations` | Send enrollment invitation | `enrollment-send-invitation` |

---

## Installed Apps

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/installed_apps/{installed_app_id}` | Get installed app details | `installed-app-get` |
| DELETE | `/installed_apps/{installed_app_id}` | Delete an installed app | `installed-app-delete` |
| POST | `/installed_apps/{installed_app_id}/update` | Update an installed app | `installed-app-update` |
| POST | `/installed_apps/{installed_app_id}/request_management` | Request management of an installed app | `installed-app-request-management` |

---

## Logs

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/logs` | List all logs | `log-list` |
| GET | `/logs/{log_id}` | Get log details | `log-get` |

---

## Profiles

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/profiles` | List all profiles | `profile-list` |
| GET | `/profiles/{profile_id}` | Get profile details | `profile-get` |
| POST | `/profiles/{profile_id}/devices/{device_id}` | Assign profile to device | `profile-assign-device` |
| DELETE | `/profiles/{profile_id}/devices/{device_id}` | Unassign profile from device | `profile-unassign-device` |

---

## Push Certificate

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/push_certificate` | Get push certificate details | `push-certificate-get` |
| GET | `/push_certificate/scsr` | Get signed CSR for push certificate | `push-certificate-scsr` |

---

## Scripts

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/scripts` | List all scripts | `script-list` |
| GET | `/scripts/{script_id}` | Get script details | `script-get` |
| PATCH | `/scripts/{script_id}` | Update a script | `script-update` |
| DELETE | `/scripts/{script_id}` | Delete a script | `script-delete` |

---

## Script Jobs

| Method | Endpoint | Description | CLI Command |
|--------|----------|-------------|-------------|
| GET | `/script_jobs` | List all script jobs | `script-job-list` |
| GET | `/script_jobs/{job_id}` | Get script job details | `script-job-get` |
| POST | `/script_jobs` | Create a script job | `script-job-create` |
| DELETE | `/script_jobs/{job_id}` | Cancel a script job | `script-job-cancel` |

---

## Summary

| Category | Endpoints |
|----------|-----------|
| Account | 2 |
| Apps | 10 |
| Assignment Groups | 19 |
| Custom Attributes | 5 |
| Custom Configuration Profiles | 5 |
| Custom Declarations | 5 |
| DEP Servers | 5 |
| Device Groups | 4 |
| Devices | 35 |
| Enrollments | 4 |
| Installed Apps | 4 |
| Logs | 2 |
| Profiles | 4 |
| Push Certificate | 2 |
| Scripts | 4 |
| Script Jobs | 4 |
| **Total** | **114** |
