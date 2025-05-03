package flags

import (
	"flag"
	"fmt"
	"os"
)

// service config
var (
	SERVICE_CONFIG_DIR string // default: home/.cache/golang-language-server
	SERVICE_DEBUG      bool
	SERVICE_STDIO      bool
	SERVICE_TCP        bool
	SERVICE_PROT       int

	VERSION      = "0.0.1"
	SERVICE_NAME = "golang-language-server"
)

func init() {
	flag.StringVar(&SERVICE_CONFIG_DIR, "config_dir", GetHome()+"/.cache/golang-language-server", "配置文件目录")
	flag.BoolVar(&SERVICE_DEBUG, "debug", false, "调试")
	flag.BoolVar(&SERVICE_STDIO, "stdio", false, "标准输出")
	flag.BoolVar(&SERVICE_TCP, "tcp", false, "rpc连接方式")
	flag.IntVar(&SERVICE_PROT, "port", 9999, "端口")
	flag.Usage = Help
	flag.Parse()
}

func Help() {
	fmt.Fprintln(os.Stderr, "This is a Go language server service. You can use the following flags to configure it:\n")
	fmt.Fprintf(os.Stderr, "Usage of %s Version:%s\n", SERVICE_NAME, VERSION)
	fmt.Fprintln(os.Stderr, "Available flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintf(os.Stderr, "  %s -port 8080 -config_dir /path/to/config -debug\n", SERVICE_NAME)
}
