package main

import (
	"fmt"
	"time"

	"github.com/minio/minio-go/v6"
)

func main() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
	} else {
		if false {
			preSignUpload(client, "hello", "path/to/myObject.png")
		} else {
			// file info
			obj, err := client.GetObject("hello", "path/to/myobject.png", minio.GetObjectOptions{})
			if nil != err {
				panic("err")
			}

			stat, _ := obj.Stat()
			fmt.Println("OBJECT: ", stat.Key, stat.ContentType, stat.Metadata, stat.Size, stat.UserTags)
		}
	}
}

func preSignUpload(client *minio.Client, bucket string, key string) {
	policy := minio.NewPostPolicy()
	policy.SetBucket(bucket)
	policy.SetKey(key)
	policy.SetExpires(time.Now().UTC().Add(4 * time.Hour)) // expires in 4 hours
	policy.SetContentType("image/png")
	policy.SetContentLengthRange(1, 1024*1024) // 1KB to 1MB
	policy.SetUserMetadata("app", "playround")

	// Get the POST form key/value object:
	url, formData, err := client.PresignedPostPolicy(policy)
	if err != nil {
		fmt.Println(err)
		return
	}

	// POST your content from the command line using `curl`
	fmt.Printf("curl -X POST ")
	for k, v := range formData {
		fmt.Printf("-F %s=%s ", k, v)
	}
	fmt.Printf("-F file=@/etc/bash.bashrc ")
	fmt.Printf("%s\n", url)
}
