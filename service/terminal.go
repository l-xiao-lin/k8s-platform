package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"k8s-platform/setting"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

const END_OF_TRANSMISSION = "\u0004"

type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		zap.L().Error("t.wsConn.ReadMessage failed", zap.Error(err))
		return copy(p, END_OF_TRANSMISSION), err
	}
	var msg TerminalMessage
	if err = json.Unmarshal(message, &msg); err != nil {
		zap.L().Error("json.Unmarshal failed", zap.Error(err))
		return copy(p, END_OF_TRANSMISSION), err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Rows, Height: msg.Cols}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		zap.L().Error("unknow message type", zap.Error(err))
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknow message type %s", msg.Operation)
	}
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		zap.L().Error("json.Marshal failed", zap.Error(err))
		return 0, err
	}
	if err = t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		zap.L().Error("t.wsConn.WriteMessage failed", zap.Error(err))
		return 0, nil
	}
	return len(p), nil
}

func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}

func (t *TerminalSession) Done() {
	close(t.doneChan)
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 10
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return upgrader
}()

// NewTerminalSession 该方法用于升级http协至websocket，并new 一个TerminalSession类型的对象返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	zap.L().Info("receive param", zap.String("namespace", namespace), zap.String("pod_name", podName), zap.String("container_name", containerName))

	conf, err := clientcmd.BuildConfigFromFlags("", setting.Conf.Kubeconfig)
	if err != nil {
		zap.L().Error("clientcmd.BuildConfigFromFlags failed", zap.Error(err))
		return
	}

	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		zap.L().Error("get pty failed", zap.Error(err))
		return
	}
	defer func() {
		zap.L().Info("close session")
		pty.Close()

	}()

	pod, err := K8s.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("获取pod失败", zap.Error(err))
		pty.Write([]byte(fmt.Sprintf("获取pod失败:%v", err)))
		pty.Done()
		return
	}

	//根据容器数据决定 是否需要指定Container字段
	var execOptions v1.PodExecOptions

	if len(pod.Spec.Containers) == 1 {
		execOptions = v1.PodExecOptions{
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
			Command: []string{"/bin/bash"},
		}
	} else {
		execOptions = v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"/bin/bash"},
		}

	}

	req := K8s.clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&execOptions, scheme.ParameterCodec)

	zap.L().Info("Request", zap.String("url", req.URL().String()))

	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}
	err = executor.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})

	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err:%v", err)
		zap.L().Error("Exec to pod error", zap.Error(err))
		pty.Write([]byte(msg))
		pty.Done()
	}
}
