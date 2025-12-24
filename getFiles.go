package iguana

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/johnbalvin/iguana/files"
)

// GetFiles returns files at a given path
func (config Config) GetFiles(workingPath string) (map[string]files.HTML, map[string]files.Static, map[string]files.SW) {
	if config.FuncReplaceRelPath == nil {
		config.FuncReplaceRelPath = skippingDefault
	}
	return config.setFiles(true, workingPath)
}

// GetFilesObf returns files at a given path filling it with it's obfuscated content
func (config Config) GetFilesObf(workingPath string) (map[string]files.HTML, map[string]files.Static, map[string]files.SW) {
	if config.FuncReplaceRelPath == nil {
		config.FuncReplaceRelPath = skippingDefault
	}
	return config.setFiles(false, workingPath)
}

func (config Config) setFiles(shouldIObfuscate bool, workingPath string) (map[string]files.HTML, map[string]files.Static, map[string]files.SW) {
	htmlFiles := make(map[string]files.HTML)
	staticFiles := make(map[string]files.Static)
	serviceWorkers := make(map[string]string)
	config.addFiles(workingPath, htmlFiles, staticFiles)
	fmt.Println("")
	if !shouldIObfuscate {
		for _, static := range staticFiles {
			if !static.Obfuscate {
				continue
			}
			code, err := config.FuncObf(static)
			if err != nil {
				log.Fatalf("------File: %s\n", static.Path)
				continue
			}
			static.ContentObf.Me = []byte(code)
			static.ContentObf.Checksum = config.FuncCheckSum(static.ContentObf.Me)
			static.ContentObf.ID, static.ContentObf.URL = config.FuncIDURLObf(static)
			staticFiles[static.Path] = static
		}
	}
	config.changePathWithURLWrapper(shouldIObfuscate, workingPath, serviceWorkers, staticFiles)
	for k, html := range htmlFiles { //check service workers definitions inside html files
		paths := files.RegServiceW.FindAll(html.Content, -1)
		for _, dpendRelPath := range paths { // -> serviceWorker.register("path.js")
			jsPathR := string(files.RegJs.Find(dpendRelPath))                 //-> "path.js"
			rawDpendRelPath := strings.TrimSpace(jsPathR[1 : len(jsPathR)-1]) //-> path.js
			join := path.Join(workingPath, rawDpendRelPath)
			serviceWorkers[join] = html.Path
			html.ServiceWorkers[join] = true
		}
		//replace all urls paths in html content by their url
		html.Content = replaceURL(shouldIObfuscate, config.SkipLogging, html.Path, html.Content, config.SWPath, serviceWorkers, staticFiles)
		html.Checksum = config.FuncCheckSum(html.Content)
		htmlFiles[k] = html
	}
	serviceWorkersToReturn := make(map[string]files.SW)
	for swPath, fileCaller := range serviceWorkers {
		if static, ok := staticFiles[swPath]; ok { //check is static file is a service worker
			static.Content.ID = config.SWPath + "/" + static.Content.Checksum + ".js"
			static.Content.URL = "/" + static.Content.ID

			static.ContentObf.ID = config.SWPath + "/" + static.ContentObf.Checksum + ".js"
			static.ContentObf.URL = "/" + static.ContentObf.ID

			serviceWorkersToReturn[static.Path] = files.SW{FileCaller: fileCaller, Static: static}
			delete(staticFiles, swPath)
		} else {
			//there  is not service work
		}
		for k, html := range htmlFiles {
			if _, ok := html.DependsFullPath[swPath]; ok {
				htmlFiles[k].ServiceWorkers[swPath] = true
			}
		}
	}
	//compressing them all
	var err error
	for path, value := range htmlFiles {
		value.ContentBR, err = files.CompressBrotli(value.Content)
		if err != nil {
			log.Println("br compression err: ", err)
		}
		value.ContentZstd, err = files.CompressZstd(value.Content)
		if err != nil {
			log.Println("Zstd compression err: ", err)
		}
		htmlFiles[path] = value
	}
	for path, value := range staticFiles {
		value.Content.ContentBR, err = files.CompressBrotli(value.Content.Me)
		if err != nil {
			log.Println("br compression err: ", err)
		}
		value.Content.ContentZstd, err = files.CompressZstd(value.Content.Me)
		if err != nil {
			log.Println("Zstd compression err: ", err)
		}
		staticFiles[path] = value
	}
	for path, value := range serviceWorkersToReturn {
		value.Content.ContentBR, err = files.CompressBrotli(value.Content.Me)
		if err != nil {
			log.Println("br compression err: ", err)
		}
		value.Content.ContentZstd, err = files.CompressZstd(value.Content.Me)
		if err != nil {
			log.Println("Zstd compression err: ", err)
		}
		serviceWorkersToReturn[path] = value
	}
	return htmlFiles, staticFiles, serviceWorkersToReturn
}
