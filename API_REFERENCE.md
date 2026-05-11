# SimpleMDM API Reference — Source de vérité (lue directement depuis https://api.simplemdm.com/v1)

## Pagination (convention globale pour les endpoints "list" qui le mentionnent)
- limit: integer, opt, default 10, max 100
- starting_after: integer, opt
- direction: enum(asc|desc), opt, default asc

## Account
- GET /account — pas de params
- PATCH /account — name:opt, apple_store_country_code:opt | form-urlencoded

## Apps
- GET /apps — include_shared:opt(false) | Supports pagination
- GET /apps/{id} — pas de params
- POST /apps — app_store_id:opt, bundle_id:opt, binary:file:opt, name:opt | multipart (curl -F)
- PATCH /apps/{id} — binary:file:opt, name:opt, deploy_to:enum(none|outdated|all):opt(none) | multipart (curl -F)
- DELETE /apps/{id}
- GET /apps/{id}/installs — Supports pagination
- POST /apps/{id}/munki_pkginfo — file:file:req | multipart (curl -F)
- DELETE /apps/{id}/munki_pkginfo

## Managed App Configs
- GET /apps/{id}/managed_configs — pas de params
- POST /apps/{id}/managed_configs — key:req, value:(pas marqué required dans la doc), value_type:(pas marqué required dans la doc) | multipart (curl -F)
- DELETE /apps/{id}/managed_configs/{config_id}
- POST /apps/{id}/managed_configs/push

## Assignment Groups
- GET /assignment_groups — Supports pagination
- GET /assignment_groups/{id}
- POST /assignment_groups — name:req, priority:opt, auto_deploy:opt(true), type:opt(standard)[deprecated], install_type:opt(managed), app_track_location:opt(true) | form-urlencoded (curl -d)
- PATCH /assignment_groups/{id} — name:opt, priority:opt, auto_deploy:opt | form-urlencoded (curl -d)
- DELETE /assignment_groups/{id}
- POST /assignment_groups/{id}/apps/{app_id} — deployment_type:opt, install_type:opt(managed) | form-urlencoded (curl -d)
- DELETE /assignment_groups/{id}/apps/{app_id}
- POST /assignment_groups/{id}/device_groups/{dg_id} [deprecated]
- DELETE /assignment_groups/{id}/device_groups/{dg_id} [deprecated]
- POST /assignment_groups/{id}/devices/{device_id} — remove_others:opt(false) | form-urlencoded (curl -d)
- DELETE /assignment_groups/{id}/devices/{device_id}
- POST /assignment_groups/{id}/push_apps
- POST /assignment_groups/{id}/update_apps
- POST /assignment_groups/{id}/profiles/{profile_id}
- DELETE /assignment_groups/{id}/profiles/{profile_id}
- POST /assignment_groups/{id}/sync_profiles
- POST /assignment_groups/{id}/clone

## Custom Attributes
- GET /custom_attributes — PAS de pagination mentionnée (pas de "Supports pagination")
- GET /custom_attributes/{id}
- POST /custom_attributes — name:req, default_value:opt | form-urlencoded
- PATCH /custom_attributes/{id} — default_value:opt ONLY | form-urlencoded
- DELETE /custom_attributes/{id}
- GET /devices/{id}/custom_attribute_values
- PUT /devices/{id}/custom_attribute_values/{name} — value:req | form-urlencoded (curl -d)
- PUT /devices/{id}/custom_attribute_values — JSON body | application/json
- PUT /custom_attribute_values/{name} — JSON body | application/json
- GET /assignment_groups/{id}/custom_attribute_values
- PUT /assignment_groups/{id}/custom_attribute_values/{id} — value:req | form-urlencoded (curl -d)
- GET /device_groups/{id}/custom_attribute_values [deprecated]
- PUT /device_groups/{id}/custom_attribute_values/{name} — value:req [deprecated] | form-urlencoded (curl -d)

## Custom Configuration Profiles
- GET /custom_configuration_profiles — search:opt | Supports pagination
- POST /custom_configuration_profiles — name:req, mobileconfig:file:req, user_scope:boolean:opt(true), attribute_support:boolean:opt(false), escape_attributes:boolean:opt, reinstall_after_os_update:boolean:opt, declarative:boolean:opt | multipart (curl -F)
- PATCH /custom_configuration_profiles/{id} — name:opt, mobileconfig:file:opt, user_scope:boolean:opt, attribute_support:boolean:opt, escape_attributes:boolean:opt, reinstall_after_os_update:boolean:opt, declarative:boolean:opt | multipart (curl -F)
- GET /custom_configuration_profiles/{id}/download
- DELETE /custom_configuration_profiles/{id}
- POST /custom_configuration_profiles/{id}/device_groups/{dg_id} [deprecated]
- DELETE /custom_configuration_profiles/{id}/device_groups/{dg_id} [deprecated]
- POST /custom_configuration_profiles/{id}/devices/{device_id}
- DELETE /custom_configuration_profiles/{id}/devices/{device_id}

## Custom Declarations
- GET /custom_declarations — search:opt | Supports pagination
- POST /custom_declarations — name:req, declaration_type:req, payload:file:req, user_scope:boolean:opt(true), attribute_support:boolean:opt(false), escape_attributes:boolean:opt, activation_predicate:string:opt | multipart
- PATCH /custom_declarations/{id} — name:opt, declaration_type:opt, payload:file:opt, user_scope:boolean:opt, attribute_support:boolean:opt, escape_attributes:boolean:opt, activation_predicate:string:opt | multipart
- GET /custom_declarations/{id}/download
- DELETE /custom_declarations/{id}
- POST /custom_declarations/{id}/devices/{device_id}
- DELETE /custom_declarations/{id}/devices/{device_id}

