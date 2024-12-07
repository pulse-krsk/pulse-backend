package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
)

type UserOauthData struct {
	OauthID   string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"default_email"`
	Fuid      string `json:"fuid"`
}

func ParseOauthToken(token string) (entity.User, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     "login.yandex.ru",
		Path:     "info",
		RawQuery: url.Values{"format": {"json"}}.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, "authUseCase.ParseOauthToken", "failed to create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, "authUseCase.ParseOauthToken", "failed to make request")
	}
	defer resp.Body.Close()

	var userOauthData UserOauthData
	err = json.NewDecoder(resp.Body).Decode(&userOauthData)
	if err != nil {
		return entity.User{}, cuserr.SystemError(err, "authUseCase.ParseOauthToken", "failed to decode response")
	}

	return userOauthDataToUser(userOauthData), nil
}

func userOauthDataToUser(userOauthData UserOauthData) entity.User {
	return entity.User{
		OauthID:   userOauthData.OauthID,
		FirstName: userOauthData.FirstName,
		LastName:  userOauthData.LastName,
		Email:     userOauthData.Email,
	}
}
