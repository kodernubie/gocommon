package conf

import (
	"os"
	"strconv"
	"strings"
)

func Str(name string, defaultVal ...string) string {

	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, ".", "_")

	ret, exist := os.LookupEnv(name)

	if !exist && len(defaultVal) > 0 {
		ret = defaultVal[0]
	}

	return ret
}

func Int(name string, defaultVal ...int) int {

	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, ".", "_")

	retStr, exist := os.LookupEnv(name)
	ret := 0

	if exist {
		ret, _ = strconv.Atoi(retStr)
	} else if len(defaultVal) > 0 {
		ret = defaultVal[0]
	}

	return ret
}

func Bool(name string, defaultVal ...bool) bool {

	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, ".", "_")

	retStr, exist := os.LookupEnv(name)
	ret := false

	if exist {
		ret = strings.ToLower(retStr) == "true"
	} else if len(defaultVal) > 0 {
		ret = defaultVal[0]
	}

	return ret
}

func Float(name string, defaultVal ...float64) float64 {

	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, ".", "_")

	retStr, exist := os.LookupEnv(name)
	var ret float64 = 0

	if exist {
		ret, _ = strconv.ParseFloat(retStr, 64)
	} else if len(defaultVal) > 0 {
		ret = defaultVal[0]
	}

	return ret
}
