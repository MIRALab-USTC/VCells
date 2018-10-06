package main

import (
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"log"
	"mime"
	"net/http"
	"os/exec"
	"path"
	"strings"
	texttmplate "text/template"
	"time"
	"io/ioutil"
	"os"
	"regexp"
	"flag"
)

func router(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.URL.Path

	if r.Method == "POST" {

		fileInfo := new(FileInfo)
		fileInfo = handleUploads(r)

		rsp := new(RspJson)
		rsp.RetCode = fileInfo.ErrorCode
		rsp.ErrMsg = fileInfo.ErrorMsg

		prefix := "upload-images/"
		fileSuffix := path.Ext(fileInfo.Name)
		filenameOnly := strings.TrimSuffix(fileInfo.Name, fileSuffix)

		rsp.SrcImgPath = prefix + fileInfo.Name
		rsp.InitSegImgPath = prefix + filenameOnly + "_init_seg" + fileSuffix
		rsp.DstImgPath = prefix + filenameOnly + "_superpixel" + fileSuffix

		// execute VCells
		cmd := exec.Command("upload_files/VCellsCpp", "upload_files/"+fileInfo.Name)
		_, err := cmd.Output()
		if err != nil {
			log.Printf("[ERROR] Execute VCellsCpp failed, err: %s", err)
		}

		rspJson, err := json.Marshal(rsp)
		check(err)
		w.Header().Set("Cache-Control", "no-cache")
		jsonType := "application/json"
		if strings.Index(r.Header.Get("Accept"), jsonType) != -1 {
			w.Header().Set("Content-Type", jsonType)
		}

		fmt.Fprintln(w, string(rspJson))

	} else if r.Method == "GET" {
		switch {
		case url == "/index.html":
			t, _ := htmltemplate.ParseFiles("template/html/index.html")
			t.Execute(w, nil)

		case strings.HasSuffix(url, ".js"):
			url = "template" + url
			w.Header().Add("Content-Type", "text/javascript")
			t, _ := texttmplate.ParseFiles(url)
			t.Execute(w, nil)

		case strings.HasSuffix(url, ".css"):
			url = "template" + url
			w.Header().Add("Content-Type", "text/css")
			t, _ := texttmplate.ParseFiles(url)
			t.Execute(w, nil)

		default:
			t, _ := htmltemplate.ParseFiles("template/html/index.html")
			t.Execute(w, nil)
		}
	}
}

func main() {
	// set mime types
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".html", "text/html")

	// ticker to clean uploaded files
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		var imgTypes = regexp.MustCompile("(gif|p?jpeg|(x-)?png)")
		for range ticker.C {
			files, err := ioutil.ReadDir("upload_files")
			if err != nil {
				log.Printf("[ERROR] Ticker read dir failed, err: %s", err)
			}

			for _, f := range files {
				fileSuffix := path.Ext("upload_files/" + f.Name())
				if imgTypes.MatchString(fileSuffix) {
					os.Remove("upload_files/" + f.Name())
				}
			}
		}
	}()

	var port = flag.String("port", "8080", "port to listen to")
	var logPath = flag.String("log", "vcells.log", "log file path")

	flag.Parse()

	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("[ERROR] Failed to open log file!\n")
	}

	log.SetOutput(logFile)

	// http handle
	http.Handle("/upload-images/", http.StripPrefix("/upload-images/", http.FileServer(http.Dir("upload_files/"))))
	http.HandleFunc("/", router)

	log.Printf("[INFO] listening to port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+ *port, nil))
}
