package models

import (
	"log"
	"server/config"
	"time"

	"gorm.io/gorm"
)

// TODO: Add foreign key logic to Service model
// GORM model for all Service records in the database
type Service struct {
	gorm.Model
	BusinessID    uint      `gorm:"column:business_id" json:"business_id"`         // ID of Business that Service is associated with
	Name          string    `gorm:"column:name" json:"name"`                       // Service name
	Description   string    `gorm:"column:desc" json:"desc"`                       // Service description
	StartDateTime time.Time `gorm:"column:start_date_time" json:"start_date_time"` // Date/time that the service starts
	Length        uint      `gorm:"column:length" json:"length"`                   // Length of time in minutes that the service will take
	Capacity      uint      `gorm:"column:capacity" json:"capacity"`               // Number of users that can sign up for the service
	CancelFee     uint      `gorm:"column:cancel_fee" json:"cancel_fee"`           // Fee (in cents) for cancelling appointment after minimum notice cutoff
	Price         uint      `gorm:"column:price" json:"price"`                     // Price (in cents) for the service being offered
}

/*
*Description*

func GetID

# Returns ID field from Service object

*Parameters*

	N/A (None)

*Returns*

	_  <uint>

		The ID of the service object
*/
func (service *Service) GetID() uint {
	return service.ID
}

/*
*Description*

func Create

Creates a new Service record in the database and returns the created record along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance where the record will be created.

*Returns*

	_  <*Service>

		The created Service record.

	_  <error>

		Encountered error (nil if no errors are encountered).
*/
func (service *Service) Create(db *gorm.DB) (map[string]Model, error) {
	// TODO: Add field validation logic (func Create) -- add as BeforeCreate gorm hook definition at the top of this file
	err := db.Create(&service).Error
	returnRecords := map[string]Model{"service": service}
	return returnRecords, err
}

/*
*Description*

func Get

Retrieves a Service record in the database by ID if it exists and returns that record along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to retrieve the specified record.

	serviceID  <uint>

		The ID of the service record being requested.

*Returns*

	_  <*Service>

		The Service record that is retrieved from the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (service *Service) Get(db *gorm.DB, serviceID uint) (map[string]Model, error) {
	err := db.First(&service, serviceID).Error
	returnRecords := map[string]Model{"service": service}
	return returnRecords, err
}

/*
*Description*

func Update

Updates the specified Service record in the database with the specified changes if the record exists.

Returns the updated record along with any errors that are thrown.

This function behaves like a PATCH method, rather than a true PUT. Any fields that aren't specified in the request body for the PUT request will not be altered for the specified record.

If a specified field's value should be deleted from the record, the appropriate null/blank should be specified for that key in the JSON request body (e.g. "type": "").

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to retrieve and update the specified record.

	serviceID  <uint>

		The ID of the service record being updated.

	updates  <map[string]interface{}>

		JSON with the fields that will be updated as keys and the updated values as values.

		Ex:
			{
				"name": "New name",
				"address": "New address"
			}

*Returns*

	_  <*Service>

		The Service record that is updated in the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (service *Service) Update(db *gorm.DB, serviceID uint, updates map[string]interface{}) (map[string]Model, error) {
	// Confirm serviceID exists in the database and get current object
	returnRecords, err := service.Get(db, serviceID)
	updateService := returnRecords["service"]

	if err != nil {
		return returnRecords, err
	}

	// TODO: Add field validation logic (func Update) -- add as BeforeUpdate gorm hook definition at the top of this file

	err = db.Model(&updateService).Where("id = ?", serviceID).Updates(updates).Error
	returnRecords = map[string]Model{"service": updateService}

	return returnRecords, err
}

// TODO: Cascade delete all records associated with service (ServiceOfferings, etc.)
/*
*Description*

func Delete

Deletes the specified Service record from the database if it exists.

Deleted record is returned along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance where the record will be created.

	serviceID  <uint>

		The ID of the service record being deleted.

*Returns*

	_  <*Service>

		The deleted Service record.

	_  <error>

		Encountered error (nil if no errors are encountered).

*/
func (service *Service) Delete(db *gorm.DB, serviceID uint) (map[string]Model, error) {
	// Confirm serviceID exists in the database and get current object
	returnRecords, err := service.Get(db, serviceID)
	deleteService := returnRecords["service"]

	if err != nil {
		return returnRecords, err
	}

	if config.Debug {
		log.Printf("\n\nService object targeted for deletion:\n\n%+v\n\n", deleteService)
	}

	// TODO:  Extend delete operations to all of the other object types associated with the Service record as is appropriate (ServiceOfferings, etc.)
	err = db.Delete(deleteService).Error
	returnRecords = map[string]Model{"service": deleteService}

	return returnRecords, err
}
