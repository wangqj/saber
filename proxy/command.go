package proxy

import "saber/proxy/redis"

//根据resp获取命令关键字
func getFlag(respArray []*redis.Resp) string {
	if len(respArray) >= 2 {
		return string(respArray[0].Value)
	}
	return ""
}
