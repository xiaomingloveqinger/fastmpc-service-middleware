/*
 *  Copyright (C) 2020-2021  AnySwap Ltd. All rights reserved.
 *  Copyright (C) 2020-2021  huangweijun@anyswap.exchange
 *
 *  This library is free software; you can redistribute it and/or
 *  modify it under the Apache License, Version 2.0.
 *
 *  This library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"github.com/anyswap/FastMulThreshold-DSA/log"
)

var (
	datadir string
)

// MakeName creates a node name that follows the ethereum convention
// for such names. It adds the operation system name and Go runtime version
// the name.
func MakeName(name, version string) string {
	return fmt.Sprintf("%s/v%s/%s/%s", name, version, runtime.GOOS, runtime.Version())
}

// FileExist checks if a file exists at filePath.
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// AbsolutePath returns datadir + filename, or filename if it is absolute.
func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}

// InitDir init data dir
func InitDir(dir string) {
	if dir == "" {
		datadir = DefaultDataDir()
		log.Info("==== InitDir() ====","datadir",datadir)
		return
	}
	if filepath.IsAbs(dir) {
		datadir = dir
	} else {
		pwdDir, _ := os.Getwd()
		datadir = filepath.Join(pwdDir, dir)
		if FileExist(datadir) != true {
		    err := os.Mkdir(datadir, os.ModePerm)
		    if err != nil {
			fmt.Printf("==== InitDir(), mk dir fail ====, datadir: %v,err: %v\n", datadir,err)
			return
		    }
		}
	}
	log.Info("==== InitDir() ====","datadir",datadir)
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	if datadir != "" {
		return datadir
	}
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "fastMPC")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "fastMPC")
		} else {
			return filepath.Join(home, ".fastMPC")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// HomeDir get home dir
func HomeDir() string {
	return homeDir()
}

// homeDir xxx
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
