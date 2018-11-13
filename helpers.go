package main
import "io/ioutil"
func fileGetContents(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}
