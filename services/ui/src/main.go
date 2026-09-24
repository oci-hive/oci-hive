package main

import (
	"codeberg.org/oci-hive/oci-hive/server"
)

func main() {
	server.Serve(":8080")
}
