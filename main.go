// https://steelseries.com/engine
// https://github.com/SteelSeries/gamesense-sdk

package main

import (
  "net/http"
  "os/exec"
  "time"

  "github.com/ztrue/tracerr"

  "github.com/ztrue/steelseries-status/steelseries"
)

const GameName = "GO_TESTS"
const EventPass = "PASS"

func main() {
  if err := run(); err != nil {
    tracerr.PrintSourceColor(err)
  }
}

func run() error {
  d := steelseries.NewDiscoverer(steelseries.CorePropsPathMacos)
  props, err := d.CoreProps()
  if err != nil {
    return tracerr.Wrap(err)
  }

  httpClient := &http.Client{
    Timeout: 30 * time.Second,
  }

  ss := steelseries.NewClient(httpClient, props.Address, GameName)

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
