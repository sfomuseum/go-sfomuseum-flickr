package s3

import (
	"bytes"
	"github.com/aaronland/go-storage"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	wof_s3 "github.com/whosonfirst/go-whosonfirst-aws/s3"
	"io"
	"io/ioutil"
)

type S3File struct {
	io.WriteCloser
	buf  *bytes.Buffer
	conn *wof_s3.S3Connection
	path string
}

// please make this better to periodically flush and append
// and all that good stuff... someone else must have done
// this by now, right? (20181120/thisisaaronland)

func NewS3File(conn *wof_s3.S3Connection, path string) (io.WriteCloser, error) {

	buf := new(bytes.Buffer)

	f := S3File{
		buf:  buf,
		conn: conn,
		path: path,
	}

	return &f, nil
}

func (f *S3File) Write(b []byte) (int, error) {
	return f.buf.Write(b)
}

func (f *S3File) WriteString(b string) (int, error) {
	return f.Write([]byte(b))
}

func (f *S3File) Close() error {

	r := bytes.NewReader(f.buf.Bytes())
	fh := ioutil.NopCloser(r)

	return f.conn.Put(f.path, fh)
}

type S3Store struct {
	storage.Store
	config *wof_s3.S3Config
	conn   *wof_s3.S3Connection
}

func NewS3Store(dsn string) (storage.Store, error) {

	cfg, err := wof_s3.NewS3ConfigFromString(dsn)

	if err != nil {
		return nil, err
	}

	conn, err := wof_s3.NewS3Connection(cfg)

	if err != nil {
		return nil, err
	}

	s := S3Store{
		config: cfg,
		conn:   conn,
	}

	return &s, nil
}

func (s *S3Store) URI(k string) string {
	return s.conn.URI(k)
}

func (s *S3Store) Get(k string) (io.ReadCloser, error) {
	return s.conn.Get(k)
}

func (s *S3Store) Create(k string, args ...interface{}) (io.WriteCloser, error) {

	return NewS3File(s.conn, k)
}

func (s *S3Store) Put(k string, in io.ReadCloser) error {

	return s.conn.Put(k, in)
}

func (s *S3Store) Delete(k string) error {

	return s.conn.Delete(k)
}

func (s *S3Store) Exists(k string) (bool, error) {

	_, err := s.conn.Head(k)

	if err != nil {

		aws_err := err.(awserr.Error)

		switch aws_err.Code() {
		case "NotFound":
			return false, nil
		case s3.ErrCodeNoSuchKey:
			return false, nil
		case s3.ErrCodeNoSuchBucket:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (s *S3Store) Walk(user_cb storage.WalkFunc) error {

	list_cb := func(obj *wof_s3.S3Object) error {
		return user_cb(obj.Key, obj)
	}

	list_opts := wof_s3.DefaultS3ListOptions()

	return s.conn.List(list_cb, list_opts)
}
