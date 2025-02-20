package server

import (
	"net/http"
	"psui"
	"psui/service"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	TMPLS_GLOB = "var/templates/*.gotmpl"
)

type Handlers struct {
	Ping                http.HandlerFunc
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
	}
}

func Devices(w http.ResponseWriter, r *http.Request, s service.Service, t *template.Template) {

	hosts, err := s.GetHosts(true) // filtered

	data := struct {
		Err   error
		Hosts []psui.Host
	}{
		Err:   err,
		Hosts: hosts,
	}

	if err := t.ExecuteTemplate(w, "devices.gotmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateDevice(w http.ResponseWriter, r *http.Request, s service.Service, t *template.Template) {
}

func NewRouter(h Handlers) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/ping", h.Ping).Methods("GET")
	router.HandleFunc("/psui/devices", h.DevicesHandler).Methods("GET")
	router.HandleFunc("/psui/devices", h.UpdateDeviceHandler).Methods("PUT")

	return router
}
