package main

import "go.uber.org/zap"

func logDomain(domain string) zap.Field {
	return zap.String("domain", domain)
}
