package log

import (
	"go.uber.org/zap"
)

func Domain(domain string) zap.Field {
	return zap.String("domain", domain)
}
