package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

func Shell() (string, error) {
	switch runtime.GOOS {
	case "linux", "openbsd", "freebsd":
		return NixShell()
	case "android":
		return AndroidShell()
	case "darwin":
		return DarwinShell()
	}

	return "", errors.New("Undefined GOOS: " + runtime.GOOS)
}

func NixShell() (string, error) {
	us, err := user.Current()
	if err != nil {
		return "/bin/bash", err
	}

	out, err := exec.Command("getent", "passwd", us.Uid).Output()
	if err != nil {
		return "/bin/bash", err
	}

	ent := strings.Split(strings.TrimSuffix(string(out), "\n"), ":")
	return ent[6], nil
}

func AndroidShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", errors.New("SHELL not defined in android.")
	}
	return shell, nil
}

func DarwinShell() (string, error) {
	dir := "Local/Default/Users/" + os.Getenv("USER")
	out, err := exec.Command("dscl", "localhost", "-read", dir, "UserShell").Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("UserShell: (/[^ ]+)\n")
	matched := re.FindStringSubmatch(string(out))
	shell := matched[1]
	if shell == "" {
		return "", errors.New(fmt.Sprintf("Invalid output: %s", string(out)))
	}

	return shell, nil
}

func ByteXor(b1, b2 []byte) (out []byte) {
	for k, v := range b1 {
		out = append(out, v^b2[k])
	}
	return
}

//func CmdRun(cmd string) (string, error) {
//	var (
//		result []byte
//		err    error
//		stderr bytes.Buffer
//	)
//	commandName , _:=Shell()
//	ccmd :=exec.Command(commandName, "-c", cmd)
//	ccmd.Stderr = &stderr
//	result, err = ccmd.Output()
//	log.Println("1",result,err,stderr.String(),cmd,ccmd.Args)
//	if err != nil {
//		return "", err
//	}
//	return strings.TrimSpace(string(result)), nil
//}
//
//func DeleteHistory(str string)bool{
//	defer func() {recover()
//	}()
//	para1 := fmt.Sprintf("cat ~/.bash_history |grep \"%s\"|awk '{print $1}'",str)
//	ids,err:=CmdRun(para1)
//	if err !=nil{
//		return false
//	}
//	idl:=strings.Split(ids, " ")
//	fmt.Println(idl,len(idl))
//	for _,id := range idl{
//		para2 := fmt.Sprintf("history -d %s",id)
//		_,err=CmdRun(para2)
//		if err !=nil{
//			return false
//		}
//	}
//	return true
//}
