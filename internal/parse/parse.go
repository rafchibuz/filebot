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

  contractDate := findDateNear(t, cfg.Patterns.ContractDatePattern, reDate)
  actDate := findDateNear(t, cfg.Patterns.ActDatePattern, reDate)

  buyerRep := findNameAfter(t, cfg.Patterns.BuyerLabelPattern, reName)
  sellerRep := findNameAfter(t, cfg.Patterns.SellerLabelPattern, reName)

  buyerPoAFrom := findDateNear(t, cfg.Patterns.PoAFromPattern, reDate)
  buyerPoATo := findDateNear(t, cfg.Patterns.PoAToPattern, reDate)
  sellerPoAFrom := findDateNear(t, cfg.Patterns.PoAFromPattern, reDate)
  sellerPoATo := findDateNear(t, cfg.Patterns.PoAToPattern, reDate)

  return models.ExtractedInfo{
    ContractDate:         contractDate,
    ActDate:              actDate,
    BuyerRepresentative:  buyerRep,
    SellerRepresentative: sellerRep,
    BuyerPoAFrom:         buyerPoAFrom,
    BuyerPoATo:           buyerPoATo,
    SellerPoAFrom:        sellerPoAFrom,
    SellerPoATo:          sellerPoATo,
  }
}

func findDateNear(text string, beforePattern string, reDate *regexp.Regexp) string {
  re := regexp.MustCompile(beforePattern + `[^\n\r\d]{0,40}` + reDate.String())
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
  re := regexp.MustCompile(labelPattern + `[^\n\r:]{0,40}[:\-–]?\s*` + reName.String())
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    name := strings.TrimSpace(m[1])
    if !regexp.MustCompile(`(?i)покупател|продав`).MatchString(name) {
      return name
    }
  }
  reLine := regexp.MustCompile(labelPattern + `([^\n\r]{0,120})`)
  if m := reLine.FindStringSubmatch(text); len(m) >= 2 {
    if nm := reName.FindString(m[1]); nm != "" {
      return strings.TrimSpace(nm)
    }
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