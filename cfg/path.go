package cfg

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	RootPath string
	CfgDir   string
	LogDir   string
	DatDir   string
	HTMLDir  string
)

func init() {
	// find "RootPath"
	initRootPath()

	// join "rootProject" and vars
	CfgDir = JoinRootPaths("cfg")
	LogDir = filepath.Join(CfgDir, "config", "logger")
	DatDir = JoinRootPaths("dat")
	HTMLDir = JoinRootPaths("server", "web", "html")
}

func JoinRootPaths(path ...string) string {
	pathLen := len(path)
	if pathLen == 0 {
		return RootPath
	}
	p := make([]string, 0, pathLen+1)
	p = append(p, RootPath)
	p = append(p, path...)

	return filepath.Join(p...)
}

func initRootPath() {

	runtimePath := getRootPathByRuntimeCaller()
	executablePath := getRootPathByExecutable()

	rootPath := executablePath
	if !isLoggerDirExist(executablePath) && isLoggerDirExist(runtimePath) {
		rootPath = runtimePath
	}

	RootPath = rootPath
}

func getRootPathByRuntimeCaller() string {
	// get current file path
	// {RootPath}\cfg\path.go
	_, p, _, _ := runtime.Caller(0)

	// find "RootPath"
	for i := 0; i < 2; i++ {
		p = filepath.Dir(p)
	}
	return p
}

func getRootPathByExecutable() string {
	// find execute path
	ex, err := os.Executable()
	if err != nil {
		return ""
	}

	// find "RootPath"
	return filepath.Dir(ex)
}

func isLoggerDirExist(rootPath string) bool {
	loggerPath := filepath.Join(rootPath, "cfg", "config", "logger")
	_, err := os.Stat(loggerPath)
	return !os.IsNotExist(err)
}
