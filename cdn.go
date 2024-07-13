package cdn

import (
	"strings"

	"github.com/goravel-kit/cdn/contracts"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/support/carbon"
)

type Cdn struct {
	drivers []contracts.Cdn
}

func NewCdn(config config.Config) *Cdn {
	name := config.GetString("cdn.driver")
	names := strings.Split(name, ",")
	var drivers []contracts.Cdn
	for _, driver := range names {
		switch driver {
		case "baishan":
			drivers = append(drivers, &BaiShan{
				Token: config.GetString("cdn.baishan.token"),
			})
		case "cloudflare":
			drivers = append(drivers, &CloudFlare{
				Key:    config.GetString("cdn.cloudflare.key"),
				Email:  config.GetString("cdn.cloudflare.email"),
				ZoneID: config.GetString("cdn.cloudflare.zone_id"),
			})
		case "huawei":
			drivers = append(drivers, &HuaWei{
				AccessKey: config.GetString("cdn.huawei.access_key"),
				SecretKey: config.GetString("cdn.huawei.secret_key"),
			})
		case "kuocai":
			drivers = append(drivers, &KuoCai{
				UserName: config.GetString("cdn.kuocai.username"),
				PassWord: config.GetString("cdn.kuocai.password"),
			})
		}
	}

	return &Cdn{
		drivers: drivers,
	}
}

func (c *Cdn) RefreshUrl(urls []string) error {
	for _, driver := range c.drivers {
		err := driver.RefreshUrl(urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cdn) RefreshPath(paths []string) error {
	for _, driver := range c.drivers {
		err := driver.RefreshPath(paths)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cdn) GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	var total uint
	for _, driver := range c.drivers {
		usage, err := driver.GetUsage(domain, startTime, endTime)
		if err != nil {
			return 0, err
		}
		total += usage
	}

	return total, nil
}
