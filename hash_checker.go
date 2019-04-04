package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	//"path/filepath"
)

func main() {
	params := os.Args
	params_len := len(params)

	fmt.Println(params_len)

	if params_len == 1 {
		fmt.Println("Поддерживаются команды calc и check")
	} else if params_len == 2 {
		fmt.Println("Невозможно создать файл")
	} else if params_len == 3 {
		if params[1] == "calc" {
			file, err := os.Create(params[2])
			if err != nil {
				fmt.Println("Невозможно создать файл")
				os.Exit(1)
			}
			defer file.Close()
			path, _ := os.Getwd()
			filelist, err := FilePathWalkDir(path)
			if err != nil {
				fmt.Println("Невозможно создать список всех файлов")
				os.Exit(1)
			}
			print(filelist)

		} else if params[1] == "check" {
			file, err := os.Open(params[2])
			if err != nil {
				fmt.Println("Невозможно открыть файл")
			}
			defer file.Close()

			data := make([]byte, 64)

			for {
				n, err := file.Read(data)
				if err == io.EOF {
					break
				}
				fmt.Println(string(data[:n]))
			}

		} else {
			fmt.Println("Поддерживаются команды calc и check")
		}

	}

}

type file struct {
	name string
	path string
	hash string
}

func (item *file) calclHash() {
	hash := sha256.New()
	hash.Write([]byte("HELLO"))
	item.hash = fmt.Sprintf("%x", hash.Sum(nil))

}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
