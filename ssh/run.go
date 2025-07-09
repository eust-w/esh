package ssh

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
	"strings"
	"path/filepath"
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

// privateKeyAuth helper
func privateKeyAuth(path string) (ssh.AuthMethod, error) {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (sv *Server) RunTerminal() error {
	var authMethods []ssh.AuthMethod
	if sv.Password != "" {
		// detect key path vs password
		if strings.HasPrefix(sv.Password, "/") || strings.HasPrefix(sv.Password, "~") {
			if keyAuth, err := privateKeyAuth(sv.Password); err == nil {
				authMethods = append(authMethods, keyAuth)
			}
		} else {
			authMethods = append(authMethods, ssh.Password(sv.Password))
		}
	}

	config := &ssh.ClientConfig{
		User: sv.Username,
		Auth: authMethods,
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
