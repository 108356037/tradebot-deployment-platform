package cli

import (
	"strings"
)

type CmdResult struct {
	Stdout string
	Stderr string
}

func ParseCreate(res *CmdResult) map[string]interface{} {
	//t := make(map[int]string)
	if res.Stdout != "" {
		resArr := strings.Split(res.Stdout, "\n")
		relName := strings.Split(resArr[0], ": ")[1]
		relStatus := strings.Split(resArr[3], ": ")[1]
		return map[string]interface{}{
			"releaseName":   relName,
			"releaseStatus": relStatus,
		}
	} else {
		errStr := strings.Trim(res.Stderr, "\n")
		errStr = strings.ReplaceAll(errStr, "\"", "")
		return map[string]interface{}{
			"error": errStr,
		}
	}
}

func ParseGet(res *CmdResult) map[string]interface{} {
	if res.Stdout != "" {

		resJson := strings.Split(res.Stdout, "\n")[0]
		resJson = strings.ReplaceAll(resJson, "\"", "")
		resJson = strings.Trim(resJson, "[]")
		releaseNames := strings.Split(resJson, ",")

		var resArr []string
		resArr = append(resArr, releaseNames...)

		if resArr[0] == "" {
			return map[string]interface{}{
				"result": "",
			}
		}

		return map[string]interface{}{
			"result": resArr,
		}
	} else {
		//resJson := strings.Split(res.Stderr, "\n")[1]
		//resJson = strings.ReplaceAll(resJson, "\"", "")
		errStr := strings.Trim(res.Stderr, "\n")
		errStr = strings.ReplaceAll(errStr, "\"", "")
		return map[string]interface{}{
			"result": errStr,
		}
	}
}

func ParseDelete(res *CmdResult) map[string]interface{} {
	if res.Stdout != "" {
		resStr := strings.Split(res.Stdout, "\n")[0]
		resStr = strings.ReplaceAll(resStr, "\"", "")
		return map[string]interface{}{
			"result": resStr,
		}
	} else {
		errStr := strings.Split(res.Stderr, "\n")[0]
		errStr = strings.ReplaceAll(errStr, "\"", "")
		return map[string]interface{}{
			"result": errStr,
		}
	}
}
