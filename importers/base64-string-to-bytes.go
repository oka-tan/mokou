package importers

import "encoding/base64"

func Base64StringToBytes(s *string) *[]byte {
	if s == nil || *s == "" {
		return nil
	}

	result, err := base64.StdEncoding.DecodeString(*s)

	if err != nil {
		return nil
	}

	return &result
}
