package nuwa

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/yumaojun03/dmidecode"
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

var aesKey = []byte("0123456789abcdef")

func (h *helperImp) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (h *helperImp) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (h *helperImp) DefaultAesEncrypt(origData []byte) ([]byte, error) {
	return h.AesEncrypt(origData, aesKey)
}

func (h *helperImp) DefaultAesDecrypt(crypted []byte) ([]byte, error) {
	return h.AesDecrypt(crypted, aesKey)
}

func (h *helperImp) AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = h.PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (h *helperImp) AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = h.PKCS7UnPadding(origData)
	return origData, nil
}

func (h *helperImp) Md5(data string) string {
	sum := md5.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

func (h *helperImp) Md5WithByte(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func (h *helperImp) PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (h *helperImp) GetSelfFilePath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return ""
	}
	path, _ = filepath.Abs(string(path[0 : i+1]))
	return path
}

func (h *helperImp) CleanSuperfluousSpace(s string) string {
	for strings.Index(s, "  ") > -1 {
		s = strings.Replace(s, "  ", " ", -1)
	}
	return strings.TrimSpace(s)
}

func (h *helperImp) Md5ToString(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (h *helperImp) GetLocalIP() string {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
		err     error
		ipv4    string
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return ""
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
			}
		}
	}
	return ipv4
}

func (h *helperImp) JsonByFile(file string, v interface{}) {
	data, _ := ioutil.ReadFile(file)
	json.Unmarshal(data, v)
}

func (h *helperImp) JsonToFile(file string, v interface{}) bool {
	data, err := json.Marshal(v)
	if err != nil {
		return false
	}
	return ioutil.WriteFile(file, data, 0644) == nil
}

func (h *helperImp) JsonEncode(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (h *helperImp) GetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func (h *helperImp) Interface2Struct(in interface{}, out interface{}) {
	raw, _ := json.Marshal(in)
	json.Unmarshal(raw, &out)
}
func (h *helperImp) Interface2Map(in interface{}, out map[string]interface{}) {

	raw, _ := json.Marshal(in)
	json.Unmarshal(raw, &out)
}

func (h *helperImp) Interface2Interface(in interface{}) (out interface{}) {
	raw, _ := json.Marshal(in)
	json.Unmarshal(raw, &out)
	return
}

func (h *helperImp) CleanExtraCharacters(a string, b string) string {
	for strings.Index(a, b+b) > -1 {
		a = strings.Replace(a, b+b, b, -1)
	}
	return a
}

func InArray(t string, arr []string) bool {

	for e := range arr {
		if arr[e] == t {
			return true
		}
	}
	return false
}

func (h *helperImp) Unzip(file string, path string) error {
	File, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer File.Close()
	for _, v := range File.File {
		info := v.FileInfo()
		fileName, _ := filepath.Abs(path + "/" + v.Name)
		_ = os.RemoveAll(fileName)
		if info.IsDir() {
			err := os.MkdirAll(fileName, 0777)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		srcFile, err := v.Open()
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer srcFile.Close()

		newFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}

func (h *helperImp) RandomInt(length int) int {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte

	for i := 0; i < length; {
		r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i) + int64(time.Now().Nanosecond())))
		b := bytes[r.Intn(len(bytes))]
		if i == 0 && b == '0' {
			continue
		}
		result = append(result, b)
		i++
	}
	num, _ := strconv.Atoi(string(result))
	return num
}

func (h *helperImp) RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < length; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i) + int64(time.Now().Nanosecond())))
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (h *helperImp) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func (h *helperImp) AbsPath(pathname string) string {
	s, _ := filepath.Abs(pathname)
	return s
}

func (h *helperImp) DirRange(pathname string, cb func(pathItem string)) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			h.DirRange(h.AbsPath(pathname+fi.Name()+"\\"), cb)
		} else {
			cb(h.AbsPath(pathname + fi.Name() + "\\"))
		}
	}
	return err
}

func (h *helperImp) ReadFileContent(path string) string {
	r, _ := ioutil.ReadFile(path)
	return string(r)
}

func (h *helperImp) WriteFileContent(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0754)
}

func (h *helperImp) MacID() string {
	dmi, err := dmidecode.New()
	if err != nil {
		return ""
	}
	infos, err := dmi.System()
	if err != nil {
		return ""
	}
	for i := range infos {
		return infos[i].UUID
	}
	return ""
}

func (h *helperImp) SendToMail(host, username, password, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", username, password, hp[0])

	e := email.NewEmail()
	e.From = username
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	return e.SendWithTLS(host, auth, &tls.Config{ServerName: hp[0]})
}

func (h *helperImp) DecompressTarGz(targzBuff *bytes.Buffer, target string) error {
	gr, err := gzip.NewReader(targzBuff)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	// 读取文件
	for {
		hn, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileTargz := filepath.Join(target, hn.Name)
		// 显示文件
		fmt.Println(fileTargz)
		// 打开文件

		fileDir := filepath.Dir(fileTargz)
		if !h.PathExists(fileDir) {
			os.MkdirAll(fileDir, os.ModeDir)
		}

		buf := bytes.NewBuffer(nil)
		// 写文件
		_, err = io.Copy(buf, tr)
		if err != nil {
			return err
		}

		ioutil.WriteFile(fileTargz, buf.Bytes(), hn.FileInfo().Mode())
	}
	return nil
}
