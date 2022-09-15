//Package asagi contains functionality for migrating Asagi data
package asagi

import "encoding/base64"

//Wraps around base64.StdEncoding.DecodeString
//to change the error handling
func base64StringToBytes(s *string) *[]byte {
	if s == nil || *s == "" {
		return nil
	}

	result, err := base64.StdEncoding.DecodeString(*s)

	if err != nil {
		return nil
	}

	return &result
}
