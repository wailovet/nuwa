package nuwa

import (
	"encoding/json"

	"fmt"

	fastJson "github.com/goccy/go-json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xujiajun/nutsdb"
)

var NutsdbDisableFastJson = false

type NutsdbImp struct {
	db     *nutsdb.DB
	err1   error
	bucket string
	prefix string
}

var defaultNutsdbDbPath = "data.nutsdb.db"
var defaultNutsdbOptions = nutsdb.DefaultOptions

func NewNutsdb(path string, opt nutsdb.Options) (ret *NutsdbImp) {
	opt.Dir = path
	db, err := nutsdb.Open(opt)
	return &NutsdbImp{
		db:     db,
		err1:   err,
		bucket: "main",
	}
}

var nutsdbImp *NutsdbImp

func NutsDB() (ret *NutsdbImp) {
	if nutsdbImp == nil {
		nutsdbImp = NewNutsdb(defaultNutsdbDbPath, defaultNutsdbOptions)
	}
	return nutsdbImp
}

func (b *NutsdbImp) Close() {
	b.db.Close()
}

func (b *NutsdbImp) Prefix(prefix string) *NutsdbImp {
	return &NutsdbImp{
		db:     b.db,
		err1:   b.err1,
		bucket: b.bucket,
		prefix: prefix,
	}
}
func (b *NutsdbImp) Bucket(bucket string) *NutsdbImp {
	return &NutsdbImp{
		db:     b.db,
		err1:   b.err1,
		bucket: bucket,
		prefix: b.prefix,
	}
}

func (b *NutsdbImp) Load(key string, val interface{}) error {
	if b.db == nil {
		return b.err1
	}
	return b.db.View(func(tx *nutsdb.Tx) error {
		v, err := tx.Get(b.bucket, []byte(b.prefix+key))
		if err != nil {
			return err
		}
		if NutsdbDisableFastJson {
			return json.Unmarshal(v.Value, val)
		}
		return fastJson.Unmarshal(v.Value, val)
	})

}

func (b *NutsdbImp) Get(key string) gjson.Result {
	data := b.GetRaw(key)
	return gjson.Parse(data)
}

func (b *NutsdbImp) GetRaw(key string) string {
	if b.db == nil {
		fmt.Println("nutsdb error:", b.err1)
		return ""
	}
	ret := ""
	b.db.View(func(tx *nutsdb.Tx) error {
		v, err := tx.Get(b.bucket, []byte(b.prefix+key))
		if err != nil {
			return err
		}
		ret = string(v.Value)
		return nil
	})
	return ret
}

func (b *NutsdbImp) Set(key string, val interface{}) error {
	if b.db == nil {
		return b.err1
	}
	return b.db.Update(func(tx *nutsdb.Tx) error {
		var err error
		var data []byte

		if NutsdbDisableFastJson {
			data, err = json.Marshal(val)
		} else {
			data, err = fastJson.Marshal(val)
		}
		if err != nil {
			return err
		}
		return tx.Put(b.bucket, []byte(b.prefix+key), data, 0)
	})
}

func (b *NutsdbImp) Page(val interface{}, page int, limit int) error {
	if page > 0 {
		page = page - 1
	}
	return b.Scan(val, page*limit, limit)
}

func (b *NutsdbImp) Delete(key string) error {
	return b.Delete(key)
}

func (b *NutsdbImp) Scan(val interface{}, offsetNum int, limitNum int) error {
	if b.db == nil {
		return b.err1
	}
	ret := "[]"
	index := 0
	err := b.db.View(func(tx *nutsdb.Tx) error {
		// Constrain 100 entries returned
		if entries, _, err := tx.PrefixScan(b.bucket, []byte(b.prefix), offsetNum, limitNum); err != nil {
			return err
		} else {
			for _, entry := range entries {
				// fmt.Println(string(entry.Key), string(entry.Value))
				reti, err := sjson.SetRaw(ret, fmt.Sprint(index), string(entry.Value))
				if err == nil {
					ret = reti
					index++
				}
			}
		}
		return nil
	})

	if NutsdbDisableFastJson {
		json.Unmarshal([]byte(ret), val)
	} else {
		fastJson.Unmarshal([]byte(ret), val)
	}
	return err

}
