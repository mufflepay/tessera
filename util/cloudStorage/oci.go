package cloudStorage

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"io"
	"mime/multipart"
)

type OCIStorage struct {
	Config     common.ConfigurationProvider
	Namespace  string
	BucketName string
}

func NewOCIStorage(config OCIStorage) *OCIStorage {
	return &OCIStorage{
		Config:     config.Config,
		Namespace:  config.Namespace,
		BucketName: config.BucketName,
	}
}

func (o *OCIStorage) OCIInit() (objectstorage.ObjectStorageClient, error) {
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(o.Config)
	if err != nil {
		return objectstorage.ObjectStorageClient{}, err
	}
	return client, nil
}

func (o *OCIStorage) Get(ctx context.Context, client objectstorage.ObjectStorageClient) (string, error) {
	request := objectstorage.GetObjectRequest{
		NamespaceName: &o.Namespace,
		BucketName:    &o.BucketName,
		//ObjectName:    &objectName, - objectName parameter
	}
	response, err := client.GetObject(ctx, request)
	if err != nil {
		return "", err
	}
	defer func(Content io.ReadCloser) {
		err := Content.Close()
		if err != nil {

		}
	}(response.Content)

	fmt.Println("get object, status code ", response.RawResponse.StatusCode)
	fmt.Println("content length ", *response.ContentLength)
	fmt.Println("content type ", *response.ContentType)

	//buf := new(strings.Builder)
	//_, err = io.Copy(buf, response.Content)
	//if err != nil {
	//	return "", err
	//}
	return "", nil
}

func (o *OCIStorage) Upload(ctx context.Context, client objectstorage.ObjectStorageClient, fileName string, fileType string, fileSize int64, fileData multipart.File) (*string, error) {

	request := objectstorage.PutObjectRequest{
		NamespaceName: &o.Namespace,
		BucketName:    &o.BucketName,
		ObjectName:    &fileName,
		ContentLength: common.Int64(fileSize),
		ContentType:   &fileType,
		PutObjectBody: fileData,
	}

	_, err := client.PutObject(ctx, request)
	if err != nil {
		return nil, err
	}

	// TODO: insert url variable on frontend when configuring the cloud storage providers variable on the frontend
	fileUrl := common.String(fmt.Sprintf("https://axrj9wzaeep0.objectstorage.af-johannesburg-1.oci.customer-oci.com/n/axrj9wzaeep0/b/bucket-20230812-0418/o/%s", fileName))
	return fileUrl, nil
}

//ociStorage := NewOCIStorage(ociConfig)
//
//ociClient, err := ociStorage.Init()
//if err != nil {
//	log.Print(err)
//}
//
//ctx := context.Background()
//
//filePath := "./loader.txt"
//
//err = ociStorage.Put(ctx, ociClient, filePath)
//if err != nil {
//	log.Print(err)
//}
//_, err = ociStorage.Get(ctx, ociClient)
//if err != nil {
//	log.Print(err)
//}
