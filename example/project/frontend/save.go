package main

import (
	"context"
	"fmt"
	"iguana"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
)

func saveServiceWorker(sw iguana.SW, normal bool) {
	ctx := context.Background()
	fmt.Print("¡")
	var id string
	var content []byte
	if normal {
		id = sw.Content.ID
		content = sw.Content.Me
	} else {
		id = sw.ContentObf.ID
		content = sw.ContentObf.Me
	}

	w := clienteCS.Object(id).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	w.ContentType = "text/plain"
	if _, err := w.Write(content); err != nil {
		log.Println("main -> saveServiceWorker:1 -> err:", err)
		return
	}
	if err := w.Close(); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if !(e.Code == 412 || e.Code == 400) {
				log.Println("main -> saveServiceWorker:2 -> err:", err)
				return
			}
		}
		return
	}
	fmt.Print("!")
}
func saveHTML(htmlFile iguana.HTML) {
	ctx := context.Background()
	fmt.Print("(")
	w := clienteCS.Object("p/" + htmlFile.Checksum).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	w.ContentType = "text/plain"
	if _, err := w.Write(htmlFile.Content); err != nil {
		log.Println("main -> saveHTML:1 -> err:", err)
		return
	}
	if err := w.Close(); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if !(e.Code == 412 || e.Code == 400) {
				log.Println("main -> saveHTML:2 -> err:", err)
				return
			}
		}
		fmt.Print(")")
		return
	}
	fmt.Print("-)")
}

func saveStatic(static iguana.Static, normal bool) error {
	ctx := context.Background()
	var id string
	var content []byte
	if normal {
		id = static.Content.ID
		content = static.Content.Me
	} else {
		id = static.ContentObf.ID
		content = static.ContentObf.Me
	}
	fmt.Print("<")
	w := clienteCS.Object(id).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx) //already check if existi
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = static.MimeType
	w.CacheControl = "must-revalidate,max-age=31536000" //agresive about caching because I'm using checkum as its name
	if _, err := w.Write(content); err != nil {
		log.Println("saveStatic -> (static Static) save:1 -> err:", err)
		return err
	}
	if err := w.Close(); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if !(e.Code == 412 || e.Code == 400) {
				log.Println("saveStatic -> (static Static) save:2 -> err:", err)
				return nil
			}
		}
		fmt.Print(">")
		return nil
	}
	fmt.Print("->")
	return nil
}
