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
  DateRegex           string `json:"dateRegex"`

  VINPattern          string `json:"vinPattern"`
  VINRegex            string `json:"vinRegex"`
  VehicleModelPattern string `json:"vehicleModelPattern"`
  SellerCompanyDKPPattern string `json:"sellerCompanyDKPPattern"`
}

// Default returns the built-in defaults used if config.json is missing or partial.
func Default() Config {
  return Config{
    UploadDir:   "uploads",
    OutputExcel: "output.xlsx",
    ExcelSheet:  "Data",
    Web: Web{ Address: ":8080" },
    Patterns: Patterns{
      ContractDatePattern: `(?i)договор[^\n\r\d]{0,80}`,
      ActDatePattern:      `(?i)акт[^\n\r\d]{0,80}`,
      DateRegex:           `(?m)(\b\d{2}[\./-]\d{2}[\./-]\d{4}\b)`,
      VINPattern:          `(?i)vin[:\s-]*`,
      VINRegex:            `(?i)\b[ABCDEFGHJKLMNPRSTUVWXYZ\d]{17}\b`,
      VehicleModelPattern: `(?i)(модель|model)[:\s-]*([^\n\r]{1,80})`,
      SellerCompanyDKPPattern: `(?i)(продавец|продавца)[^\n\r]{0,10}[:\s-]*([^\n\r]{1,160})`,
    },
  }
}

// Load reads configuration from config.json in the working directory.
// If the file is missing, defaults are returned.
func Load(path string) (Config, error) {
  cfg := Default()

  if path == "" { path = "config.json" }
  if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) { return cfg, nil }

  f, err := os.Open(filepath.Clean(path))
  if err != nil { return cfg, err }
  defer f.Close()

  dec := json.NewDecoder(f)
  var fileCfg Config
  if err := dec.Decode(&fileCfg); err != nil { return cfg, err }

  if fileCfg.UploadDir != "" { cfg.UploadDir = fileCfg.UploadDir }
  if fileCfg.OutputExcel != "" { cfg.OutputExcel = fileCfg.OutputExcel }
  if fileCfg.ExcelSheet != "" { cfg.ExcelSheet = fileCfg.ExcelSheet }
  if fileCfg.Web.Address != "" { cfg.Web.Address = fileCfg.Web.Address }

  if v := fileCfg.Patterns.ContractDatePattern; v != "" { cfg.Patterns.ContractDatePattern = v }
  if v := fileCfg.Patterns.ActDatePattern; v != "" { cfg.Patterns.ActDatePattern = v }
  if v := fileCfg.Patterns.DateRegex; v != "" { cfg.Patterns.DateRegex = v }
  if v := fileCfg.Patterns.VINPattern; v != "" { cfg.Patterns.VINPattern = v }
  if v := fileCfg.Patterns.VINRegex; v != "" { cfg.Patterns.VINRegex = v }
  if v := fileCfg.Patterns.VehicleModelPattern; v != "" { cfg.Patterns.VehicleModelPattern = v }
  if v := fileCfg.Patterns.SellerCompanyDKPPattern; v != "" { cfg.Patterns.SellerCompanyDKPPattern = v }

  return cfg, nil
} 