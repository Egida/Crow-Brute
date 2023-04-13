package crow

import (
	"bufio"
	"os"
)

func ReadServers() ([]string, error) {

	var servers = make([]string, 0)

	file, err := os.Open("servers.txt")
	if err != nil {
		return nil, err
	}

	fscanner := bufio.NewScanner(file)

	fscanner.Split(bufio.ScanLines)

	for fscanner.Scan() {
		if len(fscanner.Text()) < 2 {
			continue
		}

		servers = append(servers, fscanner.Text())
	}
	return servers, nil

}
