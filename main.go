package main

import (
	"gl-source/database"
	"gl-source/sources/gov"
)

func init() {
	database.InitDB()
}

func main() {
	gov.ProcessDownload()
}
