package flags

import (
	"flag"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/conts"
	"os"
)

// service config
var (
	SERVICE_CONFIG_DIR string // default: home/.cache/golang-language-server
	SERVICE_DEBUG      bool
	SERVICE_STDIO      bool
	SERVICE_TCP        bool
	SERVICE_PROT       int
)

func init() {
	flag.StringVar(&SERVICE_CONFIG_DIR, "config_dir", GetHome()+"/.cache/golang-language-server", "配置文件目录")
	flag.BoolVar(&SERVICE_DEBUG, "debug", false, "调试")
	flag.BoolVar(&SERVICE_STDIO, "stdio", false, "标准输出")
	flag.BoolVar(&SERVICE_TCP, "tcp", false, "rpc连接方式")
	flag.IntVar(&SERVICE_PROT, "port", 9999, "端口")
	flag.Usage = Help
	flag.Parse()

	if _, err := os.Stat(SERVICE_CONFIG_DIR); os.IsNotExist(err) {
		err = os.Mkdir(SERVICE_CONFIG_DIR, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func Help() {
	fmt.Fprintln(os.Stderr, "This is a Go language server service. You can use the following flags to configure it:\n")
	fmt.Fprintf(os.Stderr, "Usage of %s Version:%s\n", conts.SERVICE_NAME, conts.VERSION)
	fmt.Fprintln(os.Stderr, "Available flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintf(os.Stderr, "  %s -port 8080 -config_dir /path/to/config -debug\n", conts.SERVICE_NAME)
}
