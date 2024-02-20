package network

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func MoveSCU(host, aet, aec, aem, filters string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "dicomizer")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("movescu")
	cmd.Args = append(cmd.Args, host)
	cmd.Args = append(cmd.Args, strconv.FormatInt(104, 10))
	cmd.Args = append(cmd.Args, "--port", strconv.FormatInt(104, 10))

	cmd.Args = append(cmd.Args, "--aetitle", aet)
	cmd.Args = append(cmd.Args, "--call", aec)
	cmd.Args = append(cmd.Args, "--move", aem)

	cmd.Args = append(cmd.Args, "-od", tmpDir)
	cmd.Args = append(cmd.Args, "-S")
	cmd.Args = append(cmd.Args, "+xi")

	cmd.Args = append(cmd.Args, "-k", "0008,0052=STUDY")
	for _, filter := range strings.Split(filters, ";") {
		cmd.Args = append(cmd.Args, "-k", filter)
	}

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run movescu: %s", err)
	}

	return tmpDir, nil
}
