package testfakes

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	//"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

var (
	s3MockBucket = "s3MockBucket"

	MockClientGetError        = errors.New("MockClientGetError")
	MockClientPutError        = errors.New("MockClientPutError")
	MockClientDeleteError     = errors.New("MockClientDeleteError")
	MockClientWaitDeleteError = errors.New("MockClientWaitDeleteError")
	MockClientListError       = errors.New("MockClientListError")

	EmptyGetObjectOutput    = &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}
	EmptyPutObjectOutput    = &s3.PutObjectOutput{}
	EmptyDeleteObjectOutput = &s3.DeleteObjectOutput{}
	EmptyListObjectsOutput  = &s3.ListObjectsV2Output{}
)

type mockS3Client struct {
	s3iface.S3API
	getObjectMock                func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	putObjectMock                func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	deleteObjectMock             func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	waitUntilObjectNotExistsMock func(*s3.HeadObjectInput) error
	listObjectsV2Mock            func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
}

func (m mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return m.putObjectMock(input)
}

func (m mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.getObjectMock(input)
}

func (m mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return m.deleteObjectMock(input)
}

func (m mockS3Client) WaitUntilObjectNotExists(input *s3.HeadObjectInput) error {
	return m.waitUntilObjectNotExistsMock(input)
}

func (m mockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return m.listObjectsV2Mock(input)
}

var ErrorS3Client = mockS3Client{
	getObjectMock: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		return EmptyGetObjectOutput, MockClientGetError
	},
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return EmptyPutObjectOutput, MockClientPutError
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return EmptyDeleteObjectOutput, MockClientDeleteError
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return MockClientWaitDeleteError
	},
	listObjectsV2Mock: func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return EmptyListObjectsOutput, MockClientListError
	},
}

var NoErrorEmptyS3Client = mockS3Client{
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return EmptyPutObjectOutput, nil
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return EmptyDeleteObjectOutput, nil
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return nil
	},
	listObjectsV2Mock: func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return EmptyListObjectsOutput, nil
	},
}

var WorkingS3Client = mockS3Client{
	getObjectMock: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if strings.HasPrefix(*input.Key, "groups") {
			return returnGroup(*input.Key)
		}
		if strings.HasPrefix(*input.Key, "profiles") {
			return returnProfile(*input.Key)
		}
		if strings.HasPrefix(*input.Key, "ignition") {
			return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewBufferString(IgnitionYAML))}, nil
		}
		if strings.HasPrefix(*input.Key, "generic") {
			return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewBufferString(Generic))}, nil
		}
		return EmptyGetObjectOutput, MockClientGetError
	},
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return EmptyPutObjectOutput, MockClientPutError
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return EmptyDeleteObjectOutput, MockClientDeleteError
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return MockClientWaitDeleteError
	},
	listObjectsV2Mock: func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		if *input.Prefix == "groups/" {
			list := []*s3.Object{
				&s3.Object{Key: aws.String("groups/" + Group.Id + ".json")},
				&s3.Object{Key: aws.String("groups/" + GroupNoMetadata.Id + ".json")},
			}
			return &s3.ListObjectsV2Output{Contents: list}, nil
		}
		if *input.Prefix == "profiles/" {
			list := []*s3.Object{
				&s3.Object{Key: aws.String("profiles/" + Profile.Id + ".json")},
			}
			return &s3.ListObjectsV2Output{Contents: list}, nil
		}
		return EmptyListObjectsOutput, MockClientListError
	},
}

func returnGroup(gId string) (*s3.GetObjectOutput, error) {

	if gId == *aws.String("groups/" + Group.Id + ".json") {
		//json marshal will give base64 encoding for matadata
		raw := []byte(`
{
	"id": "test-group",
	"Name": "test group",
	"profile": "g1h2i3j4",
	"selector": {
		"uuid": "a1b2c3d4"
	},
	"metadata": {
		"pod_network": "10.2.0.0/16",
		"service_name":"etcd2"
	}
}
`)
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(raw))}, nil
	}
	if gId == *aws.String("groups/" + GroupNoMetadata.Id + ".json") {
		data, err := json.MarshalIndent(GroupNoMetadata, "", "\t")
		if err != nil {
			return EmptyGetObjectOutput, err
		}
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(data))}, nil
	}
	return EmptyGetObjectOutput, MockClientGetError
}

func returnProfile(pId string) (*s3.GetObjectOutput, error) {
	if pId == *aws.String("profiles/" + Profile.Id + ".json") {
		data, err := json.MarshalIndent(Profile, "", "\t")
		if err != nil {
			return EmptyGetObjectOutput, err
		}
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(data))}, nil
	}
	return EmptyGetObjectOutput, MockClientGetError
}
