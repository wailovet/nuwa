package nuwa

import (
	"fmt"
	_ "github.com/iamacarpet/go-sqlite3-dynamic"
	"github.com/wailovet/nuwa/nuwares"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	sqlitedPreFunc = func() {
		_, err := getLibraryPath()
		if err != nil {
			bPath, dllName := basePath(), "sqlite3.dll"
			resdata, _ := nuwares.Asset("nuwares/static/sqlite3.dll")

			if exist, _ := exists(bPath + "support"); !exist {
				_ = os.Mkdir(bPath+"support", os.ModeDir)
			}
			filePath := bPath + "support" + string(os.PathSeparator) + dllName
			_ = ioutil.WriteFile(filePath, resdata, 0644)
		}
	}
}

func getLibraryPath() (string, error) {
	bPath, dllName := basePath(), "sqlite3.dll"

	if exist, _ := exists(dllName); exist {
		return dllName, nil
	}

	filePath := bPath + dllName
	if exist, _ := exists(filePath); exist {
		return filePath, nil
	}

	filePath = bPath + "support" + string(os.PathSeparator) + dllName
	if exist, _ := exists(filePath); exist {
		return filePath, nil
	}

	return "", fmt.Errorf("%s not found.", dllName)
}

func basePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	return dir + string(os.PathSeparator)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
