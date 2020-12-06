package iguana

import (
	"fmt"
	"strings"

	"github.com/johnbalvin/iguana/files"
	"github.com/johnbalvin/iguana/utils"
)

//changePathWithURLWrapper changes all the paths of file content with URLs
//so each file content that have relative path dependencies
//are gonna get change by an URL as defined with the function config.FunIDURL
func (config Config) changePathWithURLWrapper(shouldIObfuscate bool, workingPath string, serviceWorkers map[string]string, staticFiles map[string]files.Static) {
	var err error
	for _, static := range staticFiles {
		if !static.ChangeContent {
			continue
		}
		var depthDependencys []string

		if depthDependencys, err = config.changePathWithURL(shouldIObfuscate, workingPath, static, serviceWorkers, staticFiles, depthDependencys, 0); err != nil {

			if config.SkipLogging {
				continue
			}
			fmt.Printf("Err: %s\n", err)
			j := 0
			for i := len(depthDependencys) - 1; i >= 0; i-- {
				fmt.Printf(" %d) %s\n", j, depthDependencys[i])
				j = j + 1
			}
			fmt.Printf("\n\n")
		}
	}
}

//changePathWithURL it's the recursive function that replace all the paths of files's content
func (config Config) changePathWithURL(shouldIObfuscate bool, workingPath string, static files.Static, sw map[string]string, staticFiles map[string]files.Static, depthDependency []string, deph int) ([]string, error) {

	if len(static.DependsFullPath) == 0 {
		return nil, nil
	}
	if deph == 30 {
		depthDependency = append(depthDependency, static.Path)
		return nil, utils.ErrDepth
	}
	deph++
	static.SetSW(workingPath, sw, staticFiles)
	for dependencPath := range static.DependsFullPath {
		dependency, ok := staticFiles[dependencPath]
		if !ok { //dependency missing, not recording directly because this it's not the principal file
			continue
		}
		if !dependency.ChangeContent { // content must no be changed, like image, fonts,  plain text
			continue
		}
		if depthDependency, err := config.changePathWithURL(shouldIObfuscate, workingPath, dependency, sw, staticFiles, depthDependency, deph); err != nil {
			depthDependency = append(depthDependency, static.Path)
			return depthDependency, err
		}
		dependency = staticFiles[dependencPath] //it may be possible the file has been change in the mean time, so I get and update
		if dependency.Content.CheckSumChanged {
			continue
		}
		dependency.Content.Me = replaceURL(shouldIObfuscate, config.SkipLogging, dependency.Path, dependency.Content.Me, config.SWPath, sw, staticFiles)
		dependency.Content.Checksum = config.FuncCheckSum(dependency.Content.Me)
		dependency.Content.ID, dependency.Content.URL = config.FuncIDURLNormal(dependency)
		dependency.Content.CheckSumChanged = true
		if !shouldIObfuscate {
			dependency.ContentObf.Me = replaceURL(shouldIObfuscate, config.SkipLogging, dependency.Path, dependency.Content.Me, config.SWPath, sw, staticFiles)
			dependency.ContentObf.Checksum = config.FuncCheckSum(dependency.Content.Me)
			dependency.ContentObf.ID, dependency.ContentObf.URL = config.FuncIDURLObf(dependency)
			dependency.ContentObf.CheckSumChanged = true
		}
		staticFiles[dependencPath] = dependency

	}
	static = staticFiles[static.Path] //updating the file content, remember this is a recursion function, so the content may have change before

	if static.Content.CheckSumChanged {
		//fmt.Printf("---CheckSumChanged file: %s\n", static.Path)
		return nil, nil
	}
	static.Content.Me = replaceURL(shouldIObfuscate, config.SkipLogging, static.Path, static.Content.Me, config.SWPath, sw, staticFiles)
	static.Content.Checksum = config.FuncCheckSum(static.Content.Me)
	static.Content.ID, static.Content.URL = config.FuncIDURLNormal(static)
	static.Content.CheckSumChanged = true
	if !shouldIObfuscate {
		static.ContentObf.Me = replaceURL(shouldIObfuscate, config.SkipLogging, static.Path, static.Content.Me, config.SWPath, sw, staticFiles)
		static.ContentObf.Checksum = config.FuncCheckSum(static.Content.Me)
		static.ContentObf.ID, static.ContentObf.URL = config.FuncIDURLObf(static)
		static.ContentObf.CheckSumChanged = true
	}
	staticFiles[static.Path] = static
	return nil, nil
}

var current int

func replaceURL(shouldIObfuscate, skipLogging bool, parentPath string, content []byte, swPath string, sw map[string]string, staticFiles map[string]files.Static) []byte {
	var inexistenDepencies []string
	tempContent := string(content)
	var dpendFullPaths = make(map[string]bool)
	paths := files.RegStatic.FindAll(content, -1)
	for _, path := range paths { //remove repeated because it will overlap the url if a file it's been use multiple times
		dpendFullPaths[string(path)] = true // ->  "path.extension"
	}
	for dpendFullPath := range dpendFullPaths { // ->  "path.extension"
		path := dpendFullPath[1 : len(dpendFullPath)-1] // ->  path.extension
		file, ok := staticFiles[path]
		if !ok {
			inexistenDepencies = append(inexistenDepencies, path)
			continue
		}
		var checkSum, url string
		if shouldIObfuscate {
			checkSum = file.Content.Checksum
			url = file.Content.URL
		} else {
			checkSum = file.ContentObf.Checksum
			url = file.ContentObf.URL
		}
		delimiter := string(dpendFullPath[0])
		var replaceWith string
		if _, ok := sw[path]; ok { //check if this is a service worker, services workers can't be in a CDN, must be serve only by your server https://github.com/w3c/ServiceWorker/issues/940
			replaceWith = delimiter + "/" + swPath + "/" + checkSum + ".js" + delimiter
		} else {
			replaceWith = delimiter + url + delimiter
		}
		tempContent = strings.ReplaceAll(tempContent, dpendFullPath, replaceWith)
	}
	if len(inexistenDepencies) != 0 {
		if skipLogging {
			return []byte(tempContent)
		}
		fmt.Printf("Dependency missing at path: %s\n", parentPath)
		for _, depencies := range inexistenDepencies {
			fmt.Printf("   %s\n", depencies)
		}
		fmt.Printf("\n")
	}
	return []byte(tempContent)
}
