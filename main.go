package main

import (
  "bytes"
  "errors"
  "fmt"
  "io"
  "log"
  "mime/multipart"
  "net/http"
  "os"
  "path/filepath"
  "regexp"
  "sort"
  "strings"
  "time"

  pdf "github.com/ledongthuc/pdf"
  "github.com/xuri/excelize/v2"
)

// ExtractedInfo holds the structured data we want to capture from PDFs.
type ExtractedInfo struct {
  SourceFiles              []string
  ContractDate             string
  ActDate                  string
  BuyerRepresentative      string
  SellerRepresentative     string
  BuyerPoAFrom             string
  BuyerPoATo               string
  SellerPoAFrom            string
  SellerPoATo              string
}

const (
  uploadDir   = "uploads"
  outputExcel = "output.xlsx"
  excelSheet  = "Data"
)

func main() {
  if err := os.MkdirAll(uploadDir, 0755); err != nil {
    log.Fatalf("cannot create upload dir: %v", err)
  }

  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/upload", handleUpload)
  http.HandleFunc("/download", handleDownload)
  http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("ok")) })

  addr := ":8080"
  log.Printf("starting server on %s", addr)
  if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatalf("server error: %v", err)
  }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  _, _ = w.Write([]byte(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>FileBot — загрузка PDF</title>
  <style>
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, sans-serif; margin: 40px; }
    .card { max-width: 760px; padding: 24px; border: 1px solid #e5e7eb; border-radius: 12px; box-shadow: 0 2px 6px rgba(0,0,0,0.06); }
    h1 { margin-top: 0; font-size: 22px; }
    input[type=file] { display: block; margin: 12px 0; }
    button { background: #111827; color: white; border: 0; border-radius: 8px; padding: 10px 16px; cursor: pointer; }
    button:hover { background: #0b1220; }
    .hint { color: #6b7280; font-size: 13px; }
    .row { margin-top: 12px; }
    .link { margin-left: 12px; }
  </style>
</head>
<body>
  <div class="card">
    <h1>Загрузка PDF документов</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
      <label>Выберите PDF файлы (можно несколько)</label>
      <input type="file" name="files" accept="application/pdf" multiple required />
      <div class="row">
        <button type="submit">Обработать и записать в Excel</button>
        <span class="hint">Результат будет сохранён в output.xlsx</span>
      </div>
    </form>
    <div class="row">
      <a class="link" href="/download">Скачать текущий Excel</a>
    </div>
  </div>
</body>
</html>`))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
  if _, err := os.Stat(outputExcel); err != nil {
    http.Error(w, "файл Excel ещё не создан", http.StatusNotFound)
    return
  }
  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", "attachment; filename=output.xlsx")
  http.ServeFile(w, r, outputExcel)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }

  if err := r.ParseMultipartForm(64 << 20); err != nil { // 64MB
    http.Error(w, fmt.Sprintf("cannot parse form: %v", err), http.StatusBadRequest)
    return
  }

  files := r.MultipartForm.File["files"]
  if len(files) == 0 {
    http.Error(w, "не выбраны файлы", http.StatusBadRequest)
    return
  }

  var sourceNames []string
  var combinedText strings.Builder

  for _, fh := range files {
    if !strings.EqualFold(filepath.Ext(fh.Filename), ".pdf") {
      continue
    }

    savedPath, err := saveUploadedFile(fh)
    if err != nil {
      http.Error(w, fmt.Sprintf("ошибка сохранения %s: %v", fh.Filename, err), http.StatusInternalServerError)
      return
    }

    text, err := extractTextFromPDF(savedPath)
    if err != nil {
      log.Printf("extract text error for %s: %v", fh.Filename, err)
    }

    if strings.TrimSpace(text) == "" {
      log.Printf("warning: no extractable text in %s (возможен скан-образ)", fh.Filename)
    }

    sourceNames = append(sourceNames, fh.Filename)
    combinedText.WriteString("\n\n==== ")
    combinedText.WriteString(fh.Filename)
    combinedText.WriteString(" ====\n")
    combinedText.WriteString(text)
  }

  info := parseFields(combinedText.String())
  info.SourceFiles = sourceNames

  if err := appendToExcel(info); err != nil {
    http.Error(w, fmt.Sprintf("не удалось записать в Excel: %v", err), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  _, _ = w.Write([]byte(`<div style="font-family:system-ui,-apple-system,Segoe UI,Roboto,sans-serif; padding:24px">` +
    `<h3>Готово</h3>` +
    `<p>Данные добавлены в <code>output.xlsx</code>.</p>` +
    `<p><a href="/download">Скачать Excel</a> · <a href="/">Назад</a></p>` +
    `</div>`))
}

func saveUploadedFile(fh *multipart.FileHeader) (string, error) {
  f, err := fh.Open()
  if err != nil {
    return "", err
  }
  defer f.Close()

  dst := filepath.Join(uploadDir, sanitizeFilename(fh.Filename))
  out, err := os.Create(dst)
  if err != nil {
    return "", err
  }
  defer out.Close()

  if _, err := io.Copy(out, f); err != nil {
    return "", err
  }
  return dst, nil
}

func sanitizeFilename(name string) string {
  name = strings.ReplaceAll(name, "\\", "_")
  name = strings.ReplaceAll(name, "/", "_")
  name = strings.ReplaceAll(name, ":", "_")
  name = strings.ReplaceAll(name, "..", ".")
  return name
}

func extractTextFromPDF(path string) (string, error) {
  file, reader, err := pdf.Open(path)
  if err != nil {
    return "", err
  }
  defer file.Close()

  var buf bytes.Buffer
  rc, err := reader.GetPlainText()
  if err != nil {
    return "", err
  }
  if _, err := io.Copy(&buf, rc); err != nil {
    return "", err
  }
  return buf.String(), nil
}

func parseFields(text string) ExtractedInfo {
  t := normalizeSpaces(text)

  contractDate := findDateNear(t, `(?i)договор[^\n\r\d]{0,80}`)
  actDate := findDateNear(t, `(?i)акт[^\n\r\d]{0,80}`)

  buyerRep := findNameAfter(t, `(?i)покупател[ья]:?`)
  sellerRep := findNameAfter(t, `(?i)продавц[а|е]:?`)

  buyerPoAFrom := findDateNear(t, `(?i)доверенн(?:ость|ости|остью)[^\n\r]{0,80}от`)
  buyerPoATo := findDateNear(t, `(?i)действи[еия][^\n\r]{0,80}до`)
  sellerPoAFrom := findDateNear(t, `(?i)доверенн(?:ость|ости|остью)[^\n\r]{0,80}от`)
  sellerPoATo := findDateNear(t, `(?i)действи[еия][^\n\r]{0,80}до`)

  return ExtractedInfo{
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

var (
  reDate = regexp.MustCompile(`(?m)(\b\d{2}[\./-]\d{2}[\./-]\d{4}\b)`) // dd.mm.yyyy or dd/mm/yyyy

  // Name pattern: simple heuristic for Russian full name or initials
  reName = regexp.MustCompile(`(?m)([А-ЯЁ][а-яё]+(?:\s+[А-ЯЁ][а-яё]+){0,2}(?:\s+[А-ЯЁ]\.[А-ЯЁ]\.)?)`)
)

func findDateNear(text string, beforePattern string) string {
  re := regexp.MustCompile(beforePattern + `[^\n\r\d]{0,40}` + reDate.String())
  if m := re.FindStringSubmatch(text); len(m) > 0 {
    d := reDate.FindString(m[0])
    if d != "" {
      return normalizeDate(d)
    }
  }
  // fallback: first date in text
  d := reDate.FindString(text)
  return normalizeDate(d)
}

func findNameAfter(text string, labelPattern string) string {
  re := regexp.MustCompile(labelPattern + `[^\n\r:]{0,40}[:\-–]?\s*` + reName.String())
  if m := re.FindStringSubmatch(text); len(m) >= 2 {
    name := strings.TrimSpace(m[1])
    // Avoid capturing the label itself if regex misfires
    if !regexp.MustCompile(`(?i)покупател|продав`).MatchString(name) {
      return name
    }
  }
  // try next names in the same line after label
  reLine := regexp.MustCompile(labelPattern + `([^\n\r]{0,120})`)
  if m := reLine.FindStringSubmatch(text); len(m) >= 2 {
    if nm := reName.FindString(m[1]); nm != "" {
      return strings.TrimSpace(nm)
    }
  }
  return ""
}

func normalizeSpaces(s string) string {
  s = strings.ReplaceAll(s, "\u00a0", " ")
  s = strings.ReplaceAll(s, "\t", " ")
  s = strings.ReplaceAll(s, "\r", " ")
  s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
  return s
}

func normalizeDate(s string) string {
  if s == "" {
    return ""
  }
  s = strings.ReplaceAll(s, "/", ".")
  s = strings.ReplaceAll(s, "-", ".")
  parts := strings.Split(s, ".")
  if len(parts) != 3 {
    return s
  }
  // zero-pad day and month if needed
  day := parts[0]
  mon := parts[1]
  year := parts[2]
  if len(day) == 1 {
    day = "0" + day
  }
  if len(mon) == 1 {
    mon = "0" + mon
  }
  // attempt to parse and reformat to dd.mm.yyyy
  if t, err := time.Parse("02.01.2006", fmt.Sprintf("%s.%s.%s", day, mon, year)); err == nil {
    return t.Format("02.01.2006")
  }
  return fmt.Sprintf("%s.%s.%s", day, mon, year)
}

func appendToExcel(info ExtractedInfo) error {
  var f *excelize.File
  var err error

  if _, err = os.Stat(outputExcel); errors.Is(err, os.ErrNotExist) {
    f = excelize.NewFile()
    idx, _ := f.NewSheet(excelSheet)
    f.SetActiveSheet(idx)
    // headers
    headers := []string{
      "Источник файлов",
      "Дата договора",
      "Дата акта",
      "Представитель покупателя",
      "Представитель продавца",
      "Доверенность покупателя от",
      "Доверенность покупателя до",
      "Доверенность продавца от",
      "Доверенность продавца до",
    }
    for i, h := range headers {
      cell, _ := excelize.CoordinatesToCellName(i+1, 1)
      _ = f.SetCellValue(excelSheet, cell, h)
    }
  } else {
    f, err = excelize.OpenFile(outputExcel)
    if err != nil {
      return err
    }
    // ensure sheet exists
    found := false
    for _, s := range f.GetSheetList() {
      if s == excelSheet {
        found = true
        break
      }
    }
    if !found {
      idx, _ := f.NewSheet(excelSheet)
      f.SetActiveSheet(idx)
    }
  }

  defer func() { _ = f.Close() }()

  rows, err := f.GetRows(excelSheet)
  if err != nil {
    return err
  }
  nextRow := len(rows) + 1

  // Join and sort source file names for stable output
  sorted := append([]string(nil), info.SourceFiles...)
  sort.Strings(sorted)
  source := strings.Join(sorted, "; ")

  values := []any{
    source,
    info.ContractDate,
    info.ActDate,
    info.BuyerRepresentative,
    info.SellerRepresentative,
    info.BuyerPoAFrom,
    info.BuyerPoATo,
    info.SellerPoAFrom,
    info.SellerPoATo,
  }

  for i, v := range values {
    cell, _ := excelize.CoordinatesToCellName(i+1, nextRow)
    if err := f.SetCellValue(excelSheet, cell, v); err != nil {
      return err
    }
  }

  if err := f.SaveAs(outputExcel); err != nil {
    return err
  }
  return nil
}