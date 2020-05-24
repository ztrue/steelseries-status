package steelseries

import (
  "bytes"
  "encoding/json"
  "net/http"
)

type HttpClient interface {
  Do(*http.Request) (*http.Response, error)
}

type Client struct {
  addr string
  gameName string
  httpClient HttpClient
}

func NewClient(httpClient HttpClient, addr string, gameName string) *Client {
  return &Client{
    addr: addr,
    gameName: gameName,
    httpClient: httpClient,
  }
}

func (c *Client) BuildBindGameEvent(eventName string) BindGameEvent {
  return BindGameEvent{
    Game: c.gameName,
    Event: eventName,
    MinValue: 0,
    MaxValue: 100,
    IconID: 1,
    Handlers: []Handler{
      {
        DeviceType: DeviceKeyboard,
        CustomZoneKeys: allKeys(),
        Color: Color{
          Gradient{
            Zero: ColorRed,
            Hundred: ColorGreen,
          },
        },
        Mode: ModeColor,
      },
    },
  }
}

func (c *Client) BuildGameEvent(eventName string, value int) GameEvent {
  return GameEvent{
    Game: c.gameName,
    Event: eventName,
    Data: Data{value},
  }
}

func (c *Client) SendBindGameEvent(event BindGameEvent) error {
  return c.send("bind_game_event", event)
}

func (c *Client) SendGameEvent(event GameEvent) error {
  return c.send("game_event", event)
}

func (c *Client) buildRequest(endpoint string, data interface{}) (*http.Request, error) {
  buf := &bytes.Buffer{}
  if err := json.NewEncoder(buf).Encode(data); err != nil {
    return nil, err
  }

  req, err := http.NewRequest("POST", c.buildURL(endpoint), buf)
  if err != nil {
    return nil, err
  }
  req.Header.Set("Content-Type", "application/json")

  return req, nil
}

func (c *Client) buildURL(endpoint string) string {
  return "http://" + c.addr + "/" + endpoint
}

func (c *Client) send(endpoint string, data interface{}) error {
  req, err := c.buildRequest(endpoint, data)
  if err != nil {
    return err
  }

  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()

  return nil
}

func allKeys() []int {
  // TODO Should be a better way to specify all keys
  var keys []int
  for i := 0; i < 500; i++ {
    keys = append(keys, i+1)
  }
  return keys
}
