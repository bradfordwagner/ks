package list

import (
	"io/ioutil"
)

// Kubeconfigs returns a list of kubeconfig files in the given directory
func Kubeconfigs(dir string) (configs []string, err error) {
	// list files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			configs = append(configs, file.Name())
		}
	}

	return
}
