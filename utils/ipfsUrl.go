package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadToPinata(fileheader *multipart.FileHeader) (string, error) {
	file, err := fileheader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileheader.Filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiJiNDQyNGUxNC00NjczLTRiNjMtYjJjMy1kMzg2YTJhMGU3MTYiLCJlbWFpbCI6InNhcnRoYWtoYXJzaDI0QGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwaW5fcG9saWN5Ijp7InJlZ2lvbnMiOlt7ImRlc2lyZWRSZXBsaWNhdGlvbkNvdW50IjoxLCJpZCI6IkZSQTEifSx7ImRlc2lyZWRSZXBsaWNhdGlvbkNvdW50IjoxLCJpZCI6Ik5ZQzEifV0sInZlcnNpb24iOjF9LCJtZmFfZW5hYmxlZCI6ZmFsc2UsInN0YXR1cyI6IkFDVElWRSJ9LCJhdXRoZW50aWNhdGlvblR5cGUiOiJzY29wZWRLZXkiLCJzY29wZWRLZXlLZXkiOiJjZTQ2ODgxZDI1ZmE4ZTY0YmM3YiIsInNjb3BlZEtleVNlY3JldCI6ImIxYTgxMjNhZjM5M2VjNjA2YjhjNWZhOWMzMTRlY2NmNzI5ODFkNTZiNTY1YTA0YTBhN2I1ZTA2ZTM0Y2JmOTUiLCJleHAiOjE3ODU5NTg3ODF9.VA9vH0FPWSsFKfM-NmU_kC6ICAycKrae8NJ3PG10bzg")
	req.Header.Set("Content-type", writer.FormDataContentType())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("The error in sending client request is : ", err)
	}
	defer resp.Body.Close()

	// Decode response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	fmt.Println("Pinata Response:", result)
	// Return IPFS hash or full gateway URL
	if hash, ok := result["IpfsHash"].(string); ok {
		return "https://gateway.pinata.cloud/ipfs/" + hash, nil
	}

	return "", fmt.Errorf("IpfsHash not found in response")
}
