# simplemdm-cli

CLI pour l'API SimpleMDM -- gerez vos appareils Apple depuis le terminal.

## Fonctionnalites

- **116 outils MCP** et **130+ commandes CLI** couvrant l'integralite de l'API SimpleMDM
- **Stockage securise de la cle API** : stockee dans le gestionnaire d'identifiants de votre systeme (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- **Multi-contexte** : gerez plusieurs comptes SimpleMDM
- **Sortie flexible** : table, JSON, YAML, CSV
- **Multi-plateforme** : macOS, Linux, Windows (amd64 et arm64)
- **Integration MCP** : 116 outils pour Claude Code, VS Code, JetBrains
- **Mise a jour automatique** : verification automatique des mises a jour et mise a jour en une commande

## Prerequis

- Un compte [SimpleMDM](https://simplemdm.com)
- Une **cle API SimpleMDM** (creee depuis [Settings > API](https://a.simplemdm.com/admin/settings/api))

## Installation

### Methode 1 : Telecharger le binaire (recommande)

Rendez-vous sur la page [Releases](https://github.com/dimer47/simplemdm-cli/releases/latest) et telechargez l'archive correspondant a votre plateforme.

Ou en une seule commande :

**macOS (Apple Silicon -- M1/M2/M3/M4) :**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_darwin_arm64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**macOS (Intel) :**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_darwin_amd64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Linux (amd64) :**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_linux_amd64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Linux (arm64 -- Raspberry Pi, etc.) :**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_linux_arm64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Windows :**

Telechargez `simplemdm-cli_windows_amd64.zip` depuis la page [Releases](https://github.com/dimer47/simplemdm-cli/releases/latest), extrayez l'archive et ajoutez le dossier a votre `PATH`.

### Methode 2 : Depuis les sources (necessite Go 1.25+)

```bash
go install github.com/dimer47/simplemdm-cli@latest
```

Le binaire sera installe dans `$GOPATH/bin/` (generalement `~/go/bin/`). Assurez-vous que ce repertoire est dans votre `PATH`.

### Methode 3 : Compilation locale

```bash
git clone https://github.com/dimer47/simplemdm-cli.git
cd simplemdm-cli
go build -o simplemdm-cli .
./simplemdm-cli version
```

### Verifier l'installation

```bash
simplemdm-cli version
# simplemdm-cli 1.0.0 (abc1234) built 2026-05-11T10:00:00Z
```

## Mise a jour

Le CLI verifie automatiquement la disponibilite de nouvelles versions au demarrage et vous avertit lorsqu'une mise a jour est disponible.

```bash
# Mettre a jour vers la derniere version
simplemdm-cli self-update
```

La mise a jour est telechargee depuis GitHub Releases et remplace le binaire actuel sur place. Si le binaire se trouve dans un repertoire protege (ex. `/usr/local/bin/`), essayez d'executer la commande avec `sudo`.

## Demarrage rapide

### 1. Obtenir une cle API SimpleMDM

1. Connectez-vous a la [console d'administration SimpleMDM](https://a.simplemdm.com)
2. Allez dans **Settings > API**
3. Copiez votre cle API

### 2. Configurer le CLI

```bash
simplemdm-cli auth login
```

Repondez aux 2 invites :
```
Context name (default): Enter          # Appuyez sur Entree pour "default"
SimpleMDM API key: ••••••••••••        # Collez votre cle API (masquee)
```

La cle API est stockee dans le **gestionnaire d'identifiants de votre systeme** (macOS Keychain, Windows Credential Manager ou Linux Secret Service) -- chiffree, jamais ecrite en clair sur le disque.

### 3. Tester

```bash
# Obtenir les details de votre compte
simplemdm-cli account get

# Lister vos appareils
simplemdm-cli device list

# Sortie JSON
simplemdm-cli device list --json
```

## Configuration

### Priorite de resolution de la cle API

| Priorite | Source | Cas d'usage |
|----------|--------|-------------|
| 1 | Flag `--api-key` / `-k` | Tests ponctuels |
| 2 | Variable d'env. `SMDM_API_KEY` | CI/CD, scripts |
| 3 | Gestionnaire d'identifiants systeme | Usage quotidien (via `auth login`) |

### Multi-contexte (plusieurs comptes SimpleMDM)

```bash
# Configurer un contexte "production"
simplemdm-cli auth login
# -> Entrez "production" comme nom de contexte

# Configurer un contexte "staging"
simplemdm-cli auth login
# -> Entrez "staging" comme nom de contexte

# Lister tous les contextes (* = actif)
simplemdm-cli auth list
# * production (****abcd1234)
#   staging    (****efgh5678)

# Changer de contexte
simplemdm-cli auth switch staging

# Utiliser un contexte pour une seule commande
simplemdm-cli device list --context production

# Verifier le statut (la cle est-elle toujours valide ?)
simplemdm-cli auth status

# Supprimer un contexte
simplemdm-cli auth remove staging
```

## Utilisation

### Compte

```bash
simplemdm-cli account get                                    # Details du compte
simplemdm-cli account update --name "My Company"             # Mettre a jour le compte
```

### Appareils

```bash
simplemdm-cli device list                                    # Lister tous les appareils
simplemdm-cli device list --json                             # Sortie JSON
simplemdm-cli device list --search "MacBook"                 # Rechercher
simplemdm-cli device get <id>                                # Details d'un appareil
simplemdm-cli device create --name "New Device"              # Creer
simplemdm-cli device update <id> --name "Renamed"            # Renommer
simplemdm-cli device delete <id>                             # Supprimer
simplemdm-cli device refresh <id>                            # Rafraichir les infos
simplemdm-cli device push-apps <id>                          # Pousser les apps
simplemdm-cli device lock <id> --message "Lost" --pin 1234   # Verrouiller
simplemdm-cli device wipe <id>                               # Effacer
simplemdm-cli device restart <id>                            # Redemarrer
simplemdm-cli device shutdown <id>                           # Eteindre
simplemdm-cli device clear-passcode <id>                     # Effacer le code
simplemdm-cli device update-os <id>                          # Mettre a jour l'OS
simplemdm-cli device unenroll <id>                           # Desenroler
```

### Appareil -- Mode perdu (iOS supervise)

```bash
simplemdm-cli device lost-mode-enable <id> --message "Call IT" --phone-number "+1234567890"
simplemdm-cli device lost-mode-disable <id>
simplemdm-cli device lost-mode-play-sound <id>
simplemdm-cli device lost-mode-update-location <id>
```

### Appareil -- Commandes avancees

```bash
simplemdm-cli device bluetooth-enable <id>                   # Activer le Bluetooth
simplemdm-cli device bluetooth-disable <id>                  # Desactiver le Bluetooth
simplemdm-cli device remote-desktop-enable <id>              # Activer le Bureau a distance
simplemdm-cli device remote-desktop-disable <id>             # Desactiver le Bureau a distance
simplemdm-cli device rotate-firmware-password <id>           # Rotation du mot de passe firmware
simplemdm-cli device rotate-recovery-lock <id>               # Rotation du verrou de recuperation
simplemdm-cli device rotate-filevault-key <id>               # Rotation de la cle FileVault
simplemdm-cli device set-admin-password <id> --new-password "..."
simplemdm-cli device rotate-admin-password <id>
simplemdm-cli device clear-firmware-password <id>
simplemdm-cli device clear-recovery-lock <id>
simplemdm-cli device clear-restrictions-password <id>
simplemdm-cli device set-timezone <id> --timezone "America/New_York"
```

### Appareil -- Attributs personnalises et utilisateurs

```bash
simplemdm-cli device custom-attributes <id>                  # Lister les valeurs des attributs
simplemdm-cli device set-custom-attribute <id> department --value "Engineering"
simplemdm-cli device profiles <id>                           # Lister les profils de l'appareil
simplemdm-cli device users <id>                              # Lister les utilisateurs de l'appareil
simplemdm-cli device delete-user <id> <user-id>              # Supprimer un utilisateur
```

### Applications

```bash
simplemdm-cli app list                                       # Lister toutes les apps
simplemdm-cli app get <id>                                   # Details d'une app
simplemdm-cli app create --name "My App" --app-store-id 123  # Creer depuis l'App Store
simplemdm-cli app create --binary ./app.ipa                  # Telecharger une app enterprise
simplemdm-cli app update <id> --binary ./app.ipa             # Mettre a jour le binaire
simplemdm-cli app delete <id>                                # Supprimer
simplemdm-cli app installs <id>                              # Lister les installations
simplemdm-cli app managed-configs <id>                       # Lister les configs gerees
simplemdm-cli app managed-config-create <id> --key "url" --value "https://..."
simplemdm-cli app managed-configs-push <id>                  # Pousser les configs vers les appareils
simplemdm-cli app managed-config-delete <id> <config-id>     # Supprimer une config
```

### Groupes d'assignation

```bash
simplemdm-cli assignment-group list                          # Lister les groupes
simplemdm-cli assignment-group get <id>                      # Details d'un groupe
simplemdm-cli assignment-group create --name "Team"          # Creer
simplemdm-cli assignment-group update <id> --name "New Name" # Mettre a jour
simplemdm-cli assignment-group delete <id>                   # Supprimer
simplemdm-cli assignment-group assign-app <id> --app-id 10   # Assigner une app
simplemdm-cli assignment-group assign-device <id> --device-id 121
simplemdm-cli assignment-group push-apps <id>                # Pousser les apps
simplemdm-cli assignment-group clone <id>                    # Cloner le groupe
```

### Profils

```bash
simplemdm-cli profile list                                   # Lister les profils
simplemdm-cli profile get <id>                               # Details d'un profil
simplemdm-cli profile assign-device <id> --device-id 121     # Assigner a un appareil
simplemdm-cli profile unassign-device <id> --device-id 121   # Desassigner
```

### Profils de configuration personnalises

```bash
simplemdm-cli custom-configuration-profile list              # Lister les profils
simplemdm-cli custom-configuration-profile delete <id>       # Supprimer
simplemdm-cli custom-configuration-profile push-device <id> --device-id 121
simplemdm-cli custom-configuration-profile remove-device <id> --device-id 121
```

### Attributs personnalises

```bash
simplemdm-cli custom-attribute list                          # Lister les attributs
simplemdm-cli custom-attribute get <id>                      # Details d'un attribut
simplemdm-cli custom-attribute create --name "department"    # Creer
simplemdm-cli custom-attribute delete <id>                   # Supprimer
```

### Declarations personnalisees

```bash
simplemdm-cli custom-declaration list                        # Lister les declarations
simplemdm-cli custom-declaration delete <id>                 # Supprimer
```

### Groupes d'appareils

```bash
simplemdm-cli device-group list                              # Lister les groupes
simplemdm-cli device-group get <id>                          # Details d'un groupe
simplemdm-cli device-group assign-device <id> --device-id 121
```

### Serveurs DEP

```bash
simplemdm-cli dep-server list                                # Lister les serveurs DEP
simplemdm-cli dep-server get <id>                            # Details d'un serveur
simplemdm-cli dep-server devices <id>                        # Lister les appareils DEP
simplemdm-cli dep-server sync <id>                           # Synchroniser le serveur
```

### Enrollements

```bash
simplemdm-cli enrollment list                                # Lister les enrollements
simplemdm-cli enrollment get <id>                            # Details d'un enrollement
simplemdm-cli enrollment delete <id>                         # Supprimer
simplemdm-cli enrollment send-invitation <id> --contact "user@company.com"
```

### Applications installees

```bash
simplemdm-cli installed-app get <id>                         # Details de l'app
simplemdm-cli installed-app delete <id>                      # Supprimer (desinstaller)
simplemdm-cli installed-app update <id>                      # Mettre a jour l'app
```

### Scripts et taches de script

```bash
simplemdm-cli script list                                    # Lister les scripts
simplemdm-cli script get <id>                                # Details d'un script
simplemdm-cli script delete <id>                             # Supprimer
simplemdm-cli script-job list                                # Lister les taches
simplemdm-cli script-job get <id>                            # Details d'une tache
simplemdm-cli script-job create --script-id 100 --device-ids "121,122"
simplemdm-cli script-job cancel <id>                         # Annuler une tache
```

### Journaux

```bash
simplemdm-cli log list                                       # Lister les journaux
simplemdm-cli log get <id>                                   # Details d'un journal
```

### Certificat Push

```bash
simplemdm-cli push-certificate get                           # Details du certificat
```

## Variables d'environnement

| Variable | Description | Valeur par defaut |
|----------|-------------|-------------------|
| `SMDM_API_KEY` | Cle API SimpleMDM | -- |
| `SMDM_OUTPUT` | Format de sortie : `table`, `json`, `yaml`, `csv` | `table` |
| `SMDM_DEBUG` | Mode debug (`true`/`false`) | `false` |
| `SMDM_CONTEXT` | Contexte actif | `default` |
| `NO_COLOR` | Desactiver les couleurs | -- |

## Completion shell

```bash
# Bash
simplemdm-cli completion bash > /etc/bash_completion.d/simplemdm-cli

# Zsh (ajoutez a votre .zshrc)
simplemdm-cli completion zsh > "${fpath[1]}/_simplemdm-cli"

# Fish
simplemdm-cli completion fish > ~/.config/fish/completions/simplemdm-cli.fish

# PowerShell
simplemdm-cli completion powershell > simplemdm-cli.ps1
```

## Integration MCP (Claude Code, VS Code, JetBrains)

Le CLI inclut un serveur [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) integre exposant **116 outils** pour les assistants IA.

### Configuration

Ajoutez ceci a vos parametres Claude Code (ou VS Code / JetBrains avec l'extension Claude) :

```json
{
  "mcpServers": {
    "simplemdm": {
      "command": "simplemdm-cli",
      "args": ["mcp-serve"],
      "env": {
        "SMDM_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

> Si vous avez deja configure la cle API via `simplemdm-cli auth login`, le serveur MCP utilisera automatiquement votre gestionnaire d'identifiants systeme -- pas besoin de la variable d'environnement `SMDM_API_KEY`.

### Outils MCP disponibles

| Categorie | Outils | Description |
|-----------|--------|-------------|
| **Compte** | `account-get`, `account-update` | Gestion des details du compte |
| **Applications** | `app-list`, `app-get`, `app-create`, `app-update`, `app-delete`, `app-installs`, `app-managed-configs`, `app-managed-config-create`, `app-managed-configs-push`, `app-managed-config-delete` | Cycle de vie complet des apps |
| **Groupes d'assignation** | `assignment-group-list`, `assignment-group-get`, `assignment-group-create`, `assignment-group-update`, `assignment-group-delete`, `assignment-group-assign-app`, `assignment-group-assign-device`, `assignment-group-push-apps`, `assignment-group-clone`, ... | Gestion des groupes |
| **Attributs personnalises** | `custom-attribute-list`, `custom-attribute-get`, `custom-attribute-create`, `custom-attribute-delete` | Attributs personnalises |
| **Profils de config. personnalises** | `custom-configuration-profile-list`, `custom-configuration-profile-delete`, `custom-configuration-profile-push-device`, `custom-configuration-profile-remove-device` | Gestion des profils |
| **Declarations personnalisees** | `custom-declaration-list`, `custom-declaration-delete` | Gestion des declarations |
| **Serveurs DEP** | `dep-server-list`, `dep-server-get`, `dep-server-devices`, `dep-server-sync` | Gestion DEP |
| **Groupes d'appareils** | `device-group-list`, `device-group-get`, `device-group-assign-device` | Gestion des groupes |
| **Appareils** | `device-list`, `device-get`, `device-create`, `device-update`, `device-delete`, `device-lock`, `device-wipe`, `device-restart`, `device-shutdown`, `device-refresh`, `device-push-apps`, `device-update-os`, `device-lost-mode-*`, ... | Gestion complete des appareils |
| **Enrollements** | `enrollment-list`, `enrollment-get`, `enrollment-delete`, `enrollment-send-invitation` | Gestion des enrollements |
| **Apps installees** | `installed-app-get`, `installed-app-delete`, `installed-app-update` | Gestion des apps |
| **Journaux** | `log-list`, `log-get` | Acces aux journaux |
| **Profils** | `profile-list`, `profile-get`, `profile-assign-device`, `profile-unassign-device` | Gestion des profils |
| **Certificat Push** | `push-certificate-get` | Infos du certificat |
| **Scripts** | `script-list`, `script-get`, `script-delete` | Gestion des scripts |
| **Taches de script** | `script-job-list`, `script-job-get`, `script-job-create`, `script-job-cancel` | Gestion des taches |

### Utilisation dans Claude Code

Une fois configure, vous pouvez simplement dire :

- *"Liste mes appareils SimpleMDM"*
- *"Verrouille l'appareil 121 avec le message 'Contactez le service IT'"*
- *"Quelles apps sont installees sur l'appareil 122 ?"*
- *"Pousse toutes les apps vers le groupe d'assignation Engineering"*
- *"Cree une tache de script pour le script 100 sur les appareils 121 et 122"*

Claude appellera automatiquement les bons outils MCP.

## Developpement

```bash
# Cloner
git clone https://github.com/dimer47/simplemdm-cli.git
cd simplemdm-cli

# Compiler
go build -o simplemdm-cli .

# Lancer les tests
go test ./...

# Analyse statique
go vet ./...
```

### Creer une nouvelle release

```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions compile et publie automatiquement
```

## Licence

MIT
