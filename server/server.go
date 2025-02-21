package server

import (
	"fmt"
	"net/http"
	"psui"
	"psui/service"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	TMPLS_GLOB = "var/templates/*.gohtml"
)

type Handlers struct {
	Ping                http.HandlerFunc
	PfHandler           http.HandlerFunc
	DevicesHandler      http.HandlerFunc
	UpdateDeviceHandler http.HandlerFunc
}

func NewHandlers(s service.Service, cfg psui.Config) Handlers {

	funcMap := template.FuncMap{}
	t := template.Must(template.New("tmpls").Funcs(funcMap).ParseGlob(TMPLS_GLOB))

	return Handlers{
		Ping: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		},
		DevicesHandler: func(w http.ResponseWriter, r *http.Request) {
			Devices(w, r, s, t)
		},
		UpdateDeviceHandler: func(w http.ResponseWriter, r *http.Request) {
			UpdateDevice(w, r, s, t)
		},
		PfHandler: func(w http.ResponseWriter, r *http.Request) {
			Pf(w, r, s)
		},
	}
}

func Devices(w http.ResponseWriter, r *http.Request, s service.Service, t *template.Template) {

	hosts, err := s.GetHosts(true) // filtered
	//hosts, err := s.GetHosts(false)

	data := struct {
		Err   error
		Hosts []psui.Host
	}{
		Err:   err,
		Hosts: hosts,
	}

	if err := t.ExecuteTemplate(w, "devices.gohtml", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateDevice(w http.ResponseWriter, r *http.Request, s service.Service, t *template.Template) {

	q := r.URL.Query()
	op := q.Get("op")
	mac := q.Get("mac")
	found := false
	var host psui.Host

	hosts, err := psui.ExecArp()
	if err != nil {
		Error(w,
			map[string]string{"status": "error", "msg": "can't retrieve the devices"},
			http.StatusInternalServerError,
		)
		fmt.Printf("can't list the devices: %s\n", err)
		return
	}

	for _, h := range hosts {
		if h.EthAddr == mac {
			found = true
			host = h
			break
		}
	}

	if found {
		ip := host.IP.String()
		fmt.Printf("op: %s; device: %s\n", op, ip)

		out, err := s.PfCommand(op, ip)
		fmt.Printf("out: %s\n", out)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			Error(w,
				map[string]string{"status": "error", "msg": "can't update device"},
				http.StatusInternalServerError,
			)
			fmt.Printf("can't update device: %s\n", err)
			return

		}
	}

	Success(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func NewRouter(h Handlers) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/ping", h.Ping).Methods("GET")
	router.HandleFunc("/psui/devices", h.DevicesHandler).Methods("GET")
	router.HandleFunc("/psui/devices", h.UpdateDeviceHandler).Methods("PUT")

	return router
}

func Pf(w http.ResponseWriter, r *http.Request, s service.Service) {
}
