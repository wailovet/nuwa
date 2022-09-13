package nuwa

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Tar struct {
	fileName string
	src      []string
}

func NewTar(fileName string) *Tar {
	return &Tar{
		fileName: fileName,
	}
}

func (t *Tar) Add(src string) {
	t.src = append(t.src, src)
}

func (t *Tar) Create2GzipMemory() (*bytes.Buffer, error) {
	fw := bytes.NewBuffer(nil)
	in, err := t.Create2Memory()
	if err != nil {
		return nil, err
	}
	w := gzip.NewWriter(fw)
	w.Write(in.Bytes())
	w.Close()
	return fw, err
}

func (t *Tar) Create2GzipFile() error {
	buf, err := t.Create2GzipMemory()
	if err != nil {
		return err
	}
	fn := t.fileName

	if filepath.Ext(t.fileName) != ".gz" {
		fn = fn + ".gz"
	}
	return ioutil.WriteFile(fn, buf.Bytes(), 0644)
}

func (t *Tar) Create2Memory() (*bytes.Buffer, error) {
	fw := bytes.NewBuffer(nil)

	// 通过fw创建一个tar.Writer
	tw := tar.NewWriter(fw)
	// 如果关闭失败会造成tar包不完整
	defer func() {
		if err := tw.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, fileName := range t.src {
		fi, err := os.Stat(fileName)
		if err != nil {
			log.Println(err)
			continue
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		hdr.Name = fileName

		// 将tar的文件信息hdr写入到tw
		err = tw.WriteHeader(hdr)
		if err != nil {
			return fw, err
		}

		// 将文件数据写入
		fs, err := os.Open(fileName)
		if err != nil {
			return fw, err
		}
		if _, err = io.Copy(tw, fs); err != nil {
			return fw, err
		}
		fs.Close()
	}
	return fw, nil
}

func (t *Tar) Create2File() error {
	buf, err := t.Create2Memory()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(t.fileName, buf.Bytes(), 0644)
}
