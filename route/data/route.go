package data

import (
	"bigagent/api"
	"bigagent/util"
)

type DataRouter struct{}

func (*DataRouter) AllRouter() *util.ConcurrentHashMap {
	chashMap := util.CreateChashMap(10, 20)
	chashMap.Set("/", api.ApiGroupApp.DataApiGrup.ShowData)
	return chashMap
}
