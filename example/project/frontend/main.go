package main

import (
	"context"
	"fmt"
	"iguana"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var normal = true
	obfuscate := os.Getenv("ObfuscateJS")
	if obfuscate == "yes" {
		normal = false
	}
	log.Println("obfuscate: ", obfuscate)
	ctx := context.Background()
	startsAt := time.Now().Unix()
	htmlFiles, serviceWorkers, err := start(normal)
	if err != nil {
		log.Fatal("main -> main:3  -> err:", err)
	}
	log.Println("----end------")
	htmlsToSave := make(map[string]html)
	for key, file := range htmlFiles {
		fileToSave := html{Checksum: file.Checksum, DataGenerate: file.DataGenerate}
		htmlsToSave[key] = fileToSave
	}
	serviceWorkersToSave := make(map[string]string)
	for _, file := range serviceWorkers {
		if normal {
			serviceWorkersToSave[file.Content.Checksum] = file.Content.ID
		} else {
			serviceWorkersToSave[file.ContentObf.Checksum] = file.ContentObf.ID
		}
	}
	batch := clienteFS.Batch()
	batch.Set(docDirectory, htmlFiles).Set(docServiceWorkers, serviceWorkersToSave)
	if _, err = batch.Commit(ctx); err != nil {
		log.Println("main -> main:4  -> err:", err)
	}
	secondsPass := time.Now().Unix() - startsAt
	fmt.Println("Ends: ", secondsPass, " seconds", " minutes: ", secondsPass/60)
}

func worker(i int, canalFiles <-chan interface{}, wg *sync.WaitGroup, normal bool) {
	for fileInterface := range canalFiles {
		switch fileInterface.(type) {
		case iguana.HTML:
			file := fileInterface.(iguana.HTML)
			saveHTML(file)
		case iguana.Static:
			file := fileInterface.(iguana.Static)
			saveStatic(file, normal)
		case iguana.SW:
			file := fileInterface.(iguana.SW)
			saveServiceWorker(file, normal)
		}
		wg.Done()
	}
}

func start(normal bool) (map[string]iguana.HTML, map[string]iguana.SW, error) {
	var htmlFiles map[string]iguana.HTML
	var staticFiles map[string]iguana.Static
	var serviceWorkers map[string]iguana.SW
	config := iguana.GetDefaultConfig()
	config.FuncIDURLNormal = idURLNormal
	config.FuncIDURLObf = idURLObf
	if normal {
		htmlFiles, staticFiles, serviceWorkers = config.GetFiles(constStartInPath)
	} else {
		htmlFiles, staticFiles, serviceWorkers = config.GetFilesObf(constStartInPath)
	}
	var wg sync.WaitGroup
	wg.Add(len(staticFiles) + len(htmlFiles) + len(serviceWorkers))
	canalFiles := make(chan interface{})
	log.Println("HTMLFiles: ", len(htmlFiles), " Static files: ", len(staticFiles))
	for i := 0; i < 30; i++ {
		go worker(i, canalFiles, &wg, normal)
	}
	for _, file := range staticFiles {
		canalFiles <- file
	}
	for _, file := range serviceWorkers {
		canalFiles <- file
	}
	for _, file := range htmlFiles {
		canalFiles <- file
	}
	fmt.Printf("waiting.....")
	wg.Wait()
	return htmlFiles, serviceWorkers, nil
}
