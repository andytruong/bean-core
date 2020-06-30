package main

import (
	"log"

	"github.com/minio/minio-go/v6"
)

func main() {
	endpoint := "127.0.0.1"
	accessKeyID := "minio"
	secretAccessKey := "minio"
	useSSL := false

	// Initialize minio client object.
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", client) // client is now setup
}
