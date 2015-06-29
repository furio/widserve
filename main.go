package main

// https://github.com/spf13/viper

import (
    "github.com/furio/widserve/server"
    "github.com/furio/widserve/refresher"
)

func main() {
    // Actually to be splitted depending on config/startup flag


    server.Main()
    //
    refresher.Main()
}
