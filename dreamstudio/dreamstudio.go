package dreamstudio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type TextToImageImage struct {
	Base64       string `json:"base64"`
	Seed         uint32 `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type TextToImageResponse struct {
	Images []TextToImageImage `json:"artifacts"`
}

// DreamStdioRequestBody 请求体
type DreamStdioRequestBody struct {
	TextPrompts        []TextPrompt `json:"text_prompts"`
	CfgScale           int          `json:"cfg_scale"`
	ClipGuidancePreset string       `json:"clip_guidance_preset"`
	Height             int          `json:"height"`
	Width              int          `json:"width"`
	Samples            int          `json:"samples"`
	Steps              int          `json:"steps"`
}

type TextPrompt struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

func TextToImage(msg string) (string, error) {
	// Build REST endpoint URL w/ specified engine
	engineId := "stable-diffusion-v1-5"
	apiHost, hasApiHost := os.LookupEnv("API_HOST")
	if !hasApiHost {
		apiHost = "https://api.stability.ai"
	}
	reqUrl := apiHost + "/v1/generation/" + engineId + "/text-to-image"

	// Acquire an API key from the environment
	//apiKey, hasAPIKey := os.LookupEnv("STABILITY_API_KEY")
	apiKey := "youKey"
	// if !hasAPIKey {
	// 	panic("Missing STABILITY_API_KEY environment variable")
	// }
	textPrompts := []TextPrompt{
		{
			Text:   msg,
			Weight: 1,
		},
	}
	requestBody := DreamStdioRequestBody{
		TextPrompts:        textPrompts,
		CfgScale:           7,
		ClipGuidancePreset: "FAST_BLUE",
		Height:             512,
		Width:              512,
		Samples:            1,
		Steps:              30,
	}

	requestData, _ := json.Marshal(requestBody)
	// if err != nil {
	// 	return nil, fmt.Errorf("json.Marshal requestBody error: %v", err)
	// }

	//log.Printf("dreamstdio request(%d) json: %s\n", runtimes, string(requestData))
	log.Printf("dreamstdio request(%d) json: %s\n", 1, string(requestData))

	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// Execute the request & read all the bytes of the body
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			panic(err)
		}
		log.Printf("Non-200 response: %s", body)
	}

	// Decode the JSON body
	var body TextToImageResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		panic(err)
	}

	// Write the images to disk
	for i, image := range body.Images {
		outFile := fmt.Sprintf("./dreamstudio/v1_txt2img_%d.png", i)
		file, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}

		imageBytes, err := base64.StdEncoding.DecodeString(image.Base64)
		if err != nil {
			panic(err)
		}

		if _, err := file.Write(imageBytes); err != nil {
			panic(err)
		}

		if err := file.Close(); err != nil {
			panic(err)
		}
	}
	return "./dreamstudio/v1_txt2img_0.png", nil
}
