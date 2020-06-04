package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ContainerName represents the name of the container
const ContainerName = "wpa_supplicant-udmpro"

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <user> <host:port> <pwd-file>", os.Args[0])
	}

	pwd, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatalf("Failed to read password from specified file: %s", os.Args[3])
	}
	pwdStr := strings.TrimSpace(string(pwd))

	config := &ssh.ClientConfig{
		User:            os.Args[1],
		Auth:            []ssh.AuthMethod{ssh.Password(string(pwdStr))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", os.Args[2], config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	running := checkContainerStatus(conn)

	if !running {
		log.Printf("%s container is not running, attempting to start it", ContainerName)
		status := startContainer(conn)

		if !status {
			log.Fatalf("Failed to start the the container: %s", ContainerName)
		} else {
			log.Printf("Successfully started %s container", ContainerName)
		}
	} else {
		log.Printf("%s container is running", ContainerName)
	}
}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)

	sessStdErr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStdErr)

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func checkContainerStatus(conn *ssh.Client) (ret bool) {
	cmd := fmt.Sprintf("podman container inspect %s", ContainerName)

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	runCommand(cmd, conn)
	return true
}

func startContainer(conn *ssh.Client) (ret bool) {
	cmd := fmt.Sprintf("podman start %s", ContainerName)

	defer func() {
		if err := recover(); err != nil {
			log.Println("failed starting container:", err)
		}
	}()

	runCommand(cmd, conn)
	return true
}
