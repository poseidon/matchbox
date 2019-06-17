package storage

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// S3Client simple s3 client to interact with a bucket.
type S3Client struct {
	svc    s3iface.S3API
	bucket string
}

// S3ClientIface implements functions to list, read, write and delete objects
// from an s3 bucket respecting a certain prefix.
type S3ClientIface interface {
	listPrefix(prefix string) ([]s3.Object, error)
	readObject(prefix, name string) ([]byte, error)
	writeObject(prefix, name string, data []byte) error
	deleteObject(prefix, name string) error
}

// NewS3Client returns a new s3 client for the given region and bucket.
func NewS3Client(region, bucket string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return &S3Client{}, err
	}
	// Create S3 service client
	svc := s3.New(sess)

	return &S3Client{
		svc:    svc,
		bucket: bucket,
	}, nil
}

// listPrefix get a prefix and returns a list of objects that are
// named under that. Errors if the request fails.
func (s *S3Client) listPrefix(prefix string) ([]*s3.Object, error) {

	resp, err := s.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		return []*s3.Object{}, err
	}

	return resp.Contents, nil
}

// readObject reads data from an object with the given name, restricted under
// a specific prefix.
func (s *S3Client) readObject(prefix, name string) ([]byte, error) {

	res, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(prefix + name),
	})

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// writeObject pushes a data object to s3 under the given prefix and name.
func (s *S3Client) writeObject(prefix, name string, data []byte) error {

	_, err := s.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(prefix + name),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}
	return nil
}

// deletObject deletes an object from the s3 bucket.
func (s *S3Client) deleteObject(prefix, name string) error {

	_, err := s.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(prefix + name),
	})
	if err != nil {
		return err
	}

	err = s.svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(prefix + name),
	})
	if err != nil {
		return err
	}

	return nil
}
