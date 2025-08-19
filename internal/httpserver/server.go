package httpserver

import (
  "fmt"
  "html/template"
  "io"
  "mime/multipart"
  "net/http"
  "os"
  "path/filepath"
  "strings"

  appcfg "filebot/internal/config"
  appexcel "filebot/internal/excel"
  appparse "filebot/internal/parse"
  apppdf "filebot/internal/pdf"
)

// Server encapsulates HTTP handling with app config.
type Server struct {
  cfg       appcfg.Config
  tmplIndex *template.Template
}

func New(cfg appcfg.Config) (*Server, error) {
  if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
    return nil, err
  }
  tmpl, err := template.ParseFiles("web/templates/index.html")
  if err != nil {
    return nil, err
  }
  return &Server{cfg: cfg, tmplIndex: tmpl}, nil
}

func (s *Server) Routes(mux *http.ServeMux) {
  mux.HandleFunc("/", s.handleIndex)
  mux.HandleFunc("/upload", s.handleUpload)
  mux.HandleFunc("/download", s.handleDownload)
  mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("ok")) })
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  _ = s.tmplIndex.Execute(w, nil)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
  if _, err := os.Stat(s.cfg.OutputExcel); err != nil {
    http.Error(w, "файл Excel ещё не создан", http.StatusNotFound)
    return
  }
  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", "attachment; filename=output.xlsx")
  http.ServeFile(w, r, s.cfg.OutputExcel)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }

  if err := r.ParseMultipartForm(64 << 20); err != nil {
    http.Error(w, fmt.Sprintf("cannot parse form: %v", err), http.StatusBadRequest)
    return
  }

  files := r.MultipartForm.File["files"]
  if len(files) == 0 {
    http.Error(w, "не выбраны файлы", http.StatusBadRequest)
    return
  }

  var firstText string
  for idx, fh := range files {
    if idx > 0 { break }
    if !strings.EqualFold(filepath.Ext(fh.Filename), ".pdf") { continue }

    savedPath, err := s.saveUploadedFile(fh)
    if err != nil {
      http.Error(w, fmt.Sprintf("ошибка сохранения %s: %v", fh.Filename, err), http.StatusInternalServerError)
      return
    }

    text, err := apppdf.ExtractTextFromPDF(savedPath)
    if err != nil {
      http.Error(w, fmt.Sprintf("ошибка извлечения текста из %s: %v", fh.Filename, err), http.StatusBadRequest)
      return
    }
    firstText = text
  }

  info := appparse.ParseFirstFileOnly(firstText, s.cfg)

  if err := appexcel.AppendToExcel(info, s.cfg); err != nil {
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

func (s *Server) saveUploadedFile(fh *multipart.FileHeader) (string, error) {
  f, err := fh.Open()
  if err != nil { return "", err }
  defer f.Close()

  dst := filepath.Join(s.cfg.UploadDir, sanitizeFilename(fh.Filename))
  out, err := os.Create(dst)
  if err != nil { return "", err }
  defer out.Close()

  if _, err := io.Copy(out, f); err != nil { return "", err }
  return dst, nil
}

func sanitizeFilename(name string) string {
  name = strings.ReplaceAll(name, "\\", "_")
  name = strings.ReplaceAll(name, "/", "_")
  name = strings.ReplaceAll(name, ":", "_")
  name = strings.ReplaceAll(name, "..", ".")
  return name
} 