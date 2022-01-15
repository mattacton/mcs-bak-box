package main

import (
	"log"
	"os"

	"github.com/mattacton/mcs-bak-box/mcsbox"
)

var (
	user   string = ""
	pass   string = ""
	host   string = ""
	port   string = ""
	bucket string = ""
)

func initEnv() {
	os.Setenv("FTP_USER", user)
	os.Setenv("FTP_PWD", pass)
	os.Setenv("FTP_HOST", host)
	os.Setenv("FTP_PORT", port)
	os.Setenv("BUCKET", bucket)
}

func main() {
	initEnv()
	log.Print("Calling backup")
	mcsbox.BackupMCSBox()
}
