package parse

import (
  "fmt"
  "regexp"
  "strings"
  "time"

  "filebot/internal/config"
  "filebot/internal/models"
)

// ParseFields extracts fields from combined text according to configured patterns.
func ParseFields(text string, cfg config.Config) models.ExtractedInfo {
  t := NormalizeSpaces(text)

  reDate := regexp.MustCompile(cfg.Patterns.DateRegex)
  reName := regexp.MustCompile(cfg.Patterns.NameRegex)
  reAmount := regexp.MustCompile(cfg.Patterns.AmountRegex)
  reVIN := regexp.MustCompile(cfg.Patterns.VINRegex)

  contractDate := findDateNear(t, cfg.Patterns.ContractDatePattern, reDate)
  actDate := findDateNear(t, cfg.Patterns.ActDatePattern, reDate)

  buyerRep := findNameAfter(t, cfg.Patterns.BuyerLabelPattern, reName)
  sellerRep := findNameAfter(t, cfg.Patterns.SellerLabelPattern, reName)

  buyerPoAFrom := findDateNear(t, cfg.Patterns.PoAFromPattern, reDate)
  buyerPoATo := findDateNear(t, cfg.Patterns.PoAToPattern, reDate)
  sellerPoAFrom := findDateNear(t, cfg.Patterns.PoAFromPattern, reDate)
  sellerPoATo := findDateNear(t, cfg.Patterns.PoAToPattern, reDate)

  // New fields
  amount := findAfterWithRegex(t, cfg.Patterns.AmountPattern, reAmount)
  vin := findAfterWithRegex(t, cfg.Patterns.VINPattern, reVIN)
  ptsDate := findDateNear(t, cfg.Patterns.PtsDatePattern, reDate)
  buyerCompany := findCompanyAfter(t, cfg.Patterns.BuyerCompanyPattern)

  return models.ExtractedInfo{
    ContractDate:             contractDate,
    ActDate:                  actDate,
    BuyerRepresentative:      buyerRep,
    SellerRepresentative:     sellerRep,
    BuyerPoAFrom:             buyerPoAFrom,
    BuyerPoATo:               buyerPoATo,
    SellerPoAFrom:            sellerPoAFrom,
    SellerPoATo:              sellerPoATo,
    ContractAmountRubles:     amount,
    VehiclePassportIssueDate: ptsDate,
    VIN:                      vin,
    BuyerCompany:             buyerCompany,
  }
}

// ParseFirstFileFields focuses on extracting items that must come from the first file only.
func ParseFirstFileFields(firstFileText string, cfg config.Config) (vin string, vehicleModel string, contractDate string, actDate string, sellerCompanyDKP string) {
  t := NormalizeSpaces(firstFileText)

  reDate := regexp.MustCompile(cfg.Patterns.DateRegex)
  reVIN := regexp.MustCompile(cfg.Patterns.VINRegex)

  vin = findAfterWithRegex(t, cfg.Patterns.VINPattern, reVIN)
  contractDate = findDateNear(t, cfg.Patterns.ContractDatePattern, reDate)
  actDate = findDateNear(t, cfg.Patterns.ActDatePattern, reDate)

  // Vehicle model: capture second group of VehicleModelPattern if present
  if m := regexp.MustCompile(cfg.Patterns.VehicleModelPattern).FindStringSubmatch(t); len(m) >= 3 {
    vehicleModel = strings.TrimSpace(m[2])
  }
  // Seller company in DKP: capture second group
  if m := regexp.MustCompile(cfg.Patterns.SellerCompanyDKPPattern).FindStringSubmatch(t); len(m) >= 3 {
    sellerCompanyDKP = strings.TrimSpace(cleanCompanyTail(m[2]))
  }
  return
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

func findNameAfter(text string, labelPattern string, reName *regexp.Regexp) string {
  re := regexp.MustCompile(labelPattern + `[^\n\r:]{0,80}[:\-–]?\s*` + reName.String())
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    name := strings.TrimSpace(m[1])
    if !regexp.MustCompile(`(?i)покупател|продав`).MatchString(name) {
      return name
    }
  }
  reLine := regexp.MustCompile(labelPattern + `([^\n\r]{0,200})`)
  if m := reLine.FindStringSubmatch(text); len(m) >= 2 {
    if nm := reName.FindString(m[1]); nm != "" {
      return strings.TrimSpace(nm)
    }
  }
  return ""
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

func findCompanyAfter(text string, labelPattern string) string {
  re := regexp.MustCompile(labelPattern + `[:\s\-–]*([^\n\r]{0,200})`)
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    line := strings.TrimSpace(m[1])
    line = strings.Trim(line, ":;–- ")
    stopIdx := len(line)
    for _, stop := range []string{"ИНН", "КПП", "ОГРН", "адрес", "тел", "e-mail", "Эл. почта"} {
      if idx := strings.Index(strings.ToLower(line), strings.ToLower(stop)); idx > 0 && idx < stopIdx {
        stopIdx = idx
      }
    }
    if stopIdx > 0 && stopIdx <= len(line) {
      return strings.TrimSpace(line[:stopIdx])
    }
    return line
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
  if s == "" {
    return ""
  }
  s = strings.ReplaceAll(s, "/", ".")
  s = strings.ReplaceAll(s, "-", ".")
  parts := strings.Split(s, ".")
  if len(parts) != 3 {
    return s
  }
  day := parts[0]
  mon := parts[1]
  year := parts[2]
  if len(day) == 1 {
    day = "0" + day
  }
  if len(mon) == 1 {
    mon = "0" + mon
  }
  if t, err := time.Parse("02.01.2006", fmt.Sprintf("%s.%s.%s", day, mon, year)); err == nil {
    return t.Format("02.01.2006")
  }
  return fmt.Sprintf("%s.%s.%s", day, mon, year)
} 