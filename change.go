package iguana

import (
	"fmt"
	"path"
	"strings"
)

func (static Static) addSw(workingPath string, sWorkers map[string]string, staticFiles map[string]Static) {
	paths := regServiceW.FindAll(static.Content.Me, -1)
	for _, dpendRelPath := range paths { // -> serviceWorker.register("path.js")
		jsPathR := string(regJs.Find(dpendRelPath))                       //-> "path.js"
		rawDpendRelPath := strings.TrimSpace(jsPathR[1 : len(jsPathR)-1]) //-> path.js
		join := path.Join(workingPath, rawDpendRelPath)
		sWorkers[join] = static.Path
	}
}

//changePathWithURL it's the recursive function that replace all the paths of files's content
func (config Config) changePathWithURL(normal bool, workingPath string, static Static, sw map[string]string, staticFiles map[string]Static, depthDependency []string, deph int) ([]string, error) {
	if len(static.DependsFullPath) == 0 {
		return nil, nil
	}
	if deph == 30 {
		depthDependency = append(depthDependency, static.Path)
		return nil, errDepth
	}
	deph++
	static.addSw(workingPath, sw, staticFiles)
	for dependencPath := range static.DependsFullPath {
		dependency, ok := staticFiles[dependencPath]
		if !ok { //dependency missing, not recording directly because this it's not the principal file
			continue
		}
		if !dependency.ChangeContent { // content must no be changed, like image, fonts,  plain text
			continue
		}
		if depthDependency, err := config.changePathWithURL(normal, workingPath, dependency, sw, staticFiles, depthDependency, deph); err != nil {
			depthDependency = append(depthDependency, static.Path)
			return depthDependency, err
		}
		dependency = staticFiles[dependencPath] //it may be possible the file has been change in the mean time, so I get and update
		if dependency.Content.CheckSumChanged {
			continue
		}
		dependency.Content.Me = replaceURL(normal, config.SkipLogging, dependency.Path, dependency.Content.Me, config.SWPath, sw, staticFiles)
		dependency.Content.Checksum = config.FuncCheckSum(dependency.Content.Me)
		dependency.Content.ID, dependency.Content.URL = config.FuncIDURLNormal(dependency)
		dependency.Content.CheckSumChanged = true
		if !normal {
			dependency.ContentObf.Me = replaceURL(normal, config.SkipLogging, dependency.Path, dependency.Content.Me, config.SWPath, sw, staticFiles)
			dependency.ContentObf.Checksum = config.FuncCheckSum(dependency.Content.Me)
			dependency.ContentObf.ID, dependency.ContentObf.URL = config.FuncIDURLObf(dependency)
			dependency.ContentObf.CheckSumChanged = true
		}
		staticFiles[dependencPath] = dependency

	}
	static = staticFiles[static.Path]
	if static.Content.CheckSumChanged {
		//fmt.Printf("---CheckSumChanged file: %s\n", static.Path)
		return nil, nil
	}
	static.Content.Me = replaceURL(normal, config.SkipLogging, static.Path, static.Content.Me, config.SWPath, sw, staticFiles)
	static.Content.Checksum = config.FuncCheckSum(static.Content.Me)
	static.Content.ID, static.Content.URL = config.FuncIDURLNormal(static)
	static.Content.CheckSumChanged = true
	if !normal {
		static.ContentObf.Me = replaceURL(normal, config.SkipLogging, static.Path, static.Content.Me, config.SWPath, sw, staticFiles)
		static.ContentObf.Checksum = config.FuncCheckSum(static.Content.Me)
		static.ContentObf.ID, static.ContentObf.URL = config.FuncIDURLObf(static)
		static.ContentObf.CheckSumChanged = true
	}
	staticFiles[static.Path] = static
	return nil, nil
}

//ChangePathWithURL changes all the paths of file content with URLs
//so each file content that have relative path dependencies
//are gonna get change by an URL as defined with the function config.FunIDURL
func (config Config) ChangePathWithURL(normal bool, workingPath string, serviceWorkers map[string]string, staticFiles map[string]Static) {
	var err error
	for _, static := range staticFiles {
		if !static.ChangeContent {
			continue
		}
		var depthDependencys []string

		if depthDependencys, err = config.changePathWithURL(normal, workingPath, static, serviceWorkers, staticFiles, depthDependencys, 0); err != nil {

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
