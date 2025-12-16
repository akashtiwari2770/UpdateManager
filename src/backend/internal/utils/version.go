package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// CompareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	// Normalize versions (remove build metadata and pre-release tags for comparison)
	v1Clean := normalizeVersion(v1)
	v2Clean := normalizeVersion(v2)

	parts1 := parseVersion(v1Clean)
	parts2 := parseVersion(v2Clean)

	// Compare major version
	if parts1[0] < parts2[0] {
		return -1
	}
	if parts1[0] > parts2[0] {
		return 1
	}

	// Compare minor version
	if parts1[1] < parts2[1] {
		return -1
	}
	if parts1[1] > parts2[1] {
		return 1
	}

	// Compare patch version
	if parts1[2] < parts2[2] {
		return -1
	}
	if parts1[2] > parts2[2] {
		return 1
	}

	return 0
}

// IsVersionNewer checks if v1 is newer than v2
func IsVersionNewer(v1, v2 string) bool {
	return CompareVersions(v1, v2) > 0
}

// IsVersionOlder checks if v1 is older than v2
func IsVersionOlder(v1, v2 string) bool {
	return CompareVersions(v1, v2) < 0
}

// normalizeVersion removes build metadata and pre-release tags for comparison
func normalizeVersion(version string) string {
	// Remove build metadata (everything after +)
	if idx := strings.Index(version, "+"); idx != -1 {
		version = version[:idx]
	}
	// Remove pre-release tags (everything after -) for basic comparison
	// We keep them for now but could enhance this later
	return version
}

// parseVersion parses a version string into [major, minor, patch]
func parseVersion(version string) [3]int {
	parts := [3]int{0, 0, 0}

	// Remove pre-release suffix if present
	if idx := strings.Index(version, "-"); idx != -1 {
		version = version[:idx]
	}

	// Split by dots
	versionParts := strings.Split(version, ".")
	for i, part := range versionParts {
		if i >= 3 {
			break
		}
		// Extract numeric part
		re := regexp.MustCompile(`\d+`)
		matches := re.FindString(part)
		if matches != "" {
			if num, err := strconv.Atoi(matches); err == nil {
				parts[i] = num
			}
		}
	}

	return parts
}

// GetVersionGapType determines the type of version gap between two versions
// Returns: "patch", "minor", "major", or "unknown"
func GetVersionGapType(current, latest string) string {
	currentParts := parseVersion(current)
	latestParts := parseVersion(latest)

	// Major version difference
	if currentParts[0] != latestParts[0] {
		return "major"
	}

	// Minor version difference
	if currentParts[1] != latestParts[1] {
		return "minor"
	}

	// Patch version difference
	if currentParts[2] != latestParts[2] {
		return "patch"
	}

	return "unknown"
}

