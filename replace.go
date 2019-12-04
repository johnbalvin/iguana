package iguana

import (
	"fmt"
	"strings"
)

func replaceURL(normal, skipLogging bool, parentPath string, content []byte, swPath string, sw map[string]string, staticFiles map[string]Static) []byte {
	var inexistenDepencies []string
	tempContent := string(content)
	var dpendFullPaths = make(map[string]bool)
	paths := regStatic.FindAll(content, -1)
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
		if normal {
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
