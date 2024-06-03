package main

import (
  "database/sql"
  "flag"
  "fmt"
  "net/url"
  "strings"

  "github.com/cloudwego/hertz/pkg/common/hlog"
  _ "github.com/go-sql-driver/mysql"
  _ "github.com/lib/pq"
  "github.com/xuri/excelize/v2"
)

var (
  dbDSN         string
  excelFilePath string
  tableName     string
)

func init() {
  flag.StringVar(&dbDSN, "dsn", "", "Database DSN")
  flag.StringVar(&excelFilePath, "excel", "", "Excel file path")
  flag.StringVar(&tableName, "table", "", "Database table name")
  flag.Parse()
}

func main() {
  hlog.Info("start")

  if dbDSN == "" || excelFilePath == "" || tableName == "" {
    hlog.Fatalf("Database DSN, table name and Excel file path are required")
  }

  dbType, err := getDBType(dbDSN)
  if err != nil {
    hlog.Fatalf("Failed to determine database type: %v", err)
  }

  if dbType == "postgres" {
    dbDSN = addSSLModeDisable(dbDSN)
  }

  db, err := sql.Open(dbType, dbDSN)
  if err != nil {
    hlog.Fatalf("Failed to connect to database: %v", err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    hlog.Fatalf("Failed to ping database: %v", err)
  }

  err = readAndInsertData(db, dbType, excelFilePath)
  if err != nil {
    hlog.Fatalf("Failed to read and insert data: %v", err)
  }

  hlog.Info("Data insertion completed")
}

func getDBType(dsn string) (string, error) {
  u, err := url.Parse(dsn)
  if err != nil {
    return "", fmt.Errorf("failed to parse DSN: %v", err)
  }

  switch u.Scheme {
  case "mysql":
    return "mysql", nil
  case "postgresql":
    return "postgres", nil
  default:
    return "", fmt.Errorf("unsupported database type: %s", u.Scheme)
  }
}

func addSSLModeDisable(dsn string) string {
  u, err := url.Parse(dsn)
  if err != nil {
    hlog.Fatalf("Failed to parse DSN: %v", err)
  }

  query := u.Query()
  query.Set("sslmode", "disable")
  u.RawQuery = query.Encode()

  return u.String()
}

func readAndInsertData(db *sql.DB, dbType, filePath string) error {
  f, err := excelize.OpenFile(filePath)
  if err != nil {
    return fmt.Errorf("failed to open excel file: %v", err)
  }
  defer f.Close()

  sheetName := f.GetSheetName(0)
  rows, err := f.GetRows(sheetName)
  if err != nil {
    return fmt.Errorf("failed to get rows: %v", err)
  }

  if len(rows) == 0 {
    return fmt.Errorf("no data found in the sheet")
  }

  headers := rows[0]
  for i, row := range rows {
    if i == 0 {
      continue // Skip header row
    }

    query, args := buildInsertQuery(dbType, headers, row)
    if query == "" {
      continue // Skip if query is empty
    }

    hlog.Infof("query:%s", query)
    hlog.Infof("args:%v", args)

    _, err := db.Exec(query, args...)
    if err != nil {
      hlog.Errorf("Failed to execute query: %v", err)
    }
  }

  return nil
}

func buildInsertQuery(dbType string, headers, row []string) (string, []interface{}) {
  var columns []string
  var placeholders []string
  var args []interface{}
  rowLen := len(row)
  for i, header := range headers {
    if i < rowLen && row[i] != "" {
      columns = append(columns, header)
      if dbType == "postgres" {
        placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
      } else {
        placeholders = append(placeholders, "?")
      }
      args = append(args, row[i])
    }
  }

  if len(columns) == 0 {
    return "", nil
  }

  query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
  return query, args
}
