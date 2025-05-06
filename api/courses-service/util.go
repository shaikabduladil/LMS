package main

import "github.com/spf13/viper"

func GetEnv(key string) string {
	return viper.GetString(key)
}
