package main

import (
  "os/exec"
  "strings"
  "time"

  "github.com/spf13/afero"
  "github.com/ztrue/tracerr"

  "github.com/ztrue/steelseries-status/steelseries"
)

type Config struct {
  Command string
  CorePropsPath string
  Developer string
  DisplayName string
  EventName string
  GameName string
  Interval time.Duration
}

type App struct {
  cfg Config
  fs afero.Fs
  httpClient steelseries.HTTPClient
  ss *steelseries.Client
  ticker *time.Ticker
}

func NewApp(cfg Config, fs afero.Fs, httpClient steelseries.HTTPClient) *App {
  return &App{
    cfg: cfg,
    fs: fs,
    httpClient: httpClient,
  }
}

func (a *App) Init() error {
  d := steelseries.NewDiscoverer(a.fs, a.cfg.CorePropsPath)
  props, err := d.CoreProps()
  if err != nil {
    return tracerr.Wrap(err)
  }

  a.ss = steelseries.NewClient(a.httpClient, props.Address, a.cfg.GameName)

  a.ticker = time.NewTicker(a.cfg.Interval)

  return nil
}

func (a *App) Listen() error {
  if err := a.Init(); err != nil {
    return tracerr.Wrap(err)
  }

  if err := a.Start(); err != nil {
    return tracerr.Wrap(err)
  }

  if err := a.Process(); err != nil {
    return tracerr.Wrap(err)
  }

  return nil
}

func (a *App) Process() error {
  for range a.ticker.C {
    segments := strings.Split(a.cfg.Command, " ")
    pass := exec.Command(segments[0], segments[1:]...).Run() == nil
    value := 100
    if !pass {
      value = 0
    }
    event := a.ss.BuildGameEvent(a.cfg.EventName, value)
    if err := a.ss.SendGameEvent(event); err != nil {
      return tracerr.Wrap(err)
    }
  }
  return nil
}

func (a *App) Start() error {
  metadata := a.ss.BuildGameMetadata(a.cfg.DisplayName, a.cfg.Developer)
  if err := a.ss.SendGameMetadata(metadata); err != nil {
    return tracerr.Wrap(err)
  }

  event := a.ss.BuildBindGameEvent(a.cfg.EventName)
  return tracerr.Wrap(a.ss.SendBindGameEvent(event))
}

func (a *App) Stop() error {
  return a.ss.SendStopGame()
}
