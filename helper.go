package nuwa

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var helper = helperImp{}

func Helper() *helperImp {
	return &helper
}

type helperImp struct {
}

func (h *helperImp) GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(fmt.Sprint(path))
	}
	return string(path[0 : i+1]), nil
}
