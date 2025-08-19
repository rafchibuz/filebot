package config

import (
  "encoding/json"
  "errors"
  "os"
  "path/filepath"
)

// Config defines application configuration loaded from config.json
// Fields are designed to be simple and overridable by the JSON file.
type Config struct {
  UploadDir   string   `json:"uploadDir"`
  OutputExcel string   `json:"outputExcel"`
  ExcelSheet  string   `json:"excelSheet"`
  Web         Web      `json:"web"`
  Patterns    Patterns `json:"patterns"`
}

type Web struct {
  Address string `json:"address"`
}

type Patterns struct {
  ContractDatePattern string `json:"contractDatePattern"`
  ActDatePattern      string `json:"actDatePattern"`
  BuyerLabelPattern   string `json:"buyerLabelPattern"`
  SellerLabelPattern  string `json:"sellerLabelPattern"`
  PoAFromPattern      string `json:"poaFromPattern"`
  PoAToPattern        string `json:"poaToPattern"`
  DateRegex           string `json:"dateRegex"`
  NameRegex           string `json:"nameRegex"`
}

// Default returns the built-in defaults used if config.json is missing or partial.
func Default() Config {
  return Config{
    UploadDir:   "uploads",
    OutputExcel: "output.xlsx",
    ExcelSheet:  "Data",
    Web: Web{
      Address: ":8080",
    },
    Patterns: Patterns{
      ContractDatePattern: `(?i)договор[^\n\r\d]{0,80}`,
      ActDatePattern:      `(?i)акт[^\n\r\d]{0,80}`,
      BuyerLabelPattern:   `(?i)покупател[ья]:?`,
      SellerLabelPattern:  `(?i)продавц[а|е]:?`,
      PoAFromPattern:      `(?i)доверенн(?:ость|ости|остью)[^\n\r]{0,80}от`,
      PoAToPattern:        `(?i)действи[еия][^\n\r]{0,80}до`,
      DateRegex:           `(?m)(\\b\\d{2}[\\./-]\\d{2}[\\./-]\\d{4}\\b)`,
      NameRegex:           `(?m)([А-ЯЁ][а-яё]+(?:\\s+[А-ЯЁ][а-яё]+){0,2}(?:\\s+[А-ЯЁ]\\.[А-ЯЁ]\\.)?)`,
    },
  }
}

// Load reads configuration from config.json in the working directory.
// If the file is missing, defaults are returned.
func Load(path string) (Config, error) {
  cfg := Default()

  if path == "" {
    path = "config.json"
  }

  // If file does not exist, return defaults without error
  if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
    return cfg, nil
  }

  f, err := os.Open(filepath.Clean(path))
  if err != nil {
    return cfg, err
  }
  defer f.Close()

  dec := json.NewDecoder(f)
  var fileCfg Config
  if err := dec.Decode(&fileCfg); err != nil {
    return cfg, err
  }

  // Merge fileCfg into defaults when non-empty
  if fileCfg.UploadDir != "" {
    cfg.UploadDir = fileCfg.UploadDir
  }
  if fileCfg.OutputExcel != "" {
    cfg.OutputExcel = fileCfg.OutputExcel
  }
  if fileCfg.ExcelSheet != "" {
    cfg.ExcelSheet = fileCfg.ExcelSheet
  }
  if fileCfg.Web.Address != "" {
    cfg.Web.Address = fileCfg.Web.Address
  }

  // Patterns: if provided, override any non-empty
  if fileCfg.Patterns.ContractDatePattern != "" {
    cfg.Patterns.ContractDatePattern = fileCfg.Patterns.ContractDatePattern
  }
  if fileCfg.Patterns.ActDatePattern != "" {
    cfg.Patterns.ActDatePattern = fileCfg.Patterns.ActDatePattern
  }
  if fileCfg.Patterns.BuyerLabelPattern != "" {
    cfg.Patterns.BuyerLabelPattern = fileCfg.Patterns.BuyerLabelPattern
  }
  if fileCfg.Patterns.SellerLabelPattern != "" {
    cfg.Patterns.SellerLabelPattern = fileCfg.Patterns.SellerLabelPattern
  }
  if fileCfg.Patterns.PoAFromPattern != "" {
    cfg.Patterns.PoAFromPattern = fileCfg.Patterns.PoAFromPattern
  }
  if fileCfg.Patterns.PoAToPattern != "" {
    cfg.Patterns.PoAToPattern = fileCfg.Patterns.PoAToPattern
  }
  if fileCfg.Patterns.DateRegex != "" {
    cfg.Patterns.DateRegex = fileCfg.Patterns.DateRegex
  }
  if fileCfg.Patterns.NameRegex != "" {
    cfg.Patterns.NameRegex = fileCfg.Patterns.NameRegex
  }

  return cfg, nil
} 