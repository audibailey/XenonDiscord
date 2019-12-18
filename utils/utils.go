package utils

import (
      "encoding/json"
)

func PrettyPrint(v interface{}) (string, error) {
      b, err := json.MarshalIndent(v, "", "  ")
      if err == nil {
		  return string(b), nil
      }
      return "", err
}
