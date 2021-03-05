package nuwa

import (
	_ "github.com/iamacarpet/go-sqlite3-dynamic"
	"xorm.io/core"
	"xorm.io/xorm"
)

var sqlitedPreFunc = func() {}

func Sqlited() *sqlite {
	sqlitedPreFunc()
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
