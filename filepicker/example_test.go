package filepicker_test

import (
	"log"
	"time"

	"github.com/filepicker/filepicker-go/filepicker"
)

const ApiKey = "put_your_api_key_here"

// This example shows how to send URLs content directly to your storage bucket.
func ExampleStoreURL() {
	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	// Store URL content using default options.
	const dataURL = "https://d3urzlae3olibs.cloudfront.net/watermark.png"
	blob, err := cl.StoreURL(dataURL, nil)
	if err != nil {
		log.Fatalf("cannot store file %q: %v\n", dataURL, err)
	}
	log.Printf("file %q stored: %q\n", dataURL, blob.URL)
}

func ExampleDownloadToFile() {
	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	blob := filepicker.NewBlob("hFHUCB3iTxyMzseuWOgG")
	if err := cl.DownloadToFile(blob, nil, "."); err != nil {
		log.Fatalf("cannot download file: %v\n", err)
	}
	log.Println("file downloaded!")
}

func ExampleStat() {
	// Create a new Filepicker.io client with S3 storage set by default.
	cl := filepicker.NewClient(ApiKey)

	options := &filepicker.StatOpts{
		Tags: []filepicker.MetaTag{filepicker.TagMd5Hash},
	}

	meta, err := cl.Stat(filepicker.NewBlob("hFHUCB3iTxyMzseuWOgG"), options)
	if err != nil {
		log.Println("cannot stat file:", err)
	}

	md5hash := "unknown"
	if md5, ok := meta.Md5Hash(); ok {
		md5hash = md5
	}
	log.Println("MD5Hash:", md5hash)
}

func ExampleMakeSecurity() {
	options := &filepicker.PolicyOpts{
		Expiry: time.Unix(1508141504, 0),
		Handle: "KW9EJhYtS6y48Whm2S6D",
	}

	policy, err := filepicker.MakePolicy(options)
	if err != nil {
		log.Fatalln("cannot create policy:", err)
	}

	security := filepicker.MakeSecurity("Z3IYZSH2UJA7VN3QYFVSVCF7PI", policy)
	log.Printf("P: %s\nS: %s\n", security.Policy, security.Signature)
}
