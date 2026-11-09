package tester

import (
	actions "Tester/backend/Actions"
	environment "Tester/backend/Environment"
	sequences "Tester/backend/Sequences"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)



type Tester struct {
	Client                  *http.Client
	SequenceState           *SequnceState
	ReqExecutionContext     map[int]RequestExecutionContext
	SequnceExecutionContext map[int][]ActionContext
	WebsocketManager	map[int64]*websocket.Conn
}



func NewTester() (*Tester, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &Tester{
		Client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		SequenceState: &SequnceState{
			State: make(map[int64][]ActionResponse),
			Order: []int64{},
		},
		ReqExecutionContext:      make(map[int]RequestExecutionContext),
		SequnceExecutionContext:  make(map[int][]ActionContext),
		WebsocketManager:    make(map[int64]*websocket.Conn),
	}, nil
}

func (t *Tester) CreateSequenceExecutionContext(db *sql.DB, sequence sequences.Sequence) (int, []int, error) {
	actionList, err := actions.GetActions(db, sequence.ActionIDs)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get sequence actions: %w", err)
	}

	var actionContexts []ActionContext

	for _, action := range actionList {
		switch action.ActionType {
		case "http":
			actionContext, err := CreateActionContext(action.ActionId, action.ActionType)
			if err != nil {
				return 0, nil, err
			}

			ctx, cancel := context.WithCancel(context.Background())

			t.ReqExecutionContext[actionContext.ReqContextId] = RequestExecutionContext{
				Ctx:    ctx,
				Cancel: cancel,
			}

			actionContexts = append(actionContexts, actionContext)

		case "websocket":
			// Handle WebSocket action

		default:
			return 0, nil, fmt.Errorf("unsupported action type: %s", action.ActionType)
		}
	}

	sequenceContextID, err := GenerateRequestID()
	if err != nil {
		return 0, nil, err
	}

	actionContextIDs := make([]int, 0, len(actionContexts))

	for _, actionContext := range actionContexts {
		actionContextIDs = append(actionContextIDs, actionContext.ReqContextId)
	}

	t.SequnceExecutionContext[sequenceContextID] = actionContexts

	return sequenceContextID, actionContextIDs, nil
}

func CreateActionContext(action_id int, action_type string) (ActionContext, error) {
	reqContextID, err := GenerateRequestID()
	if err != nil {
		return ActionContext{}, err
	}

	return ActionContext{
		ReqContextId: reqContextID,
		ActionId:     action_id,
		ActionType:   action_type,
	}, nil
}

func GenerateRequestID() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<31-1))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

func (t *Tester) ExecuteSequnce(db *sql.DB, sequence sequences.Sequence, sequence_context_id int) error {
	actionContexts, exists := t.SequnceExecutionContext[sequence_context_id]
	if !exists {
		return fmt.Errorf("sequence execution context %d not found", sequence_context_id)
	}

	if t.SequenceState.State == nil {
		t.SequenceState.State = make(map[int64][]ActionResponse)
	}

	for _, actionContext := range actionContexts {
		requestContext, exists := t.ReqExecutionContext[actionContext.ReqContextId]
		if !exists {
			return fmt.Errorf("request execution context %d not found", actionContext.ReqContextId)
		}

		response, err := t.ExecuteAction(requestContext.Ctx, db, int64(actionContext.ActionId), actionContext.ActionType)
		if err != nil {
			return fmt.Errorf("failed to execute action %d: %w", actionContext.ActionId, err)
		}

		t.SequenceState.State[int64(actionContext.ActionId)] = append(
			t.SequenceState.State[int64(actionContext.ActionId)],
			*response,
		)

		t.SequenceState.Order = append(t.SequenceState.Order, int64(actionContext.ActionId))
	}

	return nil
}

