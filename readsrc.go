package main

import "io/ioutil"

func readFile(fileName string) (string, error) {
	c, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(c), nil
}
