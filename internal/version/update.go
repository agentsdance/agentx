package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	latestReleaseURL  = "https://api.github.com/repos/agentsdance/agentx/releases/latest"
	updateCacheMaxAge = 12 * time.Hour
	releaseNotesURL   = "https://github.com/agentsdance/agentx/releases/latest"
)

// UpdateNotice describes an available update and how to upgrade.
type UpdateNotice struct {
	Latest  string
	Command string
}

type updateCache struct {
	CheckedAt time.Time `json:"checked_at"`
	Latest    string    `json:"latest"`
	Ignored   string    `json:"ignored,omitempty"`
}

// CheckForUpdate returns an update notice when a newer release is available.
func CheckForUpdate(ctx context.Context) (UpdateNotice, bool, error) {
	current, _, ok := parseSemver(Version)
	if !ok {
		if !forceUpdateCheck() {
			return UpdateNotice{}, false, nil
		}
		current = semVersion{major: 0, minor: 0, patch: 0}
	}

	cache, cacheOk := loadUpdateCache()
	latest := cache.Latest
	cacheFresh := cacheOk && latest != "" && updateCacheMaxAge > 0 && time.Since(cache.CheckedAt) <= updateCacheMaxAge
	if !cacheFresh {
		fetched, err := fetchLatestReleaseTag(ctx)
		if err != nil {
			if latest == "" {
				return UpdateNotice{}, false, err
			}
		} else {
			latest = fetched
			cache.Latest = fetched
			cache.CheckedAt = time.Now()
			_ = saveUpdateCache(cache)
		}
	}

	latestParsed, latestNormalized, ok := parseSemver(latest)
	if !ok {
		return UpdateNotice{}, false, nil
	}

	if cache.Ignored != "" && cache.Ignored == latestNormalized {
		return UpdateNotice{}, false, nil
	}

	if compareSemver(latestParsed, current) <= 0 {
		return UpdateNotice{}, false, nil
	}

	notice := UpdateNotice{
		Latest:  latestNormalized,
		Command: upgradeCommand(),
	}
	return notice, true, nil
}

// ReleaseNotesURL returns the latest release notes URL.
func ReleaseNotesURL() string {
	return releaseNotesURL
}

func forceUpdateCheck() bool {
	return strings.TrimSpace(os.Getenv("AGENTX_FORCE_UPDATE_CHECK")) != ""
}

func fetchLatestReleaseTag(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	// GitHub requires a User-Agent header.
	req.Header.Set("User-Agent", "agentx")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("update check returned status %d", resp.StatusCode)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.TagName == "" {
		return "", fmt.Errorf("update check returned empty tag")
	}

	return payload.TagName, nil
}

var semverPattern = regexp.MustCompile(`v?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z\.-]+)?`)

type semVersion struct {
	major int
	minor int
	patch int
	pre   []string
}

func parseSemver(value string) (semVersion, string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return semVersion{}, "", false
	}

	match := semverPattern.FindString(value)
	if match == "" {
		return semVersion{}, "", false
	}

	if !strings.HasPrefix(match, "v") {
		match = "v" + match
	}

	parsed, ok := parseSemverCore(match[1:])
	if !ok {
		return semVersion{}, "", false
	}

	return parsed, match, true
}

func parseSemverCore(value string) (semVersion, bool) {
	value = strings.SplitN(value, "+", 2)[0]
	core := value
	pre := ""
	if parts := strings.SplitN(value, "-", 2); len(parts) == 2 {
		core = parts[0]
		pre = parts[1]
	}

	nums := strings.Split(core, ".")
	if len(nums) != 3 {
		return semVersion{}, false
	}

	major, ok := parsePositiveInt(nums[0])
	if !ok {
		return semVersion{}, false
	}
	minor, ok := parsePositiveInt(nums[1])
	if !ok {
		return semVersion{}, false
	}
	patch, ok := parsePositiveInt(nums[2])
	if !ok {
		return semVersion{}, false
	}

	var preIDs []string
	if pre != "" {
		preIDs = strings.Split(pre, ".")
	}

	return semVersion{
		major: major,
		minor: minor,
		patch: patch,
		pre:   preIDs,
	}, true
}

