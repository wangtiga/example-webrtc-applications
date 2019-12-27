package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/pion/webrtc/v2"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {

	// Generate pem file for https
	genPem()

	// Create a MediaEngine object to configure the supported codec
	m = webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	// Create the API object with the MediaEngine
	api = webrtc.NewAPI(webrtc.WithMediaEngine(m))

}

type Flags struct {
	ListenAddr string
	CertPath   string
	KeyPath    string
}

func main() {
	if err := prometheus.Register(prommod.NewCollector("sfu_ws")); err != nil {
		panic(err)
	}

	param := Flags{
		ListenAddr: ":443",
		CertPath:   "cert.pem",
		KeyPath:    "key.pem",
	}
	flag.StringVar(&param.ListenAddr, "l", param.ListenAddr, "https ip:port")
	flag.StringVar(&param.CertPath, "cert", param.CertPath, "cert file path")
	flag.StringVar(&param.KeyPath, "key", param.KeyPath, "key file path")
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())

	// Websocket handle func
	http.HandleFunc("/ws", room)

	// Html handle func
	http.HandleFunc("/", web)

	// Support https, so we can test by lan
	fmt.Println("Web listening " + param.ListenAddr)
	panic(http.ListenAndServeTLS(
		param.ListenAddr,
		param.CertPath,
		param.KeyPath,
		nil,
	))
}
