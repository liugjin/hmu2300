/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/08/28
 * Despcription: upgrade
 *
 */

package public

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadUpgradePackage download upgrade package
func DownloadUpgradePackage(url, savepath string) error {
	// open save file
	f, err := os.Create(savepath)
	if err != nil {
		return fmt.Errorf("open save file {%s} fail, errmsg {%v}", savepath, err)
	}
	defer f.Close()

	// get resource
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get upgrade package fail, errmsg {%v}", err)
	}
	defer resp.Body.Close()

	// check status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get upgrade package fail, errmsg {%s}", resp.Status)
	}

	// copy data
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("save file fail, write bytes: %d, errmsg {%v}", n, err)
	}

	return nil
}

// Decompress decompress, src: file to decompress; dest: path to save
func Decompress(src, dest string) error {
	// open file
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	// new reader
	gr, err := gzip.NewReader(sf)
	if err != nil {
		return err
	}
	defer gr.Close()

	// create files
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// log.Println(hdr.Name)

		filename := dest + hdr.Name

		// directory, skip
		if filename[len(filename)-1] == '/' {
			continue
		}

		// common file, create
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}

	return os.Create(name)
}

// CheckFileSHA256 check file's sha256
func CheckFileSHA256(filepath, sum string) bool {
	newsum, err := FileSHA256(filepath)
	if err != nil {
		return false
	}

	if newsum != sum {
		return false
	}

	return true
}
