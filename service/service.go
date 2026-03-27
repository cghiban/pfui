package service

import (
	"errors"
	"fmt"
	"net"
	"pfui"
	"slices"
)

type Service struct {
	Cfg        pfui.Config
	devicesMap map[string]string
}

func NewService(cfg pfui.Config) Service {
	dm := map[string]string{}
	for _, d := range cfg.Devices {
		dm[d.Mac] = d.Name
	}
	return Service{
		Cfg:        cfg,
		devicesMap: dm,
	}
}

func (s Service) GetHosts(filtered bool) ([]pfui.Host, error) {
	var err error
	hosts := []pfui.Host{}
	if false && filtered {
		//for k, v := range s.Cfg.Devices {
		for _, d := range s.Cfg.Devices {
			hosts = append(hosts, pfui.Host{
				Name:    d.Name,
				EthAddr: d.Mac,
			})
		}
	} else {
		hosts, err = pfui.ExecArp()
	}
	if err != nil {
		return []pfui.Host{}, err
	}

	leases := pfui.LoadLeases()
	pf := pfui.PF{}
	banned_ips := []string{}
	banned_ips, err = pf.TableShow(s.Cfg.PFTable)

	out := []pfui.Host{}
	for _, h := range hosts {
		name, exists := s.devicesMap[h.EthAddr]
		if exists {
			h.Name = name
		} else {
			l := leases.Find(map[string]string{"hw": h.EthAddr})
			if l != nil {
				h.Name = l.Name
			}
		}

		//func Index[S ~[]E, E comparable](s S, v E) int
		idx := slices.Index(banned_ips, h.IP.String())
		if idx >= 0 {
			h.Banned = true
			banned_ips = slices.Delete(banned_ips, idx, idx+1)
		}
		if filtered {
			if exists {
				out = append(out, h)
			}
		} else {
			out = append(out, h)
		}
	}

	// add any left over banned ips
	for _, bip := range banned_ips {

		host := pfui.Host{
			Name:   "????",
			IP:     net.ParseIP(bip),
			Banned: true,
		}
		l := leases.Find(map[string]string{"ip": bip})
		if l != nil {
			host.Name = l.Name
		}
		out = append(out, host)
	}

	return out, nil
}

func ValidateCommand(cmd string) bool {

	validCommands := []string{
		"tables",
		"table",
		"add",
		"delete",
	}

	return slices.Contains(validCommands, cmd)
}

func (s Service) PfCommand(cmd string, args ...string) ([]string, error) {

	if !ValidateCommand(cmd) {
		return []string{}, fmt.Errorf("invalid command received: %s", cmd)
	}

	pf := pfui.PF{}
	var err error

	output := []string{}
	switch cmd {
	case "tables":
		output, err = pf.Tables()
	case "table":
		output, err = pf.TableShow(s.Cfg.PFTable)
	case "add":
		if len(args) != 1 {
			err = errors.New("needs and arg: ip")
		} else {
			err = pf.TableAddEntry(s.Cfg.PFTable, args[0])
		}
	case "delete":
		if len(args) != 1 {
			err = errors.New("needs an arg: ip")
		} else {
			err = pf.TableDeleteEntry(s.Cfg.PFTable, args[0])
		}
	default:
		err = errors.New("invalid PF command")
	}

	return output, err
}
