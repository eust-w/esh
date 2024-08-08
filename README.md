[🇨🇳](README_CN.md)[中文](README_CN.md)

## 🎉Introduction

esh is a cross-platform SSH link management tool designed to simplify and streamline the process of managing multiple SSH connections.It is simple yet powerful!

## ⚡Usage
It is particularly useful for developers and system administrators who frequently connect to various remote servers.

### Installation

#### 📜 Install from Source

The binary files are generated in the `out` directory.

#### 📦 Download and Install

- x86-64 Linux version: [esh-linux-amd64](https://github.com/eust-w/esh/releases)
- ARM-64 Linux version: [esh-linux-arm64](https://github.com/eust-w/esh/releases)
- x86-64 Mac version: [esh-mac-amd64](https://github.com/eust-w/esh/releases)
- x86-64 Windows version: [esh.exe](https://github.com/eust-w/esh/releases)

After downloading, you can run it directly from the command line. 

Note: Please run it via the command line!

## 🌱Interaction

esh command descriptions:

```
sql复制代码  
  add         Add a new remote SSH connection
  cluster     Use this command to connect to remote SSH or execute commands across multiple SSH sessions
  completion  Generate autocompletion scripts for the specified shell esh commands
  del         Delete an existing SSH connection using its name.
  help        Get detailed information about any esh command.
  list        List remote SSH sessions
  run         Connect to a remote SSH or run a command on it.
  set         Configure global settings for esh.
```

## ➕Development

1. Read information from `Home/esh_config.yaml`.
2. Passwords and usernames/IPs can be encrypted using AES. There should be at least two AES keys for encryption and decryption. Randomly select one (using the current time as the random seed) for encryption. Decryption is determined by the initial identifier. A root account can view the plaintext password, and the password is a salted value compiled at build time.
3. There should be login auto-completion functionality and a feature that requires entering a key to log in.
4. It should be able to execute remote commands like SSH and support cluster functionality.
