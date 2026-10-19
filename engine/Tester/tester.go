package tester

import (
	"bytes"
	"database/sql"
	"engine/Actions"
	sequences "engine/Sequences"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

)

type Tester struct {
	Client         http.Client
	Prev_res       ActionResponse
	Sequence_state SequnceState
}

func CreateTesterObj() Tester {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return Tester{
		Client:   http.Client{Jar: jar},
		Prev_res: ActionResponse{},
		Sequence_state: SequnceState{
			State:      make(map[string][]ActionResponse),
			Order:      make([]string, 0),
			CustomData: make(map[string]any),
		},
	}
}

func (tester *Tester) SetUpTester() {
	*tester = CreateTesterObj()
}

func (state SequnceState) ToJSON() map[string]any {
	return map[string]any{
		"state":       state.State,
		"order":       state.Order,
		"custom_data": state.CustomData,
	}
}

func (tester *Tester) ShowState() map[string]any {
	return map[string]any{
		"previous_response": tester.Prev_res,
		"sequence_state":    tester.Sequence_state.ToJSON(),
	}
}

func (tester *Tester) ExecuteHTTPAction(action *Actions.HTTPAction) *ActionResponse {
	opts := action.Opts

	reqURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil
	}

	query := reqURL.Query()
	for _, item := range opts.Query {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			query.Set(parts[0], parts[1])
		}
	}
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest(opts.Method, reqURL.String(), bytes.NewReader(opts.Body))
	if err != nil {
		return nil
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	resp, err := tester.Client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	response := &ActionResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		Data:       make(map[string]any),
	}

	tester.Prev_res = *response

	return response
}

func (tester *Tester) ExecuteSequence(db *sql.DB, sequence sequences.Sequence) map[string]any {
	tester.SetUpTester()

	for _, actionID := range sequence.ActionIDs {
		httpAction, err := Actions.GetHTTPAction(db, actionID, sequence.AppID)
		if err == nil {
			actionResponse := tester.ExecuteHTTPAction(&httpAction, )
			if actionResponse == nil {
				continue
			}

			actionName := httpAction.ActionName
			tester.Sequence_state.State[actionName] = append(tester.Sequence_state.State[actionName], *actionResponse)
			tester.Sequence_state.Order = append(tester.Sequence_state.Order, actionName)
			continue
		}

		wsActionExists, err := Actions.WSActionExistsByID(db, sequence.AppID, actionID)
		if err != nil || !wsActionExists {
			continue
		}

		wsAction, err := Actions.GetWSAction(db, actionID, sequence.AppID)
		if err != nil {
			continue
		}

		// Execute the WebSocket action using wsAction.
		// Record its response in Sequence_state.State and Sequence_state.Order.
		_ = wsAction
	}

	return tester.Sequence_state.ToJSON()
}