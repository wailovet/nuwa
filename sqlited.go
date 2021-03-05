package nuwa

import (
	"fmt"
	_ "github.com/iamacarpet/go-sqlite3-dynamic"
	"github.com/wailovet/nuwa/nuwares"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"xorm.io/core"
	"xorm.io/xorm"
)

func Sqlited() *sqlite {
	println(Helper().JsonEncode(nuwares.AssetNames()))
	if runtime.GOOS == "windows" {
		_, err := getLibraryPath()
		if err != nil {
			bPath, dllName := basePath(), "sqlite3.dll"
			resdata, _ := nuwares.Asset("nuwares/static/sqlite3.dll")
			_ = ioutil.WriteFile(bPath+dllName, resdata, 0644)
		}
	}
	return &_sqlite
}

var _sqlite sqlite

type sqlite struct {
	engine   *xorm.Engine
	filename string
	prefix   string
	isLog    bool
}

func (s *sqlite) Config(filename string, prefix string, isLogs ...bool) {
	s.filename = filename
	s.prefix = prefix
	if len(isLogs) > 0 {
		s.isLog = isLogs[0]
	}
}

func (s *sqlite) Xorm() *xorm.Engine {
	if s.engine == nil {
		var err error
		s.engine, err = xorm.NewEngine("sqlite3", s.filename)
		if err != nil {
			panic("数据库访问错误")
		}
		if s.prefix != "" {
			tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, s.prefix)
			s.engine.SetTableMapper(tbMapper)
		}

		if s.isLog {
			s.engine.ShowSQL(true)
		}
	}
	return s.engine
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
