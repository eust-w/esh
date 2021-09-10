package ssh

import (
	"errors"
	"esh/utils"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const (
	PasswordFlag string = "Password:"
	FailFlag     string = "denied"
	SureFlag     string = "sure"
)

func Run(ip,user, password, port,cmds string, runFlag bool) error {
	cmd := "ssh "+user+"@"+ip+" -p "+port+" "+cmds
	shell, err := utils.Shell()

	c := exec.Command(shell)

	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	defer func() { _ = ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

	if _, err := ptmx.Write([]byte(cmd + "; exit\n")); err != nil {
		return err
	}
	_, err = enterPassword(
		ptmx,
		password,
		runFlag,
	)
	if err != nil {
		return err
	}

	return nil
}

func enterPassword(ptmx *os.File, password string, runFlag bool) (string, error) {
	errChan := make(chan error)
	pwdChan := make(chan string)
	go func() {
		data := ""
		buf := make([]byte, 4096000)
		entered := false
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				errChan <- err
				break
			}
			if n == 0 {
				continue
			}
			data += string(buf[:n])
			if !entered && strings.Contains(data, PasswordFlag) {
				entered = true
				data = ""
				_, err := ptmx.Write([]byte(password + "\n"))
				if err != nil {
					return
				}
			} else if entered && len(data) > 5 {
				if strings.Contains(data, PasswordFlag) || strings.Contains(data, FailFlag) {
					errChan <- errors.New("connect fail")
					break
				}
				pwdChan <- data
				break
			} else if !entered && strings.Contains(data,SureFlag){
				data = ""
				_, err := ptmx.Write([]byte("yes" + "\n"))
				if err != nil {
					return
				}
			}
		}
	}()
	select {
	case newBuffered := <-pwdChan:
		os.Stdout.WriteString(strings.TrimPrefix(newBuffered,"\r\n"))
		go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
		if runFlag{
			_, _ = io.Copy(os.Stdout, ptmx)
		}
		return "", nil
	case err := <-errChan:
		return "", err
	}
}


func MultiRun(name,ip,user, password, port,cmds string, runFlag bool, outChan chan<- [2]string) error {
	cmd := "ssh "+user+"@"+ip+" -p "+port+" "+cmds
	shell, err := utils.Shell()

	c := exec.Command(shell)

	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	defer func() { _ = ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

	if _, err := ptmx.Write([]byte(cmd + "; exit\n")); err != nil {
		return err
	}
	_, err = multiEnterPassword(
		ptmx,
		name,
		password,
		runFlag,
		outChan,
	)
	if err != nil {
		return err
	}

	return nil
}

func multiEnterPassword(ptmx *os.File, name,password string, runFlag bool,outChan chan<- [2]string) (string, error) {
	errChan := make(chan error)
	pwdChan := make(chan string)
	go func() {
		data := ""
		buf := make([]byte, 4096)
		entered := false
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				errChan <- err
				break
			}
			if n == 0 {
				continue
			}
			data += string(buf[:n])
			log.Println("yes:",data,[]byte(data),"end")
			if !entered && strings.Contains(data, PasswordFlag) {
				entered = true
				data = ""
				_, err := ptmx.Write([]byte(password + "\n"))
				if err != nil {
					return
				}
			} else if entered && len(data) > 5 {
				if strings.Contains(data, PasswordFlag) || strings.Contains(data, FailFlag) {
					errChan <- errors.New("connect fail")
					break
				}
				pwdChan <- data
				break
			} else if !entered && strings.Contains(data,SureFlag){
				data = ""
				_, err := ptmx.Write([]byte("yes" + "\n"))
				if err != nil {
					return
				}
			}
		}
	}()
	select {
	case newBuffered := <-pwdChan:
		outChan <- [2]string{name,strings.TrimPrefix(newBuffered,"\r\n")}
		//os.Stdout.WriteString(strings.TrimPrefix(newBuffered,"\r\n"))
		go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
		if runFlag{
			_, _ = io.Copy(os.Stdout, ptmx)
		}
		return "", nil
	case err := <-errChan:
		return "", err
	}
}
