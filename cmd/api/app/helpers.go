package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type envelope map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func parseForm(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	const maxBytes = 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Parse the form data from the request
	if err := r.ParseForm(); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return fmt.Errorf("form data must not be larger than %d bytes", maxBytes)
		}

		switch {
		case strings.Contains(err.Error(), "non-form content type"):
			return errors.New("request Content-Type is not form-compatible")
		case strings.HasPrefix(err.Error(), "mime:"):
			return fmt.Errorf("invalid form data format: %v", err)
		default:
			return fmt.Errorf("error parsing form: %v", err)
		}
	}

	// Decode the form data into the destination struct
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false) // Disallow unknown fields

	if err := decoder.Decode(dst, r.PostForm); err != nil {
		var schemaUnknownKey *schema.UnknownKeyError
		if errors.As(err, &schemaUnknownKey) {
			return fmt.Errorf("form contains unknown key %q", schemaUnknownKey.Key)
		}

		var multiErr schema.MultiError
		if errors.As(err, &multiErr) {
			// Return the first error encountered (map iteration order is random)
			for field, err := range multiErr {
				if conversionErr, ok := err.(schema.ConversionError); ok {
					return fmt.Errorf("invalid value for field %q", field)
				}
				return fmt.Errorf("error in field %q: %v", field, err)
			}
		}

		return fmt.Errorf("error decoding form data: %v", err)
	}

	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func readIDParam(r *http.Request, param string) (int64, error) {
	idStr := r.PathValue(param)
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}

func readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}
