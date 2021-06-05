package main

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/qq51529210/cloud-service/util"
	"github.com/qq51529210/log"
)

const (
	defaultuploadFileBuffer         = 1024 * 64
	defaultuploadFileTickDuration   = 30
	defaultdownloadFileBuffer       = 1024 * 64
	defaultdownloadFileTickDuration = 30
)

var (
	uploadFilePool   sync.Pool
	downloadFilePool sync.Pool
)

func init() {
	uploadFilePool.New = func() interface{} {
		p := new(UploadFile)
		p.Hash = sha1.New()
		p.buff = make([]byte, defaultuploadFileBuffer)
		return p
	}
	downloadFilePool.New = func() interface{} {
		p := new(downloadFile)
		return p
	}
}

// To handle upload file.
type UploadFile struct {
	*os.File
	hash.Hash
	buff []byte
}

// Read data from r.
func (f *UploadFile) Save(r io.Reader, dir, namespace, name string, rate, dur int) (err error) {
	// Make sure there is a directory for new file.
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	// Create file.
	filePath := filepath.Join(dir, f.TempName(namespace, name))
	f.File, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			// If there is any error, remove file.
			re := os.Remove(filePath)
			if re != nil {
				log.Error(re)
			}
		}
	}()
	// Save data to file.
	f.Hash.Reset()
	_, err = util.LimitRateCopy(f, r, f.buff, rate, dur)
	if err != nil {
		f.File.Close()
		return
	}
	err = f.File.Close()
	if err != nil {
		return
	}
	// New file name by file hex hash value.
	hashFilePath := filepath.Join(dir, hex.EncodeToString(f.Hash.Sum(f.buff[:0])))
	// Rename file temp name if not exist.
	_, err = os.Stat(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Rename(filePath, hashFilePath)
			if err != nil {
				return
			}
		}
	} else {
		// File extists, remove temp file.
		err = os.Remove(filePath)
		if err != nil {
			return
		}
	}
	// Create a symbolic link if not exist.
	linkName := filepath.Join(dir, name)
	_, err = os.Stat(linkName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Symlink(hashFilePath, linkName)
		}
	}
	return
}

// Return file temp name.
func (f *UploadFile) TempName(namespace, name string) string {
	var str strings.Builder
	str.WriteString(namespace)
	str.WriteByte('.')
	str.WriteString(name)
	str.WriteByte('.')
	str.WriteString(time.Now().Format("20060102150405.000"))
	return str.String()
}

// For io.Copy
func (f *UploadFile) ReadFrom(r io.Reader) (n int64, err error) {
	var m int
	for {
		m, err = r.Read(f.buff)
		if err != nil {
			if err == io.EOF {
				f.Write(f.buff[:m])
				n += int64(m)
				err = nil
				return
			}
		}
		f.Write(f.buff[:m])
		n += int64(m)
	}
}

func (f *UploadFile) Write(b []byte) (int, error) {
	f.Hash.Write(b)
	return f.File.Write(b)
}

// To handle download file.
type downloadFile struct {
}

func (f *downloadFile) WriteTo(w io.Writer, dir, namespace, name string, dur int) (err error) {
	return nil
}
