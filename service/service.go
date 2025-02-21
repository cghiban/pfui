package service

import (
	"errors"
	"fmt"
	"psui"
	"slices"
)

type Service struct {
	cfg        psui.Config
	devicesMap map[string]string
}

func NewService(cfg psui.Config) Service {
	dm := map[string]string{}
	for _, d := range cfg.Devices {
		dm[d.Mac] = d.Name
	}
	return Service{
		cfg:        cfg,
		devicesMap: dm,
	}
}

func (s Service) GetHosts(filtered bool) ([]psui.Host, error) {
	all, err := psui.ExecArp()
	if err != nil {
		return []psui.Host{}, err
	}

	pf := psui.PF{}
	banned_ips := []string{}
	banned_ips, err = pf.TableShow(s.cfg.PFTable)

	out := []psui.Host{}
	var name string
	var exists bool
	for _, h := range all {
		name, exists = s.devicesMap[h.EthAddr]
		if exists {
			h.Name = name
		}
		// func Contains[S ~[]E, E comparable](s S, v E) bool
		if slices.Contains(banned_ips, h.IP.String()) {
			h.Banned = true
		}
		if filtered {
			if exists {
				out = append(out, h)
			}
		} else {
			out = append(out, h)
		}
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

	pf := psui.PF{}
	var err error

	output := []string{}
	switch cmd {
	case "tables":
		output, err = pf.Tables()
	case "table":
		output, err = pf.TableShow(s.cfg.PFTable)
	case "add":
		if len(args) != 1 {
			err = errors.New("needs and arg: ip")
		} else {
			err = pf.TableAddEntry(s.cfg.PFTable, args[0])
		}
	case "delete":
		if len(args) != 1 {
			err = errors.New("needs an arg: ip")
		} else {
			err = pf.TableDeleteEntry(s.cfg.PFTable, args[0])
		}
	default:
		err = errors.New("invalid PF command")
	}

	return output, err
}
