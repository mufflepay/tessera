package cloudStorage

import (
	"errors"
	"github.com/oracle/oci-go-sdk/common"
)

type StorageClient interface {
	GetStorageClient(storageType string, accountId string, r2BucketName string, accessKeyId string, accessKeySecret string, tenancyOCID string, userOCID string, region string, fingerprint string, privateKey string, namespace string, ociBucketName string) (*CloudflareR2Storage, *OCIStorage, error)
}

func GetStorageClient(storageType string, accountId string, r2BucketName string, accessKeyId string, accessKeySecret string, tenancyOCID string, userOCID string, region string, fingerprint string, privateKey string, namespace string, ociBucketName string) (*CloudflareR2Storage, *OCIStorage, error) {
	switch storageType {
	case "r2":
		r2Config := CloudflareR2Storage{
			AccountId:       accountId,
			BucketName:      r2BucketName,
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
		}
		r2Storage := NewCloudflareR2Storage(r2Config)

		return r2Storage, nil, nil
	case "oci":
		config := common.NewRawConfigurationProvider(tenancyOCID, userOCID, region, fingerprint, privateKey, nil)
		ociConfig := OCIStorage{
			Config:     config,
			Namespace:  namespace,
			BucketName: ociBucketName,
		}

		ociStorage := NewOCIStorage(ociConfig)

		return nil, ociStorage, nil
	default:
		return nil, nil, errors.New("unsupported storage type")
	}
}
