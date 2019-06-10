package get

import (
	"context"
	"io/ioutil"
	"log"
	"sync"
	"text/template"
)

//Get returns all html and static files and service workers from database
func Get() (*template.Template, map[string][]byte, map[string]SWInfo) {
	var htmlConData = template.New("default")
	var htmlSinData = make(map[string][]byte)
	var htmlAux = make(map[string]HTMLInfo)
	ctx := context.Background()
	doc, err := clienteFS.Collection("frontEnd").Doc("Directory").Get(ctx)
	if err != nil {
		log.Fatal("serviruser -> Get:1 -> err:", err)
	}
	if err = doc.DataTo(&htmlAux); err != nil {
		log.Fatal("serviruser -> Get:2 -> err:", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(htmlAux))
	for path := range htmlAux {
		go func(path string) {
			htmlInfo := htmlAux[path]
			if err := htmlInfo.getHTMLSW(); err != nil {
				log.Fatal("serviruser -> Get:3 -> ruta:", htmlInfo.Path, " err:", err)
			}
			if htmlInfo.DataGenerate {
				template.Must(htmlConData.New(htmlInfo.Path).Parse(string(htmlInfo.Content)))
			} else {
				htmlSinData[htmlInfo.Path] = htmlInfo.Content
			}
			wg.Done()
		}(path)
	}

	var serviceWorkers = make(map[string]SWInfo)
	var serviceWorkersAux = make(map[string]string)
	doc, err = clienteFS.Collection("frontEnd").Doc("SW").Get(ctx)
	if err != nil {
		log.Fatal("serviruser -> Get:4 -> err:", err)
	}
	if err = doc.DataTo(&serviceWorkersAux); err != nil {
		log.Fatal("serviruser -> Get:5 -> err:", err)
	}
	wg.Add(len(serviceWorkersAux))
	for checksum := range serviceWorkersAux {
		go func(checksum string) {
			id := serviceWorkersAux[checksum]
			swinfo := SWInfo{ID: id, Checksum: checksum}
			if err := getSW(id, &swinfo); err != nil {
				log.Fatal("serviruser -> Get:6 -> err:", err)
			}
			serviceWorkers[checksum+".js"] = swinfo
			wg.Done()
		}(checksum)
	}
	//-----------------------------------------
	wg.Wait()
	return htmlConData, htmlSinData, serviceWorkers
}

func getSW(id string, static *SWInfo) error {
	ctx := context.Background()
	reader, err := clienteCS.Object(id).NewReader(ctx)
	if err != nil {
		log.Fatal("serviruser ->  getSW:1 -> err:", err)
		return nil
	}
	static.Me, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("serviruser -> getSW:2 -> err:", err)
		return nil
	}
	if err = reader.Close(); err != nil {
		log.Fatal("serviruser -> getSW:3 -> err:", err)
		return nil
	}
	return nil
}
func (htmlInfo *HTMLInfo) getHTMLSW() error {
	ctx := context.Background()
	reader, err := clienteCS.Object("p/" + htmlInfo.Checksum).NewReader(ctx)
	if err != nil {
		log.Fatal("serviruser -> (htmlInfo *HTMLInfo) getHTMLSW:1 -> err:", err)
		return err
	}
	htmlInfo.Content, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("serviruser -> (htmlInfo *HTMLInfo) getHTMLSW:2 -> err:", err)
		return err
	}
	if err = reader.Close(); err != nil {
		log.Fatal("serviruser -> (htmlInfo *HTMLInfo) getHTMLSW:3 -> err:", err)
		return err
	}
	return nil
}
