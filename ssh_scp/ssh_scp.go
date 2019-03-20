package ssh_scp

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/xclpkg/clcolor"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"sync"
	"time"
)

type HostInfo struct {
	Host   string
	Port   int
	User   string
	Passwd string
}

type SSH struct {
	HostInfo
	Cmd string
}

type SCP struct {
	HostInfo
	SourceFile string
	RemoteFile string
}

func (s HostInfo) client() (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Passwd))
	clientConfig = &ssh.ClientConfig{
		User:            s.User,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr = fmt.Sprintf("%s:%d", s.Host, s.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return client, nil
}

func (s SSH) SSHFun(group *sync.WaitGroup) {
	defer group.Done()
	var (
		client  *ssh.Client
		session *ssh.Session
		err     error
	)
	if client, err = s.client(); err != nil {
		fmt.Println(clcolor.Red(err.Error()))
		return
	}
	if session, err = client.NewSession(); err != nil {
		fmt.Println(clcolor.Red(err.Error()))
		return
	}
	defer session.Close()
	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	session.Run(s.Cmd)
	fmt.Printf("%s\n%s%s\n", clcolor.Cyan(s.Host+" output:"), stdOut.String(), stdErr.String())

}

func (s SCP) SCPFun(group *sync.WaitGroup) {
	defer group.Done()
	var (
		client     *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	if client, err = s.client(); err != nil {
		fmt.Println(clcolor.Red(err.Error()))
		return
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(client); err != nil {
		fmt.Println(clcolor.Red(err.Error()))
		return
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(s.SourceFile)
	if err != nil {
		fmt.Println(clcolor.Red(err.Error()))
		return
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(s.SourceFile)
	fmt.Println(path.Join(s.RemoteFile, remoteFileName))
	dstFile, err := sftpClient.Create(path.Join(s.RemoteFile, remoteFileName))
	if err != nil {
		fmt.Println(clcolor.Blue(err.Error()))
		return
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}

}
