package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getFilesystemType(dev string) (string, error) {
	out, err := exec.Command("blkid", "-s", "TYPE", "-o", "value", dev).CombinedOutput()

	if err != nil {
		if len(out) == 0 {
			return "", nil
		}

		return "", errors.New(string(out))
	}

	return string(out), nil
}

func formatFilesystem(dev string, label string) error {
	out, err := exec.Command("mkfs.ext4", "-L", label, dev).CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// look for a device which name contains id, under dir
// and return the full path+filename
func waitForDevice(dir string, id string) (string, error) {

	for i := 0; i <= 10; i++ {

        files, err := os.ReadDir(dir)
        if err != nil {
            return "", err
        }
    
        for _, file := range files {
            if strings.Contains(file.Name(), id) {
				return fmt.Sprintf("%s/%s", dir, file.Name()), nil
            }
        }

		time.Sleep(500 * time.Millisecond)
	}

	return "", fmt.Errorf("Timeout waiting for file: %s", id)
}

func isDirectoryPresent(path string) (bool, error) {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return stat.IsDir(), nil
	}
}
