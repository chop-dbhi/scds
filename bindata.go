package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _email_changed_object_email_body_txt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xf2\x4e\xad\xb4\x52\xa8\xae\xd6\x03\xd2\xb5\xb5\x5c\x61\xa9\x45\xc5\x99\xf9\x79\x60\x11\x28\x1b\x28\x1a\x92\x99\x9b\x0a\x16\x02\x31\x80\xfc\xd0\x20\x1f\x30\x17\x48\x23\xf4\x28\xc0\x44\xa1\x7c\x88\x24\x57\x75\x75\x66\x9a\x82\x9e\x73\x46\x62\x5e\x7a\x6a\x31\x50\x40\x59\x01\xca\x06\x49\x21\xc4\xab\xab\x53\xf3\x52\x40\x14\x48\xb5\x63\x4a\x4a\x66\x09\xd0\x08\x88\x7a\x38\x0f\xac\x03\x49\x0e\x45\x4f\x50\x6a\x6e\x7e\x59\x62\x0e\x44\x0b\x8c\x03\xd6\x81\x90\x81\x6a\xe0\x02\x04\x00\x00\xff\xff\x3f\x78\xb2\x1d\xf4\x00\x00\x00")

func email_changed_object_email_body_txt_bytes() ([]byte, error) {
	return bindata_read(
		_email_changed_object_email_body_txt,
		"email/changed_object_email_body.txt",
	)
}

func email_changed_object_email_body_txt() (*asset, error) {
	bytes, err := email_changed_object_email_body_txt_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "email/changed_object_email_body.txt", size: 244, mode: os.FileMode(420), modTime: time.Unix(1438253661, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _email_new_object_email_body_txt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xf2\x4e\xad\xb4\x52\xa8\xae\xd6\x03\xd2\xb5\xb5\x5c\x61\xa9\x45\xc5\x99\xf9\x79\x60\x11\x28\x1b\x28\x1a\x92\x99\x9b\x0a\x16\x02\x31\x80\xfc\xd0\x20\x1f\x30\x17\x48\x23\xf4\x28\xc0\x44\xa1\x7c\x88\x24\x97\xb2\x82\x7f\x52\x56\x6a\x72\x09\x17\x17\x50\x0a\xc2\x04\x0a\x03\x02\x00\x00\xff\xff\x5e\xe9\x44\x05\x76\x00\x00\x00")

func email_new_object_email_body_txt_bytes() ([]byte, error) {
	return bindata_read(
		_email_new_object_email_body_txt,
		"email/new_object_email_body.txt",
	)
}

func email_new_object_email_body_txt() (*asset, error) {
	bytes, err := email_new_object_email_body_txt_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "email/new_object_email_body.txt", size: 118, mode: os.FileMode(420), modTime: time.Unix(1438253632, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"email/changed_object_email_body.txt": email_changed_object_email_body_txt,
	"email/new_object_email_body.txt": email_new_object_email_body_txt,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"email": &_bintree_t{nil, map[string]*_bintree_t{
		"changed_object_email_body.txt": &_bintree_t{email_changed_object_email_body_txt, map[string]*_bintree_t{
		}},
		"new_object_email_body.txt": &_bintree_t{email_new_object_email_body_txt, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