func parsePositiveInt(value string) (int, bool) {
	if value == "" {
		return 0, false
	}
	n := 0
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
	}
	return n, true
}

func compareSemver(a, b semVersion) int {
	if a.major != b.major {
		return compareInt(a.major, b.major)
	}
	if a.minor != b.minor {
		return compareInt(a.minor, b.minor)
	}
	if a.patch != b.patch {
		return compareInt(a.patch, b.patch)
	}

	return comparePrerelease(a.pre, b.pre)
}

func comparePrerelease(a, b []string) int {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	if len(a) == 0 {
		return 1
	}
	if len(b) == 0 {
		return -1
	}

	for i := 0; i < len(a) || i < len(b); i++ {
		if i >= len(a) {
			return -1
		}
		if i >= len(b) {
			return 1
		}

		ai := a[i]
		bi := b[i]
		an, aNum := parseNumericIdentifier(ai)
		bn, bNum := parseNumericIdentifier(bi)

		switch {
		case aNum && bNum:
			if an != bn {
				return compareInt(an, bn)
			}
		case aNum && !bNum:
			return -1
		case !aNum && bNum:
			return 1
		default:
			if cmp := strings.Compare(ai, bi); cmp != 0 {
				return cmp
			}
		}
	}

	return 0
}

func parseNumericIdentifier(value string) (int, bool) {
	if value == "" {
		return 0, false
	}
	n := 0
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
	}
	return n, true
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func loadUpdateCache() (updateCache, bool) {
	cachePath, err := getUpdateCachePath()
	if err != nil {
		return updateCache{}, false
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return updateCache{}, false
	}

	var cache updateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return updateCache{}, false
	}
	return cache, true
}

// SkipVersion records a version to ignore until a newer release is available.
func SkipVersion(version string) error {
	_, normalized, ok := parseSemver(version)
	if !ok {
		return fmt.Errorf("invalid version")
	}

	cache, _ := loadUpdateCache()
	cache.Ignored = normalized
	if cache.Latest == "" {
		cache.Latest = normalized
	}
	if cache.CheckedAt.IsZero() {
		cache.CheckedAt = time.Now()
	}

	return saveUpdateCache(cache)
}

func saveUpdateCache(cache updateCache) error {
	cachePath, err := getUpdateCachePath()
	if err != nil {
		return err
	}

	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func getUpdateCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentx", "cache", "update.json"), nil
}

func upgradeCommand() string {
	if cmd := strings.TrimSpace(os.Getenv("AGENTX_UPGRADE_COMMAND")); cmd != "" {
		return cmd
	}

	exe, err := os.Executable()
	if err == nil {
		if isHomebrewInstall(exe) {
			return "brew upgrade agentx"
		}
		if isGoInstall(exe) {
			return "go install github.com/agentsdance/agentx@latest"
		}
	}

	return "go install github.com/agentsdance/agentx@latest"
}

func isHomebrewInstall(exe string) bool {
	brewMarkers := []string{
		"/Cellar/agentx/",
		"/Homebrew/Cellar/agentx/",
		"/home/linuxbrew/.linuxbrew/Cellar/agentx/",
	}
	for _, marker := range brewMarkers {
		if strings.Contains(exe, marker) {
			return true
		}
	}
	return false
}

func isGoInstall(exe string) bool {
	if gobin := strings.TrimSpace(os.Getenv("GOBIN")); gobin != "" {
		if strings.HasPrefix(exe, filepath.Clean(gobin)+string(os.PathSeparator)) {
			return true
		}
	}

	gopath := strings.TrimSpace(os.Getenv("GOPATH"))
	if gopath == "" {
		return false
	}

	for _, path := range strings.Split(gopath, string(os.PathListSeparator)) {
		if path == "" {
			continue
		}
		binDir := filepath.Join(path, "bin")
		if strings.HasPrefix(exe, binDir+string(os.PathSeparator)) {
			return true
		}
	}

	return false
}
