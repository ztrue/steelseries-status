package steelseries

import (
  "encoding/json"
  "errors"
  "os"
)

const CorePropsPathMacos = "/Library/Application Support/SteelSeries Engine 3/coreProps.json"
const CorePropsPathWindows = "%PROGRAMDATA%/SteelSeries/SteelSeries Engine 3/coreProps.json"

var ErrNoCorePropsFile = errors.New("coreProps.json file not found, make sure SteelSeries Engine installed")

type CoreProps struct {
  Address string `json:"address"`
  EncryptedAddress string `json:"encrypted_address"`
}

type Discoverer struct {
  corePropsPath string
}

func NewDiscoverer(corePropsPath string) *Discoverer {
  return &Discoverer{
    corePropsPath: corePropsPath,
  }
}

func (d *Discoverer) CorePropsFileExists() bool {
  _, err := os.Stat(d.corePropsPath)
  return err == nil
}

func (d *Discoverer) CoreProps() (CoreProps, error) {
  if !d.CorePropsFileExists() {
    return CoreProps{}, ErrNoCorePropsFile
  }

  return d.readCoreProps()
}

func (d *Discoverer) readCoreProps() (CoreProps, error) {
  var props CoreProps

  f, err := os.Open(d.corePropsPath)
  if err != nil {
    return props, err
  }
  defer f.Close()

  err = json.NewDecoder(f).Decode(&props)

  return props, err
}
