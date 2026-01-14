package utils

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/ArturCassu/klee/constants"
	"github.com/spf13/viper"
)

func StartDocker() error {
	wrappers := constants.AvailableWrappersOnOS()
	if len(wrappers) == 0 {
		return fmt.Errorf("no available Docker wrappers found")
	}

	// Check if wrapper is already configured
	configuredWrapper := viper.GetString("docker.wrapper")
	var selectedWrapper constants.DockerWrapper

	if configuredWrapper != "" {
		selectedWrapper = constants.DockerWrapper(configuredWrapper)
		fmt.Printf("Using configured wrapper: %s\n", selectedWrapper)
	} else {
		// Convert []DockerWrapper to []string for prompt
		wrapperNames := make([]string, len(wrappers))
		for i, w := range wrappers {
			wrapperNames[i] = string(w)
		}

		selected := PromptSelect("Select your preferred Docker wrapper:", wrapperNames, 0)
		selectedWrapper = constants.DockerWrapper(selected)

		// Save selection
		viper.Set("docker.wrapper", string(selectedWrapper))
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Warning: Could not save wrapper preference: %v\n", err)
		}
		fmt.Printf("✓ Saved preference: %s\n", selectedWrapper)
	}

	// Start the selected Docker wrapper
	return startWrapper(selectedWrapper)
}

func startWrapper(wrapper constants.DockerWrapper) error {
	currentOS := constants.CurrentOS()

	// Handle based on OS and wrapper combination
	switch currentOS {
	case constants.MacOS:
		return startWrapperMacOS(wrapper)
	case constants.Linux:
		return startWrapperLinux(wrapper)
	case constants.Windows:
		return startWrapperWindows(wrapper)
	default:
		return fmt.Errorf("unsupported operating system: %s", currentOS)
	}
}

// startWrapperMacOS starts Docker wrapper on macOS
func startWrapperMacOS(wrapper constants.DockerWrapper) error {
	switch wrapper {
	case constants.DockerDesktop:
		return startApp("Docker Desktop")

	case constants.RancherDesktop:
		return startApp("Rancher Desktop")

	case constants.Colima:
		fmt.Println("Starting Colima...")
		cmd := exec.Command("colima", "start")
		return cmd.Run()

	case constants.OrbStack:
		return startApp("OrbStack")

	case constants.PodmanDesktop:
		return startApp("Podman Desktop")

	default:
		return fmt.Errorf("unsupported wrapper for macOS: %s", wrapper)
	}
}

// startWrapperLinux starts Docker wrapper on Linux
func startWrapperLinux(wrapper constants.DockerWrapper) error {
	switch wrapper {
	case constants.DockerEngine:
		fmt.Println("Starting Docker Engine...")
		// Try systemctl first
		cmd := exec.Command("sudo", "systemctl", "start", "docker")
		if err := cmd.Run(); err == nil {
			return nil
		}
		// Try service command as fallback
		cmd = exec.Command("sudo", "service", "docker", "start")
		return cmd.Run()

	case constants.Podman:
		fmt.Println("Starting Podman...")
		cmd := exec.Command("sudo", "systemctl", "start", "podman")
		return cmd.Run()

	case constants.DockerDesktop:
		// Docker Desktop for Linux
		return startApp("docker-desktop")

	default:
		return fmt.Errorf("unsupported wrapper for Linux: %s", wrapper)
	}
}

// startWrapperWindows starts Docker wrapper on Windows
func startWrapperWindows(wrapper constants.DockerWrapper) error {
	switch wrapper {
	case constants.DockerDesktop:
		fmt.Println("Starting Docker Desktop...")
		// Try common Docker Desktop paths
		paths := []string{
			`C:\Program Files\Docker\Docker\Docker Desktop.exe`,
			`C:\Program Files\Docker\Docker Desktop\Docker Desktop.exe`,
		}

		for _, path := range paths {
			cmd := exec.Command(path)
			if err := cmd.Start(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("could not find Docker Desktop executable")

	case constants.RancherDesktop:
		fmt.Println("Starting Rancher Desktop...")
		cmd := exec.Command("cmd", "/C", "start", "Rancher Desktop")
		return cmd.Run()

	case constants.PodmanDesktop:
		fmt.Println("Starting Podman Desktop...")
		cmd := exec.Command("cmd", "/C", "start", "Podman Desktop")
		return cmd.Run()

	default:
		return fmt.Errorf("unsupported wrapper for Windows: %s", wrapper)
	}
}

// startApp starts a macOS application
func startApp(appName string) error {
	fmt.Printf("Starting %s...\n", appName)
	cmd := exec.Command("open", "-a", appName)
	return cmd.Start()
}

// IsDockerRunning checks if Docker is accessible
func IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// WaitForDocker waits for Docker to be ready
func WaitForDocker(maxRetries int, retryDelay time.Duration) error {
	fmt.Print("Waiting for Docker to be ready")

	for i := 0; i < maxRetries; i++ {
		if IsDockerRunning() {
			fmt.Println("\n✓ Docker is ready")
			return nil
		}
		fmt.Print(".")
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("\nDocker failed to start within timeout")
}

// CheckAndStartDocker checks if Docker is running, starts it if needed
func CheckAndStartDocker() error {
	if IsDockerRunning() {
		return nil
	}

	fmt.Println("Docker is not running. Starting...")
	if err := StartDocker(); err != nil {
		return fmt.Errorf("failed to start Docker: %w", err)
	}

	// Wait for Docker to be ready
	return WaitForDocker(60, 3*time.Second)
}
