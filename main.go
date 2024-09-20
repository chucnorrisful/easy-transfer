package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

const maxUploadSize = 8 * 1024 * 1024 * 1024 // 8 GB
const targetFolder = "data"
const port = "8080"

var ip net.IP

//go:embed assets/index.html
var index []byte

func main() {
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		err = os.Mkdir(targetFolder, 0750)
		if err != nil {
			panic(err)
		}
	}
	go launchServer()

	ip = GetOutboundIP()
	fmt.Printf("Hosting on http://%v:%v\n", ip, port)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf("writing uploaded data to %s\n", wd+`\`+targetFolder)
	cmd := exec.Command(`explorer`, `/open,`, wd+`\`+targetFolder)
	cmd.Run()

	fmt.Println("Press ENTER to exit or close this terminal")
	_, _ = fmt.Scanln()
}

func launchServer() {
	// todo: https

	// hosting the files again -> todo: optional feature with flag?g
	fs := http.FileServer(http.Dir(targetFolder))
	http.Handle("/"+targetFolder+"/", http.StripPrefix("/files", fs))

	// the data receive endpoint
	http.HandleFunc("/upload", uploadFileHandler())

	// hosting the website
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(bytes.Replace(index, []byte("{{}}"), []byte("\"http://"+ip.String()+":8080\""), 1))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderError(w, "only post allowed", http.StatusBadRequest)
			return
		}
		fmt.Println("Request parsing...")

		if err := r.ParseForm(); err != nil {
			fmt.Printf("Could not parse form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["file"]
		fmt.Println("Receiving", len(files), "images")

		for _, file := range files {
			out, err := os.Create("./" + targetFolder + "/" + file.Filename)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			if file.Size > maxUploadSize {
				renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
				return
			}

			readerFile, _ := file.Open()
			_, err = io.Copy(out, readerFile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(".")
		}
		fmt.Println("done!")
		w.Write([]byte(`<div>UPLOAD SUCCESSFUL <button onClick="window.location.reload();">⟳</button></div>`))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func GetOutboundIP() net.IP {
	// todo: debug for mixed networks/host via wifi
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
