package crow

import (
	"fmt"
	random "math/rand"
	"os"
	"strings"
	"time"
)

func (c *Core) GenPassword() string {
	config := c.Config.RandomPassword

	if !config.Enabled {
		data, err := os.ReadFile(c.Config.Bruteforce.DictionaryPath)
		if err != nil {
			fmt.Println("Can't find: ", c.Config.Bruteforce.DictionaryPath)
			os.Exit(0)
		}

		strs := strings.Split(string(data), "\n")
		cStrs := len(strs)

		s1 := random.NewSource(time.Now().UnixNano())
		r1 := random.New(s1)

		pass := strs[r1.Intn(cStrs)]
		return pass
	}

	return genRand(config.RandomCustomPassword, config.RandomPasswordLen)
}

func genRand(ltrs string, length int) string {
	var letters = []rune(ltrs)

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[random.Intn(len(letters))]
	}
	return string(s)
}
