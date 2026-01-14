package constants

import "runtime"

type OperatingSystem string

const (
	Windows OperatingSystem = "windows"
	Linux   OperatingSystem = "linux"
	MacOS   OperatingSystem = "darwin"
	Unknown OperatingSystem = "unknown"
)

func CurrentOS() OperatingSystem {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return MacOS
	default:
		return Unknown
	}
}

func IsWindows() bool {
	return CurrentOS() == Windows
}

func IsLinux() bool {
	return CurrentOS() == Linux
}

func IsMacOS() bool {
	return CurrentOS() == MacOS
}
