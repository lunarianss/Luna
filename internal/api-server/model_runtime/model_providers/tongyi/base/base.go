package base

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

// const (
// 	wsURL      = "wss://dashscope.aliyuncs.com/api-ws/v1/inference/"
// 	outputFile = "output.mp3"
// )

type Header struct {
	Action       string                 `json:"action"`
	TaskID       string                 `json:"task_id"`
	Streaming    string                 `json:"streaming"`
	Event        string                 `json:"event"`
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type Payload struct {
	TaskGroup  string     `json:"task_group"`
	Task       string     `json:"task"`
	Function   string     `json:"function"`
	Model      string     `json:"model"`
	Parameters Params     `json:"parameters"`
	Resources  []Resource `json:"resources"`
	Input      Input      `json:"input"`
}

type Params struct {
	TextType   string `json:"text_type"`
	Voice      string `json:"voice"`
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	Volume     int    `json:"volume"`
	Rate       int    `json:"rate"`
	Pitch      int    `json:"pitch"`
}

type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
	Text string `json:"text"`
}

type Event struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

type TongyiTTSSDK struct {
	apiKey            string
	wsUrl             string
	conn              *websocket.Conn
	binaryAudioQueue  chan []byte
	eventMessageQueue chan *Event
	errorQueue        chan error
	done              chan struct{}
	taskStart         chan struct{}
	taskID            string
	dialer            *websocket.Dialer
	model             string
	voice             string
	format            string
}

func NewTongyiTTSSDK(apiKey, wsUrl, model, voice, format string, dialer *websocket.Dialer) *TongyiTTSSDK {
	return &TongyiTTSSDK{
		apiKey:            apiKey,
		wsUrl:             wsUrl,
		done:              make(chan struct{}),
		taskStart:         make(chan struct{}),
		binaryAudioQueue:  make(chan []byte, 15),
		eventMessageQueue: make(chan *Event, 15),
		errorQueue:        make(chan error, 1),
		dialer:            dialer,
		model:             model,
		voice:             voice,
		format:            format,
	}
}

func (ts *TongyiTTSSDK) Close() {
	close(ts.binaryAudioQueue)
	close(ts.eventMessageQueue)
	close(ts.done)
	close(ts.taskStart)
	// log.Info("close all channel")
}

func (ts *TongyiTTSSDK) connectWebSocket() error {
	header := make(http.Header)
	header.Add("X-DashScope-DataInspection", "enable")
	header.Add("Authorization", fmt.Sprintf("bearer %s", ts.apiKey))

	dialer := ts.dialer

	if dialer == nil {
		dialer = websocket.DefaultDialer
	}

	conn, _, err := dialer.Dial(ts.wsUrl, header)

	if err != nil {
		return errors.WithCode(code.ErrTTSWebSocket, fmt.Sprintf("failed to connect tongyi websocket %s", ts.wsUrl))
	}

	ts.conn = conn
	return nil
}

func (ts *TongyiTTSSDK) CloseConnection() {
	if ts.conn != nil {
		ts.conn.Close()
		// log.Info("close websocket connection")
	}
}

func (ts *TongyiTTSSDK) PushError(err error) {
	defer close(ts.errorQueue)
	ts.errorQueue <- err
}

func (ts *TongyiTTSSDK) Done() {
	ts.done <- struct{}{}
	// log.Info("done finish!")
}

func (ts *TongyiTTSSDK) startResultReceiver() {

	go func() {
		defer ts.Done()
		errorTimes := 0

		for {
			msgType, message, err := ts.conn.ReadMessage()
			if err != nil {
				log.Errorf(fmt.Sprintf("error when parse tongyi websocket message: %+v", err))
				ts.PushError(err)
				return
			}

			if msgType == websocket.BinaryMessage {
				ts.binaryAudioQueue <- message
				// audio storage
				// if err := ts.writeBinaryDataToFile(message, outputFile); err != nil {
				// 	log.Errorf(fmt.Sprintf("error when write tongyi websocket message to file: %+v", err))
				// 	ts.PushError(err)
				// 	return
				// }
			} else {
				var event Event
				err = json.Unmarshal(message, &event)
				if err != nil {
					if errorTimes > 3 {
						ts.PushError(errors.WithCode(code.ErrEncodingJSON, fmt.Sprintf("error when parse tongyi websocket event message: %s", err.Error())))
						return
					}

					log.Errorf(fmt.Sprintf("error when parse tongyi websocket event message: %s, errorTim %d", err.Error(), errorTimes))
					errorTimes += 1
					continue
				}

				ts.eventMessageQueue <- &event
				if ts.handleEvent(&event) {
					return
				}
			}
		}
	}()

}

