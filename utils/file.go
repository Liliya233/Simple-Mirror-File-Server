package utils

import (
	"encoding/json"
	"io"
	"os"
)

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil && stat.IsDir())
}

func MkDir(path string) bool {
	if !IsDir(path) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return false
		}
	}
	return true
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil && !stat.IsDir())
}

func GetFileData(fname string) ([]byte, error) {
	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	return buf, err
}

func WriteFileData(fname string, data []byte) error {
	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fp.Write(data); err != nil {
		return err
	}
	return nil
}

func WriteJsonData(fname string, data interface{}) error {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func GetJsonData(fname string, ptr interface{}) error {
	data, err := GetFileData(fname)
	if err != nil {
		return err
	}
	if len(data) == 0 || data == nil {
		return nil
	}
	err = json.Unmarshal(data, ptr)
	if err != nil {
		return err
	}
	return nil
}
