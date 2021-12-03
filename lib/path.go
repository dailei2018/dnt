package lib

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func Mk_dir_if_not(path string) {
	s, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			logrus.Fatalf("mkdir %s failed:%v\n", path, err)
		}
		fmt.Printf("mkdir %s\n", path)
	} else {
		if !s.IsDir() {
			logrus.Fatalf("%s exits not dir\n", path)
		}
	}

}
