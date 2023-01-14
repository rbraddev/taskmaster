package main

type siteConfiguration struct {
	tasksHeaders []byte
}

var funkoConfig = siteConfiguration{
	tasksHeaders: []byte("first_name,last_name,address_1,address_2"),
}

var Configuration = map[string]siteConfiguration{
	"funko": funkoConfig,
}
