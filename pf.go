package psui

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

const (
	DOAS_CMD = "/usr/bin/doas"
	PF_CMD   = "/sbin/pfctl"
)

type PF struct {
}

// Tables - `pfctl -s Tables`
func (pf PF) Tables() (tables []string, err error) {
	var out []byte
	out, err = exec.Command("/usr/bin/doas", PF_CMD, "-s", "Tables").Output()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		tables = append(tables, scanner.Text())
	}

	return tables, err
}

// TableShow - returns blocked IPs
// pfctl -t nointernet  -T show
func (pf PF) TableShow(t string) (hosts []string, err error) {
	var out []byte
	out, err = exec.Command(DOAS_CMD, PF_CMD, "-t", t, "-T", "show").Output()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		hosts = append(hosts, strings.TrimSpace(scanner.Text()))
	}

	return hosts, err

}

// TableDeleteEntry - `pfctl -t <table> -T delete <ip>`
func (pf PF) TableDeleteEntry(table, ipstr string) (err error) {
	ip := net.ParseIP(ipstr)
	if ip.String() != ipstr || !ip.IsPrivate() {
		return errors.New("was expecting a private IPv4 IP address")
	}
	cmd := exec.Command(DOAS_CMD, PF_CMD, "-t", table, "-T", "delete", ip.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if !strings.Contains(string(out), `/1 addresses deleted.`) {
		return fmt.Errorf("expected `1/1 addresses deleted.`, but got %s instead\n", out)
	}
	return nil
}

// TableAddEntry - `pfctl -t <table> -T add <ipstr>`
func (pf PF) TableAddEntry(table, ipstr string) error {
	// 1 table created. /// -- optional, when we need to create the table
	// 1/1 addresses added.
	ip := net.ParseIP(ipstr)
	if ip.String() != ipstr || !ip.IsPrivate() {
		return errors.New("was expecting a private IPv4 IP address")
	}
	//fmt.Println(strings.Join([]string{DOAS_CMD, PF_CMD, "-t", table, "-T", "add", ip.String()}, " "))
	cmd := exec.Command(DOAS_CMD, PF_CMD, "-t", table, "-T", "add", ip.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if !strings.Contains(string(out), `/1 addresses added.`) {
		return fmt.Errorf("expected `1/1 addresses added.`, but got %s instead\n", out)
	}
	return nil
}
