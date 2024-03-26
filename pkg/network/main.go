package network

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func findSCU(host, aet, aec, filters string) (string, error) {
  tmpDir, err := os.MkdirTemp("", "dicomizer")
  if err != nil {
    return "", err
  }

  cmd := exec.Command("findscu")
  cmd.Args = append(cmd.Args, host)
  cmd.Args = append(cmd.Args, strconv.FormatInt(104, 10))

  cmd.Args = append(cmd.Args, "--aetitle", aet)
  cmd.Args = append(cmd.Args, "--call", aec)

  cmd.Args = append(cmd.Args, "-od", tmpDir)
  cmd.Args = append(cmd.Args, "-S", "-X")

  cmd.Args = append(cmd.Args, "-k", "0008,0052=STUDY")
  cmd.Args = append(cmd.Args, "-k", "0020,000D") // include StudyInstanceUID in response DCM file
	for _, filter := range strings.Split(filters, ";") {
		cmd.Args = append(cmd.Args, "-k", filter)
	}

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run findscu: %s", err)
	}

  return tmpDir, nil
}

func MoveSCU(host, aet, aec, aem, filters string) (string, error) {
  tmpDir, err := findSCU(host, aet, aec, filters)
  if err != nil {
    return "", err
  }

  dirEntries, err := os.ReadDir(tmpDir)
  if err != nil {
    return "", err
  }

  for i, f := range dirEntries {
    if f.IsDir() {
      panic("Found directory in newly created temporary directory")
    }

    fmt.Printf("Processing batch: %d\n", i+1)

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

    cmd.Args = append(cmd.Args, filepath.Join(tmpDir, f.Name()))

    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
      return "", fmt.Errorf("failed to run movescu: %s", err)
    }
  }

	return tmpDir, nil
}
