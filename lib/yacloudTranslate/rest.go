package yacloudTranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// REST implementation of YaTranslate interface
type RestYaTranslate struct {
	// YaCloud folder ID
	FolderId string

	// API key for simple authentication.
	// Will you Api-Key auth scheme when specified
	ApiKey string

	// IAM token for authentication.
	// Will use Bearer auth scheme when specified and ApiKey is empty
	IAMToken string

	// default - translate.api.cloud.yandex.net
	Domain string

	// default - translate/v2
	BaseUrl string

	// logger for debugging
	Logger *log.Logger
}

type detectLanguageRequest struct {
	DetectLanguageRequest
	FolderId string `json:"folderId"`
}

type listLanguagesRequest struct {
	ListLanguagesRequest
	FolderId string `json:"folderId"`
}

type translateRequest struct {
	TranslateRequest
	FolderId string `json:"folderId"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func stringOrDefault(s1, s2 string) string {
	if len(s1) > 0 {
		return s1
	}
	return s2
}

func (s RestYaTranslate) callRestApi(method string, params any) ([]byte, error) {
	url := fmt.Sprintf("https://%s/%s/%s",
		stringOrDefault(s.Domain, "translate.api.cloud.yandex.net"),
		stringOrDefault(s.BaseUrl, "translate/v2"),
		method)

	if s.Logger != nil {
		s.Logger.Printf("yacloud translate: %s", url)
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if s.Logger != nil {
		s.Logger.Printf("yacloud translate: %s", string(body))
	}

	const maxRetries = 5
	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		req.Header.Set("content-type", "application/json")
		if len(s.ApiKey) > 0 {
			req.Header.Set("authorization", fmt.Sprintf("Api-Key %s", s.ApiKey))
		} else if len(s.IAMToken) > 0 {
			req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.IAMToken))
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if s.Logger != nil {
				s.Logger.Printf("yacloud translate: attempt %d failed: %s", attempt+1, err)
			}
			continue
		}

		d, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if s.Logger != nil {
			s.Logger.Printf("yacloud translate: %s", string(d))
		}

		if resp.StatusCode != http.StatusOK {
			var apiErr apiError
			if err = json.Unmarshal(d, &apiErr); err != nil {
				lastErr = err
				continue
			}
			lastErr = fmt.Errorf("api error %d: %s", apiErr.Code, apiErr.Message)
			continue
		}

		return d, nil
	}

	return nil, lastErr
}

func (s RestYaTranslate) DetectLanguage(req DetectLanguageRequest) (res DetectLanguageResponse, err error) {
	data, err := s.callRestApi("detect", detectLanguageRequest{
		DetectLanguageRequest: req,
		FolderId:              s.FolderId,
	})
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

func (s RestYaTranslate) ListLanguages(req ListLanguagesRequest) (res ListLanguagesResponse, err error) {
	data, err := s.callRestApi("languages", listLanguagesRequest{
		ListLanguagesRequest: req,
		FolderId:             s.FolderId,
	})
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

func (s RestYaTranslate) Translate(req TranslateRequest) (res TranslateResponse, err error) {
	data, err := s.callRestApi("translate", translateRequest{
		TranslateRequest: req,
		FolderId:         s.FolderId,
	})
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