func (t *Tester) ExecuteAction(ctx context.Context, db *sql.DB, action_id int64, action_type string) (*ActionResponse, error) {
    switch action_type {
    case "http":
        action, err := actions.GetHttpAction(db, action_id)
        if err != nil {
            return nil, fmt.Errorf("failed to get HTTP action: %w", err)
        }

        response, err := t.ExecuteHttpAction(ctx, action, db)
        if err != nil {
            return response, fmt.Errorf("failed to execute HTTP action: %w", err)
        }

        return response, nil

    case "websocket":
        action, err := actions.GetWSAction(db, action_id)
        if err != nil {
            return nil, fmt.Errorf("failed to get WebSocket action: %w", err)
        }

        conn := t.WebsocketManager[action_id]

        if conn == nil {
            var response *http.Response
            conn, response, err = t.CreateWebSocketConnection(ctx, action, db)
            if err != nil {
                return nil, err
            }

            t.WebsocketManager[action_id] = conn

            _ = response
        }

        return t.ExecuteWsAction(ctx, action, conn, db)

    default:
        return nil, fmt.Errorf("unsupported action type: %s", action_type)
    }
}

func (t *Tester) CreateWebSocketConnection(ctx context.Context, action actions.WSAction, db *sql.DB) (*websocket.Conn, *http.Response, error) {
    env, err := environment.GetEnvironment(db, int(action.EnvironmentID))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get environment: %w", err)
    }

    opts := action.Opts
    opts.URL = ReplaceEnvironmentVariables(opts.URL, env.KeysValues)

    wsURL, err := url.Parse(opts.URL)
    if err != nil {
        return nil, nil, fmt.Errorf("invalid websocket URL: %w", err)
    }

    query := wsURL.Query()
    for key, value := range opts.QueryParams {
        query.Set(key, ReplaceEnvironmentVariables(value, env.KeysValues))
    }
    wsURL.RawQuery = query.Encode()

    headers := http.Header{}
    for key, value := range opts.Headers {
        headers.Set(key, ReplaceEnvironmentVariables(value, env.KeysValues))
    }

    dialer := websocket.Dialer{
        HandshakeTimeout: opts.HandshakeTimeout,
        Subprotocols:     opts.Subprotocols,
    }

    conn, resp, err := dialer.DialContext(ctx, wsURL.String(), headers)
    if err != nil {
        return nil, resp, fmt.Errorf("failed to establish websocket connection: %w", err)
    }

    return conn, resp, nil
}


func (t *Tester) ExecuteHttpAction(ctx context.Context, action actions.HttpAction, db *sql.DB) (*ActionResponse, error) {
    env, err := environment.GetEnvironment(db, int(action.EnvironmentId))
    if err != nil {
        return &ActionResponse{Success: false, Message: fmt.Sprintf("failed to get environment: %v", err)}, fmt.Errorf("failed to get environment: %w", err)
    }

    opts := action.Opts

    opts.URL = ReplaceEnvironmentVariables(opts.URL, env.KeysValues)

    reqURL, err := url.Parse(opts.URL)
    if err != nil {
        return &ActionResponse{Success: false, Message: fmt.Sprintf("invalid request URL: %v", err)}, fmt.Errorf("invalid request URL: %w", err)
    }

    query := reqURL.Query()
    for _, item := range opts.Query {
        item = ReplaceEnvironmentVariables(item, env.KeysValues)
        parts := strings.SplitN(item, "=", 2)
        if len(parts) == 2 {
            query.Set(parts[0], parts[1])
        }
    }
    reqURL.RawQuery = query.Encode()

    body := ReplaceEnvironmentVariables(string(opts.Body), env.KeysValues)

    req, err := http.NewRequestWithContext(ctx, opts.Method, reqURL.String(), strings.NewReader(body))
    if err != nil {
        return &ActionResponse{Success: false, Message: fmt.Sprintf("failed to create HTTP request: %v", err)}, fmt.Errorf("failed to create HTTP request: %w", err)
    }

    for key, value := range opts.Headers {
        req.Header.Set(key, ReplaceEnvironmentVariables(value, env.KeysValues))
    }

    resp, err := t.Client.Do(req)
    if err != nil {
        return &ActionResponse{Success: false, Message: fmt.Sprintf("HTTP request failed: %v", err)}, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return &ActionResponse{
            Success:    false,
            Message:    fmt.Sprintf("failed to read response body: %v", err),
            StatusCode: resp.StatusCode,
            Status:     resp.Status,
            Headers:    resp.Header,
            Cookies:    resp.Cookies(),
        }, fmt.Errorf("failed to read response body: %w", err)
    }

    response := &ActionResponse{
        Success:    true,
        Message:    "HTTP request executed successfully",
        Body:       string(responseBody),
        StatusCode: resp.StatusCode,
        Status:     resp.Status,
        Headers:    resp.Header,
        Cookies:    resp.Cookies(),
        Data:       make(map[string]any),
    }

    return response, nil
}

