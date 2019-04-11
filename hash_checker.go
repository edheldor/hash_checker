package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Hashmap struct {
	mx sync.Mutex
	m  map[string]string
}

func newHashmap() *Hashmap {
	return &Hashmap{
		m: make(map[string]string),
	}
}

func (h *Hashmap) Load(key string) (string, bool) {
	h.mx.Lock()
	defer h.mx.Unlock()
	val, ok := h.m[key]
	return val, ok
}

func (h *Hashmap) Store(key string, value string) {
	h.mx.Lock()
	defer h.mx.Unlock()
	h.m[key] = value
}

func main() {

	params := os.Args
	params_len := len(params)
	hashes := newHashmap()
	var wg sync.WaitGroup

	if params_len == 1 {
		fmt.Println("Поддерживается команда calc имя_файла")
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

			//полуаем путь директории в которой запускается файл
			path, _ := os.Getwd()
			filelist, err := FilePathWalkDir(path)
			if err != nil {
				fmt.Println("Невозможно создать список всех файлов")
				os.Exit(1)
			}

			for _, filepath := range filelist {
				wg.Add(1)
				go hashcalc(filepath, hashes, &wg)
			}
			// ждем завершения горутин
			wg.Wait()
			//Неожиданно чтение из мапа производится в случайном порядке, поэтому чтобы вывести в итоговый файл сортировку по алфавиту делаем следующее безобразие: итерируемся по слайсу filelist
			//там все по алфавиту

			for i := range filelist {

				hash, _ := hashes.Load(filelist[i])
				towrite := fmt.Sprintf("%s    %s \n", filelist[i], hash) //4 пробела между путем к файлу и хэшем
				file.WriteString(towrite)
			}

		} else {
			fmt.Println("Поддерживается команда calc имя_файла")
		}

	}

}

func hashcalc(filepath string, hashes *Hashmap, wg *sync.WaitGroup) {

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
	hashes.Store(filepath, hash)
	wg.Done()
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
