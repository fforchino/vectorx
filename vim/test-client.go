package main

import (
	"vectorx/pkg/vim-client"
)

func main() {
	vim_client.SendMessageAndGo("localhost:8070", "Tom", "005070ac", "Human", "00000000", "Hello")
}
