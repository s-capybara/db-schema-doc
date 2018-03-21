package dbsd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Auth struct {
	Username string
	Password string
	Database string
}

type ColumnCollection = []string
type RawRecordCollection = []map[string]string
type RecordCollection = [][]string

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Root(auth Auth, table string, columns []string) {
	fullColumns, rawData := ReadData(auth, table)
	data := FormatData(columns, fullColumns, rawData)
	RenderTable(columns, data)
}

func DataSourceName(auth Auth) string {
	return fmt.Sprintf("%s:%s@/%s", auth.Username, auth.Password, auth.Database)
}

func ReadData(auth Auth, table string) (ColumnCollection, RawRecordCollection) {
	db, err := sql.Open("mysql", DataSourceName(auth))
	checkError(err)
	defer db.Close()

	rows, err := db.Query("SHOW FULL COLUMNS FROM " + table)
	checkError(err)

	fullColumns, err := rows.Columns()
	checkError(err)

	records := make(RawRecordCollection, 0)
	values := make([]sql.RawBytes, len(fullColumns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		checkError(err)

		record := make(map[string]string)

		for i, value := range values {
			column := fullColumns[i]
			record[column] = string(value)
		}

		records = append(records, record)
	}

	return fullColumns, records
}

func FormatData(columns ColumnCollection, fullColumns ColumnCollection, rawRecords RawRecordCollection) RecordCollection {
	records := make(RecordCollection, len(rawRecords))

	for i := range records {
		for _, column := range columns {
			records[i] = append(records[i], rawRecords[i][column])
		}
	}

	return records
}

func RenderTable(columns ColumnCollection, records RecordCollection) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(records)
	table.Render()
}
