package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	"github.com/sirupsen/logrus"
)

var picSizes = []string{"28x28", "34x34", "42x42", "50x50", "56x56", "68x68", "75x75", "84x84", "100x100", "200x200"}

type UserOauthData struct {
	OauthID         string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"default_email"`
	IsAvatarEmpty   bool   `json:"is_avatar_empty"`
	DefaultAvatarID string `json:"default_avatar_id"`
}

func ParseOauthToken(token string) (entity.User, error) {
	const op = "authUseCase.ParseOauthToken"
	url := url.URL{
		Scheme:   "https",
		Host:     "login.yandex.ru",
		Path:     "info",
		RawQuery: url.Values{"format": {"json"}}.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, op, "failed to create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, op, "failed to make request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, op, "failed to read response")
	}

	var userOauthData UserOauthData
	err = json.Unmarshal(data, &userOauthData)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, op, "failed to decode response")
	}

	return userOauthDataToUser(userOauthData), nil
}

func userOauthDataToUser(userOauthData UserOauthData) entity.User {
	var avatar string
	if !userOauthData.IsAvatarEmpty {
		for _, size := range picSizes {
			avatarURL := fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/%s", userOauthData.DefaultAvatarID, size)

			resp, err := http.Head(avatarURL)
			if err != nil {
				logrus.WithError(err).Warnf("Failed to fetch avatar for size %s", size)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				avatar = avatarURL
				break
			}
		}
	}

	return entity.User{
		OauthID:   userOauthData.OauthID,
		FirstName: userOauthData.FirstName,
		LastName:  userOauthData.LastName,
		Email:     userOauthData.Email,
		Avatar:    avatar,
	}
}
