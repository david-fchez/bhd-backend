package ver

import "strconv"

const (
	AppVersion       = 0.1
	AppName          = "BHD BunnyHedger "
	AppNameSuffix    = " - there is no winter as crypto winter"
	BuildTypeRelease = "RELEASE"
	BuildTypeDebug   = "DEBUG"
	BuildType        = BuildTypeRelease
)

// IsNewerVersion returns true if the version
// number is higher than the current app version
func IsNewerVersion(version float64) bool {
	return version > AppVersion
}

// GetVersionString returns the current app version string
func GetVersionString() string {
	return AppName + strconv.FormatFloat(AppVersion, 'f', 4, 64) + AppNameSuffix
}

// ToVersionString returns the current app version string
func ToVersionString(version float64) string {
	return strconv.FormatFloat(version, 'f', 4, 64)
}

// FromVersionString returns the current app version string
func FromVersionString(version string) (float64, error) {
	return strconv.ParseFloat(version, 64)
}

// GetBuildType returns the current build type
func GetBuildType() string {
	return BuildType
}
