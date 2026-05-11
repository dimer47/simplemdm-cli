package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dimer47/simplemdm-cli/internal/update"
	"github.com/spf13/cobra"
)

func newSelfUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update simplemdm-cli to the latest version",
		Long:  `Check for a new version of simplemdm-cli on GitHub Releases and install it.`,
		RunE:  runSelfUpdate,
	}
	return cmd
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := Version

	fmt.Fprintf(os.Stderr, "Current version: %s\n", currentVersion)
	fmt.Fprintf(os.Stderr, "Checking for updates...\n")

	result := update.CheckForce(currentVersion)

	if result.LatestVersion == "" {
		return fmt.Errorf("could not determine the latest version (network error or dev build)")
	}

	if !result.UpdateAvailable {
		fmt.Fprintf(os.Stderr, "You are already running the latest version (%s).\n", currentVersion)
		return nil
	}

	fmt.Fprintf(os.Stderr, "New version available: %s -> %s\n", currentVersion, result.LatestVersion)

	if result.DownloadURL == "" {
		return fmt.Errorf("no compatible binary found for %s/%s.\nPlease download manually from: https://github.com/dimer47/simplemdm-cli/releases/latest", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Fprintf(os.Stderr, "Downloading %s...\n", result.DownloadURL)

	archiveData, err := downloadArchive(result.DownloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Extracting binary...\n")

	binaryData, err := extractBinary(archiveData, result.DownloadURL)
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("could not resolve executable path: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Replacing binary at %s...\n", execPath)

	if err := replaceBinary(execPath, binaryData); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("permission denied: try running with sudo:\n  sudo %s self-update", os.Args[0])
		}
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Clear the update cache since we just updated
	update.ClearCache()

	fmt.Fprintf(os.Stderr, "Successfully updated to version %s!\n", result.LatestVersion)
	return nil
}

// downloadArchive fetches the archive from the given URL.
func downloadArchive(url string) ([]byte, error) {
	client := &http.Client{Timeout: 120 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// extractBinary extracts the simplemdm-cli binary from the archive.
func extractBinary(archiveData []byte, archiveURL string) ([]byte, error) {
	if strings.HasSuffix(archiveURL, ".zip") {
		return extractFromZip(archiveData)
	}
	return extractFromTarGz(archiveData)
}

// extractFromTarGz extracts the binary from a .tar.gz archive.
func extractFromTarGz(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		name := filepath.Base(header.Name)
		if isBinaryName(name) && header.Typeflag == tar.TypeReg {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("binary not found in tar.gz archive")
}

// extractFromZip extracts the binary from a .zip archive.
func extractFromZip(data []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	for _, f := range zipReader.File {
		name := filepath.Base(f.Name)
		if isBinaryName(name) && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open zip entry: %w", err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("binary not found in zip archive")
}

// isBinaryName checks if the file name matches the expected binary name.
func isBinaryName(name string) bool {
	return name == "simplemdm-cli" || name == "simplemdm-cli.exe"
}

// replaceBinary atomically replaces the binary at the given path.
func replaceBinary(execPath string, newBinary []byte) error {
	dir := filepath.Dir(execPath)

	// Write new binary to a temp file in the same directory (for atomic rename)
	tmpFile, err := os.CreateTemp(dir, "simplemdm-cli-update-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on failure
	defer func() {
		// Remove temp file if it still exists (means we failed)
		if _, err := os.Stat(tmpPath); err == nil {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(newBinary); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	// Make executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return err
	}

	// Rename (atomic on most systems when same filesystem)
	if err := os.Rename(tmpPath, execPath); err != nil {
		return err
	}

	// Prevent defer from removing the file we just renamed
	// (Stat on tmpPath will now fail since it was renamed)

	return nil
}
