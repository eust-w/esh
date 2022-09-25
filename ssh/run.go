package ssh

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type Server struct {
	Username string //Username
	IP       string //IP Address
	Password string //Password
	Port     string //Port
	Client   *Client
}

func NewServer(user, password, ip string, port string) *Server {
	return &Server{Username: user, IP: ip, Port: port, Password: password}
}
func (sv *Server) Run(args []string) (string, error) {
	command := strings.Join(args, " ")
	runFlag := strings.Trim(command, "") == ""
	if runFlag {
		err := sv.RunTerminal()
		return "", err
	} else {
		out, err := sv.RunCommand(command)
		return out, err
	}
}

func (sv *Server) RunTerminal() error {
	config := &ssh.ClientConfig{
		User: sv.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sv.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", sv.IP+":"+sv.Port, config)
	Check(err)
	defer client.Close()

	session, err := client.NewSession()
	Check(err)
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	Check(err)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	fd2 := int(os.Stdout.Fd())
	termWidth, termHeight, err := terminal.GetSize(fd)
	termWidth, termHeight, err2 := terminal.GetSize(fd2)
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
	Check(err)

	err = session.Shell()
	Check(err)

	err = session.Wait()
	Check(err)
	return nil
}

func (sv *Server) RunCommand(cmd string) (string, error) {
	if cmd == "" {
		return "", errors.New("no cmd run")
	}
	var err error
	sv.Client, err = NewClient(sv.IP, sv.Port, sv.Username, sv.Password)
	output, err := sv.Client.Output(cmd)
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
