package route

import "bigagent/route/data"

type RouterGroup struct {
	DataRouter data.DataRouter
}

var RouterGroupApp = new(RouterGroup)
