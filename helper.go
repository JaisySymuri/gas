package gas

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var (
	RemoteServerURL string
	BackendName     string
	BackendType     string
	Host            string
	Port            string
	DatabaseType    string
	DatabaseName    string
)

func SetRemoteServerURL(url string) {
	RemoteServerURL = url
}

func SetBackendName(name string) {
	BackendName = name
}

func SetBackendType(typ string) {
	BackendType = typ
}

func SetHost(h string) {
	Host = h
}

func SetPort(p string) {
	Port = p
}

func SetDatabaseType(dbType string) {
	DatabaseType = dbType
}

func SetDatabaseName(dbName string) {
	DatabaseName = dbName
}

// Initialization function to set Server constants.
func InitializeServer(name, typ, h, p, dbType, dbName string) {
	// SetRemoteServerURL(remoteURL)
	SetBackendName(name)
	SetBackendType(typ)
	SetHost(h)
	SetPort(p)
	SetDatabaseType(dbType)
	SetDatabaseName(dbName)
}

// Initialization function to set Server constants.
func InitializeClient(remoteURL, name, typ, h, p string) {
	SetRemoteServerURL(remoteURL)
	SetBackendName(name)
	SetBackendType(typ)
	SetHost(h)
	SetPort(p)
}

// Function to get the name of the calling function using the Go runtime
func GetCallingFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func GetCurrentDirectoryName() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extract the directory's name using filepath.Base
	dirName := filepath.Base(cwd)
	return dirName, nil
}

func PerformHTTPRequest(method, url string, payload []byte) ([]byte, error) {
	fullURL := RemoteServerURL + url

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP request failed with status: %s", response.Status)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func GetCurrentURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
}

func CreateRequest(endpoint string, method string, body io.Reader) (*http.Request, error) {
	url := RemoteServerURL + endpoint
	req, err := http.NewRequest(method, url, body)
	return req, err
}
