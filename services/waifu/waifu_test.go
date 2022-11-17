package waifu

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
)

func TestGetPixivMoe(t *testing.T) {
	pixivmoe := &PixelMoe{}
	ids, err := pixivmoe.getPixivIdsByKeyword("草神", 0, 5, false)
	if err != nil {
		t.Skip(err)
	}
	for _, id := range ids {
		t.Logf("https://pixiv.net/artworks/%d", id)
		data, err := getIllust(id)
		if err != nil {
			t.Log(err)
			continue
		}
		t.Logf("title: %s, tags: %s, url: %s", data.Title, strings.Join(pixivmoe.toArr(data.Tags), ", "), data.Images.Original)
	}
}

func TestGetPixivIcon(t *testing.T) {
	url := "https://i.pximg.net/user-profile/img/2022/09/26/02/35/44/23383020_ad04155d3b239285249e6d0837123609_50.jpg"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Referer", "https://pixiv.net")
	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Skip(err)
	}

	if res.StatusCode == 403 {
		t.Fatal(err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Skip(err)
	}

	t.Log(string(b))
}

func TestGetLolicron(t *testing.T) {

	loli := &Lolicron{}

	data, err := loli.GetImages(NewOptions(
		WithKeyword("草神"),
		WithAmount(5),
		WithR18(false),
	))

	if err != nil {

		if e, ok := err.(*request.HttpError); ok {
			defer e.Response.Body.Close()
			if b, err := io.ReadAll(e.Response.Body); err == nil {
				t.Log(string(b))
			}
		}

		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("%+v\n", d)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
}
