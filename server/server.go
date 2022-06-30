/*
Copyright © 2022 John Hooks

*/

package server

import (
	"fmt"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

const (
	bucket = "pages"
)

var docsTemplate = `
<h1>JetDocs</h1>
<h3>Documents</h3>
<ul>
  {{ range .  }}
	<li><a href={{ .URL }}>{{ .Name }}</a></li>
  {{ end }}
</ul>
`

type Server struct {
	Conn   *nats.Conn
	Router *mux.Router
	Port   int
}

type Link struct {
	URL  string
	Name string
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetNatsConn(nc *nats.Conn) *Server {
	s.Conn = nc
	return s
}

func (s *Server) SetRouter(r *mux.Router) *Server {
	s.Router = r
	return s
}

func (s *Server) getPage(id string) ([]byte, error) {
	js, err := s.Conn.JetStream()
	if err != nil {
		return nil, err
	}

	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, err
	}

	data, err := kv.Get(id)
	if err != nil {
		return nil, err
	}

	return data.Value(), err
}

func (s *Server) getPages() ([]string, error) {
	js, err := s.Conn.JetStream()
	if err != nil {
		return nil, err
	}

	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, err
	}

	return kv.Keys()
}

func (s *Server) GetPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	data, err := s.getPage(vars["id"])
	if err != nil && err == nats.ErrKeyNotFound {
		http.Error(w, "page not found", 404)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", 500)
		return
	}

	w.Write(markdown.ToHTML(data, nil, nil))
}

func (s *Server) GetPages(w http.ResponseWriter, r *http.Request) {
	pages, err := s.getPages()
	if err != nil && err == nats.ErrNoKeysFound {
		http.Error(w, "no documents found", 404)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}

	var urls []Link
	for _, v := range pages {
		link := Link{
			URL:  fmt.Sprintf("http://127.0.0.1:%d/pages/%s", s.Port, v),
			Name: v,
		}
		urls = append(urls, link)
	}

	tmpl := template.Must(template.New("").Parse(docsTemplate))

	if err := tmpl.Execute(w, urls); err != nil {
		log.Println(err)
	}

}

func StartEmbeddedNATS(nc *nats.Conn) (*nats.Conn, error) {
	sopts := natsserver.Options{
		JetStream: true,
		StoreDir:  viper.GetString("store-dir"),
		Host:      "127.0.0.1",
		Port:      44566,
	}

	ns, err := natsserver.NewServer(&sopts)
	if err != nil {
		return nil, err
	}

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		return nil, fmt.Errorf("NATS was not able to start")
	}

	return nats.Connect(ns.ClientURL())
}

func InitializeBucket(nc *nats.Conn) error {
	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	config := nats.KeyValueConfig{
		Bucket:  bucket,
		History: 10,
	}

	_, err = js.KeyValue(bucket)
	if err != nil && err != nats.ErrBucketNotFound {
		return err
	}

	_, err = js.CreateKeyValue(&config)
	if err != nil {
		return err
	}

	return nil

}
