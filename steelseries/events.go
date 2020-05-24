package steelseries

const DeviceKeyboard = "keyboard"

const ModeColor = "color"

var ColorGreen = RGB{0, 255, 0}
var ColorRed = RGB{255, 0, 0}

type Blank struct {
  Game string `json:"game"`
}

type BindGameEvent struct {
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
  Color Color `json:"color"`
  Mode string `json:"mode"`
}

type Color struct {
  Gradient Gradient `json:"gradient"`
}

type Gradient struct {
  Zero RGB `json:"zero"`
  Hundred RGB `json:"hundred"`
}

type RGB struct {
  Red int `json:"red"`
  Green int `json:"green"`
  Blue int `json:"blue"`
}

type GameEvent struct {
  Game string `json:"game"`
  Event string `json:"event"`
  Data Data `json:"data"`
}

type Data struct {
  Value int `json:"value"`
}

type GameMetadata struct {
  Game string `json:"game"`
  GameDisplayName string `json:"game_display_name"`
  Developer string `json:"developer"`
  DeinitializeTimerLengthMS int `json:"deinitialize_timer_length_ms"`
}
