package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type server struct {
	serverName string
	router    *mux.Router
	logger *log.Logger
	logInfo *logStruct
	counterFile string
}

type logStruct struct {
	ServerStartedAt time.Time `json:"server_started"`
	Requests int `json:"requests"`
	LastRequestAt time.Time `json:"lastrequest"`
	Routes []*logRoute `json:"routes"`
	isDirty bool
}

type logRoute struct {
	RouteName string `json:"route"`
	Requests int `json:"requests"`
	LastRequestAt time.Time `json:"lastrequest"`
}

func NewServer(serverName, counterFile string) *server {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	s := &server{
		serverName:  serverName,
		router:      mux.NewRouter(),
		logger:      logger,
		logInfo:     nil,
		counterFile: counterFile,
	}
	s.InitLogStruct()
	s.initLogWriter()
	return s
}

func (s *server) InitLogStruct() {
	s.logInfo = &logStruct{
		ServerStartedAt: time.Now().UTC(),
		Requests: 0,
		Routes:   nil,
	}
	s.readCounterFile()
}

func (s *server) initLogWriter() {
	go func() {
		for true {
			if s.logInfo.isDirty {
				s.saveCounterFile()
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func (s *server) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logRouteRequest(r.URL.Path)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
}

func (s *server) logRouteRequest(routeName string) {
	found := false
	for _, r := range s.logInfo.Routes {
		if r.RouteName == routeName {
			r.Requests++
			r.LastRequestAt = time.Now().UTC()
			found = true
			break
		}
	}
	if found == false {
		s.logInfo.Routes = append(s.logInfo.Routes, &logRoute{
			RouteName: routeName,
			Requests:  1,
			LastRequestAt: time.Now().UTC(),
		})
	}
	s.logInfo.Requests++
	s.logInfo.LastRequestAt = time.Now().UTC()
	s.logInfo.isDirty = true
}

func (s *server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logRouteRequest("/notfound")
		response := struct {
			Code string `json:code`
			Message string `json:message`
		}{
			Code: "404",
			Message: fmt.Sprintf("Route %s not found", r.URL.Path),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}
}

func (s *server) readCounterFile() {
	info, err := os.Stat(s.counterFile)
	if os.IsNotExist(err) {
		s.logger.Printf("No counter file found - initialized new log structure")
		return
	}
	if !info.IsDir() {
		file, _ := ioutil.ReadFile(s.counterFile)
		err = json.Unmarshal([]byte(file), &s.logInfo)
		if err != nil {
			return
		}
		s.logger.Printf("Initialized log structure from existing counter file %s", s.counterFile)
	}
	return
}

func (s *server) saveCounterFile() {
	s.logger.Printf("writing stats to counter file %s", s.counterFile)
	file, _ := json.Marshal(s.logInfo)
	_ = ioutil.WriteFile(s.counterFile, file, 0644)
	s.logInfo.isDirty = false
}