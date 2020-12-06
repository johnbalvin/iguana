package iguana

import (
	"fmt"
	"iguana/files"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

//addFiles add files to htmlFiles,staticFiles , also changes the inmediate relative path to a relative path base on the full path
func (config Config) addFiles(workingPath string, htmlFiles map[string]files.HTML, staticFiles map[string]files.Static) {
	if !config.SkipLogging {
		fmt.Print(".")
	}
	subdirectories, err := ioutil.ReadDir(workingPath)
	if err != nil {
		log.Fatal("iguana -> addFiles:1 -> err:", err)
	}
	var folders []string
	for _, subdirectory := range subdirectories {
		name := subdirectory.Name()
		if subdirectory.IsDir() { //check if it's a file or directory
			folders = append(folders, workingPath+"/"+name)
			continue
		}
		ext := filepath.Ext(name)
		if _, ok := files.MimeTypeAllow[ext]; !ok {
			continue
		}
		fullPath := workingPath + "/" + name
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Fatal("iguana -> addFiles:2 -> err:", err)
		}
		if ext == ".css" {
			content = files.CorrectFormatCss(content)
		}
		var dpendRelPaths = make(map[string]bool)
		paths := files.RegStatic.FindAll(content, -1)
		for _, path := range paths { //remove repeated depencies like for example multiple images... if not, it could cause an overlay when replacing it's path
			dpendRelPaths[string(path)] = true
		}
		var dependsFullPath = make(map[string]bool) //all dependecies the file have map[fullPath]true
		tempContent := string(content)
		for dpendRelPath := range dpendRelPaths {
			delimiter := string(dpendRelPath[0])
			if dpendRelPath[0] != dpendRelPath[len(dpendRelPath)-1] { //this avoids false positives like string: I'm a false postive string.js" <- it maches because it has ' and "
				continue
			}
			rawDpendRelPath := strings.TrimSpace(dpendRelPath[1 : len(dpendRelPath)-1]) //removing delimiters and spaces to use path.Join
			join := path.Join(workingPath, rawDpendRelPath)                             //this gives me full path
			if strings.HasPrefix(workingPath, "./") && !strings.HasPrefix(join, "./") { //path.join removes "./" so I will just add it
				join = "./" + join
			}
			if config.FuncReplaceRelPath(fullPath, rawDpendRelPath) {
				continue
			}
			dependecyFullPath := delimiter + join + delimiter
			tempContent = strings.ReplaceAll(tempContent, dpendRelPath, dependecyFullPath)
			dependsFullPath[join] = true
		}
		if ext == ".html" {
			var html files.HTML
			html.DependsFullPath = dependsFullPath
			html.Path = fullPath
			if strings.Contains(tempContent, "{{") { // just to check if your html is completly static or you it generates data with templates
				html.DataGenerate = true
			}
			html.Content = []byte(tempContent)
			html.Checksum = config.FuncCheckSum(html.Content)
			html.ServiceWorkers = make(map[string]bool)
			writeToMapHTML(htmlFiles, html.Path, html)
			continue
		}

		var static files.Static
		static.Name = name
		static.Extension = ext
		static.Path = fullPath
		static.DependsFullPath = dependsFullPath
		static.Content.Me = []byte(tempContent)
		static.Content.Checksum = config.FuncCheckSum(static.Content.Me)
		static.Content.ID, static.Content.URL = config.FuncIDURLNormal(static)
		static.MimeType = files.MimeTypeAllow[ext]

		if strings.Contains(static.Name, ".min.") { //minified files most not be obfuscated, I got headache debuging this
			static.ContentObf.Me = static.Content.Me
			static.ContentObf.Checksum = static.Content.Checksum
			static.ContentObf.URL = static.Content.URL
		} else {
			static.Obfuscate = true
		}
		if _, ok := files.MimeTypeReplace[ext]; ok {
			static.ChangeContent = true //verify is file can have dependencies
		}
		if strings.HasPrefix(static.Path, "./") {
			writeToMapStatic(staticFiles, static.Path, static)
		}
	}
	var wg sync.WaitGroup
	wg.Add(len(folders))
	for _, folder := range folders {
		go func(folder string) {
			defer wg.Done()
			config.addFiles(folder, htmlFiles, staticFiles)
		}(folder)
	}
	wg.Wait()
}
func writeToMapHTML(htmlFiles map[string]files.HTML, key string, content files.HTML) {
	lockHTML.Lock()
	defer lockHTML.Unlock()
	htmlFiles[key] = content
}
func writeToMapStatic(staticFiles map[string]files.Static, key string, content files.Static) {
	lockStatic.Lock()
	defer lockStatic.Unlock()
	staticFiles[key] = content
}
