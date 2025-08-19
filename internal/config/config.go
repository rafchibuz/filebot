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
  Ocr         Ocr      `json:"ocr"`
}

type Web struct {
  Address string `json:"address"`
}

type Patterns struct {
  ContractDateContextPattern string `json:"contractDateContextPattern"`
  // Date regex that matches dd.mm.yyyy, dd/mm/yyyy, and "23 августа 2024"
  DateRegex string `json:"dateRegex"`

  VINPattern       string `json:"vinPattern"`
  VINRegex         string `json:"vinRegex"`
  CommercialNamePattern string `json:"commercialNamePattern"`
  SellerCompanyDKPPattern string `json:"sellerCompanyDKPPattern"`
}

type Ocr struct {
  Enabled       bool   `json:"enabled"`
  Lang          string `json:"lang"`          // e.g. "rus+eng"
  Dpi           int    `json:"dpi"`           // image DPI for pdftoppm
  PdftoppmPath  string `json:"pdftoppmPath"`  // optional full path to pdftoppm
}

// Default returns the built-in defaults used if config.json is missing or partial.
func Default() Config {
  return Config{
    UploadDir:   "uploads",
    OutputExcel: "output.xlsx",
    ExcelSheet:  "Data",
    Web: Web{ Address: ":8080" },
    Patterns: Patterns{
      // Контекст вокруг даты ДКП (ищем дату рядом с упоминанием договора лизинга или ДКП)
      ContractDateContextPattern: `(?i)(договор[а-я\s-]*лизинга|дкп|договор\s+купли-?продажи)[^\n\r\d]{0,120}`,
      DateRegex: `(?m)(\b\d{1,2}[\./-]\d{1,2}[\./-]\d{2,4}\b|\b\d{1,2}\s+(?:января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября|декабря)\s+\d{4}\b)`,

      VINPattern: `(?i)vin[:\s-]*`,
      VINRegex:   `(?i)\b[ABCDEFGHJKLMNPRSTUVWXYZ\d]{17}\b`,

      CommercialNamePattern: `(?i)(коммерческое\s+наименование)[:\s-]*([^\n\r]{1,120})`,
      SellerCompanyDKPPattern: `(?i)(продавец|продавца)[^\n\r]{0,10}[:\s-]*([^\n\r]{1,200})`,
    },
    Ocr: Ocr{
      Enabled:      false,
      Lang:         "rus+eng",
      Dpi:          300,
      PdftoppmPath: "",
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

  if v := fileCfg.Patterns.ContractDateContextPattern; v != "" { cfg.Patterns.ContractDateContextPattern = v }
  if v := fileCfg.Patterns.DateRegex; v != "" { cfg.Patterns.DateRegex = v }
  if v := fileCfg.Patterns.VINPattern; v != "" { cfg.Patterns.VINPattern = v }
  if v := fileCfg.Patterns.VINRegex; v != "" { cfg.Patterns.VINRegex = v }
  if v := fileCfg.Patterns.CommercialNamePattern; v != "" { cfg.Patterns.CommercialNamePattern = v }
  if v := fileCfg.Patterns.SellerCompanyDKPPattern; v != "" { cfg.Patterns.SellerCompanyDKPPattern = v }

  // Merge OCR
  cfg.Ocr.Enabled = fileCfg.Ocr.Enabled
  if fileCfg.Ocr.Lang != "" { cfg.Ocr.Lang = fileCfg.Ocr.Lang }
  if fileCfg.Ocr.Dpi != 0 { cfg.Ocr.Dpi = fileCfg.Ocr.Dpi }
  if fileCfg.Ocr.PdftoppmPath != "" { cfg.Ocr.PdftoppmPath = fileCfg.Ocr.PdftoppmPath }

  return cfg, nil
} 