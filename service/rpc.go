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

func (service *ServiceMiddleWare) GetGroupIdByRawData(raw string) map[string]interface{} {
	if data, err := getGroupIdByRawData(raw); err != nil {
		log.Error("GetGroupIdByRawData", "error", err.Error())
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

func (service *ServiceMiddleWare) KeyGenByRawData(raw string) map[string]interface{} {
	if data, err := doKeyGenByRawData(raw); err != nil {
		log.Error("KeyGenByRawData", "error", err.Error())
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

func (service *ServiceMiddleWare) GetReqAddrStatus(keyId string) map[string]interface{} {
	if data, err := getReqAddrStatus(keyId); err != nil {
		log.Error("getReqAddrStatus", "error", err.Error())
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

func (service *ServiceMiddleWare) GetAccountList(userAccount string) map[string]interface{} {
	if data, err := getAccountList(userAccount); err != nil {
		log.Error("getAccountList", "error", err.Error())
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
