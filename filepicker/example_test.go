package filepicker_test

import (
	"log"

	"github.com/filepicker/filepicker-go/filepicker"
)

// This example shows how to send URLs content directly to your storage bucket.
func ExampleStoreURL() {
	const (
		apiKey  = "put_your_api_key_here"
		dataUrl = "https://d3urzlae3olibs.cloudfront.net/watermark.png"
	)

	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(apiKey)

	// Store URL content using default options.
	blob, err := cl.StoreURL(dataUrl, nil)
	if err != nil {
		log.Fatalf("cannot store file %q: %v\n", dataUrl, err)
	}
	log.Printf("file %q stored: %q\n", dataUrl, blob.Url)
}
