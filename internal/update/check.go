package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	repoOwner = "dimer47"
	repoName  = "simplemdm-cli"
	apiURL    = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases/latest"
	cacheFile = "update-check.json"
	cacheTTL  = 24 * time.Hour
)

// CheckResult contains the result of an update check.
type CheckResult struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	UpdateAvailable bool   `json:"update_available"`
	DownloadURL     string `json:"download_url"`
}

// cachedCheck is the on-disk cache format.
type cachedCheck struct {
	CheckedAt time.Time    `json:"checked_at"`
	Result    *CheckResult `json:"result"`
}

// githubRelease represents the relevant fields from the GitHub API response.
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

// githubAsset represents a release asset.
type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// cacheDir returns the path to the cache directory (~/.simplemdm-cli).
func cacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".simplemdm-cli")
}

// cachePath returns the full path to the cache file.
func cachePath() string {
	return filepath.Join(cacheDir(), cacheFile)
}

// loadCache tries to load a cached check result.
func loadCache() *cachedCheck {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil
	}
	var c cachedCheck
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}

// saveCache writes the check result to disk.
func saveCache(result *CheckResult) {
	c := cachedCheck{
		CheckedAt: time.Now(),
		Result:    result,
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(cacheDir(), 0700)
	_ = os.WriteFile(cachePath(), data, 0600)
}

// trimVersion removes a leading "v" from a version string.
func trimVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}

// findAssetURL looks for the appropriate download URL for the current OS/arch.
func findAssetURL(assets []githubAsset) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Build expected suffixes based on goreleaser naming convention:
	// simplemdm-cli_<os>_<arch>.tar.gz  (linux/darwin)
	// simplemdm-cli_<os>_<arch>.zip     (windows)
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}

	expected := fmt.Sprintf("%s_%s_%s%s", repoName, goos, goarch, ext)

	for _, a := range assets {
		if strings.EqualFold(a.Name, expected) {
			return a.BrowserDownloadURL
		}
	}

	// Fallback: partial match
	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if strings.Contains(lower, goos) && strings.Contains(lower, goarch) {
			return a.BrowserDownloadURL
		}
	}

	return ""
}

// Check queries GitHub for the latest release and returns a CheckResult.
// Results are cached for 24h to avoid spamming the API.
func Check(currentVersion string) *CheckResult {
	// Skip check for dev builds
	if currentVersion == "" || currentVersion == "dev" {
		return &CheckResult{
			CurrentVersion:  currentVersion,
			LatestVersion:   "",
			UpdateAvailable: false,
		}
	}

	// Try cache first
	if cached := loadCache(); cached != nil {
		if time.Since(cached.CheckedAt) < cacheTTL && cached.Result != nil {
			// Update current version in case it changed (re-build)
			cached.Result.CurrentVersion = currentVersion
			cached.Result.UpdateAvailable = trimVersion(cached.Result.LatestVersion) != trimVersion(currentVersion)
			return cached.Result
		}
	}

	result := checkRemote(currentVersion)
	saveCache(result)
	return result
}

// CheckForce always queries GitHub, ignoring the cache.
func CheckForce(currentVersion string) *CheckResult {
	result := checkRemote(currentVersion)
	saveCache(result)
	return result
}

// checkRemote performs the actual HTTP request to the GitHub API.
func checkRemote(currentVersion string) *CheckResult {
	client := &http.Client{Timeout: 2 * time.Second}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &CheckResult{CurrentVersion: currentVersion}
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", repoName+"/"+currentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return &CheckResult{CurrentVersion: currentVersion}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &CheckResult{CurrentVersion: currentVersion}
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return &CheckResult{CurrentVersion: currentVersion}
	}

	latestVersion := trimVersion(release.TagName)
	currentClean := trimVersion(currentVersion)

	return &CheckResult{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestVersion,
		UpdateAvailable: latestVersion != currentClean,
		DownloadURL:     findAssetURL(release.Assets),
	}
}

// ClearCache removes the update check cache file.
func ClearCache() {
	_ = os.Remove(cachePath())
}
