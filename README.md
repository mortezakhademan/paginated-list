# Paginated List MongoDB 

This repo use for get data from mongoDB as paginated list (pagination + filter + sort)

## How to use:
1. import package
```go
import componentsList "github.com/mortezakhademan/paginated-list"
```
2. prepare list variable:
   1. set filters
   2. set sort field
   3. set result type
   4. set page-no
   5. set page size
```go
    list := componentsList.NewList()
    list.Filters = nil
    list.Sort = nil
    if pageNo > 0 {
        list.PageNo = pageNo
    }
    if pageSize > 0 {
        list.PageSize = pageSize
    }  
    list.ResultType = 0
```
3. add pipeline stages if you need stages (Optional)
```go
list.AddPipelineStage(bson.D{{"$match", bson.D{{"status", bson.M{"$ne": user.USER_STATUS_USER_STATUS_DELETED}}}}})
```
4. call RunQuery function
```go
paginatedItems:= []*ResultModel{}
list.RunQuery(collection, map[string]*models.ColumnInfo{
		"id":              models.NewObjectIDColumnInfo("_id", "in"),
		"type":            models.NewIntColumnInfo("type", "type"),
		"status":          models.NewIntColumnInfo("status", "status"),
		"lastName":        models.NewColumnInfo("last_name", "like"),
		"firstName":       models.NewColumnInfo("first_name", "like"),
		"mobile":          models.NewColumnInfo("mobile.mobile", "like"),
		"createdAt":       models.NewDateColumnInfo("created_at", "created_at"),
	}, "-createdAt", &paginatedItems)
```
Now, list variable contains paginated list info (total items count) + page Items