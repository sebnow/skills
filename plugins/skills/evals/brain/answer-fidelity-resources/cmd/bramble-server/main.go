// Command bramble-server runs the bramble server: the GUI, the API and
// the endpoint runners pull tasks from.
package main

import (
	"flag"
	"log/slog"
	"os"

	"bramble/internal/server"
)

func main() {
	db := flag.String("db", "bramble.db", "SQLite database file")
	listen := flag.String("listen", "127.0.0.1:8080", "address to serve on")
	flag.Parse()
	if err := server.Run(*db, *listen); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
