package iguana

import (
	"fmt"
	"log"
	"path"
	"strings"
)

//GetFiles returns files at a given path
func (config Config) GetFiles(workingPath string) (map[string]HTML, map[string]Static, map[string]SW) {
	if config.FuncReplaceRelPath == nil {
		config.FuncReplaceRelPath = skippingDefault
	}
	return config.getFiles(true, workingPath)
}

//GetFilesObf returns files at a given path filling it with it's obfuscated content
func (config Config) GetFilesObf(workingPath string) (map[string]HTML, map[string]Static, map[string]SW) {
	if config.FuncReplaceRelPath == nil {
		config.FuncReplaceRelPath = skippingDefault
	}
	return config.getFiles(false, workingPath)
}

func (config Config) getFiles(normal bool, workingPath string) (map[string]HTML, map[string]Static, map[string]SW) {
	htmlFiles := make(map[string]HTML)
	staticFiles := make(map[string]Static)
	serviceWorkers := make(map[string]string)
	config.AddFiles(normal, workingPath, htmlFiles, staticFiles)
	fmt.Println("")
	if !normal {
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
	config.ChangePathWithURL(normal, workingPath, serviceWorkers, staticFiles)
	for k, html := range htmlFiles { //check service workers definitions inside html files
		paths := regServiceW.FindAll(html.Content, -1)
		for _, dpendRelPath := range paths { // -> serviceWorker.register("path.js")
			jsPathR := string(regJs.Find(dpendRelPath))                       //-> "path.js"
			rawDpendRelPath := strings.TrimSpace(jsPathR[1 : len(jsPathR)-1]) //-> path.js
			join := path.Join(workingPath, rawDpendRelPath)
			serviceWorkers[join] = html.Path
			html.ServiceWorkers[join] = true
		}
		//replace all urls paths in html content by their url
		html.Content = replaceURL(normal, config.SkipLogging, html.Path, html.Content, config.SWPath, serviceWorkers, staticFiles)
		html.Checksum = config.FuncCheckSum(html.Content)
		htmlFiles[k] = html
	}
	serviceWorkersToReturn := make(map[string]SW)
	for swPath, fileCaller := range serviceWorkers {
		if static, ok := staticFiles[swPath]; ok { //check is static file is a service worker
			static.Content.ID = config.SWPath + "/" + static.Content.Checksum + ".js"
			static.Content.URL = "/" + static.Content.ID

			static.ContentObf.ID = config.SWPath + "/" + static.ContentObf.Checksum + ".js"
			static.ContentObf.URL = "/" + static.ContentObf.ID

			serviceWorkersToReturn[static.Path] = SW{FileCaller: fileCaller, Static: static}
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
	return htmlFiles, staticFiles, serviceWorkersToReturn
}
