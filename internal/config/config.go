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

  // New patterns
  AmountPattern       string `json:"amountPattern"`
  AmountRegex         string `json:"amountRegex"`
  VINPattern          string `json:"vinPattern"`
  VINRegex            string `json:"vinRegex"`
  PtsDatePattern      string `json:"ptsDatePattern"`
  BuyerCompanyPattern string `json:"buyerCompanyPattern"`
  VehicleModelPattern string `json:"vehicleModelPattern"`
  SellerCompanyDKPPattern string `json:"sellerCompanyDKPPattern"`
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
      DateRegex:           `(?m)(\b\d{2}[\./-]\d{2}[\./-]\d{4}\b)`,
      NameRegex:           `(?m)([А-ЯЁ][а-яё]+(?:\s+[А-ЯЁ][а-яё]+){0,2}(?:\s+[А-ЯЁ]\.[А-ЯЁ]\.)?)`,

      AmountPattern:       `(?i)общая\s+сумма\s+договора[^\n\r\d]{0,80}`,
      AmountRegex:         `(?m)(\b\d{1,3}(?:[\s\u00a0]\d{3})*(?:[\.,]\d{2})\b)`,
      VINPattern:          `(?i)vin[:\s-]*`,
      VINRegex:            `(?i)\b[ABCDEFGHJKLMNPRSTUVWXYZ\d]{17}\b`,
      PtsDatePattern:      `(?i)паспорт\s+транспортного\s+средства[^\n\r\d]{0,80}(?:выдан|от)`,
      BuyerCompanyPattern: `(?i)(покупатель|покупателя)[^\n\r]{0,120}`,
      VehicleModelPattern: `(?i)(модель|model)[:\s-]*([^\n\r]{1,80})`,
      SellerCompanyDKPPattern: `(?i)(продавец|продавца)[^\n\r]{0,10}[:\s-]*([^\n\r]{1,160})`,
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

  // Merge simple fields
  if fileCfg.UploadDir != "" { cfg.UploadDir = fileCfg.UploadDir }
  if fileCfg.OutputExcel != "" { cfg.OutputExcel = fileCfg.OutputExcel }
  if fileCfg.ExcelSheet != "" { cfg.ExcelSheet = fileCfg.ExcelSheet }
  if fileCfg.Web.Address != "" { cfg.Web.Address = fileCfg.Web.Address }

  // Merge patterns selectively if provided
  if v := fileCfg.Patterns.ContractDatePattern; v != "" { cfg.Patterns.ContractDatePattern = v }
  if v := fileCfg.Patterns.ActDatePattern; v != "" { cfg.Patterns.ActDatePattern = v }
  if v := fileCfg.Patterns.BuyerLabelPattern; v != "" { cfg.Patterns.BuyerLabelPattern = v }
  if v := fileCfg.Patterns.SellerLabelPattern; v != "" { cfg.Patterns.SellerLabelPattern = v }
  if v := fileCfg.Patterns.PoAFromPattern; v != "" { cfg.Patterns.PoAFromPattern = v }
  if v := fileCfg.Patterns.PoAToPattern; v != "" { cfg.Patterns.PoAToPattern = v }
  if v := fileCfg.Patterns.DateRegex; v != "" { cfg.Patterns.DateRegex = v }
  if v := fileCfg.Patterns.NameRegex; v != "" { cfg.Patterns.NameRegex = v }

  if v := fileCfg.Patterns.AmountPattern; v != "" { cfg.Patterns.AmountPattern = v }
  if v := fileCfg.Patterns.AmountRegex; v != "" { cfg.Patterns.AmountRegex = v }
  if v := fileCfg.Patterns.VINPattern; v != "" { cfg.Patterns.VINPattern = v }
  if v := fileCfg.Patterns.VINRegex; v != "" { cfg.Patterns.VINRegex = v }
  if v := fileCfg.Patterns.PtsDatePattern; v != "" { cfg.Patterns.PtsDatePattern = v }
  if v := fileCfg.Patterns.BuyerCompanyPattern; v != "" { cfg.Patterns.BuyerCompanyPattern = v }
  if v := fileCfg.Patterns.VehicleModelPattern; v != "" { cfg.Patterns.VehicleModelPattern = v }
  if v := fileCfg.Patterns.SellerCompanyDKPPattern; v != "" { cfg.Patterns.SellerCompanyDKPPattern = v }

  return cfg, nil
} 