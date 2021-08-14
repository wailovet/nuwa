package nuwa

import (
	"fmt"
	"testing"
	"time"
)

func TestNutsdbImp_Scan(t *testing.T) {
	var ts []time.Time
	b := NutsDB().Prefix("test")
	for i := 0; i < 100; i++ {
		b.Set(fmt.Sprint(i), time.Now())
	}

	b.Page(&ts, 1, 30)
	fmt.Println(len(ts))
	b.Page(&ts, 2, 30)
	fmt.Println(len(ts))
	b.Page(&ts, 3, 30)
	fmt.Println(len(ts))
	b.Page(&ts, 4, 30)
	fmt.Println(len(ts))
}
