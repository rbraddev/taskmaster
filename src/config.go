package main

type siteConfiguration struct {
	tasksHeaders []byte
}

var funkoConfig = siteConfiguration{
	tasksHeaders: []byte("item,qty,first_name,last_name,address_1,address_2,city,postcode,username,password,guest"),
}

var Configuration = map[string]siteConfiguration{
	"funko": funkoConfig,
}
