package nuwa

import (
	"bytes"
	"encoding/json"

	fastJson "github.com/goccy/go-json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/errors/fmt"
)

var BoltDisableFastJson = false

type BoltImp struct {
	db     *bolt.DB
	err1   error
	bucket string
	prefix string
}

var defaultBoltDbPath = "data.bolt.db"
var defaultBoltOptions *bolt.Options = nil

func NewBolt(path string, opt *bolt.Options) (ret *BoltImp) {
	db, err := bolt.Open(path, 0666, opt)
	return &BoltImp{
		db:     db,
		err1:   err,
		bucket: "main",
	}
}

var boltImp *BoltImp

func Bolt() (ret *BoltImp) {
	if boltImp == nil {
		db, err := bolt.Open(defaultBoltDbPath, 0666, defaultBoltOptions)
		boltImp = &BoltImp{
			db:     db,
			err1:   err,
			bucket: "main",
		}
	}
	return boltImp
}

func (b *BoltImp) DB() *bolt.DB {
	return b.db
}

func (b *BoltImp) Close() {
	b.db.Close()
}

func (b *BoltImp) Prefix(prefix string) *BoltImp {
	return &BoltImp{
		db:     b.db,
		err1:   b.err1,
		bucket: b.bucket,
		prefix: prefix,
	}
}

func (b *BoltImp) Bucket(bucket string) *BoltImp {
	return &BoltImp{
		db:     b.db,
		err1:   b.err1,
		bucket: bucket,
	}
}

func (b *BoltImp) Load(key string, val interface{}) error {
	if b.db == nil {
		return b.err1
	}
	return b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(key))
		if BoltDisableFastJson {
			return json.Unmarshal(v, val)
		}
		return fastJson.Unmarshal(v, val)
	})

}

func (b *BoltImp) Get(key string) gjson.Result {
	data := b.GetRaw(key)
	return gjson.Parse(data)
}

func (b *BoltImp) Delete(key string) error {
	return b.Delete(key)
}

func (b *BoltImp) GetRaw(key string) string {
	if b.db == nil {
		fmt.Println("bolt error:", b.err1)
		return ""
	}
	ret := ""
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		if b != nil {
			v := b.Get([]byte(key))
			ret = string(v)
		}
		return nil
	})
	return ret
}

func (b *BoltImp) Range(cb func(k, v string) error) {
	if b.db == nil {
		fmt.Println("bolt error:", b.err1)
		return
	}
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		if b != nil {
			b.ForEach(func(k, v []byte) error {
				return cb(string(k), string(v))
			})
		}
		return nil
	})
}

func (b *BoltImp) Set(key string, val interface{}) error {
	if b.db == nil {
		return b.err1
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(b.bucket))
		if err != nil {
			return err
		}
		var data []byte
		if BoltDisableFastJson {
			data, err = json.Marshal(val)
		} else {
			data, err = fastJson.Marshal(val)
		}
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})
}

func (b *BoltImp) Page(val interface{}, page int, limit int) error {
	if page > 0 {
		page = page - 1
	}
	return b.Scan(val, page*limit, limit)
}

func (b *BoltImp) Scan(val interface{}, offsetNum int, limitNum int) error {
	if b.db == nil {
		return b.err1
	}
	ret := "[]"
	index := 0
	err := b.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(b.bucket))
		if bk != nil {
			c := bk.Cursor()
			for k, v := c.Seek([]byte(b.prefix)); k != nil && bytes.HasPrefix(k, []byte(b.prefix)); k, v = c.Next() {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				reti, err := sjson.SetRaw(ret, fmt.Sprint(index), string(v))
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
