package filepicker_test

import (
	"log"

	"github.com/filepicker/filepicker-go/filepicker"
)

const ApiKey = "put_your_api_key_here"

// This example shows how to send URLs content directly to your storage bucket.
func ExampleStoreURL() {
	const dataUrl = "https://d3urzlae3olibs.cloudfront.net/watermark.png"

	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	// Store URL content using default options.
	blob, err := cl.StoreURL(dataUrl, nil)
	if err != nil {
		log.Fatalf("cannot store file %q: %v\n", dataUrl, err)
	}
	log.Printf("file %q stored: %q\n", dataUrl, blob.Url)
}

func ExampleDownloadToFile() {
	const fileHandle = "hFHUCB3iTxyMzseuWOgG"

	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	blob := filepicker.NewBlob(fileHandle)
	if err := cl.DownloadToFile(blob, nil, "."); err != nil {
		log.Fatalf("cannot download file: %v\n", err)
	}
	log.Println("file downloaded!")
}

func ExampleStat() {
	const fileHandle = "hFHUCB3iTxyMzseuWOgG"

	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	options := &filepicker.MetaOpts{
		Tags: []filepicker.MetaTag{filepicker.TagMd5Hash},
	}

	meta, err := cl.Stat(filepicker.NewHandle(fileHandle), options)
	if err != nil {
		log.Println("cannot stat file:", err)
	}

	md5hash := "unknown"
	if md5, ok := meta.Md5Hash(); ok {
		md5hash = md5
	}
	log.Println("MD5Hash:", md5hash)
}
