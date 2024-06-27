package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/philippgille/chromem-go"
	"github.com/siuyin/aigotut/emb"
)

var db *chromem.DB

func main() {
	collection := initDB()
	fmt.Println(collection.Count())

	exportDB()
}

func initDB() *chromem.Collection {
	docs := loadDocuments()

	db = chromem.NewDB()
	c, err := db.CreateCollection("aigogo", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	c.AddDocuments(ctx, docs, runtime.NumCPU())
	return c

}

func exportDB() {
	compressWithGZIP := false
	encryptKey32Bytes := ""
	if err := db.Export("./db.gob", compressWithGZIP, encryptKey32Bytes); err != nil {
		log.Fatal(err)
	}
}

func loadDocuments() []chromem.Document {
	f, err := os.Open("embeddings.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var rec emb.Rec
	dec := gob.NewDecoder(f)
	docs := []chromem.Document{}
	for {
		if err := dec.Decode(&rec); err != nil {
			break
		}
		docs = addDoc(docs, &rec)
	}
	return docs
}

func addDoc(docs []chromem.Document, rec *emb.Rec) []chromem.Document {
	d := chromem.Document{
		ID:      rec.ID,
		Content: rec.Title + rec.Content,
	}
	d.Embedding = append(d.Embedding, rec.Embedding...)
	docs = append(docs, d)
	return docs
}
