package pdf

import (
  "bytes"
  "errors"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "runtime"
  "strings"
  "strconv"

  appcfg "filebot/internal/config"

  pdf "github.com/ledongthuc/pdf"
)

// ExtractText returns plain text content of a PDF file.
// If OCR is enabled and no text is extracted, attempts OCR via pdftoppm + tesseract.
func ExtractText(path string, cfg appcfg.Config) (string, error) {
  // Try text layer first
  if text, err := extractTextFromPDF(path); err == nil && strings.TrimSpace(text) != "" {
    return text, nil
  }
  if !cfg.Ocr.Enabled {
    // Return any text we got (possibly empty) or the error from first attempt
    return extractTextFromPDF(path)
  }
  // OCR fallback
  return ocrPDF(path, cfg)
}

func extractTextFromPDF(path string) (string, error) {
  f, r, err := pdf.Open(path)
  if err != nil {
    return "", err
  }
  defer f.Close()

  var buf bytes.Buffer
  rc, err := r.GetPlainText()
  if err != nil {
    return "", err
  }
  if _, err := io.Copy(&buf, rc); err != nil {
    return "", err
  }
  return buf.String(), nil
}

func ocrPDF(path string, cfg appcfg.Config) (string, error) {
  // Use pdftoppm to render pages to PNGs
  pdftoppm := cfg.Ocr.PdftoppmPath
  if pdftoppm == "" {
    pdftoppm = "pdftoppm"
    if runtime.GOOS == "windows" {
      pdftoppm = "pdftoppm.exe"
    }
  }

  tmpDir, err := os.MkdirTemp("", "filebot_ocr_*")
  if err != nil { return "", err }
  defer os.RemoveAll(tmpDir)

  base := filepath.Join(tmpDir, "page")
  dpi := cfg.Ocr.Dpi
  if dpi <= 0 { dpi = 300 }

  cmd := exec.Command(pdftoppm, "-r",  strconv.Itoa(dpi), "-png", path, base)
  // For Windows PATH resolution, rely on environment
  out, err := cmd.CombinedOutput()
  if err != nil {
    return "", errors.New("pdftoppm error: " + string(out))
  }

  // Collect generated PNGs
  entries, err := os.ReadDir(tmpDir)
  if err != nil { return "", err }

  var pages []string
  for _, e := range entries {
    name := e.Name()
    if strings.HasPrefix(name, "page-") && strings.HasSuffix(name, ".png") {
      pages = append(pages, filepath.Join(tmpDir, name))
    }
  }
  if len(pages) == 0 {
    return "", errors.New("no images produced by pdftoppm")
  }

  // OCR each page with tesseract
  tesseract := "tesseract"
  if runtime.GOOS == "windows" {
    tesseract = "tesseract.exe"
  }
  var b strings.Builder
  for _, img := range pages {
    // tesseract <img> stdout -l rus+eng
    args := []string{img, "stdout", "-l", cfg.Ocr.Lang}
    cmd := exec.Command(tesseract, args...)
    out, err := cmd.CombinedOutput()
    if err != nil {
      // continue but record error context
      b.WriteString("\n[OCR error: ")
      b.WriteString(filepath.Base(img))
      b.WriteString(": ")
      b.WriteString(string(out))
      b.WriteString("]\n")
      continue
    }
    b.Write(out)
    b.WriteString("\n\n")
  }
  return b.String(), nil
} 