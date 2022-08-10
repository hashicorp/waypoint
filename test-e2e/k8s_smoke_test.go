package test

import (
	"fmt"
	"strings"
	"testing"
)

var (
	kubernetesTestDir = fmt.Sprintf("%s/kubernetes/nodejs", examplesRootDir)
)

func TestWaypointKubernetesInstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("install", "-platform=kubernetes", "-accept-tos", fmt.Sprintf("-k8s-server-image=%s", wpServerImage), fmt.Sprintf("-k8s-odr-image=%s", wpOdrImage))

	if err != nil {
		t.Errorf("unexpected error installing server to kubernetes: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output installing server to kubernetes: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully installed and configured!") {
		t.Errorf("No success message detected after kubernetes server install:\n%s", stdout)
	}
}

func TestWaypointKubernetesUp(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("init")

	if err != nil {
		t.Errorf("unexpected error initializing waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output initializing waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Project initialized!") {
		t.Errorf("No success message detected after initializing project:\n%s", stdout)
	}

	stdout, stderr, err = wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointKubernetesUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("server", "upgrade", "-platform=kubernetes", "-auto-approve", fmt.Sprintf("-k8s-server-image=%s", wpServerImageUpgrade), fmt.Sprintf("-k8s-odr-image=%s", wpOdrImageUpgrade), "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error upgrading server in kubernetes: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output upgrading server in kubernetes: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint has finished upgrading the server") {
		t.Errorf("No success message detected after kubernetes server install:\n%s", stdout)
	}
}

func TestWaypointKubernetesUpAfterUpgrade(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("up")

	if err != nil {
		t.Errorf("unexpected error deploying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output deploying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "The deploy was successful!") {
		t.Errorf("No success message detected after deploying project:\n%s", stdout)
	}
}

func TestWaypointKubernetesDestroy(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("destroy")

	if err != nil {
		t.Errorf("unexpected error destroying waypoint project: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output destroying waypoint project: %s", stderr)
	}

	if !strings.Contains(stdout, "Destroy successful!") {
		t.Errorf("No success message detected after destroying project:\n%s", stdout)
	}
}

func TestWaypointKubernetesUninstall(t *testing.T) {
	wp := NewBinary(t, wpBinary, kubernetesTestDir)
	stdout, stderr, err := wp.RunRaw("server", "uninstall", "-platform=kubernetes", "-auto-approve", "-snapshot=false")

	if err != nil {
		t.Errorf("unexpected error uninstalling waypoint server: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output uninstalling waypoint server: %s", stderr)
	}

	if !strings.Contains(stdout, "Waypoint server successfully uninstalled") {
		t.Errorf("No success message detected after uninstalling server:\n%s", stdout)
	}
}
