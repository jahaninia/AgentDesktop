package jolClientPush

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type Client struct {
	token     string
	extension string
	address   string
	wg        *sync.WaitGroup
	ctx       context.Context
	debug     bool
	popup     string
	fields    []string
}

func NewClient(token, address, popup, extension string, fields []string, debug bool, wg *sync.WaitGroup, ctx context.Context) *Client {
	return &Client{
		token:     token,
		extension: extension,
		address:   address,
		debug:     debug,
		popup:     popup,
		fields:    fields,
		wg:        wg,
		ctx:       ctx,
	}
}

// func (c *Client) Connect() {
// 	defer c.wg.Done()

// 	backoff := time.Second
// 	for {
// 		select {
// 		case <-c.ctx.Done():
// 			fmt.Println("Client stopped")
// 			return
// 		default:
// 			err := c.ConnectOnce()
// 			if err != nil {
// 				fmt.Println("connect error:", err)
// 				time.Sleep(backoff)
// 				if backoff < 30*time.Second {
// 					backoff *= 2
// 				}
// 				continue
// 			}
// 			backoff = time.Second // reset بعد از موفقیت
// 		}
// 	}
// }

func (c *Client) ConnectOnce() error {
	// مسیر جدید سمت سرور (GET برای اتصال SSE)
	url := fmt.Sprintf("%s/connect?token=%s&ext=%s", c.address, c.token, c.extension)

	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("SSE Connected to server:", url)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			c.processMessage([]byte(strings.TrimSpace(data)))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
func (c *Client) processMessage(data []byte) {
	jsonData := map[string]string{}
	if c.debug {
		fmt.Printf("Received message: %s\n", string(data))
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	query := c.popup
	for _, v := range c.fields {
		query = strings.ReplaceAll(query, fmt.Sprintf("{%s}", v), jsonData[v])
	}

	c.open(query)
}

func (c *Client) open3(url string) error {

	fmt.Println(url)
	cmdStr := fmt.Sprintf(`start "" "%s"`, url)
	fmt.Println(cmdStr)
	return exec.Command("cmd", "/C", cmdStr).Start()
}
func (c *Client) open1(url string) error {
	var cmd string
	var args []string
	fmt.Println(url)
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", ""}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
func (c *Client) open(url string) error {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}