func ReplaceEnvironmentVariables(input string, values map[string]any) string {
    for key, value := range values {
        placeholder := "{" + key + "}"
        input = strings.ReplaceAll(input, placeholder, fmt.Sprint(value))
    }

    return input
}

func (t *Tester) CancelRequest(reqContextID int) error {
	requestContext, exists := t.ReqExecutionContext[reqContextID]
	if !exists {
		return fmt.Errorf("request execution context %d not found", reqContextID)
	}

	requestContext.Cancel()

	return nil
}

func (t *Tester) CancelSequence(sequenceContextID int) error {
	actionContexts, exists := t.SequnceExecutionContext[sequenceContextID]
	if !exists {
		return fmt.Errorf("sequence execution context %d not found", sequenceContextID)
	}

	for _, actionContext := range actionContexts {
		if requestContext, exists := t.ReqExecutionContext[actionContext.ReqContextId]; exists {
			requestContext.Cancel()
		}
	}

	return nil
}

func (t *Tester) ExecuteWsAction(ctx context.Context, action actions.WSAction, conn *websocket.Conn, db *sql.DB) (*ActionResponse, error) {
    env, err := environment.GetEnvironment(db, int(action.EnvironmentID))
    if err != nil {
        return &ActionResponse{Success: false, Message: err.Error()}, err
    }

    opts := action.Opts

    response := &ActionResponse{
        Success: true,
        Message: "WebSocket connection established successfully",
        Data:    make(map[string]any),
    }

    messages := append([]actions.WSMessage(nil), opts.Messages...)
    sort.Slice(messages, func(i, j int) bool {
        return messages[i].Sequence < messages[j].Sequence
    })

    for _, message := range messages {
        payload := ReplaceEnvironmentVariables(string(message.Payload), env.KeysValues)

        messageType := websocket.TextMessage

        if strings.EqualFold(message.MessageType, "BINARY") {
            messageType = websocket.BinaryMessage
        } else if !strings.EqualFold(message.MessageType, "TEXT") {
            err := fmt.Errorf("unsupported websocket message type: %s", message.MessageType)
            response.Success = false
            response.Message = err.Error()
            return response, err
        }

        if opts.WriteTimeout > 0 {
            conn.SetWriteDeadline(time.Now().Add(opts.WriteTimeout))
        }

        if err := conn.WriteMessage(messageType, []byte(payload)); err != nil {
            response.Success = false
            response.Message = fmt.Sprintf("failed to send websocket message: %v", err)
            return response, err
        }
    }

    response.Message = "WebSocket messages sent successfully"

    return response, nil
}

func (t *Tester) ListenWebSocket(ctx context.Context, actionID int64, conn *websocket.Conn) {
    go func() {
        <-ctx.Done()
        conn.Close()
    }()

    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        messageType, message, err := conn.ReadMessage()
        if err != nil {
            if ctx.Err() != nil {
                return
            }

            fmt.Printf("WebSocket Action ID %d read error: %v\n", actionID, err)
            return
        }

        messageTypeName := "TEXT"
        if messageType == websocket.BinaryMessage {
            messageTypeName = "BINARY"
        }

        fmt.Printf("Action ID: %d\n", actionID)
        fmt.Printf("Message Type: %s\n", messageTypeName)
        fmt.Printf("Message: %s\n", string(message))
    }
}