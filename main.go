package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/skip2/go-qrcode"
)

const maxUploadSize = 8 * 1024 * 1024 * 1024 // 8 GiB
const targetFolder = "data"

//go:embed assets/index.html
var indexPage []byte

func main() {
	var tlsEnabled bool
	flag.BoolVar(&tlsEnabled, "secure", false, "enables HTTPS encryption with self-signed certificate")
	flag.BoolVar(&tlsEnabled, "s", false, "shorthand for -secure")
	flag.Parse()

	ip := GetOutboundIP()
	go launchServer(ip, tlsEnabled)

	cancelChan := make(chan os.Signal, 1)
	endedChan := make(chan struct{})
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		link := fmt.Sprintf("http://%v:8080", ip)
		if tlsEnabled {
			link = fmt.Sprintf("https://%v", ip)
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		qr, err := qrcode.New(link, qrcode.Medium)
		if err != nil {
			panic(err)
		}

		fmt.Printf(`
                      Easy Transfer
        high-speed local ad-hoc file transmission
                  by chucnorrisful, 2025

hosting the upload website on %v
uploaded data will be written to:
    %v\%v

            %v

       press ENTER to exit or close this terminal
`, link, wd, targetFolder, strings.ReplaceAll(qr.ToSmallString(false), "\n", "\n            "))
		_, _ = fmt.Scanln()
		endedChan <- struct{}{}
	}()

	select {
	case <-cancelChan:
	case <-endedChan:
	}

	wd, _ := os.Getwd() //would have paniced earier

	// todo: think about changing UX, open website on receiver and add button to open target folder
	cmd := exec.Command(`explorer`, `/open,`, wd+`\`+targetFolder)
	_ = cmd.Run()
}

func launchServer(ip net.IP, tlsEnabled bool) {
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		err = os.Mkdir(targetFolder, 0750)
		if err != nil {
			panic(err)
		}
	}

	// hosting the files again -> todo: optional feature with flag?g
	fs := http.FileServer(http.Dir(targetFolder))
	http.Handle("/"+targetFolder+"/", http.StripPrefix("/"+targetFolder+"/", fs))

	// the data receive endpoint
	http.HandleFunc("/upload", uploadFileHandler())

	// hosting the website
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(indexPage)
	})

	if !tlsEnabled {
		log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
		return
	}

	tlsConf, err := createSelfSignedCertificate()
	if err != nil {
		panic(err)
	}

	server := http.Server{
		TLSConfig: tlsConf,
		Handler:   http.DefaultServeMux,
		Addr:      ":443",
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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
	_, _ = w.Write([]byte(message))
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

func createSelfSignedCertificate() (*tls.Config, error) {
	// code from https://shaneutt.com/blog/golang-ca-and-signed-cert-go/

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumberCa, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}
	serialNumberCert, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial2 number: %v", err)
	}

	ca := &x509.Certificate{
		SerialNumber: serialNumberCa,
		Subject: pkix.Name{
			Organization: []string{"FOSS"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	cert := &x509.Certificate{
		SerialNumber: serialNumberCert,
		Subject: pkix.Name{
			Organization: []string{"FOSS"},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{serverCert}}, err
}
