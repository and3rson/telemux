module github.com/and3rson/telemux/gormpersistence

go 1.16

replace github.com/and3rson/telemux => ../

require (
	github.com/and3rson/telemux v1.3.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.6
)
