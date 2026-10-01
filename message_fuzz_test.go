package convertago_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/JSYoo5B/convertago/internal/benchmarksource"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
)

func FuzzGeneratedReflectionConversion(f *testing.F) {
	f.Add(uint8(0), "Notice", "Hello ", "Jane", "https://example.com/image.png", int64(123), true)
	f.Add(uint8(1), "", "", "", "relative", int64(0), false)
	f.Add(uint8(2), "알림", "본문 ", "사용자", "http://example.com/image.png", int64(-1), true)
	f.Add(uint8(3), "Notice", "<b>&", "\xff", "https:///image", int64(1), false)
	v := validator.New(validator.WithRequiredStructEnabled())
	kakaowork.RegisterValidation(v)
	slack.RegisterValidation(v)
	googlechat.RegisterValidation(v)
	f.Fuzz(func(t *testing.T, platform uint8, title, prefix, name, url string, count int64, present bool) {
		flat := benchmarksource.FlatMessage{Title: title, Prefix: prefix, Name: name, Count: int(count)}
		nested := benchmarksource.NestedMessage{
			Title: title, Items: []benchmarksource.Line{{Prefix: prefix, Name: name}},
		}
		if present {
			nested.Picture = &benchmarksource.Picture{URL: url, Alt: name}
			nested.Footer = &prefix
		}
		dynamic := benchmarksource.DynamicMessage{Content: &nested}
		convert := []func(any, ...convertago.Option) (any, error){
			func(input any, options ...convertago.Option) (any, error) {
				return convertago.ToKakaoworkMessage(input, options...)
			},
			func(input any, options ...convertago.Option) (any, error) {
				return convertago.ToSlackMessage(input, options...)
			},
			func(input any, options ...convertago.Option) (any, error) {
				return convertago.ToGoogleChatMessage(input, options...)
			},
		}[platform%3]
		for _, inputs := range []struct{ generated, reflected any }{
			{&flat, (*benchmarkReflectionFlatMessage)(&flat)},
			{&nested, (*benchmarkReflectionNestedMessage)(&nested)},
			{&dynamic, (*benchmarkReflectionDynamicMessage)(&dynamic)},
		} {
			got, gotErr := convert(inputs.generated, convertago.WithStrict())
			want, wantErr := convert(inputs.reflected, convertago.WithStrict())
			if !reflect.DeepEqual(gotErr, wantErr) {
				t.Fatalf("diagnostics differ: generated=%v reflected=%v", gotErr, wantErr)
			}
			if gotErr != nil {
				continue
			}
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			wantJSON, err := json.Marshal(want)
			if err != nil || !bytes.Equal(gotJSON, wantJSON) {
				t.Fatalf("JSON differs: generated=%s reflected=%s (%v)", gotJSON, wantJSON, err)
			}
			if err := v.Struct(got); err != nil {
				t.Fatalf("successful conversion failed native validation: %v", err)
			}
		}
	})
}
