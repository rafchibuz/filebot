package pdf

import (
  "bytes"
  "io"

  pdf "github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF returns the plain text content of a PDF file.
func ExtractTextFromPDF(path string) (string, error) {
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