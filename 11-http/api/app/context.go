// Package app provides application support for context and MongoDB access.
// Current Status Codes:
// 		200 OK           : StatusOK                  : Call is success and returning data.
// 		400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
// 		401 Unauthorized : StatusUnauthorized        : Authentication failure.
// 		404 Not Found    : StatusNotFound            : Invalid URL or identifier.
// 		500 Internal     : StatusInternalServerError : Application specific beyond scope of user.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
)

// Context contains data associated with a single request.
type Context struct {
	Session *mgo.Session
	http.ResponseWriter
	Request   *http.Request
	Params    map[string]string
	SessionID string
}

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

type jsonError struct {
	Error  string    `json:"error"`
	Fields []Invalid `json:"fields,omitempty"`
}

// Authenticate handles the authentication of each request.
func (c *Context) Authenticate() error {
	log.Println(c.SessionID, ": api : Authenticate : Started")

	// ServeError(w, errors.New("Auth Error"), http.StatusUnauthorized)

	log.Println(c.SessionID, ": api : Authenticate : Completed")
	return nil
}

// Respond sends JSON to the client.
//
// If code is StatusNoContent, v is expected to be nil.
func (c *Context) Respond(v interface{}, code int) {
	log.Printf("%v : api : Respond [%d] : Started", c.SessionID, code)

	if code == http.StatusNoContent {
		c.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		// v failed to marshal (programmer error), so panic
		log.Panicf("%v : api : Respond [%d] : Failed: %v", c.SessionID, code, err)
	}

	datalen := len(data) + 1 // account for trailing LF
	h := c.Header()
	h.Set("Content-Type", "application/json")
	h.Set("Content-Length", strconv.Itoa(datalen))
	c.WriteHeader(code)
	fmt.Fprintf(c, "%s\n", data)

	log.Printf("%v : api : Respond [%d] : Completed", c.SessionID, code)
}

// RespondInvalid sends JSON describing field validation errors.
func (c *Context) RespondInvalid(fields []Invalid) {
	v := jsonError{
		Error:  "field validation failure",
		Fields: fields,
	}
	c.Respond(v, http.StatusBadRequest)
}

// RespondError sends JSON describing the error
func (c *Context) RespondError(error string, code int) {
	c.Respond(jsonError{Error: error}, code)
}
