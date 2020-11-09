package s3cli

import (
	"github.com/qiaogw/log"
	"github.com/qiaogw/pkg/config"

	"github.com/minio/minio-go/v6"
)

func ListObj() {
	m := config.Config.S3
	// Create a done channel to control 'ListObjects' go routine.
	client, err := minio.New(m.Endpoint, m.AccessKeyID, m.SecretAccessKey, m.Secure)
	if err != nil {
		log.Error(err)
		return
	}
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	isRecursive := true
	list := 0
	objectCh := client.ListObjects(m.Bucket, "", isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			log.Error(object.Err)
			return
		}
		list++
	}
}
