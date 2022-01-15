package mcsbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jlaffaye/ftp"
)

type Config struct {
	User, Pass, Host, Port, Bucket string
}

func BackupMCSBox() error {
	config := Config{os.Getenv("FTP_USER"), os.Getenv("FTP_PWD"), os.Getenv("FTP_HOST"), os.Getenv("FTP_PORT"), os.Getenv("BUCKET")}

	configFilesToBackup := []string{
		"ops.json",
		"server.properties",
		"whitelist.json",
		"banned-ips.json",
		"banned-players.json"}

	worldFileName := "world.zip"

	c, err := ftp.Dial(fmt.Sprintf("%s:%s", config.Host, config.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(config.User, config.Pass)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Problem getting client Error: %s", err)
	}
	defer client.Close()

	err = backupToCloudStorage(c, client, ctx, config.Bucket, worldFileName)
	if err != nil {
		log.Fatalf("Could not backup world. Error: %s", err)
	}

	for _, filePath := range configFilesToBackup {
		err := backupToCloudStorage(c, client, ctx, config.Bucket, filePath)
		if err != nil {
			fmt.Printf("Skipping file '%s'", err)
			continue
		}
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func backupToCloudStorage(ftpConn *ftp.ServerConn, storageClient *storage.Client,
	ctx context.Context, bucket string, filePath string) error {

	log.Printf("Backup up %s", filePath)

	_, err := ftpConn.NameList(filePath)
	if err != nil {
		log.Printf("Problem finding file '%s' Error: %s", filePath, err)
		return errors.New(filePath)
	}

	source, err := ftpConn.Retr(filePath)
	if err != nil {
		log.Printf("Could not read file '%s' Error: %s", filePath, err)
		return errors.New(filePath)
	}
	defer source.Close()

	sourceData, err := ioutil.ReadAll(source)
	if err != nil {
		log.Printf("Problem saving source to memory %s", err)
		return errors.New(filePath)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*8)
	defer cancel()

	wc := storageClient.Bucket(bucket).Object(filePath).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(sourceData)); err != nil {
		log.Printf("Problem copying %s to cloud storage. Error: %s", filePath, err)
		return errors.New(filePath)
	}

	if err := wc.Close(); err != nil {
		log.Printf("Problem closing cloud storate writer. Error: %s", err)
		return errors.New(filePath)
	}

	return nil
}
