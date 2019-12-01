package stats

import "io/ioutil"

func openFile(path string) (string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
