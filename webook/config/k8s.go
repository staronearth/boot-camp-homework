//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:xayy@dev.123@tcp(k8s-mysql-svc:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "k8s-redis-svc:6380",
	},
}
