package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	tls "github.com/refraction-networking/utls"
	"golang.org/x/exp/slices"
)

//go:embed deviceids.txt
var deviceFile embed.FS

func deviceCheck() error {
	file, _ := deviceFile.ReadFile("deviceids.txt")
	deviceIds := strings.Split(string(file), "\n")
	id, err := machineid.ProtectedID("**TaskMaster**")
	if err != nil {
		return err
	}
	if !slices.Contains(deviceIds, id) {
		return ErrUnauth
	}
	return nil
}

func createSiteDir(site string) error {
	if err := os.Mkdir(site, 0755); err != nil {
		return err
	}
	if err := createTaskTemplate(site, false); err != nil {
		return err
	}
	return nil
}

func getFolderList() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	dirList := []string{}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			dirList = append(dirList, fileInfo.Name())
		}
	}
	return dirList, nil
}

func initialCheck() error {
	updated := false

	dirList, err := getFolderList()
	if err != nil {
		return err
	}

	for _, site := range SITES {
		if !slices.Contains(dirList, site) {
			if err = createSiteDir(site); err != nil {
				return err
			}
			updated = true
		}
	}

	if updated {
		return ErrAppUpdated
	}

	return nil
}

func createTaskTemplate(site string, new bool) error {
	files, err := os.ReadDir(fmt.Sprintf("./%s", site))
	if err != nil {
		return err
	}

	fileList := []string{}
	for _, fileInfo := range files {
		fileList = append(fileList, fileInfo.Name())
	}

	if slices.Contains(fileList, "tasks.csv") && new {
		return ErrTasksFileExists
	}

	if !slices.Contains(fileList, "tasks.csv") {
		err = os.WriteFile(fmt.Sprintf("./%s/tasks.csv", site), Configuration[site].tasksHeaders, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadCsv(site string, file string) ([][]string, error) {
	f, err := os.Open(fmt.Sprintf("./%s/%s.csv", site, file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	tasks, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func getTransport(s string, p Proxy) (*http.Transport, error) {
	tr := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			User:   url.UserPassword(p.username, p.password),
			Host:   fmt.Sprintf("%s:%s", p.host, p.port),
		}),
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			tcpConn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
			config := tls.Config{ServerName: s}
			tlsConn := tls.UClient(tcpConn, &config, tls.HelloFirefox_Auto)

			err = tlsConn.Handshake()
			if err != nil {
				fmt.Printf("uTlsConn.Handshake() error: %v", err)
			}

			return tlsConn, nil
		},
	}
	return tr, nil
}

func getClient(s string, p Proxy) (*http.Client, error) {
	tr, err := getTransport(s, p)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Jar:       jar,
		Transport: tr,
	}, err
}

func getCapKey(url string, key string, p Proxy) (string, error) {
	cap_values := map[string]any{
		"clientKey": "70215859748a404f4978c6a0f54c11e7",
		"task": map[string]any{
			"type":          "NoCaptchaTask",
			"websiteURL":    "https://uk.topps.com/customer/account/login",
			"websiteKey":    key,
			"proxyType":     "http",
			"proxyAddress":  p.host,
			"proxyPort":     p.port,
			"proxyLogin":    p.username,
			"proxyPassword": p.password,
			"userAgent":     "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:70.0) Gecko/20100101 Firefox/70.0",
		},
	}
	json_data, err := json.Marshal(cap_values)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("https://api.capmonster.cloud/createTask", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return "", err
	}

	var res map[string]any
	d := json.NewDecoder(resp.Body)
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		fmt.Print(err)
	}

	cap_values = map[string]any{
		"clientKey": "70215859748a404f4978c6a0f54c11e7",
		"taskId":    res["taskId"],
	}
	json_data, err = json.Marshal(cap_values)
	if err != nil {
		return "", err
	}

	var recap string
	for {
		resp, err = http.Post("https://api.capmonster.cloud/getTaskResult", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			return "", err
		}

		d = json.NewDecoder(resp.Body)
		d.UseNumber()
		err = d.Decode(&res)
		if err != nil {
			return "", err
		}
		fmt.Printf("\nerrorId: %v", res["errorId"])
		error_id, err := res["errorId"].(json.Number).Int64()
		if err != nil {
			return "", err
		}
		if error_id != 0 {
			panic("error with cap")
		}
		if res["status"] == "ready" {
			solution, ok := res["solution"].(map[string]any)
			if ok {
				recap = solution["gRecaptchaResponse"].(string)
			} else {
				fmt.Print("failed to map solution")
			}
			break
		}
		fmt.Print("checking task again...")
		time.Sleep(500 * time.Millisecond)
	}
	return recap, err
}
