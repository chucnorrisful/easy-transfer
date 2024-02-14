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
)

const maxUploadSize = 8 * 1024 * 1024 * 1024 // 8 GB

var ip net.IP

//go:embed index.html
var index []byte

func main() {
	targetFs := createFolder()
	ip = GetOutboundIP()
	go launchServer(targetFs)

	fmt.Printf("Hosting on http://%v:8080\n", ip)
	fmt.Printf("writing uploaded data to %s\n", targetFs)
	fmt.Println("Press ENTER to exit or close this terminal")
	_, _ = fmt.Scanln()
}

func launchServer(targetFs string) {
	// hosting the files again -> todo: optional feature with flag?
	fs := http.FileServer(http.Dir(targetFs))
	http.Handle("/"+targetFs+"/", http.StripPrefix("/files", fs))

	// the data receive endpoint
	http.HandleFunc("/upload", uploadFileHandler(targetFs))

	// hosting the website
	http.HandleFunc("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(bytes.Replace(index, []byte("{{}}"), []byte("\"http://"+ip.String()+":8080\""), 1))
	}))

	//todo: make port configurable
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createFolder() string {
	var dataDir = "data"
	if len(os.Args) >= 2 {
		dataDir = os.Args[1]
		if len(os.Args) > 2 {
			fmt.Println("WARNING: additional options specified which are ignored")
		}
	} else {
		fmt.Printf("WARNING: no folder name given, defaulting to '%s'\n", dataDir)
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0750)
	}
	return dataDir
}

func uploadFileHandler(targetFs string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderError(w, "only post allowed", http.StatusBadRequest)
			return
		}
		fmt.Println("Request parsing...")

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["images"]
		fmt.Println("Receiving", len(files), "images")

		for _, file := range files {
			out, err := os.Create("./" + targetFs + "/" + file.Filename)
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
