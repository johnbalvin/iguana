package get

import (
	"html/template"
	"iguana/example/project/clients"
	"time"
)

//ServirHTMLS has all html files that have generate data in it
var ServirHTMLS *template.Template

//ServirSW has all service workers
var ServirSW = make(map[string][]byte)

//SevirHTMLSinData has all html files that have generate data in it
var SevirHTMLSinData = make(map[string][]byte)

var clienteFS = clients.ClientFirestore()
var clienteCS = clients.ClientStorage()

//HTMLInfo contains HTML's info
type HTMLInfo struct {
	Path            string
	Content         []byte
	Checksum        string
	DependsFullPath map[string]bool
	DataGenerate    bool
	LastModify      time.Time
}

//SWInfo contains service worker's info
type SWInfo struct {
	ID       string
	Checksum string
	Me       []byte
}
