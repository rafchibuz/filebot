package parse

import (
  "fmt"
  "regexp"
  "strings"
  "time"

  "filebot/internal/config"
  "filebot/internal/models"
)

// ParseFirstFileOnly extracts VIN, CommercialName, ContractDateUS (MM/DD/YYYY), SellerCompanyDKP from first file text.
func ParseFirstFileOnly(firstFileText string, cfg config.Config) models.ExtractedInfo {
  t := NormalizeSpaces(firstFileText)

  reDate := regexp.MustCompile(cfg.Patterns.DateRegex)
  reVIN := regexp.MustCompile(cfg.Patterns.VINRegex)

  vin := findAfterWithRegex(t, cfg.Patterns.VINPattern, reVIN)

  // Contract date: prefer date near context (лизинга/ДКП), else first date
  contractDateRaw := findDateNear(t, cfg.Patterns.ContractDateContextPattern, reDate)
  contractDateUS := toUSDate(contractDateRaw)

  // Commercial name: capture second regex group
  var commercialName string
  if m := regexp.MustCompile(cfg.Patterns.CommercialNamePattern).FindStringSubmatch(t); len(m) >= 3 {
    commercialName = strings.TrimSpace(m[2])
  }

  // Seller company (DKP): capture second group and trim tail
  var sellerCompanyDKP string
  if m := regexp.MustCompile(cfg.Patterns.SellerCompanyDKPPattern).FindStringSubmatch(t); len(m) >= 3 {
    sellerCompanyDKP = strings.TrimSpace(cleanCompanyTail(m[2]))
  }

  return models.ExtractedInfo{
    VIN:              vin,
    CommercialName:   commercialName,
    ContractDateUS:   contractDateUS,
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
  pattern := beforePattern + `[^\n\r\d]{0,120}` + reDate.String()
  if m := regexp.MustCompile(pattern).FindStringSubmatch(text); len(m) > 0 {
    if d := reDate.FindString(m[0]); d != "" { return strings.TrimSpace(d) }
  }
  if d := reDate.FindString(text); d != "" { return strings.TrimSpace(d) }
  return ""
}

func findAfterWithRegex(text string, labelPattern string, target *regexp.Regexp) string {
  re := regexp.MustCompile(labelPattern + `([^\n\r]{0,200})`)
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    if v := target.FindString(m[1]); v != "" {
      return strings.TrimSpace(v)
    }
  }
  if v := target.FindString(text); v != "" { return strings.TrimSpace(v) }
  return ""
}

func NormalizeSpaces(s string) string {
  s = strings.ReplaceAll(s, "\u00a0", " ")
  s = strings.ReplaceAll(s, "\t", " ")
  s = strings.ReplaceAll(s, "\r", " ")
  s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
  return s
}

func toUSDate(s string) string {
  if s == "" { return "" }
  s = strings.TrimSpace(s)

  // Try numeric formats first
  for _, layout := range []string{"02.01.2006", "2.1.2006", "02/01/2006", "2/1/2006", "02-01-2006", "2-1-2006"} {
    if t, err := time.Parse(layout, s); err == nil {
      return t.Format("01/02/2006")
    }
  }

  // Try Russian month form: "23 августа 2024"
  ruMonths := map[string]time.Month{
    "января": 1, "февраля": 2, "марта": 3, "апреля": 4, "мая": 5, "июня": 6,
    "июля": 7, "августа": 8, "сентября": 9, "октября": 10, "ноября": 11, "декабря": 12,
  }
  parts := strings.Fields(s)
  if len(parts) == 3 {
    day := parts[0]
    monName := strings.ToLower(parts[1])
    year := parts[2]
    if mon, ok := ruMonths[monName]; ok {
      // Normalize day to two digits
      if len(day) == 1 { day = "0" + day }
      // Build time and format
      composed := fmt.Sprintf("%s.%02d.%s", day, mon, year)
      if t, err := time.Parse("02.01.2006", composed); err == nil {
        return t.Format("01/02/2006")
      }
    }
  }
  return s // fallback: return as-is
} 