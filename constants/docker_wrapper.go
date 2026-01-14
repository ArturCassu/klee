package constants

type DockerWrapper string

const (
	// Cross-platform wrappers (macOS, Linux, Windows)
	DockerDesktop  DockerWrapper = "docker-desktop"
	RancherDesktop DockerWrapper = "rancher-desktop"
	PodmanDesktop  DockerWrapper = "podman-desktop"

	// macOS-specific wrappers
	Colima   DockerWrapper = "colima"
	OrbStack DockerWrapper = "orbstack"

	// Linux-specific wrappers
	DockerEngine DockerWrapper = "docker-engine"
	Podman       DockerWrapper = "podman"
)

func AvailableWrappersOnOS() []DockerWrapper {
	var wrappers []DockerWrapper

	switch CurrentOS() {
	case MacOS:
		wrappers = append(wrappers, DockerDesktop, RancherDesktop, PodmanDesktop, Colima, OrbStack)
	case Linux:
		wrappers = append(wrappers, DockerEngine, Podman)
	case Windows:
		wrappers = append(wrappers, DockerDesktop)
	}

	return wrappers
}
