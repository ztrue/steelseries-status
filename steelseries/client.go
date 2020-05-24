package steelseries

import (
  "bytes"
  "encoding/json"
  "net/http"
)

const DeviceKeyboard = "keyboard"

const ModeColor = "color"

var ColorGreen = RGB{0, 255, 0}
var ColorRed = RGB{255, 0, 0}

type Event struct {
  Game string `json:"game"`
  Event string `json:"event"`
  Data Data `json:"data"`
}

type Data struct {
  Value int `json:"value"`
}

type RGB struct {
  Red int `json:"red"`
  Green int `json:"green"`
  Blue int `json:"blue"`
}

type BindEvent struct {
  Game string `json:"game"`
  Event string `json:"event"`
  MinValue int `json:"min_value"`
  MaxValue int `json:"max_value"`
  IconID int `json:"icon_id"`
  Handlers []Handler `json:"handlers"`
}

type Handler struct {
  DeviceType string `json:"device-type"`
  // Zone string `json:"zone"`
  CustomZoneKeys []int `json:"custom-zone-keys"`
  Color Gradient `json:"color"`
  Mode string `json:"mode"`
}

type Gradient struct {
  Gradient GradientValues `json:"gradient"`
}

type GradientValues struct {
  Zero RGB `json:"zero"`
  Hundred RGB `json:"hundred"`
}

type HttpClient interface {
  Do(*http.Request) (*http.Response, error)
}

type Client struct {
  addr string
  gameName string
  httpClient HttpClient
}

func NewClient(addr string, httpClient HttpClient, gameName string) *Client {
  return &Client{
    addr: addr,
    gameName: gameName,
    httpClient: httpClient,
  }
}

func (c *Client) Register(eventName string) error {
  // TODO Specify all keys normally
  var keys []int
  for i := -1000; i < 1000; i++ {
    keys = append(keys, i+1)
  }

  event := BindEvent{
    Game: c.gameName,
    Event: eventName,
    MinValue: 0,
    MaxValue: 100,
    IconID: 1,
    Handlers: []Handler{
      {
        DeviceType: DeviceKeyboard,
        // Zone: "main-keyboard",
        CustomZoneKeys: keys,
        Color: Gradient{
          GradientValues{
            Zero: ColorRed,
            Hundred: ColorGreen,
          },
        },
        Mode: ModeColor,
      },
    },
  }

  buf := &bytes.Buffer{}
  if err := json.NewEncoder(buf).Encode(event); err != nil {
    return err
  }

  req, err := http.NewRequest("POST", "http://" + c.addr + "/bind_game_event", buf)
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")

  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()

  return nil
}

func (c *Client) Update(eventName string, value int) error {
  event := Event{
    Game: c.gameName,
    Event: eventName,
    Data: Data{value},
  }

  buf := &bytes.Buffer{}
  if err := json.NewEncoder(buf).Encode(event); err != nil {
    return err
  }

  req, err := http.NewRequest("POST", "http://" + c.addr + "/game_event", buf)
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")

  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()

  return nil
}
