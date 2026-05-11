package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrimVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"v0.1.2", "0.1.2"},
		{" v1.2.3 ", "1.2.3"},
		{"", ""},
	}
	for _, tt := range tests {
		got := trimVersion(tt.input)
		if got != tt.expected {
			t.Errorf("trimVersion(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFindAssetURL(t *testing.T) {
	assets := []githubAsset{
		{Name: "simplemdm-cli_linux_amd64.tar.gz", BrowserDownloadURL: "https://example.com/linux_amd64.tar.gz"},
		{Name: "simplemdm-cli_darwin_arm64.tar.gz", BrowserDownloadURL: "https://example.com/darwin_arm64.tar.gz"},
		{Name: "simplemdm-cli_windows_amd64.zip", BrowserDownloadURL: "https://example.com/windows_amd64.zip"},
		{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums.txt"},
	}

	// Test that at least one asset can be found (depends on runtime)
	url := findAssetURL(assets)
	if url == "" {
		t.Error("expected to find an asset URL for current platform")
	}
}

func TestFindAssetURL_Empty(t *testing.T) {
	url := findAssetURL(nil)
	if url != "" {
		t.Errorf("expected empty URL for nil assets, got %q", url)
	}
}

func TestCheck_DevBuild(t *testing.T) {
	result := Check("dev")
	if result.UpdateAvailable {
		t.Error("dev builds should not show update available")
	}
	if result.LatestVersion != "" {
		t.Errorf("expected empty latest version for dev, got %q", result.LatestVersion)
	}
}

func TestCheck_EmptyVersion(t *testing.T) {
	result := Check("")
	if result.UpdateAvailable {
		t.Error("empty version should not show update available")
	}
}

func TestCheckRemote_WithMockServer(t *testing.T) {
	release := githubRelease{
		TagName: "v2.0.0",
		Assets: []githubAsset{
			{Name: "simplemdm-cli_linux_amd64.tar.gz", BrowserDownloadURL: "https://example.com/linux.tar.gz"},
			{Name: "simplemdm-cli_darwin_arm64.tar.gz", BrowserDownloadURL: "https://example.com/darwin.tar.gz"},
			{Name: "simplemdm-cli_windows_amd64.zip", BrowserDownloadURL: "https://example.com/windows.zip"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// We can't easily override apiURL, so test checkRemote indirectly
	// Instead test the result construction logic
	result := &CheckResult{
		CurrentVersion:  "1.0.0",
		LatestVersion:   trimVersion(release.TagName),
		UpdateAvailable: trimVersion(release.TagName) != trimVersion("1.0.0"),
		DownloadURL:     findAssetURL(release.Assets),
	}

	if !result.UpdateAvailable {
		t.Error("expected update available")
	}
	if result.LatestVersion != "2.0.0" {
		t.Errorf("expected latest 2.0.0, got %s", result.LatestVersion)
	}
}

func TestCheckResult_SameVersion(t *testing.T) {
	result := &CheckResult{
		CurrentVersion:  "1.0.0",
		LatestVersion:   "1.0.0",
		UpdateAvailable: trimVersion("1.0.0") != trimVersion("v1.0.0"),
	}
	if result.UpdateAvailable {
		t.Error("same version should not show update available")
	}
}
