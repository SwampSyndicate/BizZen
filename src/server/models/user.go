package models

import (
	"errors"
	"fmt"
	"log"
	"server/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TODO: Add foreign key logic to User model
// TODO: Add constraint for AccountType column to limit user types
// GORM model for all User records in the database
type User struct {
	gorm.Model
	Email       string `gorm:"not null;unique;column:email" json:"email"`                             // User's email address
	Password    string `gorm:"not null;column:password" json:"password"`                              // User's hashed password
	AccountType string `gorm:"not null;column:account_type" json:"account_type"`                      // Account type of the User record (Individual, Business, System)
	FirstName   string `gorm:"not null;column:first_name" json:"first_name"`                          // User's first name
	LastName    string `gorm:"not null;column:last_name" json:"last_name"`                            // User's last name
	BusinessID  *uint  `gorm:"column:business_id;default:null" json:"business_id" sql:"DEFAULT:NULL"` // ID of the Business record associated with the User record
}

/*
*Description*

func HashPassword

Generates a hash from the provided password string and assigns it to the calling User's Password attribute.

This conforms to best practice of storing hashed passwords in the application database, rather than plain text.

*Parameters*

	password  <string>

		The plain text password that will be hashed.

*Returns*

	_  <error>

		Encountered error (nil if no errors encountered).
*/
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

/*
*Description*

func CheckPassword

Checks if a given password matches the hashed password associated with the calling User record's account.

This function uses the bcrypt algorithm to compare the given password with the hashed password stored in the calling User struct.

If the given password matches the hashed password, nil is returned.

*Parameters*

	password  <string>

		The password to be checked against the calling User's hashed password.

*Returns*

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

/*
*Description*

func GetID

# Returns ID field from User object

*Parameters*

	N/A (None)

*Returns*

	_  <uint>

		The ID of the User object
*/
func (user *User) GetID() uint {
	return user.ID
}

/*
*Description*

func Create

Creates a new User record in the database and returns the created record along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance where the record will be created.

*Returns*

	_  <*User>

		The created User record.

	_  <error>

		Encountered error (nil if no errors are encountered).
*/
func (user *User) Create(db *gorm.DB) (map[string]Model, error) {
	// TODO: Add field validation logic (func Create) -- add as BeforeCreate gorm hook definition at the top of this file
	err := db.Create(&user).Error
	returnRecords := map[string]Model{"user": user}
	return returnRecords, err
}

/*
*Description*

func Get

Retrieves a User record in the database by ID if it exists and returns that record along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to retrieve the specified record.

	userID <uint>

		The ID of the User record being requested.

*Returns*

	_  <*User>

		The User record that is retrieved from the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) Get(db *gorm.DB, userID uint) (map[string]Model, error) {
	err := db.First(&user, userID).Error
	returnRecords := map[string]Model{"user": user}
	return returnRecords, err
}

/*
*Description*

func GetAll

Retrieves all User records from the database.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that the records will be retrieved from.

*Returns*

	_  <[]User>

		The list of User records that are retrieved from the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) GetAll(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error

	return users, err
}

/*
*Description*

func GetRecordsByPrimaryIDs

Retrieves a list of User records from the database using their IDs (primary key).

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that the records will be retrieved from.

	ids  <[]uint>

		The list of User IDs that will be used to retrieve User records.

*Returns*

	_  <[]User>

		The list of User records that are retrieved from the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) GetRecordsByPrimaryIDs(db *gorm.DB, ids []uint) ([]User, error) {
	var users []User

	err := db.Where(ids).Find(&users).Error
	return users, err
}

/*
*Description*

func GetAppointments

Retrieves the list of all Appointments that are associated with the specified User.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that the records will be retrieved from.

	userID  <uint>

		The User ID that will be used to retrieve the list of Appointment records.

*Returns*

	_  <[]Appointment>

		The list of Appointment records that are retrieved from the database that are associated with the specified User.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) GetAppointments(db *gorm.DB, userID uint) ([]Appointment, error) {
	var appt Appointment
	var appts []Appointment

	appts, err := appt.GetRecordsBySecondaryID(db, "user_id", userID)
	return appts, err
}

/*
*Description*

func GetServiceAppointments

Retrieves the list of all Appointments (and the Service each Appointment is for) that are associated with the specified User.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that the records will be retrieved from.

	userID  <uint>

		The User ID that will be used to retrieve the list of Appointment/Service records.

*Returns*

	_  <[]map[string]interface{}>

		A list of JSON objects that each have an "appointment" key and a "service" key with the respective Appointment/Service record.

		Ex:
			[
				{
					"appointment": {
						"ID":11,
						"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"DeletedAt": null,
						"service_id":22,
						"user_id":33,
						"cancel_date_time":null,
						"active":true
					},
					"service": {
						"ID": 22,
						"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"DeletedAt": null,
						"business_id":42,
						"name":"Yoga class",
						"desc":"30 minute beginner yoga class",
						"start_date_time":"2023-05-31T14:30:00.0000000-05:00",
						"length":30,
						"capacity":20,
						"price":2000,
						"cancel_fee":0
					}
				},
				{
					"appointment": {
						"ID":44,
						"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
						"DeletedAt": null,
						"service_id":55,
						"user_id":66,
						"cancel_date_time":null,
						"active":true
					},
					"service": {
						"ID": 55,
						"CreatedAt": "2020-02-05T01:23:45.6789012-05:00",
						"UpdatedAt": "2020-02-05T01:23:45.6789012-05:00",
						"DeletedAt": null,
						"business_id":99,
						"name":"Spin class",
						"desc":"60 minute intermediate spin class",
						"start_date_time":"2023-04-20T10:00:00.0000000-05:00",
						"length":60,
						"capacity":10,
						"price":5000,
						"cancel_fee":1000
					}
				},
				...
			]

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) GetServiceAppointments(db *gorm.DB, userID uint) ([]map[string]interface{}, error) {
	var appts []Appointment
	var apptServiceID uint
	var apptService Service
	var serviceAppointments []map[string]interface{}

	// Get list of appointments for specified UserID
	appts, apptErr := user.GetAppointments(db, userID)
	if apptErr != nil {
		var errorMessage string = fmt.Sprintf("User ID (%d) does not have any appointment records in the database.  [%s]", userID, apptErr)
		return serviceAppointments, errors.New(errorMessage)
	}

	// Get list of ServiceIDs from user's appointments
	for _, appt := range appts {
		// Get Service associated with each of the user's appointments
		apptServiceID = appt.GetServiceID()
		returnedRecords, svcErr := apptService.Get(db, apptServiceID)
		if svcErr != nil {
			var errorMessage string = fmt.Sprintf("Service ID (%d) does not exist in the database, but is associated with Appointment ID (%d).  [%s]", userID, appt.GetID(), apptErr)
			return serviceAppointments, errors.New(errorMessage)
		}

		// Structure JSON appropriately and append to list of service appointments
		var svcAppt map[string]interface{} = map[string]interface{}{"appointment": appt, "service": returnedRecords["service"]}
		serviceAppointments = append(serviceAppointments, svcAppt)
	}

	return serviceAppointments, nil
}

/*
*Description*

func HasServiceAppointment

Returns whether the specified User has an Appointment scheduled for the specified Service.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to verify the User's Appointment status.

	userID  <uint>

		The ID of the User record

	serviceID  <uint>

		The ID of the Service record

*Returns*

	_  <bool>

		'True' if User has an Appointment for the specfied Service. 'False' if not.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) HasServiceAppointment(db *gorm.DB, userID uint, serviceID uint) (bool, error) {
	var appts []Appointment
	var apptServiceID uint

	// Get list of appointments for specified UserID
	appts, apptErr := user.GetAppointments(db, userID)
	if apptErr != nil {
		var errorMessage string = fmt.Sprintf("User ID (%d) does not have any appointment records in the database.  [%s]", userID, apptErr)
		return false, errors.New(errorMessage)
	}

	for _, appt := range appts {
		// Get Service associated with each of the user's appointments
		apptServiceID = appt.GetServiceID()
		if serviceID == apptServiceID {
			return true, nil
		}
	}

	return false, nil
}

/*
*Description*

func GetUserByEmail

Retrieves a User record in the database by email if it exists and returns that record along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to retrieve the specified record.

	userEmail  <string>

		The email of the User record being requested.

*Returns*

	_  <*User>

		The User record that is retrieved from the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) GetUserByEmail(db *gorm.DB, userEmail string) (*User, error) {
	err := db.First(&user, userEmail).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

/*
*Description*

func Update

Updates the specified User record in the database with the specified changes if the record exists.

Returns the updated record along with any errors that are thrown.

This function behaves like a PATCH method, rather than a true PUT. Any fields that aren't specified in the request body for the PUT request will not be altered for the specified record.

If a specified field's value should be deleted from the record, the appropriate null/blank should be specified for that key in the JSON request body (e.g. "type": "").

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance that will be used to retrieve and update the specified record.

	userID  <uint>

		The ID of the User record being updated.

	updates  <map[string]interface{}>

		JSON with the fields that will be updated as keys and the updated values as values.

		Ex:
			{
				"name": "New name",
				"address": "New address"
			}

*Returns*

	_  <*User>

		The User record that is updated in the database.

	_  <error>

		Encountered error (nil if no errors are encountered)
*/
func (user *User) Update(db *gorm.DB, userID uint, updates map[string]interface{}) (map[string]Model, error) {
	// Confirm userID exists in the database and get current object
	returnRecords, err := user.Get(db, userID)
	updateUser := returnRecords["user"]

	if err != nil {
		return returnRecords, err
	}

	// TODO: Add field validation logic (func Update) -- add as BeforeUpdate gorm hook definition at the top of this file
	err = db.Model(&updateUser).Where("id = ?", userID).Updates(updates).Error
	returnRecords = map[string]Model{"user": updateUser}

	return returnRecords, err
}

/*
*Description*

func Delete

Deletes the specified User record from the database by ID if it exists.

Deleted record is returned along with any errors that are thrown.

*Parameters*

	db  <*gorm.DB>

		A pointer to the database instance where the record will be deleted from.

	userID  <uint>

		The ID of the User record being deleted.

*Returns*

	_  <*User>

		The deleted User record.

	_  <error>

		Encountered error (nil if no errors are encountered).
*/
func (user *User) Delete(db *gorm.DB, userID uint) (map[string]Model, error) {
	// Confirm userID exists in the database and get current object
	returnRecords, err := user.Get(db, userID)
	deleteUser := returnRecords["user"]

	if err != nil {
		return returnRecords, err
	}

	if config.Debug {
		log.Printf("\n\nUser object targeted for deletion:\n\n%+v\n\n", deleteUser)
	}

	err = db.Delete(deleteUser).Error
	returnRecords = map[string]Model{"user": deleteUser}

	return returnRecords, err
}
