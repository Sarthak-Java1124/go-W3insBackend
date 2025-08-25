package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func UploadJSONToPinata(metadata map[string]interface{}) (string, error) {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(metadata); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinJSONToIPFS", payloadBuf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiJiNDQyNGUxNC00NjczLTRiNjMtYjJjMy1kMzg2YTJhMGU3MTYiLCJlbWFpbCI6InNhcnRoYWtoYXJzaDI0QGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwaW5fcG9saWN5Ijp7InJlZ2lvbnMiOlt7ImRlc2lyZWRSZXBsaWNhdGlvbkNvdW50IjoxLCJpZCI6IkZSQTEifSx7ImRlc2lyZWRSZXBsaWNhdGlvbkNvdW50IjoxLCJpZCI6Ik5ZQzEifV0sInZlcnNpb24iOjF9LCJtZmFfZW5hYmxlZCI6ZmFsc2UsInN0YXR1cyI6IkFDVElWRSJ9LCJhdXRoZW50aWNhdGlvblR5cGUiOiJzY29wZWRLZXkiLCJzY29wZWRLZXlLZXkiOiJjZTQ2ODgxZDI1ZmE4ZTY0YmM3YiIsInNjb3BlZEtleVNlY3JldCI6ImIxYTgxMjNhZjM5M2VjNjA2YjhjNWZhOWMzMTRlY2NmNzI5ODFkNTZiNTY1YTA0YTBhN2I1ZTA2ZTM0Y2JmOTUiLCJleHAiOjE3ODU5NTg3ODF9.VA9vH0FPWSsFKfM-NmU_kC6ICAycKrae8NJ3PG10bzg")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if hash, ok := result["IpfsHash"].(string); ok {
		return "https://gateway.pinata.cloud/ipfs/" + hash, nil
	}

	return "", fmt.Errorf("IpfsHash not found in response")
}
