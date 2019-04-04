package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type hashmap map[string]string

func main() {

	params := os.Args
	params_len := len(params)

	if params_len == 1 {
		fmt.Println("Поддерживаются команды calc имя_файла и check имя_файла")
	} else if params_len == 2 {
		fmt.Println("Нужно больше аргументов")
	} else if params_len == 3 {
		if params[1] == "calc" {
			file, err := os.Create(params[2])
			if err != nil {
				fmt.Println("Невозможно создать файл")
				os.Exit(1)
			}
			defer file.Close()
			hashes := make(hashmap)
			//полуаем путь директории в которой запускается файл
			path, _ := os.Getwd()
			filelist, err := FilePathWalkDir(path)
			if err != nil {
				fmt.Println("Невозможно создать список всех файлов")
				os.Exit(1)
			}

			for _, filepath := range filelist {
				hashcalc(filepath, hashes)
			}

			for key, _ := range hashes {
				towrite := key + " " + hashes[key] + "\n"
				file.WriteString(towrite)
			}

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
			fmt.Println("Поддерживаются команды calc имя_файла и check имя_файла")
		}

	}

}

func hashcalc(filepath string, hashes hashmap) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	hashes[filepath] = hash

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
