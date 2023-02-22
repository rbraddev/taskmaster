package main

type siteConfiguration struct {
	tasksHeaders []byte
	capAPIKey    string
}

var baseConfig = siteConfiguration{
	capAPIKey: "70215859748a404f4978c6a0f54c11e7",
}

var funkoConfig = siteConfiguration{
	capAPIKey:    baseConfig.capAPIKey,
	tasksHeaders: []byte("item,qty,first_name,last_name,address_1,address_2,city,postcode,username,password,guest"),
}

var toppsConfig = siteConfiguration{
	capAPIKey:    baseConfig.capAPIKey,
	tasksHeaders: []byte("item,qty,first_name,last_name,address_1,address_2,city,postcode,username,password"),
}

var Configuration = map[string]siteConfiguration{
	"funko": funkoConfig,
	"topps": toppsConfig,
}
