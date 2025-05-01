package flags

import (
	"flag"
	"fmt"
	"os"
)

// service config
var (
	SERVICE_PROT       int
	SERVICE_CONFIG_DIR string // default: home/.cache/golang-language-server
	SERVICE_DEBUG      bool

	VERSION      = "0.0.1"
	SERVICE_NAME = "golang-language-server"
)

func init() {
	flag.IntVar(&SERVICE_PROT, "port", 9999, "")
	flag.StringVar(&SERVICE_CONFIG_DIR, "config_dir", GetHome()+"/.cache/golang-language-server", "config dir")
	flag.BoolVar(&SERVICE_DEBUG, "debug", false, "")
	flag.Usage = Help
	flag.Parse()
}

func Help() {
	fmt.Fprintf(os.Stderr, "Usage of %s Version:%s\n", SERVICE_NAME, VERSION)
	fmt.Fprintln(os.Stderr, "This is a Go language server service.\nYou can use the following flags to configure it:")
	fmt.Fprintln(os.Stderr, "Available flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintf(os.Stderr, "  %s -port 8080 -config_dir /path/to/config -debug\n", SERVICE_NAME)
}
