//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:example@tcp(localhost:3306)/webook",
		//DSN: "root:example@tcp(172.17.32.102:30708)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
		//Addr: "172.17.32.102:32381",
	},
}
