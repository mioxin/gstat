package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

const DATA string = `{"bin":"90032396","name":"ИП \"Аутсорсинговая компания Аксултан\"","registerDate":"","okedCode":"69202","okedName":"Деятельность в области составления счетов и бухгалтерского учета","krpCode":"105","krpName":"Малые предприятия (\u003c= 5)","krpBfCode":"","krpBfName":"","kseCode":"1122","kseName":"Национальные частные нефинансовые корпорации – ОПП","katoAddress":"ЖАМБЫЛСКАЯ ОБЛАСТЬ, ТАРАЗ Г.А., Г.ТАРАЗ","fio":"АБВА БАН РЕВНА","ip":true}`
const IIN string = "90032396"

type BodyStr struct {
	data io.Reader
}

func (bs BodyStr) Close() error {
	return nil
}
func (bs BodyStr) Read(p []byte) (int, error) {
	return bs.data.Read(p)
}

type ErrTimeout struct{}

func (e *ErrTimeout) Error() string {
	return "timeout"
}

func (e *ErrTimeout) Temporary() bool {
	return true
}

type MockHelper struct {
	param    string
	workTime int
}

func (mh MockHelper) Param(k, v string) *HttpHelper {
	return &HttpHelper{}
}

func (mh MockHelper) Get(ctx context.Context) *HttpHelperResponse {
	var (
		body   io.ReadCloser
		err    error
		code   int
		status string
	)
	switch mh.param {
	case "0":
		err = new(ErrTimeout)
	case "429":
		code = 429
		status = "429 Too Many Requests."
	case IIN:
		code = 200
		status = "200 OK."
		body = BodyStr{strings.NewReader(fmt.Sprintf(`{"success":true,"obj":%s}`, DATA))}
	}
	timer := time.NewTimer(time.Duration(mh.workTime) * time.Second)

	select {
	case <-ctx.Done(): // cancel
	case <-timer.C: //work time
	}

	return &HttpHelperResponse{
		StatusCode: code,
		Status:     status,
		err:        err,
		Body:       body,
	}
}

func TestHttpConn(t *testing.T) {
	cases := []struct {
		nameCase          string
		mHelper           *MockHelper
		cancelTime        int
		wantresp, wanterr string
	}{
		{
			nameCase:   "timeout work",
			mHelper:    &MockHelper{param: "0", workTime: 4},
			cancelTime: 5,
			wanterr:    "tamporary err",
		},
		{
			nameCase:   "timeout cancel",
			mHelper:    &MockHelper{param: "0", workTime: 4},
			cancelTime: 1,
			wanterr:    "stop by context cancel",
		},
		{
			nameCase:   "get IIN",
			mHelper:    &MockHelper{param: IIN, workTime: 1},
			cancelTime: 5,
			wantresp:   `true`,
		},
		{
			nameCase:   "get timeout status 429",
			mHelper:    &MockHelper{param: "429", workTime: 1},
			cancelTime: 5,
			wanterr:    "429 Too Many Requests.",
		},
	}
	for _, test := range cases {
		t.Run(test.nameCase, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), time.Duration(test.cancelTime)*time.Second)
			resp, err := HttpConn(ctx, *test.mHelper, IIN, test.nameCase)
			if test.wantresp != fmt.Sprintf("%v", resp.Success) {
				t.Errorf("%s:\n, want: %v\t\t%v\n get:  %v\t\t%v\n", test.nameCase, test.wantresp, test.wanterr, resp, err)
			}
		})
	}
}
