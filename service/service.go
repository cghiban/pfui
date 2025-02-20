package service

import (
	"psui"
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
	if filtered {
		out := []psui.Host{}
		for _, h := range all {
			if name, exists := s.devicesMap[h.EthAddr]; exists {
				h.Name = name
				out = append(out, h)
			}
		}

		return out, nil
	}
	return all, nil
}
