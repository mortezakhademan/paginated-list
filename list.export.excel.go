package componentsList

import (
	"context"
	"fmt"
	"github.com/mortezakhademan/paginated-list/mongodb/models"
	"github.com/tealeg/xlsx"
	ptime "github.com/yaa110/go-persian-calendar"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (list *List) prepareExcelResult(cursor *mongo.Cursor, context context.Context) error {

	list.ExcelFile = xlsx.NewFile()

	sheet, err := list.ExcelFile.AddSheet("sheet 1")
	if err != nil {
		fmt.Printf("can't create sheet in export excel : %s", err)
		return err
	}

	columns := list.setExcelHeaderRow(sheet)

	for cursor.Next(context) {
		documentData := map[string]interface{}{}
		cursor.Decode(documentData)

		list.setExcelDataRow(sheet, columns, documentData)

	}

	return nil
}

func (list *List) setExcelHeaderRow(sheet *xlsx.Sheet) []string {
	headerRow := sheet.AddRow()
	columns := []string{}
	for column, columnInfo := range list.fromCollectionColumnsMap {
		if !columnInfo.ExportSettings.Show {
			continue
		}
		title := columnInfo.Title
		if title == "" {
			title = column
		}
		headerRow.AddCell().Value = title
		columns = append(columns, column)
	}
	return columns
}

func (list *List) setExcelDataRow(sheet *xlsx.Sheet, columns []string, documentData map[string]interface{}) {
	excelRow := sheet.AddRow()

	// Create our map, and retrieve the value for each column from the pointers slice,
	//storing it in the map with the name of the column as the key.
	for _, col := range columns {
		colValStr := ""
		columnInfo := list.fromCollectionColumnsMap[col]
		val, ok := documentData[columnInfo.Column]
		if ok && val != nil {
			if columnInfo.PrepareValueFunc != nil {
				colValStr = fmt.Sprint(columnInfo.PrepareValueFunc(val))
			}
			switch columnInfo.DataType {
			case models.COLUMN_INFO_DATA_TYPE_TEXT:
				if data, ok := val.(primitive.DateTime); ok {
					colValStr = ptime.New(data.Time()).Format("yyyy/MM/dd HH:mm")
				}
			default:

				b, ok := val.([]byte)
				if ok {
					colValStr = string(b)
				} else {
					colValStr = fmt.Sprint(val)
				}
			}
		}
		excelRow.AddCell().Value = fmt.Sprint(colValStr)
	}

}
