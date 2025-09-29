package pkg

import (
	"context"

	"github.com/stephenafamo/bob/types/pgtypes"
)

type clientIpKeyType struct{}

var clientIpKey = clientIpKeyType{}

func SetCtxClientIp(c context.Context, ip string) context.Context {
	inet := pgtypes.Inet{}
	inet.Scan(ip)

	ctx := context.WithValue(c, clientIpKey, inet)

	return ctx
}

func GetCtxClientIp(c context.Context) (pgtypes.Inet, bool) {
	val := c.Value(clientIpKey)

	ip, ok := val.(pgtypes.Inet)

	return ip, ok
}
