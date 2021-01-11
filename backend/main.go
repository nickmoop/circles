package main


import (
    "flag"

    "backend/handlers"
)



func main() {
    var hostname = flag.String("address", "localhost:8080", "http service address")
    handlers.ApplyHandlers(hostname)
}
