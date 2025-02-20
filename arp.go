package psui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Host struct {
	IP      net.IP
	Name    string
	EthAddr string
	If      string
	Expire  string
	Flags   string
}

func parseLine(line string) (Host, error) {

	fields := strings.Fields(line)
	//if len(fields) < 3

	return Host{
		IP:      net.ParseIP(fields[0]),
		EthAddr: fields[1],
		If:      fields[2],
		Expire:  fields[3],
	}, nil
}
func Parse(s string) ([]Host, error) {

	hosts := []Host{}
	scanner := bufio.NewScanner(strings.NewReader(s))

	i := 1
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Index(line, "Host") == 0 {
			continue
		}
		h, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing line %d: %s", i, err)
			continue
		}
		hosts = append(hosts, h)
		//fmt.Println(i, line, h) // Println will add back the final '\n'
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return hosts, nil
}

func ExecArp() ([]Host, error) {
	var out string
	//var err error
	//out, err = exec.Command("/usr/sbin/arp", "-an").Output()
	// return Parse(string(out))

	out = `Host                                 Ethernet Address    Netif Expire    Flags
141.157.237.1                        b2:a8:6e:fd:5a:22     em0 17m0s
141.157.237.108                      00:e0:67:18:19:0c     em0 permanent l
192.168.9.1                          00:e0:67:18:19:0d     em1 permanent l
192.168.9.4                          fc:4d:d4:31:4a:42     em1 17m10s
192.168.9.14                         b8:27:eb:81:14:22     em1 14m5s
192.168.9.19                         bc:d7:d4:5b:e7:69     em1 19m50s
192.168.9.20                         88:66:5a:59:28:26     em1 18m54s
192.168.9.21                         9e:2f:a9:11:b0:3d     em1 19m35s
192.168.9.24                         da:58:41:aa:e5:cd     em1 19m51s
192.168.9.26                         b8:27:eb:b9:5d:d5     em1 1m32s
192.168.9.27                         78:20:a5:6f:72:e4     em1 1m5s
192.168.9.30                         a8:b5:7c:20:39:6a     em1 19m54s    `
	return Parse(out)
}
