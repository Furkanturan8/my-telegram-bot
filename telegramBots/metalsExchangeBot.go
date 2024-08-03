package telegramBots

type MetalInfo struct {
	Name    string `json:"name"`
	Buying  string `json:"Buying"`
	Type    string `json:"Type"`
	Selling string `json:"Selling"`
	Change  string `json:"Change"`
}

type MetalsExchangeResp map[string]MetalInfo

func GetMetalsExchange() (string, error) {

	return "", nil
}
