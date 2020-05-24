package main

import (
  "encoding/json"
  "net/http"
  "os"
  "os/exec"
  "time"

  "github.com/ztrue/tracerr"

  "github.com/ztrue/steelseries-status/steelseries"
)

const corePropsPath = "/Library/Application Support/SteelSeries Engine 3/coreProps.json"

const GameName = "GO_TESTS"
const EventPass = "PASS"

func main() {
  if err := run(corePropsPath); err != nil {
    tracerr.PrintSourceColor(err)
  }
}

type CoreProps struct {
  Address string `json:"address"`
  EncryptedAddress string `json:"encrypted_address"`
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

  httpClient := &http.Client{
    Timeout: 30 * time.Second,
  }

  ss := steelseries.NewClient(httpClient, addr, GameName)

  event := ss.BuildBindGameEvent(EventPass)
  if err := ss.SendBindGameEvent(event); err != nil {
    return tracerr.Wrap(err)
  }

  for range time.NewTicker(time.Second).C {
    pass := exec.Command("go", "test", "./tests").Run() == nil
    value := 100
    if !pass {
      value = 0
    }
    event := ss.BuildGameEvent(EventPass, value)
    ss.SendGameEvent(event)
  }

  return nil
}
