package psui

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Host struct {
	IP      net.IP
	Name    string
	EthAddr string
	If      string
	Expire  string
	Flags   string
	Banned  bool
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
func Parse(arp_out []byte) ([]Host, error) {

	hosts := []Host{}
	scanner := bufio.NewScanner(bytes.NewReader(arp_out))

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
		fmt.Println("»", i, line, h) // Println will add back the final '\n'
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return hosts, nil
}

func ExecArp() ([]Host, error) {
	out, err := exec.Command("/usr/sbin/arp", "-an").Output()
	if err != nil {
		return []Host{}, err
	}

	return Parse(out)
}
