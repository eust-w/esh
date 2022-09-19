package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"strconv"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type Connect struct {
	serverName string //Server name
}
type Server struct {
	ServerName string //Server name
	Username   string //Username
	IP         string //IP Address
	Password   string //Password
	Port       int    //Port
}

func (conn *Connect) Run() {
	sv := Server{ServerName: "lt", Username: "root", IP: "172.20.1.34", Port: 22, Password: "password"}
	config := &ssh.ClientConfig{
		User: sv.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sv.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", sv.IP+":"+strconv.Itoa(sv.Port), config)
	fmt.Println("Dial:", err)
	Check(err)
	defer client.Close()

	session, err := client.NewSession()
	fmt.Println("NewSession:", err)
	Check(err)
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	fmt.Println("MakeRaw:", err)
	fmt.Println("oldState:", oldState)
	Check(err)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	fd2 := int(os.Stdout.Fd())
	termWidth, termHeight, err := terminal.GetSize(fd)
	fmt.Println("GetSize:", termWidth, termHeight, err)
	termWidth, termHeight, err2 := terminal.GetSize(fd2)
	fmt.Println("GetSize:", termWidth, termHeight, err2)
	Check(err2)

	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	//err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	//err = session.RequestPty("ms-terminal", termHeight, termWidth, modes)
	err = session.RequestPty("", termHeight, termWidth, modes)
	fmt.Println("RequestPty:", err)
	Check(err)

	err = session.Shell()
	fmt.Println("Shell:", err)
	Check(err)

	err = session.Wait()
	fmt.Println("Wait:", err)
	Check(err)
}

func main() {
	h := Connect{}
	h.Run()
}
