package models

import (
	"github.com/mortezakhademan/paginated-list/mongodb/dataType"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

const (
	COLUMN_INFO_DATA_TYPE_TEXT      = "text"
	COLUMN_INFO_DATA_TYPE_DATE      = "date"
	COLUMN_INFO_DATA_TYPE_INT       = "int"
	COLUMN_INFO_DATA_TYPE_BOOL      = "bool"
	COLUMN_INFO_DATA_TYPE_FLOAT64   = "float64"
	COLUMN_INFO_DATA_TYPE_OBJECT_ID = "objectID"

	COLUMN_INFO_FILTER_OPERATOR_NONE = "-"
)

type ColumnInfo struct {
	///
	Column string
	/// input and output data used ColumnAlias and Column used
	ColumnAlias      string
	FilterOperator   string
	DataType         string
	Title            string
	PrepareValueFunc func(interface{}) interface{}
	ExportSettings   *ExportSettings
}

type ExportSettings struct {
	/// show or hide column in export file
	Show bool
}

func NewColumnInfo(column string, filterOperator string) *ColumnInfo {
	return NewTextColumnInfo(column, filterOperator, column)
}
func newColumnInfo(column string, filterOperator string, title string, dataType string) *ColumnInfo {
	return newColumnInfoWithExportSettings(column, filterOperator, title, dataType, &ExportSettings{Show: true})
}
func newColumnInfoWithExportSettings(column string, filterOperator string, title string, dataType string, exportSettings *ExportSettings) *ColumnInfo {
	return &ColumnInfo{
		Column:         column,
		FilterOperator: filterOperator,
		DataType:       dataType,
		Title:          title,
		ExportSettings: exportSettings,
	}
}
func NewTextColumnInfo(column string, filterOperator string, title string) *ColumnInfo {
	return newColumnInfo(column, filterOperator, title, COLUMN_INFO_DATA_TYPE_TEXT)
}

func NewIntColumnInfo(column string, title string) *ColumnInfo {
	return newColumnInfo(column, "=", title, COLUMN_INFO_DATA_TYPE_INT)
}

func NewIntFilterColumnInfo(column string, filterOperator string, title string) *ColumnInfo {
	return newColumnInfo(column, filterOperator, title, COLUMN_INFO_DATA_TYPE_INT)
}

func NewBoolColumnInfo(column string, title string) *ColumnInfo {
	return newColumnInfo(column, "=", title, COLUMN_INFO_DATA_TYPE_BOOL)
}

func NewFloat64ColumnInfo(column string, title string) *ColumnInfo {
	return newColumnInfo(column, "=", title, COLUMN_INFO_DATA_TYPE_FLOAT64)
}

func NewObjectIDColumnInfo(column string, filterOperator string) *ColumnInfo {
	return newColumnInfoWithExportSettings(column, filterOperator, "", COLUMN_INFO_DATA_TYPE_OBJECT_ID, &ExportSettings{Show: false})
}

func NewDateColumnInfo(column string, title string) *ColumnInfo {
	return newColumnInfo(column, "between", title, COLUMN_INFO_DATA_TYPE_DATE)
}

func (columnInfo *ColumnInfo) SetPrepareValueFunc(prepareValueFunc func(interface{}) interface{}) *ColumnInfo {
	columnInfo.PrepareValueFunc = prepareValueFunc
	return columnInfo
}
func (columnInfo *ColumnInfo) HideInExport() *ColumnInfo {
	columnInfo.ExportSettings.Show = false
	return columnInfo
}

func (columnInfo *ColumnInfo) GetFilterVal(filterVal string) interface{} {
	switch strings.ToLower(columnInfo.FilterOperator) {
	case ">":
		return bson.M{"$gt": columnInfo.convertValueToDataType(filterVal)}
	case ">=":
		return bson.M{"$gte": columnInfo.convertValueToDataType(filterVal)}
	case "<":
		return bson.M{"$lt": columnInfo.convertValueToDataType(filterVal)}
	case "<=":
		return bson.M{"$lte": columnInfo.convertValueToDataType(filterVal)}
	case "=":
		return columnInfo.convertValueToDataType(filterVal)
	case "between":
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(values) == 1 || values[1] == nil {
			return bson.M{"$gte": values[0]}
		} else if values[0] == nil {
			return bson.M{"$lte": values[1]}
		} else {
			return bson.M{
				"$gt": values[0],
				"$lt": values[1],
			}
		}
	case "in":
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(values) == 1 {
			return values[0]
		} else {
			return bson.M{"$in": values}
		}
	case "not in":
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(values) == 1 {
			return bson.M{"$ne": values[0]}
		} else {
			return bson.M{"$nin": values}
		}
	case "like":
		return primitive.Regex{Pattern: filterVal, Options: "i"}
	case "is":
		if filterVal == "0" {
			return bsontype.Null
		}
		return bson.M{"$ne": bsontype.Null}
	case COLUMN_INFO_FILTER_OPERATOR_NONE:
		break
	}
	return nil
}

func (columnInfo *ColumnInfo) convertValueToDataType(value string) interface{} {
	switch columnInfo.DataType {
	case COLUMN_INFO_DATA_TYPE_DATE:
		return dataType.ParseDateTime(value)
	case COLUMN_INFO_DATA_TYPE_OBJECT_ID:
		objectId, _ := primitive.ObjectIDFromHex(value)
		return objectId
	case COLUMN_INFO_DATA_TYPE_INT:
		if value == "" {
			return nil
		}
		i, _ := strconv.Atoi(value)
		if i == 0 {
			return bson.M{"$in": bson.A{nil, 0}}
		}
		return i
	case COLUMN_INFO_DATA_TYPE_FLOAT64:
		if value == "" {
			return nil
		}
		i, _ := strconv.ParseFloat(value, 64)
		return i
	case COLUMN_INFO_DATA_TYPE_BOOL:
		if value == "" {
			return nil
		}
		i, _ := strconv.ParseBool(value)
		if !i {
			return bson.M{"$in": bson.A{nil, false}}
		}
		return i
	}
	return value
}

func (columnInfo *ColumnInfo) convertValueToArrayDataType(value string) []interface{} {

	strValues := strings.Split(strings.TrimSpace(value), ",")
	values := []interface{}{}
	for _, strValue := range strValues {
		values = append(values, columnInfo.convertValueToDataType(strValue))
	}
	return values
}
