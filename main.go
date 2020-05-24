// https://steelseries.com/engine
// https://github.com/SteelSeries/gamesense-sdk

package main

import (
  "flag"
  "net/http"
  "time"

  "github.com/spf13/afero"
  "github.com/ztrue/tracerr"

  "github.com/ztrue/steelseries-status/steelseries"
)

const Developer = "github.com/ztrue"

const DisplayName = "Build Status"

const EventName = "PASS"

const GameName = "BUILD_STATUS"

func main() {
  cfg := parseConfig()

  fs := afero.NewOsFs()

  httpClient := &http.Client{
    Timeout: 30 * time.Second,
  }

  if err := NewApp(cfg, fs, httpClient).Listen(); err != nil {
    tracerr.PrintSourceColor(err)
  }
}

func parseConfig() Config {
  command := flag.String("c", "go test", "command to run")
  corePropsPath := flag.String("p", "", "coreProps.json path")
  intervalMS := flag.Int("t", 1000, "interval in ms")
  // TODO help command

  flag.Parse()

  if *corePropsPath == "" {
    // TODO Detect OS
    *corePropsPath = steelseries.CorePropsPathMacos
  }

  return Config{
    Command: *command,
    CorePropsPath: *corePropsPath,
    Developer: Developer,
    DisplayName: DisplayName,
    EventName: EventName,
    GameName: GameName,
    Interval: time.Duration(*intervalMS) * time.Millisecond,
  }
}
