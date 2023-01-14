package main

import "errors"

var (
	ErrInvalidSite     = errors.New("")
	ErrUnauth          = errors.New("you are unauthorised to use this app")
	ErrAppUpdated      = errors.New("app has been updated, please run again")
	ErrSiteRequired    = errors.New("site required if creating new tasks template")
	ErrTasksFileExists = errors.New("tasks.csv already exists, please rename and try again")
)
