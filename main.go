// https://steelseries.com/engine
// https://github.com/SteelSeries/gamesense-sdk

package main

import (
  "flag"
  "net/http"
  "os/exec"
  "strings"
  "time"

  "github.com/spf13/afero"
  "github.com/ztrue/tracerr"

  "github.com/ztrue/steelseries-status/steelseries"
)

const GameName = "BUILD_STATUS"
const DisplayName = "Build Status"
const Developer = "github.com/ztrue"
const EventPass = "PASS"

type Config struct {
  Command string
  CorePropsPath string
  Interval time.Duration
}

func main() {
  if err := run(parseConfig()); err != nil {
    tracerr.PrintSourceColor(err)
  }
}

func parseConfig() Config {
  command := flag.String("c", "go test", "command to run")
  corePropsPath := flag.String("p", "", "coreProps.json path")
  intervalMS := flag.Int("t", 1000, "interval in ms")

  flag.Parse()

  if *corePropsPath == "" {
    // TODO Detect OS
    *corePropsPath = steelseries.CorePropsPathMacos
  }

  return Config{
    Command: *command,
    CorePropsPath: *corePropsPath,
    Interval: time.Duration(*intervalMS) * time.Millisecond,
  }
}

func run(cfg Config) error {
  ss, err := start(cfg)
  if err != nil {
    return tracerr.Wrap(err)
  }

  return process(ss, cfg)
}

func start(cfg Config) (*steelseries.Client, error) {
  fs := afero.NewOsFs()

  d := steelseries.NewDiscoverer(fs, cfg.CorePropsPath)
  props, err := d.CoreProps()
  if err != nil {
    return nil, tracerr.Wrap(err)
  }

  httpClient := &http.Client{
    Timeout: 30 * time.Second,
  }

  ss := steelseries.NewClient(httpClient, props.Address, GameName)

  metadata := ss.BuildGameMetadata(DisplayName, Developer)
  if err := ss.SendGameMetadata(metadata); err != nil {
    return nil, tracerr.Wrap(err)
  }

  event := ss.BuildBindGameEvent(EventPass)
  if err := ss.SendBindGameEvent(event); err != nil {
    return nil, tracerr.Wrap(err)
  }

  return ss, nil
}

func process(ss *steelseries.Client, cfg Config) error {
  for range time.NewTicker(cfg.Interval).C {
    segments := strings.Split(cfg.Command, " ")
    pass := exec.Command(segments[0], segments[1:]...).Run() == nil
    value := 100
    if !pass {
      value = 0
    }
    event := ss.BuildGameEvent(EventPass, value)
    if err := ss.SendGameEvent(event); err != nil {
      return tracerr.Wrap(err)
    }
  }
  return nil
}