## DEP Servers
- GET /dep_servers — Supports pagination
- GET /dep_servers/{id}
- POST /dep_servers/{id}/sync
- GET /dep_servers/{id}/dep_devices — Supports pagination
- GET /dep_servers/{id}/dep_devices/{dep_device_id}

## Devices
- GET /devices — search:opt, include_awaiting_enrollment:opt(false), include_secret_custom_attributes:opt(false) | Supports pagination
- GET /devices/{id} — include_secret_custom_attributes:opt(false)
- POST /devices — name:req, group_id:opt[deprecated], static_group_ids:opt(array) | form-urlencoded (curl -d)
- PATCH /devices/{id} — name:opt, device_name:opt | form-urlencoded (curl -d)
- DELETE /devices/{id}
- GET /devices/{id}/profiles — Supports pagination
- GET /devices/{id}/installed_apps — Supports pagination
- GET /devices/{id}/users — Supports pagination (macOS only)
- DELETE /devices/{id}/users/{user_id}
- POST /devices/{id}/push_apps
- POST /devices/{id}/refresh
- POST /devices/{id}/restart — rebuild_kernel_cache:boolean:opt(false), notify_user:boolean:opt(false) | form-urlencoded
- POST /devices/{id}/shutdown
- POST /devices/{id}/lock — message:opt, phone_number:opt, pin:req(macOS) | form-urlencoded (curl -d)
- POST /devices/{id}/clear_passcode
- POST /devices/{id}/clear_firmware_password
- POST /devices/{id}/rotate_firmware_password
- POST /devices/{id}/clear_recovery_lock_password
- POST /devices/{id}/clear_restrictions_password (iOS/iPad)
- POST /devices/{id}/rotate_recovery_lock_password
- POST /devices/{id}/rotate_filevault_key
- POST /devices/{id}/set_admin_password — new_password:req | form-urlencoded (curl -d)
- POST /devices/{id}/rotate_admin_password
- POST /devices/{id}/wipe — pin:opt, clear_custom_attributes:opt(false), disable_activation_lock:opt(true), preserve_data_plan:opt(false), disallow_proximity_setup:opt(false), return_to_service:opt(false), wifi_network_id:opt(integer), obliteration_behavior:enum(ObliterateWithWarning|DoNotObliterate):opt, unassign_direct_profiles:opt(false) | form-urlencoded
- POST /devices/{id}/update_os — os_update_mode:req(macOS):enum(smart_update|download_only|notify_only|install_asap|force_update), version_type:enum(latest_minor_version|latest_major_version):opt(latest_major_version) | form-urlencoded
- POST /devices/{id}/remote_desktop
- DELETE /devices/{id}/remote_desktop
- POST /devices/{id}/bluetooth
- DELETE /devices/{id}/bluetooth
- POST /devices/{id}/set_time_zone — time_zone:req | form-urlencoded (curl -d)
- POST /devices/{id}/unenroll — unassign_direct_profiles:opt(false) | form-urlencoded
- POST /devices/{id}/lost_mode — message:opt, phone_number:opt, footnote:opt | form-urlencoded (at least message or phone_number required)
- DELETE /devices/{id}/lost_mode
- POST /devices/{id}/lost_mode/play_sound
- POST /devices/{id}/lost_mode/update_location

## Device Groups [deprecated]
- GET /device_groups — Supports pagination [deprecated]
- GET /device_groups/{id} [deprecated]
- POST /device_groups/{id}/devices/{device_id} [deprecated]
- POST /device_groups/{id}/clone [deprecated]

## Enrollments
- GET /enrollments — Supports pagination
- GET /enrollments/{id}
- POST /enrollments/{id}/invitations — contact:req | form-urlencoded (curl -d)
- DELETE /enrollments/{id}

## Installed Apps
- GET /installed_apps/{id}
- POST /installed_apps/{id}/request_management
- POST /installed_apps/{id}/update
- DELETE /installed_apps/{id}

## Logs
- GET /logs — serial_number:opt | Supports pagination
- GET /logs/{id}

## Profiles
- GET /profiles — search:opt | Supports pagination
- GET /profiles/{id}
- POST /profiles/{id}/device_groups/{dg_id} [deprecated]
- DELETE /profiles/{id}/device_groups/{dg_id} [deprecated]
- POST /profiles/{id}/devices/{device_id}
- DELETE /profiles/{id}/devices/{device_id}

## Push Certificate
- GET /push_certificate
- PUT /push_certificate — file:file:req, apple_id:opt | multipart
- GET /push_certificate/scsr

## Scripts
- GET /scripts — Supports pagination
- GET /scripts/{id}
- POST /scripts — name:req, variable_support:opt(0), file:file:req | multipart (curl -F)
- PATCH /scripts/{id} — name:opt, variable_support:opt, file:file:opt | multipart
- DELETE /scripts/{id}

## Script Jobs
- GET /script_jobs — Supports pagination
- GET /script_jobs/{id}
- POST /script_jobs — script_id:req, device_ids:opt, group_ids:opt[deprecated], assignment_group_ids:opt, custom_attribute:opt, custom_attribute_regex:opt | multipart (curl -F)
- DELETE /script_jobs/{id}
