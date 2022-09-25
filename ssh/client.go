package ssh

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"strconv"
	"time"
)

const DefaultTimeout = 30 * time.Second

type Client struct {
	*Config
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
}

func AuthWithAgent() (ssh.AuthMethod, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errors.New("Agent Disabled")
	}
	socks, err := net.Dial("unix", sock)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_client := agent.NewClient(socks)
	signers, err := _client.Signers()
	return ssh.PublicKeys(signers...), nil
}

func (c *Client) Output(cmd string) ([]byte, error) {
	session, err := c.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.Output(cmd)
}

func (c *Config) Connect() (client *Client, err error) {

	return nil, nil
}

func New(cnf *Config) (client *Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User:            cnf.User,
		Timeout:         DefaultTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if cnf.Port == 0 {
		cnf.Port = 22
	}

	if cnf.Password != "" {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(cnf.Password))
	}
	if auth, err := AuthWithAgent(); err == nil {
		clientConfig.Auth = append(clientConfig.Auth, auth)
	}
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(cnf.Host, strconv.Itoa(cnf.Port)), clientConfig)

	if err != nil {
		return client, errors.New("Failed to dial ssh: " + err.Error())
	}
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	return &Client{SSHClient: sshClient, SSHSession: session}, nil
}

func NewClient(host, port, user, password string) (client *Client, err error) {
	p, _ := strconv.Atoi(port)
	if user == "" {
		user = "root"
	}
	var config = &Config{
		Host:     host,
		Port:     p,
		User:     user,
		Password: password,
	}
	return New(config)
}
