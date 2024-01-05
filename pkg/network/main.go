package network

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/suyashkumar/dicom"
)

func findscu(tmpDir, host string, port uint, aet, aec string, dataset *dicom.Dataset) error {
	cmd := exec.Command("findscu")
	cmd.Args = append(cmd.Args, host)
	cmd.Args = append(cmd.Args, strconv.FormatUint(uint64(port), 10))

	cmd.Args = append(cmd.Args, "-aet "+aet)
	cmd.Args = append(cmd.Args, "-aec "+aec)

	cmd.Args = append(cmd.Args, "-P")

	cmd.Args = append(cmd.Args, "-od "+tmpDir)
	cmd.Args = append(cmd.Args, "-X")
	for _, el := range dataset.Elements {
		elStr := fmt.Sprintf("%d,%d", el.Tag.Group, el.Tag.Element)
		if v := el.Value.String(); v != "" {
			elStr += ("=" + v)
		}
		cmd.Args = append(cmd.Args, "-k")
		cmd.Args = append(cmd.Args, elStr)
	}

	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

func movescu(tmpDir, host string, port uint, aet, aec, aem string) (string, error) {
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return "", err
	}

	outDir := filepath.Join(tmpDir, "out")
	if err := os.Mkdir(outDir, 0755); err != nil {
		return "", err
	}

	for _, file := range files {
		cmd := exec.Command("movescu")
		cmd.Args = append(cmd.Args, host)
		cmd.Args = append(cmd.Args, strconv.FormatUint(uint64(port), 10))

		cmd.Args = append(cmd.Args, "-aet "+aet)
		cmd.Args = append(cmd.Args, "-aec "+aec)
		cmd.Args = append(cmd.Args, "-aem "+aem)

		cmd.Args = append(cmd.Args, "--port "+strconv.FormatInt(int64(104), 10))

		cmd.Args = append(cmd.Args, "-od "+tmpDir)

		cmd.Args = append(cmd.Args, "-k 0008,0052=IMAGE")

		cmd.Args = append(cmd.Args, filepath.Join(outDir, file.Name()))

		if err := cmd.Run(); err != nil {
			return "", err
		}
	}

	return outDir, nil
}
