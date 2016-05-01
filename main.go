package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":8080", "address to listen on")
)

func main() {
	flag.Parse()

	if *listen == "" {
		log.Fatal("listen must be provided")
	}

	http.HandleFunc("/", fetchHashicorpPackage)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func fetchHashicorpPackage(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://releases.hashicorp.com" + r.RequestURI)
	if err != nil {
		if resp.StatusCode >= 400 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := bytes.NewReader(body)
	zipReader, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tw := tar.NewWriter(w)

	for _, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fc, err := ioutil.ReadAll(rc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hdr := &tar.Header{
			Name: f.Name,
			Mode: int64(f.Mode()),
			Size: int64(len(fc)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := tw.Write(fc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := tw.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
