package main

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/exec"
  "time"

  "github.com/ztrue/tracerr"
)

const corePropsPath = "/Library/Application Support/SteelSeries Engine 3/coreProps.json"

const GameName = "GO_TESTS"
const EventPass = "PASS"
const EventFail = "FAIL"

var ColorGreen = RGB{0, 255, 0}
var ColorRed = RGB{255, 0, 0}

func main() {
  if err := run(corePropsPath); err != nil {
    tracerr.PrintSourceColor(err)
  }
}

type CoreProps struct {
  Address string `json:"address"`
  EncryptedAddress string `json:"encrypted_address"`
}

type Event struct {
  Game string `json:"game"`
  Event string `json:"event"`
  Data Data `json:"data"`
}

type Data struct {
  Value int `json:"value"`
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

type RGB struct {
  Red int `json:"red"`
  Green int `json:"green"`
  Blue int `json:"blue"`
}

type Gradient struct {
  Gradient GradientValues `json:"gradient"`
}

type GradientValues struct {
  Zero RGB `json:"zero"`
  Hundred RGB `json:"hundred"`
}

func run(path string) error {
  // TODO Check file exists

  f, err := os.Open(path)
  if err != nil {
    return tracerr.Wrap(err)
  }
  defer f.Close()

  var props CoreProps

  if err := json.NewDecoder(f).Decode(&props); err != nil {
    return tracerr.Wrap(err)
  }

  addr := props.Address

  client := &http.Client{
    Timeout: 30 * time.Second,
  }

  {
    // TODO Specify all keys normally
    var keys []int
    for i := -1000; i < 1000; i++ {
      keys = append(keys, i+1)
    }

    event := BindEvent{
      Game: GameName,
      Event: EventPass,
      MinValue: 0,
      MaxValue: 100,
      IconID: 1,
      Handlers: []Handler{
        {
          DeviceType: "keyboard",
          // Zone: "main-keyboard",
          CustomZoneKeys: keys,
          Color: Gradient{
            GradientValues{
              Zero: ColorRed,
              Hundred: ColorGreen,
            },
          },
          Mode: "color",
        },
      },
    }


    buf := &bytes.Buffer{}
    if err := json.NewEncoder(buf).Encode(event); err != nil {
      return tracerr.Wrap(err)
    }

    req, err := http.NewRequest("POST", "http://" + addr + "/bind_game_event", buf)
    if err != nil {
      return tracerr.Wrap(err)
    }
    req.Header.Set("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
      return tracerr.Wrap(err)
    }
    defer res.Body.Close()

    content, err := ioutil.ReadAll(res.Body)
    if err != nil {
      return tracerr.Wrap(err)
    }
    log.Println(res.StatusCode)
    log.Println(string(content))
  }

  for range time.NewTicker(time.Second).C {
    pass := exec.Command("go", "test", "./tests").Run() == nil
    update(client, addr, pass)
  }

  return nil
}

func update(client *http.Client, addr string, pass bool) error {
  value := 100
  if !pass {
    value = 0
  }

  event := Event{
    Game: GameName,
    Event: EventPass,
    Data: Data{value},
  }

  buf := &bytes.Buffer{}
  if err := json.NewEncoder(buf).Encode(event); err != nil {
    return tracerr.Wrap(err)
  }

  req, err := http.NewRequest("POST", "http://" + addr + "/game_event", buf)
  if err != nil {
    return tracerr.Wrap(err)
  }
  req.Header.Set("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    return tracerr.Wrap(err)
  }
  defer res.Body.Close()

  content, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return tracerr.Wrap(err)
  }

  log.Println(res.StatusCode)
  log.Println(string(content))
  return nil
}
