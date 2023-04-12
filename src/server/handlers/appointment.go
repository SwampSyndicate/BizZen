package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"server/utils"
	_ "time"
)

/*
*Description*

func CreateAppointment

Creates a new appointment record in the database.

*Parameters*

	writer  <http.ResponseWriter>

		The HTTP response writer

	request  <*http.Request>

		The HTTP request

*Returns*

	None

*Expected request format*

	Type:  POST

	Route:	/appointment

	Body:
		Format: JSON

		Required fields:

			service_id  <uint>

				ID of Service record Appointment is associated with

			user_id  <uint>

				ID of User that booked the appointment

		Optional fields:

			cancel_date_time  <time.Time>

				Date/time when appointment was cancelled (if cancelled, else null)

			active  <bool>

				Status flag (true for Active, false for Cancelled)

*Example request(s)*

	POST /appointment
	{
		"service_id":123
		"user_id":123
	}

*Response format*

	Success:

		HTTP/1.1 201 Created
		Content-Type: application/json

		{
			"ID": 123456,
			"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"DeletedAt": null,
			"service_id":123,
			"user_id":123,
			"cancel_date_time":null,
			"active":true
		}

	Failure:

		-- Case = Bad request body
		HTTP/1.1 400 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}

		-- Case = Database operation error
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}
*/
func (app *Application) CreateAppointment(writer http.ResponseWriter, request *http.Request) {
	appt := models.Appointment{}

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&appt); err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			err.Error())

		return
	}

	defer request.Body.Close()

	returnedRecords, err := appt.Create(app.AppDB)
	createdAppointment := returnedRecords["appointment"]
	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusInternalServerError,
			err.Error())

		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusCreated,
		createdAppointment)
}

/*
*Description*

func GetAppointment

Get appointment record from the database by ID.

*Parameters*

	writer  <http.ResponseWriter>

		The HTTP response writer

	request  <*http.Request>

		The HTTP request

*Returns*

	None

*Expected request format*

	Type:	GET

	Route:	/appointment/{id}

	Body:

		None

*Example request(s)*

	GET /appointment/123456

*Response format*

	Success:

		HTTP/1.1 200 OK
		Content-Type: application/json

		{
			"ID": 123456,
			"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"DeletedAt": null,
			"service_id":123,
			"user_id":123,
			"cancel_date_time":null,
			"active":true
		}

	Failure:
		-- Case = ID missing from or incorrectly formatted in request url
		HTTP/1.1 400 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}

		HTTP/1.1 404 Resource Not Found
		Content-Type: application/json

		{
			"error":"ERROR MESSAGE TEXT HERE"
		}
*/
func (app *Application) GetAppointment(writer http.ResponseWriter, request *http.Request) {
	appt := models.Appointment{}
	apptID, err := utils.ParseRequestID(request)

	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			err.Error())

		return
	}

	returnedAppointment, err := appt.Get(app.AppDB, apptID)
	if err != nil {
		var errorMessage string = fmt.Sprintf("Appointment ID (%d) does not exist in the database.  [%s]", apptID, err)

		utils.RespondWithError(
			writer,
			http.StatusNotFound,
			errorMessage)

		log.Printf("ERROR:  %s", errorMessage)

		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		returnedAppointment)
}

/*
*Description*

func UpdateAppointment

Updates the appointment record associated with the specified appointment ID in the database.

This function behaves like a PATCH method, rather than a true PUT. Any fields that aren't specified in the request body for the PUT request will not be altered for the specified record.

If a specified field's value should be deleted from the record, the appropriate null/blank should be specified for that key in the JSON request body (e.g. "address2": "").

*Parameters*

	writer  <http.ResponseWriter>

		The HTTP response writer

	request  <*http.Request>

		The HTTP request

*Returns*

	None

*Expected request format*

	Type:   PUT

	Route:  /appointment/{id}

	Body:
		Format: JSON

		Required fields:

			N/A  --  At least one field should be present in the request body, but no fields are specifically required to be present in the request body.

		Optional fields:

			service_id  <uint>

				ID of Service record Appointment is associated with

			user_id  <uint>

				ID of User that booked the appointment

			cancel_date_time  <time.Time>

				Date/time when appointment was cancelled (if cancelled, else null)

			active  <bool>

				Status flag (true for Active, false for Cancelled)

*Example request(s)*

	PUT /appointment/123456
	{
		"active":false
	}

*Response format*

	Success:

		HTTP/1.1 200 OK
		Content-Type: application/json

		{
			"ID": 123456,
			"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"UpdatedAt": "2020-02-13T04:20:12.6789012-05:00",
			"DeletedAt": null,
			"service_id":123,
			"user_id":123,
			"cancel_date_time":"2020-01-31T04:20:12.6789012-05:00",
			"active":false
		}

	Failure:
		-- Case = Bad request body or missing/misformatted ID in request URL
		HTTP/1.1 400 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}

		-- Case = Database operation error
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}
*/
func (app *Application) UpdateAppointment(writer http.ResponseWriter, request *http.Request) {
	appt := models.Appointment{}
	apptID, err := utils.ParseRequestID(request)

	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			err.Error())

		return
	}

	var updates map[string]interface{}

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&updates); err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			err.Error())

		return
	}

	defer request.Body.Close()

	returnedRecords, err := appt.Update(app.AppDB, apptID, updates)
	updatedAppointment := returnedRecords["appointment"]
	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusInternalServerError,
			err.Error())

		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		updatedAppointment)
}

/*
*Description*

func DeleteAppointment

Delete an appointment record from the database by appointment ID if the ID exists in the database.

Deleted appointment record is returned in the response body if the operation is sucessful.

*Parameters*

	writer  <http.ResponseWriter>

		The HTTP response writer

	request  <*http.Request>

		The HTTP request

*Returns*

	None

*Expected request format*

	Type:  	DELETE

	Route:	/appointment/{id}

	Body:

		None

*Example request(s)*

	DELETE /appointment/123456

*Response format*

	Success:

		HTTP/1.1 200 OK
		Content-Type: application/json

		{
			"ID": 123456,
			"CreatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"UpdatedAt": "2020-01-01T01:23:45.6789012-05:00",
			"DeletedAt": "2022-06-31T04:20:12.6789012-05:00",
			"service_id":123,
			"user_id":123,
			"cancel_date_time":null,
			"active":true
		}

	Failure:
		-- Case = ID missing from or incorrectly formatted in request url
		HTTP/1.1 400 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}

		-- Case = Database operation error
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{
		"error":"ERROR MESSAGE TEXT HERE"
		}
*/
func (app *Application) DeleteAppointment(writer http.ResponseWriter, request *http.Request) {
	appt := models.Appointment{}
	apptID, err := utils.ParseRequestID(request)

	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			err.Error())

		return
	}

	returnedRecords, err := appt.Delete(app.AppDB, apptID)
	deletedAppointment := returnedRecords["appointment"]
	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusInternalServerError,
			err.Error())

		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		deletedAppointment)

}