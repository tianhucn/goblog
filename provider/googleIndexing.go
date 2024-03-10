package provider

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"
)

func (w *Website) SaveGoogleIndexingAccount(content string) error {
	return w.SaveSettingValueRaw(GoogleIndexingJsonKey, content)
}

func (w *Website) GetGoogleIndexingAccess() (*indexing.Service, error) {
	content := w.GetSettingValue(GoogleIndexingJsonKey)

	if len(content) == 0 {
		return nil, errors.New("账号错误")
	}

	ctx := context.Background()

	client, err := indexing.NewService(ctx, option.WithCredentialsJSON([]byte(content)))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (w *Website) PushGoogleIndexing(client *indexing.Service, domain string) (int, error) {
	notification := indexing.UrlNotification{
		Type: "URL_UPDATED",
		Url:  domain,
	}
	res, err := client.UrlNotifications.Publish(&notification).Do()

	if err != nil {
		return -1, err
	}

	w.logPushResult("google", fmt.Sprintf("%v, %s", domain, res.HTTPStatusCode))

	return res.HTTPStatusCode, nil
}

func (w *Website) PushGoogle(list []string) error {
	client, err := w.GetGoogleIndexingAccess()
	if err != nil {
		return err
	}
	for _, domain := range list {
		_, err := w.PushGoogleIndexing(client, domain)
		if err != nil {
			return err
		}
	}

	return nil
}
