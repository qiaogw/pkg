package s3cli

import (
	"github.com/minio/minio-go/v6"
	"github.com/qiaogw/pkg/logs"
)

func ListObj(endpoint, accessKeyID, secretAccessKey,bucket string, secure bool) error {
	//m := config.Config.S3
	// Create a done channel to control 'ListObjects' go routine.
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		logs.Error(err)
		return err
	}
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	isRecursive := true
	list := 0
	objectCh := client.ListObjects(bucket, "", isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			logs.Error(object.Err)
			return object.Err
		}
		list++
	}
	return nil
}
