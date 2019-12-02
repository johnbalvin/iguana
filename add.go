package iguana

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

//AddFiles add files to htmlFiles,staticFiles , also changes the inmediate relative path to a relative path base on the full path
func (config Config) AddFiles(normal bool, workingPath string, htmlFiles map[string]HTML, staticFiles map[string]Static) {
	fmt.Print(".")
	subdirectories, err := ioutil.ReadDir(workingPath)
	if err != nil {
		log.Fatal("iguana -> AddFiles:1 -> err:", err)
	}
	var folders []string
	for _, subdirectory := range subdirectories {
		name := subdirectory.Name()
		if subdirectory.IsDir() { //check if it's a file or directory
			folders = append(folders, workingPath+"/"+name)
			continue
		}
		separator := strings.Split(name, ".")
		ext := separator[len(separator)-1]
		if _, ok := mimeTypeAllow[ext]; !ok {
			continue
		}
		fullPath := workingPath + "/" + name
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Fatal("iguana -> AddFiles:2 -> err:", err)
		}
		var dpendRelPaths = make(map[string]bool)
		paths := regStatic.FindAll(content, -1)
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
			dependecyFullPath := delimiter + join + delimiter
			tempContent = strings.ReplaceAll(tempContent, dpendRelPath, dependecyFullPath)
			dependsFullPath[join] = true
		}
		if ext == "html" {
			var html HTML
			html.DependsFullPath = dependsFullPath
			html.Path = fullPath
			if strings.Contains(tempContent, "{{") { // just to check if your html is completly static or you it generates data with templates
				html.DataGenerate = true
			}
			html.Content = []byte(tempContent)
			html.Checksum = config.FuncCheckSum(html.Content)
			key := strings.TrimLeft(html.Path, "./")
			html.ServiceWorkers = make(map[string]bool)
			htmlFiles[key] = html
			continue
		}
		var static Static
		static.Name = name
		static.Extension = ext
		static.Path = fullPath
		static.DependsFullPath = dependsFullPath
		static.Content.Me = []byte(tempContent)
		static.Content.Checksum = config.FuncCheckSum(static.Content.Me)
		static.Content.ID, static.Content.URL = config.FuncIDURLNormal(static)
		static.MimeType = mimeTypeAllow[ext]

		if strings.Contains(static.Name, ".min.") { //minified files most not be obfuscated, I got headache debuging this
			static.ContentObf.Me = static.Content.Me
			static.ContentObf.Checksum = static.Content.Checksum
			static.ContentObf.URL = static.Content.URL
		} else {
			static.Obfuscate = true
		}
		if _, ok := mimeTypeReplace[ext]; ok {
			static.ChangeContent = true //verify is file can have dependencies
		}
		staticFiles[static.Path] = static
		if strings.HasPrefix(static.Path, "./") {
			staticFiles[static.Path[2:]] = static
		}
	}
	for _, folder := range folders {
		config.AddFiles(normal, folder, htmlFiles, staticFiles)
	}
}
