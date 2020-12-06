package files

import (
	"path"
	"strings"
)

func (static Static) SetSW(workingPath string, sWorkers map[string]string, staticFiles map[string]Static) {
	paths := RegServiceW.FindAll(static.Content.Me, -1)
	for _, dpendRelPath := range paths { // -> serviceWorker.register("path.js")
		jsPathR := string(RegJs.Find(dpendRelPath))                       //-> "path.js"
		rawDpendRelPath := strings.TrimSpace(jsPathR[1 : len(jsPathR)-1]) //-> path.js
		join := path.Join(workingPath, rawDpendRelPath)
		sWorkers[join] = static.Path
	}
}
