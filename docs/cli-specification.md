# simplemdm-cli -- Specification CLI

## 1. Vue d'ensemble

| Champ | Valeur |
|-------|--------|
| **Nom** | `simplemdm-cli` |
| **Objectif** | CLI non-officielle pour [SimpleMDM](https://simplemdm.com) permettant de gerer les appareils Apple (iOS/macOS/tvOS), apps, profils, scripts et plus, depuis le terminal. |
| **Pattern de commande** | `simplemdm-cli [FLAGS_GLOBAUX] RESSOURCE ACTION [FLAGS_LOCAUX] [ARGS]` |
| **API ciblee** | SimpleMDM REST API v1 -- `https://a.simplemdm.com/api/v1` |
| **Authentification** | HTTP Basic Auth (API key comme username, password vide) |
| **Inspirations** | `tailscale-cli`, `gh` (GitHub CLI), `kubectl` |
| **Langage** | Go 1.25+ |
| **Framework CLI** | [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper) |
| **Licence** | Proprietaire |
| **Repository** | `github.com/dimer47/simplemdm-cli` |

---

## 2. Flags globaux

| Flag | Short | Variable d'env | Description | Default |
|------|-------|----------------|-------------|---------|
| `--api-key` | `-k` | `SMDM_API_KEY` | Cle API SimpleMDM | `""` |
| `--output` | `-o` | `SMDM_OUTPUT` | Format de sortie : `table`, `json`, `yaml`, `csv` | `table` |
| `--json` | | | Raccourci pour `--output json` | `false` |
| `--quiet` | `-q` | | Supprime les sorties non-essentielles | `false` |
| `--debug` | | `SMDM_DEBUG` | Mode debug (affiche version, commit, date) | `false` |
| `--no-color` | | `NO_COLOR` | Desactive la couleur | `false` |
| `--context` | `-c` | `SMDM_CONTEXT` | Contexte d'authentification actif | `""` |

**Priorite de resolution de la cle API :**
1. Flag `--api-key` / variable d'env `SMDM_API_KEY`
2. Keychain systeme (via le contexte actif)

---

## 3. Configuration

Le fichier de configuration est stocke dans `~/.simplemdm-cli/config.json`.

```json
{
  "default_context": "production",
  "contexts": {
    "production": {
      "name": "production"
    },
    "staging": {
      "name": "staging"
    }
  }
}
```

Les cles API sont stockees dans le **keychain systeme** (macOS Keychain, Linux Secret Service, Windows Credential Manager) via la librairie `go-keyring`. Elles ne sont jamais ecrites en clair sur disque.

---

## 4. Commandes d'administration

### 4.1 auth -- Gestion de l'authentification

| Commande | Description |
|----------|-------------|
| `auth login` | Ajouter un nouveau contexte d'authentification (interactif : nom + cle API) |
| `auth status` | Afficher le contexte actif, la cle masquee et valider aupres de l'API |
| `auth switch <context>` | Basculer vers un autre contexte |
| `auth list` | Lister tous les contextes configures (`*` = actif) |
| `auth remove <context>` | Supprimer un contexte et sa cle API du keychain |

### 4.2 completion -- Autocompletion shell

| Commande | Description |
|----------|-------------|
| `completion bash` | Generer le script d'autocompletion pour Bash |
| `completion zsh` | Generer le script d'autocompletion pour Zsh |
| `completion fish` | Generer le script d'autocompletion pour Fish |
| `completion powershell` | Generer le script d'autocompletion pour PowerShell |

**Exemple d'installation (zsh) :**
```bash
simplemdm-cli completion zsh > "${fpath[1]}/_simplemdm-cli"
```

### 4.3 version

| Commande | Description |
|----------|-------------|
| `version` | Affiche la version, le commit et la date de build |

**Sortie :** `simplemdm-cli v1.2.3 (abc1234) built 2026-05-10T12:00:00Z`

### 4.4 self-update

| Commande | Description |
|----------|-------------|
| `self-update` | Verifie les mises a jour sur GitHub Releases et installe la derniere version |

Le CLI effectue egalement une verification automatique en arriere-plan (non-bloquante, timeout 1s) apres chaque commande et affiche un message si une nouvelle version est disponible.

---

## 5. Commandes par ressource

### 5.1 account -- Gestion du compte

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `account get` | `GET /account` | Obtenir les details du compte |
| `account update` | `PATCH /account` | Mettre a jour le compte |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--name` | Nom du compte | `update` |
| `--apple-store-country-code` | Code pays Apple Store | `update` |

---

### 5.2 app -- Gestion des applications

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `app list` | `GET /apps` | Lister toutes les apps |
| `app get <app-id>` | `GET /apps/{id}` | Obtenir les details d'une app |
| `app create` | `POST /apps` | Creer une app (App Store, enterprise ou custom) |
| `app update <app-id>` | `PATCH /apps/{id}` | Mettre a jour une app |
| `app delete <app-id>` | `DELETE /apps/{id}` | Supprimer une app |
| `app installs <app-id>` | `GET /apps/{id}/installs` | Lister les installations d'une app |
| `app managed-configs <app-id>` | `GET /apps/{id}/managed_configs` | Lister les configurations gerees |
| `app managed-config-create <app-id>` | `POST /apps/{id}/managed_configs` | Creer une configuration geree |
| `app managed-configs-push <app-id>` | `POST /apps/{id}/managed_configs/push` | Pousser les configs gerees sur les appareils |
| `app managed-config-delete <app-id> <config-id>` | `DELETE /apps/{id}/managed_configs/{cid}` | Supprimer une configuration geree |
| `app munki-pkginfo-update <app-id>` | `POST /apps/{id}/munki_pkginfo` | Uploader un fichier Munki pkginfo |
| `app munki-pkginfo-delete <app-id>` | `DELETE /apps/{id}/munki_pkginfo` | Supprimer le Munki pkginfo |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--direction` | Tri : `asc` ou `desc` | `list` |
| `--include-shared` | Inclure les apps partagees : `true` / `false` | `list` |
| `--name` | Nom de l'app | `create`, `update` |
| `--app-store-id` | App Store ID | `create` |
| `--bundle-id` | Bundle ID | `create` |
| `--binary` | Chemin vers le fichier binaire de l'app | `create`, `update` |
| `--deploy-to` | Deployer vers : `none`, `outdated`, `all` | `update` |
| `--key` | Cle de configuration (requis) | `managed-config-create` |
| `--value` | Valeur de configuration (requis) | `managed-config-create` |
| `--value-type` | Type de valeur : `string`, `integer`, `boolean` | `managed-config-create` |
| `--file` | Chemin vers le fichier pkginfo (requis) | `munki-pkginfo-update` |

---

### 5.3 assignment-group (alias: ag) -- Groupes d'assignation

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `assignment-group list` | `GET /assignment_groups` | Lister tous les groupes d'assignation |
| `assignment-group get <id>` | `GET /assignment_groups/{id}` | Obtenir les details d'un groupe |
| `assignment-group create` | `POST /assignment_groups` | Creer un groupe d'assignation |
| `assignment-group update <id>` | `PATCH /assignment_groups/{id}` | Mettre a jour un groupe |
| `assignment-group delete <id>` | `DELETE /assignment_groups/{id}` | Supprimer un groupe |
| `assignment-group assign-app <gid> <aid>` | `POST /assignment_groups/{gid}/apps/{aid}` | Assigner une app au groupe |
| `assignment-group unassign-app <gid> <aid>` | `DELETE /assignment_groups/{gid}/apps/{aid}` | Retirer une app du groupe |
| `assignment-group assign-device <gid> <did>` | `POST /assignment_groups/{gid}/devices/{did}` | Assigner un appareil au groupe |
| `assignment-group unassign-device <gid> <did>` | `DELETE /assignment_groups/{gid}/devices/{did}` | Retirer un appareil du groupe |
| `assignment-group assign-device-group <gid> <dgid>` | `POST /assignment_groups/{gid}/device_groups/{dgid}` | Assigner un device group |
| `assignment-group unassign-device-group <gid> <dgid>` | `DELETE /assignment_groups/{gid}/device_groups/{dgid}` | Retirer un device group |
| `assignment-group assign-profile <gid> <pid>` | `POST /assignment_groups/{gid}/profiles/{pid}` | Assigner un profil au groupe |
| `assignment-group unassign-profile <gid> <pid>` | `DELETE /assignment_groups/{gid}/profiles/{pid}` | Retirer un profil du groupe |
| `assignment-group push-apps <id>` | `POST /assignment_groups/{id}/push_apps` | Pousser les apps sur tous les appareils |
| `assignment-group update-apps <id>` | `POST /assignment_groups/{id}/update_apps` | Mettre a jour les apps sur tous les appareils |
| `assignment-group sync-profiles <id>` | `POST /assignment_groups/{id}/sync_profiles` | Synchroniser les profils du groupe |
| `assignment-group clone <id>` | `POST /assignment_groups/{id}/clone` | Cloner un groupe d'assignation |
| `assignment-group custom-attributes <id>` | `GET /assignment_groups/{id}/custom_attribute_values` | Obtenir les attributs personnalises |
| `assignment-group set-custom-attribute <gid> <name>` | `PUT /assignment_groups/{gid}/custom_attribute_values/{name}` | Definir un attribut personnalise |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--name` | Nom du groupe (requis a la creation) | `create`, `update` |
| `--auto-deploy` | Deploiement automatique des apps | `create`, `update` |
| `--priority` | Priorite du groupe | `create`, `update` |
| `--type` | Type de groupe | `create` |
| `--install-type` | Type d'installation | `create` |
| `--app-track-location` | Emplacement du track d'app | `create` |
| `--remove-others` | Retirer l'appareil des autres groupes | `assign-device` |
| `--value` | Valeur de l'attribut (requis) | `set-custom-attribute` |

---

### 5.4 custom-attribute (alias: ca) -- Attributs personnalises

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `custom-attribute list` | `GET /custom_attributes` | Lister tous les attributs |
| `custom-attribute get <id>` | `GET /custom_attributes/{id}` | Obtenir les details d'un attribut |
| `custom-attribute create` | `POST /custom_attributes` | Creer un attribut |
| `custom-attribute update <id>` | `PATCH /custom_attributes/{id}` | Mettre a jour un attribut |
| `custom-attribute delete <id>` | `DELETE /custom_attributes/{id}` | Supprimer un attribut |
| `custom-attribute set-value <name>` | `PUT /custom_attribute_values/{name}` | Definir la valeur globale d'un attribut |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--name` | Nom de l'attribut (requis) | `create` |
| `--default-value` | Valeur par defaut | `create`, `update` |
| `--value` | Valeur (requis) | `set-value` |

---

### 5.5 custom-configuration-profile (alias: ccp) -- Profils de configuration personnalises

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `custom-configuration-profile list` | `GET /custom_configuration_profiles` | Lister tous les profils |
| `custom-configuration-profile create` | `POST /custom_configuration_profiles` | Creer un profil (upload .mobileconfig) |
| `custom-configuration-profile update <id>` | `PATCH /custom_configuration_profiles/{id}` | Mettre a jour un profil |
| `custom-configuration-profile delete <id>` | `DELETE /custom_configuration_profiles/{id}` | Supprimer un profil |
| `custom-configuration-profile download <id>` | `GET /custom_configuration_profiles/{id}/download` | Telecharger le .mobileconfig |
| `custom-configuration-profile push-to-device <pid> <did>` | `POST /custom_configuration_profiles/{pid}/devices/{did}` | Pousser le profil sur un appareil |
| `custom-configuration-profile remove-from-device <pid> <did>` | `DELETE /custom_configuration_profiles/{pid}/devices/{did}` | Retirer le profil d'un appareil |
| `custom-configuration-profile assign-device-group <pid> <dgid>` | `POST /custom_configuration_profiles/{pid}/device_groups/{dgid}` | Assigner a un device group |
| `custom-configuration-profile unassign-device-group <pid> <dgid>` | `DELETE /custom_configuration_profiles/{pid}/device_groups/{dgid}` | Retirer d'un device group |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--search` | Recherche par nom | `list` |
| `--name` | Nom du profil (requis a la creation) | `create`, `update` |
| `--mobileconfig` | Chemin vers le fichier .mobileconfig (requis a la creation) | `create`, `update` |
| `--user-scope` | Scope utilisateur | `create`, `update` |
| `--attribute-support` | Activer le support des attributs | `create`, `update` |
| `--escape-attributes` | Echapper les attributs | `create`, `update` |
| `--reinstall-after-os-update` | Reinstaller apres mise a jour OS | `create`, `update` |
| `--declarative` | Profil declaratif | `create`, `update` |
| `-o, --output` | Chemin du fichier de sortie | `download` |

---

### 5.6 custom-declaration (alias: cd) -- Declarations personnalisees (DDM)

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `custom-declaration list` | `GET /custom_declarations` | Lister toutes les declarations |
| `custom-declaration create` | `POST /custom_declarations` | Creer une declaration (upload payload) |
| `custom-declaration update <id>` | `PATCH /custom_declarations/{id}` | Mettre a jour une declaration |
| `custom-declaration delete <id>` | `DELETE /custom_declarations/{id}` | Supprimer une declaration |
| `custom-declaration download <id>` | `GET /custom_declarations/{id}/download` | Telecharger le payload |
| `custom-declaration push-to-device <did> <devid>` | `POST /custom_declarations/{did}/devices/{devid}` | Pousser sur un appareil |
| `custom-declaration remove-from-device <did> <devid>` | `DELETE /custom_declarations/{did}/devices/{devid}` | Retirer d'un appareil |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--search` | Recherche par nom | `list` |
| `--name` | Nom de la declaration (requis) | `create`, `update` |
| `--declaration-type` | Type de declaration (requis a la creation) | `create`, `update` |
| `--payload` | Chemin vers le fichier payload (requis a la creation) | `create`, `update` |
| `--user-scope` | Scope utilisateur | `create`, `update` |
| `--attribute-support` | Activer le support des attributs | `create`, `update` |
| `--escape-attributes` | Echapper les attributs | `create`, `update` |
| `--activation-predicate` | Predicat d'activation | `create`, `update` |
| `-o, --output` | Chemin du fichier de sortie | `download` |

---

### 5.7 dep-server (alias: dep) -- Serveurs DEP (Device Enrollment Program)

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `dep-server list` | `GET /dep_servers` | Lister tous les serveurs DEP |
| `dep-server get <id>` | `GET /dep_servers/{id}` | Obtenir les details d'un serveur |
| `dep-server devices <id>` | `GET /dep_servers/{id}/dep_devices` | Lister les appareils DEP du serveur |
| `dep-server device-get <sid> <did>` | `GET /dep_servers/{sid}/dep_devices/{did}` | Obtenir les details d'un appareil DEP |
| `dep-server sync <id>` | `POST /dep_servers/{id}/sync` | Synchroniser le serveur DEP |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list`, `devices` |
| `--starting-after` | Curseur de pagination | `list`, `devices` |

---

### 5.8 device -- Gestion des appareils

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `device list` | `GET /devices` | Lister tous les appareils |
| `device get <id>` | `GET /devices/{id}` | Obtenir les details d'un appareil |
| `device create` | `POST /devices` | Creer un appareil |
| `device update <id>` | `PATCH /devices/{id}` | Mettre a jour un appareil |
| `device delete <id>` | `DELETE /devices/{id}` | Supprimer un appareil |
| `device installed-apps <id>` | `GET /devices/{id}/installed_apps` | Lister les apps installees |
| `device push-apps <id>` | `POST /devices/{id}/push_apps` | Pousser les apps assignees |
| `device refresh <id>` | `POST /devices/{id}/refresh` | Rafraichir les informations |
| `device lock <id>` | `POST /devices/{id}/lock` | Verrouiller l'appareil |
| `device wipe <id>` | `POST /devices/{id}/wipe` | Effacer l'appareil |
| `device restart <id>` | `POST /devices/{id}/restart` | Redemarrer l'appareil |
| `device shutdown <id>` | `POST /devices/{id}/shutdown` | Eteindre l'appareil |
| `device clear-passcode <id>` | `POST /devices/{id}/clear_passcode` | Supprimer le code de deverrouillage |
| `device clear-firmware-password <id>` | `POST /devices/{id}/clear_firmware_password` | Supprimer le mot de passe firmware |
| `device update-os <id>` | `POST /devices/{id}/update_os` | Lancer une mise a jour OS |
| `device bluetooth-enable <id>` | `POST /devices/{id}/bluetooth` | Activer le Bluetooth |
| `device bluetooth-disable <id>` | `DELETE /devices/{id}/bluetooth` | Desactiver le Bluetooth |
| `device remote-desktop-enable <id>` | `POST /devices/{id}/remote_desktop` | Activer Remote Desktop |
| `device remote-desktop-disable <id>` | `DELETE /devices/{id}/remote_desktop` | Desactiver Remote Desktop |
| `device rotate-firmware-password <id>` | `POST /devices/{id}/rotate_firmware_password` | Rotation du mot de passe firmware |
| `device clear-recovery-lock <id>` | `POST /devices/{id}/clear_recovery_lock_password` | Supprimer le recovery lock |
| `device rotate-recovery-lock <id>` | `POST /devices/{id}/rotate_recovery_lock_password` | Rotation du recovery lock |
| `device rotate-filevault-key <id>` | `POST /devices/{id}/rotate_filevault_key` | Rotation de la cle FileVault |
| `device set-admin-password <id>` | `POST /devices/{id}/set_admin_password` | Definir le mot de passe admin |
| `device rotate-admin-password <id>` | `POST /devices/{id}/rotate_admin_password` | Rotation du mot de passe admin |
| `device clear-restrictions-password <id>` | `POST /devices/{id}/clear_restrictions_password` | Supprimer le mot de passe restrictions |
| `device profiles <id>` | `GET /devices/{id}/profiles` | Lister les profils de l'appareil |
| `device users <id>` | `GET /devices/{id}/users` | Lister les utilisateurs de l'appareil |
| `device delete-user <did> <uid>` | `DELETE /devices/{did}/users/{uid}` | Supprimer un utilisateur de l'appareil |
| `device set-timezone <id>` | `POST /devices/{id}/set_time_zone` | Definir le fuseau horaire |
| `device unenroll <id>` | `POST /devices/{id}/unenroll` | Desenroler l'appareil |
| `device custom-attributes <id>` | `GET /devices/{id}/custom_attribute_values` | Lister les attributs personnalises |
| `device set-custom-attribute <did> <name>` | `PUT /devices/{did}/custom_attribute_values/{name}` | Definir un attribut personnalise |
| `device set-custom-attributes <id>` | `PUT /devices/{id}/custom_attribute_values` | Definir plusieurs attributs (JSON) |
| `device lost-mode-enable <id>` | `POST /devices/{id}/lost_mode` | Activer le mode perdu |
| `device lost-mode-disable <id>` | `DELETE /devices/{id}/lost_mode` | Desactiver le mode perdu |
| `device lost-mode-play-sound <id>` | `POST /devices/{id}/lost_mode/play_sound` | Jouer le son du mode perdu |
| `device lost-mode-update-location <id>` | `POST /devices/{id}/lost_mode/update_location` | Demander une mise a jour de localisation |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--search` | Recherche textuelle | `list` |
| `--include-awaiting-enrollment` | Inclure les appareils en attente : `true` / `false` | `list` |
| `--include-secret-custom-attributes` | Inclure les attributs secrets : `true` / `false` | `list`, `get` |
| `--name` | Nom SimpleMDM de l'appareil | `create`, `update` |
| `--group-id` | ID du groupe (deprecie, utiliser `--static-group-ids`) | `create` |
| `--static-group-ids` | IDs de groupes statiques separes par virgule | `create` |
| `--device-name` | Nom de l'appareil (necessite supervision, async) | `update` |
| `--message` | Message a afficher | `lock`, `lost-mode-enable` |
| `--phone-number` | Numero de telephone a afficher | `lock`, `lost-mode-enable` |
| `--pin` | Code PIN | `lock`, `wipe` |
| `--clear-custom-attributes` | Effacer les attributs personnalises | `wipe` |
| `--disable-activation-lock` | Desactiver le verrouillage d'activation | `wipe` |
| `--preserve-data-plan` | Preserver le forfait data (iOS) | `wipe` |
| `--disallow-proximity-setup` | Interdire le proximity setup (iOS) | `wipe` |
| `--return-to-service` | Retour en service (iOS 17+/tvOS 18+) | `wipe` |
| `--wifi-network-id` | ID du reseau WiFi pour retour en service | `wipe` |
| `--obliteration-behavior` | Comportement d'obliteration (macOS 12+) | `wipe` |
| `--unassign-direct-profiles` | Retirer les profils directs | `wipe`, `unenroll` |
| `--rebuild-kernel-cache` | Reconstruire le kernel cache (macOS 11+) | `restart` |
| `--notify-user` | Notifier l'utilisateur (macOS 11.3+) | `restart` |
| `--os-update-mode` | Mode de mise a jour OS | `update-os` |
| `--version-type` | Type de version | `update-os` |
| `--version` | Version cible | `update-os` |
| `--new-password` | Nouveau mot de passe admin (requis) | `set-admin-password` |
| `--timezone` | Fuseau horaire TZ (requis) | `set-timezone` |
| `--value` | Valeur de l'attribut (requis) | `set-custom-attribute` |
| `--data` | Donnees JSON cle-valeur (requis) | `set-custom-attributes` |
| `--footnote` | Note de bas de page | `lost-mode-enable` |

---

### 5.9 device-group (alias: dg) -- Groupes d'appareils

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `device-group list` | `GET /device_groups` | Lister tous les groupes |
| `device-group get <id>` | `GET /device_groups/{id}` | Obtenir les details d'un groupe |
| `device-group assign-device <gid> <did>` | `POST /device_groups/{gid}/devices/{did}` | Assigner un appareil |
| `device-group clone <id>` | `POST /device_groups/{id}/clone` | Cloner un groupe |
| `device-group custom-attributes <id>` | `GET /device_groups/{id}/custom_attribute_values` | Obtenir les attributs personnalises |
| `device-group set-custom-attribute <gid> <name>` | `PUT /device_groups/{gid}/custom_attribute_values/{name}` | Definir un attribut |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--value` | Valeur de l'attribut (requis) | `set-custom-attribute` |

---

### 5.10 enrollment -- Enrollements

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `enrollment list` | `GET /enrollments` | Lister tous les enrollements |
| `enrollment get <id>` | `GET /enrollments/{id}` | Obtenir les details |
| `enrollment delete <id>` | `DELETE /enrollments/{id}` | Supprimer un enrollement |
| `enrollment send-invitation <id>` | `POST /enrollments/{id}/invitations` | Envoyer une invitation |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--contact` | Email ou numero de telephone (requis) | `send-invitation` |

---

### 5.11 installed-app -- Applications installees

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `installed-app get <id>` | `GET /installed_apps/{id}` | Obtenir les details d'une app installee |
| `installed-app delete <id>` | `DELETE /installed_apps/{id}` | Supprimer une app installee |
| `installed-app update <id>` | `POST /installed_apps/{id}/update` | Mettre a jour une app installee |
| `installed-app request-management <id>` | `POST /installed_apps/{id}/request_management` | Demander la gestion d'une app |

---

### 5.12 log -- Journaux d'evenements

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `log list` | `GET /logs` | Lister tous les logs |
| `log get <id>` | `GET /logs/{id}` | Obtenir les details d'un log |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--serial-number` | Filtrer par numero de serie | `list` |

---

### 5.13 profile -- Profils

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `profile list` | `GET /profiles` | Lister tous les profils |
| `profile get <id>` | `GET /profiles/{id}` | Obtenir les details d'un profil |
| `profile assign-device <pid> <did>` | `POST /profiles/{pid}/devices/{did}` | Assigner un profil a un appareil |
| `profile unassign-device <pid> <did>` | `DELETE /profiles/{pid}/devices/{did}` | Retirer un profil d'un appareil |
| `profile assign-device-group <pid> <dgid>` | `POST /profiles/{pid}/device_groups/{dgid}` | Assigner a un device group |
| `profile unassign-device-group <pid> <dgid>` | `DELETE /profiles/{pid}/device_groups/{dgid}` | Retirer d'un device group |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--search` | Recherche par nom ou type | `list` |

---

### 5.14 push-certificate -- Certificat push

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `push-certificate get` | `GET /push_certificate` | Obtenir les details du certificat push |
| `push-certificate scsr` | `GET /push_certificate/scsr` | Obtenir le CSR signe |
| `push-certificate update` | `PUT /push_certificate` | Mettre a jour le certificat push |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--file` | Chemin vers le certificat push Apple (requis) | `update` |
| `--apple-id` | Apple ID associe | `update` |

---

### 5.15 script -- Scripts

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `script list` | `GET /scripts` | Lister tous les scripts |
| `script get <id>` | `GET /scripts/{id}` | Obtenir les details d'un script |
| `script create` | `POST /scripts` | Creer un script (upload fichier) |
| `script update <id>` | `PATCH /scripts/{id}` | Mettre a jour un script |
| `script delete <id>` | `DELETE /scripts/{id}` | Supprimer un script |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--name` | Nom du script (requis a la creation) | `create`, `update` |
| `--file` | Chemin vers le fichier script (requis a la creation) | `create`, `update` |
| `--variable-support` | Activer le support des variables | `create`, `update` |

---

### 5.16 script-job -- Taches de script

| Commande | API Endpoint | Description |
|----------|-------------|-------------|
| `script-job list` | `GET /script_jobs` | Lister toutes les taches |
| `script-job get <id>` | `GET /script_jobs/{id}` | Obtenir les details d'une tache |
| `script-job create` | `POST /script_jobs` | Creer une tache de script |
| `script-job cancel <id>` | `DELETE /script_jobs/{id}` | Annuler une tache |

**Flags specifiques :**

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--limit` | Limiter le nombre de resultats | `list` |
| `--starting-after` | Curseur de pagination | `list` |
| `--script-id` | ID du script (requis) | `create` |
| `--device-ids` | IDs d'appareils separes par virgule | `create` |
| `--assignment-group-ids` | IDs de groupes d'assignation separes par virgule | `create` |
| `--group-ids` | IDs de device groups separes par virgule (deprecie) | `create` |
| `--custom-attribute` | Filtre par attribut personnalise | `create` |
| `--custom-attribute-regex` | Filtre par regex d'attribut personnalise | `create` |

---

## 6. Arborescence complete des commandes

```
simplemdm-cli
|
|-- auth
|   |-- login
|   |-- status
|   |-- switch <context>
|   |-- list
|   +-- remove <context>
|
|-- account
|   |-- get
|   +-- update
|
|-- app
|   |-- list
|   |-- get <app-id>
|   |-- create
|   |-- update <app-id>
|   |-- delete <app-id>
|   |-- installs <app-id>
|   |-- managed-configs <app-id>
|   |-- managed-config-create <app-id>
|   |-- managed-configs-push <app-id>
|   |-- managed-config-delete <app-id> <config-id>
|   |-- munki-pkginfo-update <app-id>
|   +-- munki-pkginfo-delete <app-id>
|
|-- assignment-group  [alias: ag]
|   |-- list
|   |-- get <group-id>
|   |-- create
|   |-- update <group-id>
|   |-- delete <group-id>
|   |-- assign-app <group-id> <app-id>
|   |-- unassign-app <group-id> <app-id>
|   |-- assign-device <group-id> <device-id>
|   |-- unassign-device <group-id> <device-id>
|   |-- assign-device-group <group-id> <device-group-id>
|   |-- unassign-device-group <group-id> <device-group-id>
|   |-- assign-profile <group-id> <profile-id>
|   |-- unassign-profile <group-id> <profile-id>
|   |-- push-apps <group-id>
|   |-- update-apps <group-id>
|   |-- sync-profiles <group-id>
|   |-- clone <group-id>
|   |-- custom-attributes <group-id>
|   +-- set-custom-attribute <group-id> <attribute-name>
|
|-- custom-attribute  [alias: ca]
|   |-- list
|   |-- get <attribute-id>
|   |-- create
|   |-- update <attribute-id>
|   |-- delete <attribute-id>
|   +-- set-value <attribute-name>
|
|-- custom-configuration-profile  [alias: ccp]
|   |-- list
|   |-- create
|   |-- update <profile-id>
|   |-- delete <profile-id>
|   |-- download <profile-id>
|   |-- push-to-device <profile-id> <device-id>
|   |-- remove-from-device <profile-id> <device-id>
|   |-- assign-device-group <profile-id> <device-group-id>
|   +-- unassign-device-group <profile-id> <device-group-id>
|
|-- custom-declaration  [alias: cd]
|   |-- list
|   |-- create
|   |-- update <declaration-id>
|   |-- delete <declaration-id>
|   |-- download <declaration-id>
|   |-- push-to-device <declaration-id> <device-id>
|   +-- remove-from-device <declaration-id> <device-id>
|
|-- dep-server  [alias: dep]
|   |-- list
|   |-- get <server-id>
|   |-- devices <server-id>
|   |-- device-get <server-id> <dep-device-id>
|   +-- sync <server-id>
|
|-- device
|   |-- list
|   |-- get <device-id>
|   |-- create
|   |-- update <device-id>
|   |-- delete <device-id>
|   |-- installed-apps <device-id>
|   |-- push-apps <device-id>
|   |-- refresh <device-id>
|   |-- lock <device-id>
|   |-- wipe <device-id>
|   |-- restart <device-id>
|   |-- shutdown <device-id>
|   |-- clear-passcode <device-id>
|   |-- clear-firmware-password <device-id>
|   |-- update-os <device-id>
|   |-- bluetooth-enable <device-id>
|   |-- bluetooth-disable <device-id>
|   |-- remote-desktop-enable <device-id>
|   |-- remote-desktop-disable <device-id>
|   |-- rotate-firmware-password <device-id>
|   |-- clear-recovery-lock <device-id>
|   |-- rotate-recovery-lock <device-id>
|   |-- rotate-filevault-key <device-id>
|   |-- set-admin-password <device-id>
|   |-- rotate-admin-password <device-id>
|   |-- clear-restrictions-password <device-id>
|   |-- profiles <device-id>
|   |-- users <device-id>
|   |-- delete-user <device-id> <user-id>
|   |-- set-timezone <device-id>
|   |-- unenroll <device-id>
|   |-- custom-attributes <device-id>
|   |-- set-custom-attribute <device-id> <attribute-name>
|   |-- set-custom-attributes <device-id>
|   |-- lost-mode-enable <device-id>
|   |-- lost-mode-disable <device-id>
|   |-- lost-mode-play-sound <device-id>
|   +-- lost-mode-update-location <device-id>
|
|-- device-group  [alias: dg]
|   |-- list
|   |-- get <group-id>
|   |-- assign-device <group-id> <device-id>
|   |-- clone <group-id>
|   |-- custom-attributes <group-id>
|   +-- set-custom-attribute <group-id> <attribute-name>
|
|-- enrollment
|   |-- list
|   |-- get <enrollment-id>
|   |-- delete <enrollment-id>
|   +-- send-invitation <enrollment-id>
|
|-- installed-app
|   |-- get <installed-app-id>
|   |-- delete <installed-app-id>
|   |-- update <installed-app-id>
|   +-- request-management <installed-app-id>
|
|-- log
|   |-- list
|   +-- get <log-id>
|
|-- profile
|   |-- list
|   |-- get <profile-id>
|   |-- assign-device <profile-id> <device-id>
|   |-- unassign-device <profile-id> <device-id>
|   |-- assign-device-group <profile-id> <device-group-id>
|   +-- unassign-device-group <profile-id> <device-group-id>
|
|-- push-certificate
|   |-- get
|   |-- scsr
|   +-- update
|
|-- script
|   |-- list
|   |-- get <script-id>
|   |-- create
|   |-- update <script-id>
|   +-- delete <script-id>
|
|-- script-job
|   |-- list
|   |-- get <job-id>
|   |-- create
|   +-- cancel <job-id>
|
|-- completion
|   |-- bash
|   |-- zsh
|   |-- fish
|   +-- powershell
|
|-- version
|-- self-update
+-- mcp-serve  [hidden]
```

---

## 7. Exemples d'utilisation

### Authentification

```bash
# Configurer un contexte
simplemdm-cli auth login

# Verifier le statut
simplemdm-cli auth status

# Basculer entre contextes
simplemdm-cli auth switch production

# Utiliser une cle API ponctuelle
simplemdm-cli --api-key "YOUR_KEY" device list

# Via variable d'environnement
export SMDM_API_KEY="YOUR_KEY"
simplemdm-cli device list
```

### Gestion des appareils

```bash
# Lister les appareils avec recherche
simplemdm-cli device list --search "MacBook" --limit 10

# Details d'un appareil en JSON
simplemdm-cli device get 123 --json

# Verrouiller un appareil avec un message
simplemdm-cli device lock 123 --message "Appareil perdu" --phone-number "+33600000000"

# Effacer un appareil avec options avancees
simplemdm-cli device wipe 123 --pin "123456" --disable-activation-lock --return-to-service

# Mettre a jour l'OS
simplemdm-cli device update-os 123 --os-update-mode "force_update" --version "18.0"

# Mode perdu
simplemdm-cli device lost-mode-enable 123 --message "Veuillez retourner cet appareil"
simplemdm-cli device lost-mode-play-sound 123

# Gestion des attributs personnalises
simplemdm-cli device custom-attributes 123
simplemdm-cli device set-custom-attribute 123 department --value "Engineering"
```

### Gestion des applications

```bash
# Lister les apps
simplemdm-cli app list --json

# Creer une app App Store
simplemdm-cli app create --name "Slack" --app-store-id "618783545"

# Uploader une app enterprise
simplemdm-cli app create --name "MonApp" --binary ./MonApp.ipa

# Gerer les managed configs
simplemdm-cli app managed-configs 456
simplemdm-cli app managed-config-create 456 --key "server_url" --value "https://api.example.com"
simplemdm-cli app managed-configs-push 456
```

### Groupes d'assignation

```bash
# Creer un groupe avec auto-deploy
simplemdm-cli ag create --name "Equipe Dev" --auto-deploy

# Assigner un appareil et une app
simplemdm-cli ag assign-device 10 123
simplemdm-cli ag assign-app 10 456

# Pousser les apps
simplemdm-cli ag push-apps 10

# Cloner un groupe
simplemdm-cli ag clone 10
```

### Profils de configuration personnalises

```bash
# Creer un profil
simplemdm-cli ccp create --name "WiFi Corp" --mobileconfig ./wifi.mobileconfig

# Pousser sur un appareil specifique
simplemdm-cli ccp push-to-device 789 123

# Telecharger un profil
simplemdm-cli ccp download 789 --output ./backup.mobileconfig
```

### Scripts et taches

```bash
# Creer un script
simplemdm-cli script create --name "Inventaire" --file ./inventaire.sh --variable-support

# Lancer une tache sur plusieurs appareils
simplemdm-cli script-job create --script-id 42 --device-ids "123,456,789"

# Lancer une tache sur un groupe d'assignation
simplemdm-cli script-job create --script-id 42 --assignment-group-ids "10,20"

# Annuler une tache
simplemdm-cli script-job cancel 99
```

### Formats de sortie

```bash
# Table (defaut)
simplemdm-cli device list

# JSON
simplemdm-cli device list --json
simplemdm-cli device list -o json

# YAML
simplemdm-cli device list -o yaml

# CSV
simplemdm-cli device list -o csv
```

---

## 8. Integration MCP (Model Context Protocol)

Le CLI embarque un serveur MCP complet, permettant aux LLM (Claude, etc.) d'interagir avec SimpleMDM via le protocole JSON-RPC sur stdio.

### Demarrage

```bash
simplemdm-cli mcp-serve
```

La commande `mcp-serve` est cachee (non affichee dans l'aide). Elle demarre un serveur MCP qui expose toutes les operations de l'API SimpleMDM en tant que tools.

### Configuration Claude Desktop

```json
{
  "mcpServers": {
    "simplemdm": {
      "command": "simplemdm-cli",
      "args": ["mcp-serve"],
      "env": {
        "SMDM_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

### Tools MCP disponibles

Le serveur MCP expose les tools suivants (nomenclature `ressource-action`) :

| Categorie | Tools |
|-----------|-------|
| **Account** | `account-get`, `account-update` |
| **Apps** | `app-list`, `app-get`, `app-create`, `app-update`, `app-delete`, `app-installs`, `app-managed-configs`, `app-managed-config-create`, `app-managed-configs-push`, `app-managed-config-delete` |
| **Assignment Groups** | `assignment-group-list`, `assignment-group-get`, `assignment-group-create`, `assignment-group-update`, `assignment-group-delete`, `assignment-group-assign-app`, `assignment-group-unassign-app`, `assignment-group-assign-device`, `assignment-group-unassign-device`, `assignment-group-assign-device-group`, `assignment-group-unassign-device-group`, `assignment-group-assign-profile`, `assignment-group-unassign-profile`, `assignment-group-push-apps`, `assignment-group-update-apps`, `assignment-group-sync-profiles`, `assignment-group-clone`, `assignment-group-custom-attributes`, `assignment-group-set-custom-attribute` |
| **Custom Attributes** | `custom-attribute-list`, `custom-attribute-get`, `custom-attribute-create`, `custom-attribute-update`, `custom-attribute-delete` |
| **Custom Config Profiles** | `custom-configuration-profile-list`, `custom-configuration-profile-update`, `custom-configuration-profile-delete`, `custom-configuration-profile-push-device`, `custom-configuration-profile-remove-device` |
| **Custom Declarations** | `custom-declaration-list`, `custom-declaration-update`, `custom-declaration-delete`, `custom-declaration-push-device`, `custom-declaration-remove-device` |
| **DEP Servers** | `dep-server-list`, `dep-server-get`, `dep-server-devices`, `dep-server-device-get`, `dep-server-sync` |
| **Device Groups** | `device-group-list`, `device-group-get`, `device-group-assign-device`, `device-group-clone` |
| **Devices** | `device-list`, `device-get`, `device-create`, `device-update`, `device-delete`, `device-installed-apps`, `device-push-apps`, `device-refresh`, `device-lock`, `device-wipe`, `device-restart`, `device-shutdown`, `device-clear-passcode`, `device-clear-firmware-password`, `device-rotate-firmware-password`, `device-clear-recovery-lock`, `device-rotate-recovery-lock`, `device-rotate-filevault-key`, `device-set-admin-password`, `device-rotate-admin-password`, `device-clear-restrictions-password`, `device-bluetooth-enable`, `device-bluetooth-disable`, `device-remote-desktop-enable`, `device-remote-desktop-disable`, `device-set-timezone`, `device-delete-user`, `device-update-os`, `device-profiles`, `device-users`, `device-unenroll`, `device-custom-attributes`, `device-set-custom-attribute`, `device-lost-mode-enable`, `device-lost-mode-disable`, `device-lost-mode-play-sound`, `device-lost-mode-update-location` |
| **Enrollments** | `enrollment-list`, `enrollment-get`, `enrollment-delete`, `enrollment-send-invitation` |
| **Installed Apps** | `installed-app-get`, `installed-app-delete`, `installed-app-update`, `installed-app-request-management` |
| **Logs** | `log-list`, `log-get` |
| **Profiles** | `profile-list`, `profile-get`, `profile-assign-device`, `profile-unassign-device`, `push-certificate-get`, `push-certificate-scsr` |
| **Scripts** | `script-list`, `script-get`, `script-update`, `script-delete` |
| **Script Jobs** | `script-job-list`, `script-job-get`, `script-job-create`, `script-job-cancel` |

> **Note :** Les operations necessitant un upload de fichier (creation d'app avec binaire, creation de script, creation de profils de configuration) ne sont pas disponibles via MCP car le protocole ne supporte pas le transfert de fichiers.

---

## 9. Stack technique

| Composant | Technologie |
|-----------|------------|
| **Langage** | Go 1.25+ |
| **Framework CLI** | [Cobra](https://github.com/spf13/cobra) v1.10 |
| **Configuration** | [Viper](https://github.com/spf13/viper) v1.21 |
| **Keychain** | [go-keyring](https://github.com/zalando/go-keyring) v0.2 |
| **Sortie table** | [tablewriter](https://github.com/olekukonko/tablewriter) v1.1 |
| **Sortie YAML** | [yaml.v3](https://gopkg.in/yaml.v3) |
| **Terminal** | [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) (lecture securisee de la cle API) |
| **Serveur MCP** | [mcp-go](https://github.com/mark3labs/mcp-go) v0.52 |
| **Distribution** | Binaires cross-platform via GitHub Releases (tar.gz / zip) |
| **Self-update** | Telechargement et remplacement atomique du binaire depuis GitHub Releases |
| **Pagination API** | Cursor-based (`starting_after` + `limit`) |
| **Upload fichiers** | Multipart form-data (apps, scripts, profils, certificats) |