func (ts *TongyiTTSSDK) handleEvent(event *Event) bool {

	switch event.Header.Event {
	case "task-started":
		ts.taskStart <- struct{}{}
	case "result-generated":
		return false
	case "task-finished":
		return true
	case "task-failed":
		ts.handleTaskFailed(event)
		return true
	default:
		log.Errorf("unexpected tongyi event: %v\n", event)
	}
	return false
}

func (ts *TongyiTTSSDK) handleTaskFailed(event *Event) {
	if event.Header.ErrorMessage != "" {
		ts.PushError(fmt.Errorf("task failed reason due to %s", event.Header.ErrorMessage))
	} else {
		ts.PushError(fmt.Errorf("task failed due to unknown reason"))

	}
}

func (ts *TongyiTTSSDK) sendRunTaskCmd() error {
	runTaskCmd, err := ts.generateRunTaskCmd()
	if err != nil {
		return err
	}
	err = ts.conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return err
}

func (ts *TongyiTTSSDK) generateRunTaskCmd() (string, error) {
	ts.taskID = uuid.NewString()

	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    ts.taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "tts",
			Function:  "SpeechSynthesizer",
			Model:     "cosyvoice-v1",
			Parameters: Params{
				TextType:   "PlainText",
				Voice:      ts.voice,
				Format:     ts.format,
				SampleRate: 22050,
				Volume:     50,
				Rate:       1,
				Pitch:      1,
			},
			Input: Input{},
		},
	}

	if runTaskCmd.Payload.Parameters.Format == "" {
		runTaskCmd.Payload.Parameters.Format = "mp3"
	}

	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	if err != nil {
		return "", errors.WithCode(code.ErrEncodingJSON, err.Error())
	}
	return string(runTaskCmdJSON), nil
}

func (ts *TongyiTTSSDK) sendContinueTaskCmd(texts []string) error {

	for _, text := range texts {
		runTaskCmd, err := ts.generateContinueTaskCmd(text, ts.taskID)
		if err != nil {
			return err
		}

		err = ts.conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
		if err != nil {
			return errors.WithCode(code.ErrTTSWebSocketWrite, err.Error())
		}
	}
	return nil
}

func (ts *TongyiTTSSDK) generateContinueTaskCmd(text string, taskID string) (string, error) {
	runTaskCmd := Event{
		Header: Header{
			Action:    "continue-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{
				Text: text,
			},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	if err != nil {
		return "", errors.WithCode(code.ErrEncodingJSON, err.Error())
	}
	return string(runTaskCmdJSON), nil
}

func (ts *TongyiTTSSDK) sendFinishTaskCmd() error {
	finishTaskCmd, err := ts.generateFinishTaskCmd(ts.taskID)
	if err != nil {
		return err
	}
	err = ts.conn.WriteMessage(websocket.TextMessage, []byte(finishTaskCmd))
	return err
}

func (ts *TongyiTTSSDK) generateFinishTaskCmd(taskID string) (string, error) {
	finishTaskCmd := Event{
		Header: Header{
			Action:    "finish-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{},
		},
	}
	finishTaskCmdJSON, err := json.Marshal(finishTaskCmd)

	if err != nil {
		return "", err
	}
	return string(finishTaskCmdJSON), nil
}

// func (ts *TongyiTTSSDK) writeBinaryDataToFile(data []byte, filePath string) error {
// 	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.Write(data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (ts *TongyiTTSSDK) GetErrorQueues() <-chan error {
	return ts.errorQueue
}

func (ts *TongyiTTSSDK) GetAudioBinaryQueues() <-chan []byte {
	return ts.binaryAudioQueue
}

func (ts *TongyiTTSSDK) GetEventQueue() <-chan *Event {
	return ts.eventMessageQueue
}

func (ts *TongyiTTSSDK) GetDone() <-chan struct{} {
	return ts.done
}

func (ts *TongyiTTSSDK) Generate(ctx context.Context, texts []string) {
	err := ts.connectWebSocket()

	if err != nil {
		ts.PushError(err)
		return
	}

	ts.startResultReceiver()

	err = ts.sendRunTaskCmd()

	if err != nil {
		ts.PushError(err)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*7)

	defer cancel()

	select {
	case <-ctx.Done():
		ts.PushError(errors.WithCode(code.ErrContextTimeout, "tongy start-run event context timeout"))
		return
	case <-ts.taskStart:
	}

	if err := ts.sendContinueTaskCmd(texts); err != nil {
		ts.PushError(err)
		return
	}

	if err := ts.sendFinishTaskCmd(); err != nil {
		ts.PushError(err)
		return
	}
}
