package structs

import "github.com/jinzhu/gorm"

type Orders struct {
	gorm.Model
	Customer_Name string
	Ordered_At string
}
type Items struct {
	gorm.Model
	Code string
	Description string
	Quantity int
	Order_Id int
}