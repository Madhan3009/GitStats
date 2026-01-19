package main

import (
	"flag"
)

func main() {
	var email string
	var folder string
	flag.StringVar(&email, "email", "your@email.com", "Email for reference")
	flag.StringVar(&folder, "file", "", "File for reference")
	flag.Parse()
	if folder != "" {
		scan(folder)
		return
	}
	stats(email)
}
