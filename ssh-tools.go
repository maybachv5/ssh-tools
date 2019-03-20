package main

import (
	"flag"
	"fmt"
	"github.com/ssh-tools/config"
	"github.com/ssh-tools/ssh_scp"
	"github.com/xclpkg/clcolor"
	"sync"
)

func main() {
	flag.Parse()
	switch {
	case *config.CommandType != "ssh" && *config.CommandType != "scp":
		return

	}
	//if (*commandType != "ssh" || *commandType != "scp") || (*commandType == "ssh" && *filePath == "") ||
	//	(*commandType == "scp" && (*remoteFile == "" || *sourceFile == "")) {
	//	return
	//}
	hosts, err := config.ParseServers(*config.FilePath)
	if err != nil {
		fmt.Println(clcolor.Red(err.Error()))
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(hosts))
	for _, host := range hosts {
		switch *config.CommandType {
		case "ssh":
			s := ssh_scp.SSH{host, *config.Command}
			go s.SSHFun(wg)
		case "scp":
			s := ssh_scp.SCP{host, *config.SourceFile, *config.RemoteFile}
			go s.SCPFun(wg)

		default:
			return
		}
	}
	wg.Wait()

}
