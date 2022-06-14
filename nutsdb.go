package nuwa

import (
	"encoding/json"
	"errors"

	"fmt"

	fastJson "github.com/goccy/go-json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xujiajun/nutsdb"
	"github.com/xujiajun/nutsdb/inmemory"
)

var NutsdbDisableFastJson = false

type NutsdbImp struct {
	db         *nutsdb.DB
	inmemorydb *inmemory.DB
	err1       error
	bucket     string
	prefix     string
}

var defaultNutsdbDbPath = "data.nutsdb.db"
var DefaultNutsdbOptions = nutsdb.Options{
	EntryIdxMode:         nutsdb.HintBPTSparseIdxMode,
	SegmentSize:          8 * 1024 * 1024,
	NodeNum:              1,
	RWMode:               nutsdb.FileIO,
	SyncEnable:           true,
	StartFileLoadingMode: nutsdb.MMap,
}

func NewNutsdb(path string, opt nutsdb.Options) (ret *NutsdbImp) {
	opt.Dir = path
	db, err := nutsdb.Open(opt)
	return &NutsdbImp{
		db:     db,
		err1:   err,
		bucket: "main",
	}
}
func NewNutsdbMemory(opt inmemory.Options) (ret *NutsdbImp) {
	db, err := inmemory.Open(opt)
	return &NutsdbImp{
		inmemorydb: db,
		err1:       err,
		bucket:     "main",
	}
}

var nutsdbImp *NutsdbImp

func NutsDB() (ret *NutsdbImp) {
	if nutsdbImp == nil {
		nutsdbImp = NewNutsdb(defaultNutsdbDbPath, DefaultNutsdbOptions)
	}
	return nutsdbImp
}

func (b *NutsdbImp) DB() *nutsdb.DB {
	return b.db
}

func (b *NutsdbImp) Close() {
	if b.inmemorydb != nil {
		return
	}
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

	if b.inmemorydb != nil {
		v, err := b.inmemorydb.Get(b.bucket, []byte(b.prefix+key))
		if err != nil {
			return err
		}
		if NutsdbDisableFastJson {
			return json.Unmarshal(v.Value, val)
		}
		return fastJson.Unmarshal(v.Value, val)
	}

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

	if b.inmemorydb != nil {
		v, err := b.inmemorydb.Get(b.bucket, []byte(b.prefix+key))
		if err != nil {
			return ""
		}
		return string(v.Value)
	}

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

func (b *NutsdbImp) Set(key string, val interface{}, ttls ...uint32) error {
	var ttl uint32
	if len(ttls) > 0 {
		ttl = ttls[0]
	}
	if b.inmemorydb != nil {
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
		return b.inmemorydb.Put(b.bucket, []byte(b.prefix+key), data, ttl)
	}

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

func (b *NutsdbImp) Delete(key string) error {

	if b.inmemorydb != nil {
		b.inmemorydb.Delete(b.bucket, []byte(b.prefix+key))
	}
	return b.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(b.bucket, []byte(b.prefix+key))
	})
}

func (b *NutsdbImp) Page(val interface{}, page int, limit int) error {
	if page > 0 {
		page = page - 1
	}
	return b.Scan(val, page*limit, limit)
}

func (b *NutsdbImp) Scan(val interface{}, offsetNum int, limitNum int) error {

	ret := "[]"
	index := 0
	if b.inmemorydb != nil {
		return errors.New("不支持")
	}

	if b.db == nil {
		return b.err1
	}
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
