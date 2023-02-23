package service

import "github.com/anyswap/FastMulThreshold-DSA/log"

type ServiceMiddleWare struct{}

func (service *ServiceMiddleWare) GetGroupId(threshold string, userAccountsAndIpPortAddr []string) map[string]interface{} {
	if data, err := getGroupId(threshold, userAccountsAndIpPortAddr); err != nil {
		log.Error("getGroupId", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetGroupIdAndEnodesByRawData(raw string) map[string]interface{} {
	if data, err := getGroupIdAndEnodesByRawData(raw); err != nil {
		log.Error("GetGroupIdAndEnodesByRawData", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) KeyGen(rsv string, msg string) map[string]interface{} {
	if data, err := doKeyGen(rsv, msg); err != nil {
		log.Error("KeyGen", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}
