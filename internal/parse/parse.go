package parse

import (
  "fmt"
  "regexp"
  "strings"
  "time"

  "filebot/internal/config"
  "filebot/internal/models"
)

// ParseFirstFileOnly extracts VIN, Model, ContractDate, ActDate, SellerCompanyDKP from first file text.
func ParseFirstFileOnly(firstFileText string, cfg config.Config) models.ExtractedInfo {
  t := NormalizeSpaces(firstFileText)

  reDate := regexp.MustCompile(cfg.Patterns.DateRegex)
  reVIN := regexp.MustCompile(cfg.Patterns.VINRegex)

  vin := findAfterWithRegex(t, cfg.Patterns.VINPattern, reVIN)
  contractDate := findDateNear(t, cfg.Patterns.ContractDatePattern, reDate)
  actDate := findDateNear(t, cfg.Patterns.ActDatePattern, reDate)

  var vehicleModel string
  if m := regexp.MustCompile(cfg.Patterns.VehicleModelPattern).FindStringSubmatch(t); len(m) >= 3 {
    vehicleModel = strings.TrimSpace(m[2])
  }

  var sellerCompanyDKP string
  if m := regexp.MustCompile(cfg.Patterns.SellerCompanyDKPPattern).FindStringSubmatch(t); len(m) >= 3 {
    sellerCompanyDKP = strings.TrimSpace(cleanCompanyTail(m[2]))
  }

  return models.ExtractedInfo{
    VIN:              vin,
    VehicleModel:     vehicleModel,
    ContractDate:     contractDate,
    ActDate:          actDate,
    SellerCompanyDKP: sellerCompanyDKP,
  }
}

func cleanCompanyTail(s string) string {
  s = strings.TrimSpace(s)
  stops := []string{"ИНН", "КПП", "ОГРН", "адрес", "тел", "e-mail", "Эл. почта"}
  stopIdx := len(s)
  for _, stop := range stops {
    if idx := strings.Index(strings.ToLower(s), strings.ToLower(stop)); idx > 0 && idx < stopIdx {
      stopIdx = idx
    }
  }
  return strings.TrimSpace(s[:stopIdx])
}

func findDateNear(text string, beforePattern string, reDate *regexp.Regexp) string {
  re := regexp.MustCompile(beforePattern + `[^\n\r\d]{0,60}` + reDate.String())
  if m := re.FindStringSubmatch(text); len(m) > 0 {
    d := reDate.FindString(m[0])
    if d != "" {
      return NormalizeDate(d)
    }
  }
  d := reDate.FindString(text)
  return NormalizeDate(d)
}

func findAfterWithRegex(text string, labelPattern string, target *regexp.Regexp) string {
  re := regexp.MustCompile(labelPattern + `([^\n\r]{0,200})`)
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    if v := target.FindString(m[1]); v != "" {
      return strings.TrimSpace(v)
    }
  }
  if v := target.FindString(text); v != "" {
    return strings.TrimSpace(v)
  }
  return ""
}

func NormalizeSpaces(s string) string {
  s = strings.ReplaceAll(s, "\u00a0", " ")
  s = strings.ReplaceAll(s, "\t", " ")
  s = strings.ReplaceAll(s, "\r", " ")
  s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
  return s
}

func NormalizeDate(s string) string {
  if s == "" { return "" }
  s = strings.ReplaceAll(s, "/", ".")
  s = strings.ReplaceAll(s, "-", ".")
  parts := strings.Split(s, ".")
  if len(parts) != 3 { return s }
  day, mon, year := parts[0], parts[1], parts[2]
  if len(day) == 1 { day = "0" + day }
  if len(mon) == 1 { mon = "0" + mon }
  if t, err := time.Parse("02.01.2006", fmt.Sprintf("%s.%s.%s", day, mon, year)); err == nil {
    return t.Format("02.01.2006")
  }
  return fmt.Sprintf("%s.%s.%s", day, mon, year)
} 