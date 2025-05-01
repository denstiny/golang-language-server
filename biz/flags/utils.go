package flags

import (
	"os"
	"os/user"
	"runtime"
)

func GetHome() string {
	if runtime.GOOS == "windows" {
		// 在 Windows 系统上，使用 USERPROFILE 环境变量
		home := os.Getenv("USERPROFILE")
		if home == "" {
			panic("could not determine home directory")
		}
		return home
	}
	// 在 Unix 系统（包括 Linux 和 macOS）上，使用 os/user 包
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.HomeDir
}
