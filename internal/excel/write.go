package excel

import (
  "errors"
  "os"

  "github.com/xuri/excelize/v2"

  "filebot/internal/config"
  "filebot/internal/models"
)

// AppendToExcel ensures the file and sheet exist, and appends a row with info.
func AppendToExcel(info models.ExtractedInfo, cfg config.Config) error {
  var f *excelize.File
  var err error

  if _, err = os.Stat(cfg.OutputExcel); errors.Is(err, os.ErrNotExist) {
    f = excelize.NewFile()
    idx, _ := f.NewSheet(cfg.ExcelSheet)
    f.SetActiveSheet(idx)
    headers := []string{
      "VIN",
      "Модель",
      "Дата ДКП",
      "Дата акта приема-передачи",
      "Наименование продавца (ДКП)",
    }
    for i, h := range headers {
      cell, _ := excelize.CoordinatesToCellName(i+1, 1)
      _ = f.SetCellValue(cfg.ExcelSheet, cell, h)
    }
  } else {
    f, err = excelize.OpenFile(cfg.OutputExcel)
    if err != nil {
      return err
    }
    exists := false
    for _, s := range f.GetSheetList() {
      if s == cfg.ExcelSheet {
        exists = true
        break
      }
    }
    if !exists {
      idx, _ := f.NewSheet(cfg.ExcelSheet)
      f.SetActiveSheet(idx)
    }
  }
  defer func() { _ = f.Close() }()

  rows, err := f.GetRows(cfg.ExcelSheet)
  if err != nil {
    return err
  }
  nextRow := len(rows) + 1

  values := []any{
    info.VIN,
    info.VehicleModel,
    info.ContractDate,
    info.ActDate,
    info.SellerCompanyDKP,
  }

  for i, v := range values {
    cell, _ := excelize.CoordinatesToCellName(i+1, nextRow)
    if err := f.SetCellValue(cfg.ExcelSheet, cell, v); err != nil {
      return err
    }
  }

  if err := f.SaveAs(cfg.OutputExcel); err != nil {
    return err
  }
  return nil
} 