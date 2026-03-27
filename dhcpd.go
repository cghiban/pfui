package pfui

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

/*
IP      HW      Lease start     Lease end       Wins
192.168.9.18    10:27:f5:68:a5:2c       2026/03/27 10:58:52 UTC 2026/03/27 22:58:52 UTC EAP610-10-27-F5-68-A5-2C
192.168.9.19    bc:d7:d4:5b:e7:69       2026/03/27 10:21:34 UTC 2026/03/27 22:21:34 UTC StreamingStick4K-2
192.168.9.20    88:66:5a:59:28:26       2026/03/27 02:08:02 UTC 2026/03/28 02:08:02 UTC vin
192.168.9.21    9e:2f:a9:11:b0:3d       2026/03/27 09:48:46 UTC 2026/03/28 09:48:46 UTC iPhone
192.168.9.24    da:58:41:aa:e5:cd       2026/03/27 01:20:49 UTC 2026/03/28 01:20:49 UTC iPad
*/

type DateTime struct {
	time.Time
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006/01/02 15:04:05 UTC", csv)
	if err != nil {
		fmt.Printf("error parsing %s -- %s\n", csv, err)
	}
	return err
}

type Lease struct {
	IP    net.IP   `csv:"IP"`
	HW    string   `csv:"HW"`
	Start DateTime `csv:"Lease start"`
	End   DateTime `csv:"Lease end"`
	Name  string   `csv:"Wins"`
}

type Leases []*Lease

func (leases Leases) Find(params map[string]string) *Lease {

	for _, l := range leases {
		switch {
		case params["ip"] == l.IP.String():
			return l
		case params["hw"] == l.HW:
			return l
		}
	}

	return nil
}

func LoadLeases() Leases {

	leases := []*Lease{}

	f, err := os.OpenFile("/var/db/dhcpd.leases.tsv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("error reading leases: %s", err)
		return leases
	}
	defer f.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	if err = gocsv.Unmarshal(f, &leases); err != nil {
		log.Printf("error unmarshalling: %s", err)
	}
	return leases
}
