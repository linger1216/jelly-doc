package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	"context"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	// This service
	pb "github.com/linger1216/jelly-doc/src/server/pb"
)

const contentType = "application/json; charset=utf-8"

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewApiClient
	_ = io.Copy
	_ = errors.Wrap
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}
	serverOptions = append(serverOptions, options...)
	m := mux.NewRouter()

	m.Methods("POST").Path("/jd/v1/api").Handler(httptransport.NewServer(
		endpoints.CreateEndpoint,
		DecodeHTTPCreateZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("GET").Path("/jd/v1/api/{ids}").Handler(httptransport.NewServer(
		endpoints.GetEndpoint,
		DecodeHTTPGetZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("HEAD").Path("/jd/v1/api").Handler(httptransport.NewServer(
		endpoints.ListEndpoint,
		DecodeHTTPListZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))
	m.Methods("GET").Path("/jd/v1/api").Handler(httptransport.NewServer(
		endpoints.ListEndpoint,
		DecodeHTTPListOneRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("PUT").Path("/jd/v1/api").Handler(httptransport.NewServer(
		endpoints.UpdateEndpoint,
		DecodeHTTPUpdateZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("DELETE").Path("/jd/v1/api/{ids}").Handler(httptransport.NewServer(
		endpoints.DeleteEndpoint,
		DecodeHTTPDeleteZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))
	return m
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(errorWrapper{Error: err.Error()})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// httpError satisfies the Headerer and StatusCoder interfaces in
// package github.com/go-kit/kit/transport/http.
type httpError struct {
	error
	statusCode int
	headers    map[string][]string
}

func (h httpError) StatusCode() int {
	return h.statusCode
}

func (h httpError) Headers() http.Header {
	return h.headers
}

// Server Decode

// DecodeHTTPCreateZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded create request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPCreateZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.CreateApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPGetZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded get request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.GetApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	IdsGetStr := pathParams["ids"]

	var IdsGet []string
	if len(IdsGetStr) > 0 {
		IdsGet = strings.Split(IdsGetStr, ",")
	}
	req.Ids = IdsGet

	return &req, err
}

// DecodeHTTPListZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded list request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPListZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ListApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if HeaderListStrArr, ok := queryParams["header"]; ok {
		HeaderListStr := HeaderListStrArr[0]
		HeaderList, err := strconv.ParseInt(HeaderListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting HeaderList from query, queryParams: %v", queryParams))
		}
		req.Header = int32(HeaderList)
	}

	if CurrentPageListStrArr, ok := queryParams["current_page"]; ok {
		CurrentPageListStr := CurrentPageListStrArr[0]
		CurrentPageList, err := strconv.ParseInt(CurrentPageListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting CurrentPageList from query, queryParams: %v", queryParams))
		}
		req.CurrentPage = int32(CurrentPageList)
	}

	if PageSizeListStrArr, ok := queryParams["page_size"]; ok {
		PageSizeListStr := PageSizeListStrArr[0]
		PageSizeList, err := strconv.ParseInt(PageSizeListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting PageSizeList from query, queryParams: %v", queryParams))
		}
		req.PageSize = int32(PageSizeList)
	}

	return &req, err
}

// DecodeHTTPListOneRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded list request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPListOneRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ListApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if HeaderListStrArr, ok := queryParams["header"]; ok {
		HeaderListStr := HeaderListStrArr[0]
		HeaderList, err := strconv.ParseInt(HeaderListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting HeaderList from query, queryParams: %v", queryParams))
		}
		req.Header = int32(HeaderList)
	}

	if CurrentPageListStrArr, ok := queryParams["current_page"]; ok {
		CurrentPageListStr := CurrentPageListStrArr[0]
		CurrentPageList, err := strconv.ParseInt(CurrentPageListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting CurrentPageList from query, queryParams: %v", queryParams))
		}
		req.CurrentPage = int32(CurrentPageList)
	}

	if PageSizeListStrArr, ok := queryParams["page_size"]; ok {
		PageSizeListStr := PageSizeListStrArr[0]
		PageSizeList, err := strconv.ParseInt(PageSizeListStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting PageSizeList from query, queryParams: %v", queryParams))
		}
		req.PageSize = int32(PageSizeList)
	}

	return &req, err
}

// DecodeHTTPUpdateZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded update request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPUpdateZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UpdateApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if ApisUpdateStrArr, ok := queryParams["Apis"]; ok {
		ApisUpdateStr := ApisUpdateStrArr[0]

		var ApisUpdate []*pb.ApiModel
		err = json.Unmarshal([]byte(ApisUpdateStr), &ApisUpdate)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't decode ApisUpdate from %v", ApisUpdateStr)
		}
		req.Apis = ApisUpdate
	}

	return &req, err
}

// DecodeHTTPDeleteZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded delete request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPDeleteZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.DeleteApiRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	IdsDeleteStr := pathParams["ids"]

	var IdsDelete []string
	if len(IdsDeleteStr) > 0 {
		IdsDelete = strings.Split(IdsDeleteStr, ",")
	}
	req.Ids = IdsDelete

	return &req, err
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	marshaller := jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}

	return marshaller.Marshal(w, response.(proto.Message))
}

// Helper functions

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// Tune specific change.
	// also add the request url
	ctx = context.WithValue(ctx, "request-url", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}
