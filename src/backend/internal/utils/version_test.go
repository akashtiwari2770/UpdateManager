package utils

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"v1 equals v2", "1.2.3", "1.2.3", 0},
		{"v1 newer patch", "1.2.4", "1.2.3", 1},
		{"v1 older patch", "1.2.2", "1.2.3", -1},
		{"v1 newer minor", "1.3.0", "1.2.3", 1},
		{"v1 older minor", "1.1.3", "1.2.3", -1},
		{"v1 newer major", "2.0.0", "1.2.3", 1},
		{"v1 older major", "0.9.9", "1.2.3", -1},
		{"different major", "3.0.0", "2.9.9", 1},
		{"same version with build", "1.2.3+build", "1.2.3", 0},
		{"same version with pre-release", "1.2.3-beta", "1.2.3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%s, %s) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestIsVersionNewer(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{"v1 is newer", "1.2.4", "1.2.3", true},
		{"v1 is older", "1.2.2", "1.2.3", false},
		{"v1 equals v2", "1.2.3", "1.2.3", false},
		{"major version newer", "2.0.0", "1.9.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsVersionNewer(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("IsVersionNewer(%s, %s) = %v, expected %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestIsVersionOlder(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{"v1 is older", "1.2.2", "1.2.3", true},
		{"v1 is newer", "1.2.4", "1.2.3", false},
		{"v1 equals v2", "1.2.3", "1.2.3", false},
		{"major version older", "1.9.9", "2.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsVersionOlder(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("IsVersionOlder(%s, %s) = %v, expected %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestGetVersionGapType(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected string
	}{
		{"patch gap", "1.2.3", "1.2.5", "patch"},
		{"minor gap", "1.2.3", "1.3.0", "minor"},
		{"major gap", "1.2.3", "2.0.0", "major"},
		{"same version", "1.2.3", "1.2.3", "unknown"},
		{"multiple minor", "1.2.3", "1.5.0", "minor"},
		{"multiple major", "1.2.3", "3.0.0", "major"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVersionGapType(tt.current, tt.latest)
			if result != tt.expected {
				t.Errorf("GetVersionGapType(%s, %s) = %s, expected %s", tt.current, tt.latest, result, tt.expected)
			}
		})
	}
}

