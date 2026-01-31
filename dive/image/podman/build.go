//go:build linux || darwin

package podman

import (
	"os"

	"github.com/wagoodman/dive/internal/log"
)

func buildImageFromCli(buildArgs []string) (string, error) {
	iidfile, err := os.CreateTemp("/tmp", "dive.*.iid")
	if err != nil {
		return "", err
	}
	defer func() {
		if removeErr := os.Remove(iidfile.Name()); removeErr != nil {
			log.WithFields("error", removeErr, "path", iidfile.Name()).Warn("failed to remove podman iidfile")
		}
	}()
	defer func() {
		if closeErr := iidfile.Close(); closeErr != nil {
			log.WithFields("error", closeErr, "path", iidfile.Name()).Warn("failed to close podman iidfile")
		}
	}()

	allArgs := append([]string{"--iidfile", iidfile.Name()}, buildArgs...)
	err = runPodmanCmd("build", allArgs...)
	if err != nil {
		return "", err
	}

	imageId, err := os.ReadFile(iidfile.Name())
	if err != nil {
		return "", err
	}

	return string(imageId), nil
}
