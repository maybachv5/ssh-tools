package config

import (
	"bufio"
	"flag"
	"github.com/ssh-tools/ssh_scp"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	CommandType = flag.String("t", "", "操作类型:(ssh|scp)")
	FilePath    = flag.String("f", "", "服务器配置文件路径")
	Command     = flag.String("c", "", "执行命令，必须与-t ssh 一起执行")
	SourceFile  = flag.String("sf", "", "源文件位置,必须与-t scp 一起执行")
	RemoteFile  = flag.String("rf", "", "目标文件位置,必须与-t scp 一起执行")
)

func ParseServers(filepath string) ([]ssh_scp.HostInfo, error) {
	var hosts []ssh_scp.HostInfo
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		info := strings.Fields(string(scanner.Text()))
		if len(info) != 4 {
			log.Fatal("error parameter:", scanner.Text())
			continue
		}
		port, err := strconv.Atoi(info[1])
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, ssh_scp.HostInfo{info[0], port, info[2], info[3]})
	}
	return hosts, nil

}
