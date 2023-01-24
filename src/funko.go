package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Funko struct {
	BaseModule
	Tasks []FunkoTask
}

type FunkoTask struct {
	Task
}

func (m *Funko) load() error {
	if err := m.loadProxies(); err != nil {
		return err
	}
	if err := m.loadTasks(); err != nil {
		return err
	}
	return nil
}

func (m *Funko) loadTasks() error {
	csvData, err := loadCsv("funko", "tasks")
	if err != nil {
		return err
	}

	for i, line := range csvData {
		if i > 0 {
			t := FunkoTask{}
			for j, field := range line {
				switch j {
				case 0:
					t.item = field
				case 1:
					t.qty, err = strconv.Atoi(field)
					if err != nil {
						return err
					}
				case 2:
					t.firstname = field
				case 3:
					t.lastname = field
				case 4:
					t.address1 = field
				case 5:
					t.address2 = field
				case 6:
					t.city = field
				case 7:
					t.postcode = field
				case 8:
					t.username = field
				case 9:
					t.password = field
				case 10:
					t.guest, err = strconv.ParseBool(field)
					if err != nil {
						return err
					}
				}
			}
			m.Tasks = append(m.Tasks, t)
		}
	}
	return nil
}

func (m *Funko) run() error {
	if err := m.initBrowser(); err != nil {
		return err
	}
	// if m.Monitor {
	// 	go m.monitor()
	// }
	var wg sync.WaitGroup
	for i, t := range m.Tasks {
		wg.Add(1)
		go func(t FunkoTask, i int) {
			defer wg.Done()
			_ = m.runTask(t, i)
		}(t, i)
	}
	wg.Wait()
	return nil
}

//	func (m *Funko) monitor() {
//		time.Sleep(time.Second * 10)
//		m.Live = true
//	}

func (t *FunkoTask) login(p playwright.Page) error {
	if _, err := p.Goto("https://funkoeurope.com/", playwright.PageGotoOptions{Timeout: playwright.Float(60000)}); err != nil {
		return err
	}

	if _, err := p.Goto("https://funkoeurope.com/account/", playwright.PageGotoOptions{Timeout: playwright.Float(60000)}); err != nil {
		return err
	}

	user, err := p.Locator(`#customer\[email\]`)
	if err != nil {
		return err
	}
	user.Type(t.username)

	pass, err := p.Locator(`#customer\[password\]`)
	if err != nil {
		return err
	}
	pass.Type(t.password)

	time.Sleep(time.Second * 2)

	l, err := p.Locator(`#customer_login > button`)
	if err != nil {
		return err
	}
	l.Click()
	p.WaitForNavigation()
	// l.Click()

	return nil
}

func (t *FunkoTask) cart(p playwright.Page) error {
	if _, err := p.Goto(fmt.Sprintf("https://funkoeurope.com/cart/%s:%v", t.item, t.qty), playwright.PageGotoOptions{Timeout: playwright.Float(60000)}); err != nil {
		return err
	}
	return nil
}

func (t *FunkoTask) checkout(p playwright.Page) error {
	if _, err := p.Goto("https://funkoeurope.com/pages/checkout", playwright.PageGotoOptions{Timeout: playwright.Float(60000)}); err != nil {
		return err
	}
	if err := p.Pause(); err != nil {
		return err
	}

	return nil
}

func (m *Funko) runTask(t FunkoTask, i int) error {
	page, err := m.Browser.NewPage(playwright.BrowserNewContextOptions{
		Proxy: &playwright.BrowserNewContextOptionsProxy{
			Server:   playwright.String(fmt.Sprintf("http://%s:%s", m.Proxies[i].host, m.Proxies[i].port)),
			Username: playwright.String(fmt.Sprint(m.Proxies[i].username)),
			Password: playwright.String(fmt.Sprint(m.Proxies[i].password)),
		},
	})
	if err != nil {
		return err
	}

	if err = t.login(page); err != nil {
		return err
	}

	if err = t.cart(page); err != nil {
		return err
	}

	if err = t.checkout(page); err != nil {
		return err
	}

	return nil
}
