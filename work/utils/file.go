package utils

// @see https://golang.org/pkg/path/filepath/

import (
	"os"
	"path/filepath"
)

//-------------------------------------------------------------------------------------------------

// FileExist file_path p
func FileExist(p string) bool {
	i, e := os.Stat(p)
	if os.IsNotExist(e) {
		return false
	}
	return !i.IsDir()
}

//-------------------------------------------------------------------------------------------------

// FileRead
//
// ex. "./aaa/bbb/ccc/target.json"
//
// utils.FileRead("./aaa/bbb/ccc","target.json")
//
// utils.FileRead("./aaa/bbb/ccc/target.json")
func FileRead(path ...string) ([]byte, error) {
	p := filepath.Join(path...)
	d, e := os.ReadFile(p)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//-------------------------------------------------------------------------------------------------

// FileWrite overwrite content
func FileWrite(dirPath, fileName string, content []byte) error {
	e := os.MkdirAll(dirPath, os.ModePerm)
	if e != nil {
		return e
	}

	p := filepath.Join(dirPath, fileName)
	return os.WriteFile(p, content, os.ModePerm)
}

// FileWriteAppend append content
func FileWriteAppend(dirPath, fileName string, content []byte) error {
	e := os.MkdirAll(dirPath, os.ModePerm)
	if e != nil {
		return e
	}

	p := filepath.Join(dirPath, fileName)
	_, e = os.Stat(p)

	if e == nil { // 檔案已存在

		f, e := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if e != nil {
			return e
		}
		defer f.Close()

		_, e = f.Write(content)
		return e

	} else if os.IsNotExist(e) { // 檔案不存在

		return os.WriteFile(p, content, os.ModePerm)
	}
	return e
}
