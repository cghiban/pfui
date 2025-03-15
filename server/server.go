package server

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"psui"
	"psui/service"
	"strings"
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

func NewHandlers(s service.Service) Handlers {

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

		_, err := s.PfCommand(op, ip)
		if err != nil {
			//fmt.Printf("err: %s\n", err)
			fmt.Printf("can't update device: %s\n", err)
			Error(w,
				map[string]string{"status": "error", "msg": "can't update device"},
				http.StatusInternalServerError,
			)
			return
		}
	}

	Success(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func NewRouter(h Handlers, cfg psui.Config) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", h.Ping).Methods("GET")

	auth := cfg.Auth
	if auth != nil && auth.User != "" && auth.Pass != "" {
		router.HandleFunc("/psui/devices", basicAuth(auth.User, auth.Pass)(h.DevicesHandler)).Methods("GET")
		router.HandleFunc("/psui/devices", basicAuth(auth.User, auth.Pass)(h.UpdateDeviceHandler)).Methods("PUT")
	} else {
		router.HandleFunc("/psui/devices", h.DevicesHandler).Methods("GET")
		router.HandleFunc("/psui/devices", h.UpdateDeviceHandler).Methods("PUT")
	}

	router.Use(loggingMiddleware)

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func basicAuth(username, password string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			payload, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 || pair[0] != username || pair[1] != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next(w, r)
		})
	}
}

func Pf(w http.ResponseWriter, r *http.Request, s service.Service) {
}
