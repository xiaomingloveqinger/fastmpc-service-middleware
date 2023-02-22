package service

type ServiceMiddleWare struct{}

func (service *ServiceMiddleWare) GetGroupIdAndEnodes(threshold string, userAccountsAndIpPortAddr []string) map[string]interface{} {
	if data, err := getGroupIdAndEnodes(threshold, userAccountsAndIpPortAddr); err != nil {
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
