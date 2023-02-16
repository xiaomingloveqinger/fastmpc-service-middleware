package service

type ServiceMiddleWare struct{}

func (service *ServiceMiddleWare) TestJsonrpc(rsv string, msg string) map[string]interface{} {
	if data, err := GetTestData(); err != nil {
		return map[string]interface{}{
			"Status": "Error",
			"Tip":    "something happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "SUCCESS",
			"Tip":    "I am tip",
			"Error":  "",
			"Data":   data,
		}
	}
}
