package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

var (
	user       string
	password   string
	host       string
	privateKey string
)

func init() {
	user = os.Getenv("SSH_USER")
	password = os.Getenv("SSH_PASSWORD")
	host = os.Getenv("SSH_HOST")
	privateKey = os.Getenv("SSH_PRIVATE_KEY")
}

func main() {
	pKey, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		log.Fatalf("could not parse private SSH key: %v", err)
	}

	log.Printf("About to connect to %s\n", host)
	sshC, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // don't do this in prod!
		BannerCallback: func(message string) error {
			log.Printf("Connected. SSH banner %s\n", message)
			return nil
		},
		HostKeyAlgorithms: []string{"ssh-rsa"},
		Timeout:           time.Minute,
	})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer sshC.Close()
	sftpC, err := sftp.NewClient(sshC)
	if err != nil {
		log.Fatalf("could not construct a SFTP client: %v", err)
	}
	defer sftpC.Close()
	_, err = sftpC.Stat(".")
	if err != nil {
		log.Fatalf("could not check if SFTP worked: %v", err)
	}
	log.Printf("All good!")
}
