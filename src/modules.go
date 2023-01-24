package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type Task struct {
	item       string
	qty        int
	firstname  string
	lastname   string
	address1   string
	address2   string
	city       string
	postcode   string
	username   string
	password   string
	guest      bool
	status     string
	carted     int
	checkedout bool
	failed     bool
}

type Proxy struct {
	host     string
	port     string
	username string
	password string
}

type BaseModule struct {
	Site    string
	Proxies []Proxy
	Monitor bool
	Live    bool
	Browser playwright.Browser
}

type SiteModule interface {
	load() error
	loadTasks() error
	loadProxies() error
	run() error
}

func (b *BaseModule) loadProxies() error {
	f, err := os.Open(fmt.Sprintf("./%s/%s.txt", b.Site, "proxies"))
	if err != nil {
		return err
	}

	// r, _ := regexp.Compile(`(\S*):(\S*)@(\S*):(\S*)`)
	proxyData := bufio.NewScanner(f)
	for proxyData.Scan() {
		l := proxyData.Text()
		if l != "" {
			res := strings.Split(l, ":")
			p := Proxy{
				username: res[2],
				password: res[3],
				host:     res[0],
				port:     res[1],
			}
			b.Proxies = append(b.Proxies, p)
		}
	}
	return nil
}

func (b *BaseModule) initBrowser() error {
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Proxy: &playwright.BrowserTypeLaunchOptionsProxy{
			Server: playwright.String("http://myproxy.local"),
		},
	}

	pw, err := playwright.Run()
	if err != nil {
		return err
	}

	b.Browser, err = pw.Chromium.Launch(launchOptions)
	if err != nil {
		return err
	}

	return nil
}

func initModule(s string) (SiteModule, error) {
	switch s {
	case "funko":
		t := &Funko{}
		t.Site = "funko"
		if err := t.load(); err != nil {
			return nil, err
		}
		return t, nil
	default:
		return nil, ErrInvalidSite
	}
}
