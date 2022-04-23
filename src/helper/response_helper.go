package helper

type ResponseHelper struct {
}

func (r *ResponseHelper) CreateResponse(result []interface{}, resErr error) (map[string]interface{}, int) {

	respMsg := make(map[string]interface{})

	if resErr != nil {

		respMsg["success"] = false
		respMsg["data"] = resErr.Error()
		respMsg["message"] = "SQL ERROR"
		respMsg["code"] = 500

		return respMsg, 500
	}

	respMsg["success"] = true
	respMsg["data"] = result
	respMsg["message"] = "Fetching News Data"
	respMsg["code"] = 200

	return respMsg, 200
}
