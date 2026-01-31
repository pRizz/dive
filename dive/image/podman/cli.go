//go:build linux || darwin

package podman

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pRizz/dive/internal/log"
	"github.com/pRizz/dive/internal/utils"
)

// runPodmanCmd runs a given Podman command in the current tty
func runPodmanCmd(cmdStr string, args ...string) error {
	if !isPodmanClientBinaryAvailable() {
		return fmt.Errorf("cannot find podman client executable")
	}

	allArgs := utils.CleanArgs(append([]string{cmdStr}, args...))

	fullCmd := strings.Join(append([]string{"docker"}, allArgs...), " ")
	log.WithFields("cmd", fullCmd).Trace("executing")

	cmd := exec.Command("podman", allArgs...)
	cmd.Env = os.Environ()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func streamPodmanCmd(args ...string) (io.Reader, error) {
	if !isPodmanClientBinaryAvailable() {
		return nil, fmt.Errorf("cannot find podman client executable")
	}

	allArgs := utils.CleanArgs(args)
	fullCmd := strings.Join(append([]string{"docker"}, allArgs...), " ")
	log.WithFields("cmd", fullCmd).Trace("executing (streaming)")

	cmd := exec.Command("podman", allArgs...)
	cmd.Env = os.Environ()

	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := writer.Close(); closeErr != nil {
			log.WithFields("error", closeErr).Warn("failed to close podman pipe writer")
		}
	}()

	cmd.Stdout = writer
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return reader, nil
}

func isPodmanClientBinaryAvailable() bool {
	_, err := exec.LookPath("podman")
	return err == nil
}
