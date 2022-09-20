package ssh

import (
	"bufio"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const DefaultTimeout = 30 * time.Second

type Client struct {
	*Config
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
	SFTPClient *sftp.Client
}

func (c *Client) Output(cmd string) ([]byte, error) {
	session, err := c.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.Output(cmd)
}

func NewDSN() (client *Client) {
	return nil
}
func Connect(cnf *Config) (client *Client, err error) {

	return nil, nil
}

func (cnf *Config) Connect() (client *Client, err error) {

	return nil, nil
}

// Close the underlying SSH connection
func (c *Client) Close() {
	c.SFTPClient.Close()
	c.SSHClient.Close()
	c.SSHSession.Close()
}

// New 创建SSH client
func New(cnf *Config) (client *Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User:            cnf.User,
		Timeout:         DefaultTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if cnf.Port == 0 {
		cnf.Port = 22
	}

	// 1. privite key file
	if len(cnf.KeyFiles) != 0 {
		if auth, err := AuthWithPrivateKeys(cnf.KeyFiles, cnf.Passphrase); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, auth)
		}

	} else {
		keypath := KeyFile()
		if FileExist(keypath) {
			if auth, err := AuthWithPrivateKey(keypath, cnf.Passphrase); err == nil {
				clientConfig.Auth = append(clientConfig.Auth, auth)
			}
		}

	}
	// 2. 密码方式 放在key之后,这样密钥失败之后可以使用Password方式
	if cnf.Password != "" {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(cnf.Password))
	}
	// 3. agent 模式放在最后,这样当前两者都不能使用时可以采用Agent模式
	if auth, err := AuthWithAgent(); err == nil {
		clientConfig.Auth = append(clientConfig.Auth, auth)
	}

	// hostPort := config.Host + ":" + strconv.Itoa(config.Port)
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(cnf.Host, strconv.Itoa(cnf.Port)), clientConfig)

	if err != nil {
		return client, errors.New("Failed to dial ssh: " + err.Error())
	}

	// create sftp client
	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return client, errors.New("Failed to conn sftp: " + err.Error())
	}

	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	// defer session.Close()

	return &Client{SSHClient: sshClient, SFTPClient: sftpClient, SSHSession: session}, nil
}

// NewClient 根据配置
func NewClient(host, port, user, password string) (client *Client, err error) {
	p, _ := strconv.Atoi(port)
	// if err != nil {
	// 	p = 22
	// }
	if user == "" {
		user = "root"
	}
	var config = &Config{
		Host:     host,
		Port:     p,
		User:     user,
		Password: password,
		// KeyFiles: []string{"~/.ssh/id_rsa"},
		Passphrase: password,
	}
	return New(config)
}

func NewWithAgent(Host, Port, User string) (client *Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User:            User,
		Timeout:         DefaultTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	auth, err := AuthWithAgent()
	if err != nil {
		return nil, err
	}
	clientConfig.Auth = append(clientConfig.Auth, auth)
	// hostPort := config.Host + ":" + strconv.Itoa(config.Port)
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(Host, Port), clientConfig)

	if err != nil {
		return client, errors.New("Failed to dial ssh: " + err.Error())
	}

	// create sftp client
	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(sshClient, sftp.MaxPacket(10240000)); err != nil {
		return client, errors.New("Failed to conn sftp: " + err.Error())
	}
	return &Client{SSHClient: sshClient, SFTPClient: sftpClient}, nil

}
func NewWithPrivateKey(Host, Port, User, Passphrase string) (client *Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User:            User,
		Timeout:         DefaultTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// 3. privite key file
	auth, err := AuthWithPrivateKey(KeyFile(), Passphrase)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	clientConfig.Auth = append(clientConfig.Auth, auth)

	// hostPort := config.Host + ":" + strconv.Itoa(config.Port)
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(Host, Port), clientConfig)

	if err != nil {
		return client, errors.New("Failed to dial ssh: " + err.Error())
	}

	// create sftp client
	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return client, errors.New("Failed to conn sftp: " + err.Error())
	}
	return &Client{SSHClient: sshClient, SFTPClient: sftpClient}, nil

}

func KeyFile() string {

	home, err := os.UserHomeDir()
	if err != nil {
		println(err.Error())
		return ""
	}
	key := filepath.ToSlash(path.Join(home, ".ssh/id_rsa"))
	return key
}
func FileExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func MkdirAll(path string) error {
	// 检测文件夹是否存在   若不存在  创建文件夹
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		}
	}
	return nil
}

//Md5File 计算md5
func Md5File(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	hash := crypto.MD5.New()
	_, err = io.Copy(hash, r)
	if err != nil {
		return "", err
	}

	out := hex.EncodeToString(hash.Sum(nil))
	return out, nil
}
