package main

import (
	"bytes"
	_ "embed"
	"errors"
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
	targetFs, err := createFolder()
	if err != nil {
		return
	}

	ip = GetOutboundIP()
	fmt.Printf("Hosting on http://%v:8080\n", ip)

	http.HandleFunc("/upload", uploadFileHandler(targetFs))

	//fs := http.FileServer(http.Dir(targetFs))
	//http.Handle("/"+targetFs+"/", http.StripPrefix("/files", fs))
	http.HandleFunc("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(bytes.Replace(index, []byte("{{}}"), []byte("\"http://"+ip.String()+":8080\""), 1))
	}))

	log.Print("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createFolder() (string, error) {
	if len(os.Args) != 2 {
		fmt.Println("no folder name given")
		return "", errors.New("no folder name given")
	}

	os.Mkdir(os.Args[1], 0750)
	return os.Args[1], nil
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
		//filePaths := []string{}

		fmt.Println("Receiving", len(files), "images")

		for _, file := range files {
			//filePath := "http://" + ip.String() + ":8080/files/" + file.Filename
			//filePaths = append(filePaths, filePath)
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
		//fmt.Println(filePaths)
		//enc := json.NewEncoder(w)
		//enc.Encode(struct {
		//	Filepaths []string `json:"filepaths"`
		//}{filePaths})
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
