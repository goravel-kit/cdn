package contracts

import "github.com/goravel/framework/support/carbon"

type Cdn interface {
	// RefreshUrl 通过URL刷新缓存
	RefreshUrl(urls []string) error
	// RefreshPath 通过路径刷新缓存
	RefreshPath(paths []string) error
	// GetUsage 获取域名请求量
	GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error)
}
