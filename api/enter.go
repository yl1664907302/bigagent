package api

import "bigagent/api/data"

type ApiGroup struct {
	DataApiGrup data.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
