package crow

import (
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func getCpuUsage() int {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0
	}
	return int(math.Ceil(percent[0]))
}

func (c *Core) Balancer() {

	for {
		cpu := getCpuUsage()
		if cpu >= c.Config.Balancer.BalancerValue {
			go func() {
				c.BalancerChan <- true
			}()
		}
	}

}
