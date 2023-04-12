package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/OmarFaruk-0x01/mongo-backup/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// var (
// 	aws_access_key = "AKIAQA7ZYSCEG4XFOUPM"
// 	aws_secret_key = "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"
// )

type Application struct {
	root_dir   string
	backup_dir string

	uri          string
	db           string
	filename     string
	archive_path string

	aws_access_key string
	aws_secret_key string

	s3_bucket string
	s3_region string
}

func main() {
	root_dir, err := os.UserHomeDir()
	if err != nil {
		root_dir = "."
	}
	root_dir = path.Join(root_dir, ".bkp")
	out_dir := path.Join(root_dir, "backups")
	uri := flag.String("uri", "", "Mongo Db URI")
	db_name := flag.String("db", "", "Mongo Db Name")

	s3_bucket := flag.String("bucket", "", "S3 bucket name")
	s3_region := flag.String("region", "", "S3 bucket region")

	aws_access_key := flag.String("aws-key", "", "AWS access key")
	aws_secret_key := flag.String("aws-sec", "", "AWS secret key")

	flag.Parse()

	filename := fmt.Sprintf("%s-%v.gz", *db_name, time.Now().Format("2006-01-02_15-04-05"))

	app := &Application{
		root_dir:     root_dir,
		backup_dir:   out_dir,
		filename:     filename,
		archive_path: path.Join(out_dir, filename),

		uri: *uri,
		db:  *db_name,

		aws_access_key: *aws_access_key,
		aws_secret_key: *aws_secret_key,

		s3_bucket: *s3_bucket,
		s3_region: *s3_region,
	}
	defer app.resetArchives()

	if err = app.Validate(); err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	err = app.Dump()
	if err != nil {
		log.Fatal(err)
	}

	err = app.SendToS3()
	if err != nil {
		log.Fatal(err)
	}

}

func (a *Application) Dump() (err error) {
	cmd := exec.Command("mongodump", "--gzip", "--uri", fmt.Sprintf("\"%s/%s\"", a.uri, a.db), fmt.Sprintf("--archive=%s", a.archive_path))

	// fmt.Println("Command ", cmd)

	output, err := cmd.CombinedOutput()

	// fmt.Println("Executed: ", string(output), err)
	log.Println(string(output))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (a *Application) SendToS3() error {
	credential := credentials.NewStaticCredentials(a.aws_access_key, a.aws_secret_key, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(a.s3_region),
		Credentials: credential,
	})
	if err != nil {
		return fmt.Errorf("failed to create aws session: %v", err)
	}

	s3Client := s3.New(sess)
	bucketName := a.s3_bucket
	backupFile := a.archive_path
	backupFileKey := fmt.Sprintf("database-backups/%s/%s", a.db, a.filename)
	_, err = os.Stat(backupFile)
	if err != nil {
		return fmt.Errorf("failed to access backup file: %v", err)
	}

	f, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer f.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(backupFileKey),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload backup file to s3: %v", err)
	}

	return nil
}

func (a *Application) Validate() error {
	err := os.MkdirAll(a.backup_dir, 0755)
	if err != nil {
		return err
	}
	db_uri, err := utils.GetURLOrigin(a.uri)

	if err != nil {
		return err
	}
	a.uri = *db_uri

	if utils.IsEmpty(a.aws_access_key) {
		return errors.New("aws access key required")
	}

	if utils.IsEmpty(a.aws_secret_key) {
		return errors.New("aws secret key required")
	}

	if utils.IsEmpty(a.db) {
		return errors.New("db name required")
	}

	if utils.IsEmpty(a.s3_bucket) {
		return errors.New("s3 bucket name required")
	}

	if utils.IsEmpty(a.s3_region) {
		return errors.New("s3 region name required")
	}

	return nil
}

func (a *Application) resetArchives() {
	err := os.Remove(fmt.Sprintf("%v", a.archive_path))
	// cmd := exec.Command("rm", "-rf", fmt.Sprintf("%v/*", a.backup_dir))
	// log.Println(cmd)
	if err != nil {
		log.Fatal(err)
	}

}
