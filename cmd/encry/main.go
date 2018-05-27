package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/SPO-Next/spo/src/util/encrypt"
)

func main() {
	var privateKey string
	var password string
	var method string
	var filename string
	flag.StringVar(&method, "method", "", "method encry|decry")
	flag.StringVar(&privateKey, "privatekey", "", "private secret")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&filename, "filename", "./encryped_privatekey.file", "storage file")
	flag.Parse()
	if method == "" {
		fmt.Printf("need -method\n")
		return
	}
	if method == "encry" {
		if privateKey == "" {
			fmt.Printf("need -privatekey\n")
			return
		}
	}
	if password == "" {
		fmt.Printf("need -password\n")
		return
	}

	if len(password) != 16 {
		fmt.Printf("password length must 16\n")
		return
	}

	switch method {
	case "encry":
		content, err := encrypt.Encrypt([]byte(password), privateKey)
		if err != nil {
			fmt.Printf("encrypted err %v\n", err)
			return
		}

		if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Printf("write error %v\n", err)
			return
		}

	case "decry":
		encryptMsg, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("read secret file failed, %+v", err)
			return
		}
		origin, err := encrypt.Decrypt([]byte(password), string(encryptMsg))
		if err != nil {
			fmt.Printf("decrypt error %v\n", err)
			return
		}

		fmt.Printf("origin:%s\n", origin)
	default:
		fmt.Printf("method should be encry|decry only\n")
		return
	}
}
