package client

import (
	"io"

	"resty.dev/v3"
)

// MultipartProgressCallback is called during multipart uploads to report progress.
type MultipartProgressCallback func(fieldName string, fileName string, bytesWritten int64, totalBytes int64)

// requestExecutor is the execution backend for a RequestBuilder.
// Transport implements it directly; tests can supply a mock via NewMockRequestBuilder.
type requestExecutor interface {
	execute(req *resty.Request, method, path string) (*resty.Response, error)
	executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error)
}

// RequestBuilder constructs a single API request using a fluent interface.
// The service layer owns the full request shape — headers, body, query params,
// result target — before handing the completed request to the executor which
// handles auth, retry, and other transport concerns.
//
// Usage:
//
//	resp, err := s.client.NewRequest(ctx).
//	    SetHeader("Accept", "application/json").
//	    SetHeader("Content-Type", "application/json").
//	    SetBody(payload).
//	    SetResult(&result).
//	    Post(endpoint)
type RequestBuilder struct {
	req      *resty.Request
	executor requestExecutor
	result   any
}

// SetHeader sets a request-level header. Empty values are ignored.
func (b *RequestBuilder) SetHeader(key, value string) *RequestBuilder {
	if value != "" {
		b.req.SetHeader(key, value)
	}
	return b
}

// SetQueryParam adds a URL query parameter. Empty values are ignored.
func (b *RequestBuilder) SetQueryParam(key, value string) *RequestBuilder {
	if value != "" {
		b.req.SetQueryParam(key, value)
	}
	return b
}

// SetQueryParams adds multiple URL query parameters in bulk. Empty values are ignored.
func (b *RequestBuilder) SetQueryParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		if v != "" {
			b.req.SetQueryParam(k, v)
		}
	}
	return b
}

// SetBody sets the request body. Nil is ignored.
func (b *RequestBuilder) SetBody(body any) *RequestBuilder {
	if body != nil {
		b.req.SetBody(body)
	}
	return b
}

// SetResult sets the target for JSON unmarshaling of a successful response.
func (b *RequestBuilder) SetResult(result any) *RequestBuilder {
	b.result = result
	b.req.SetResult(result)
	return b
}

// SetFormData sets URL-encoded form data for the request.
func (b *RequestBuilder) SetFormData(formData map[string]string) *RequestBuilder {
	if len(formData) > 0 {
		b.req.SetFormData(formData)
	}
	return b
}

// SetMultipartFile configures the request for a multipart file upload.
// Execute with Post after setting any additional form fields or headers.
// Content-Type is managed automatically by resty.
func (b *RequestBuilder) SetMultipartFile(fileField, fileName string, fileReader io.Reader, fileSize int64, callback MultipartProgressCallback) *RequestBuilder {
	if fileReader != nil && fileName != "" && fileField != "" {
		field := &resty.MultipartField{
			Name:        fileField,
			FileName:    fileName,
			ContentType: "application/octet-stream",
			Reader:      fileReader,
			FileSize:    fileSize,
		}
		if callback != nil {
			field.ProgressCallback = func(p resty.MultipartFieldProgress) {
				callback(p.Name, p.FileName, p.Written, p.FileSize)
			}
		}
		b.req.SetMultipartFields(field)
	}
	return b
}

// SetMultipartFormData adds additional form fields to a multipart request.
func (b *RequestBuilder) SetMultipartFormData(formFields map[string]string) *RequestBuilder {
	if len(formFields) > 0 {
		b.req.SetMultipartFormData(formFields)
	}
	return b
}

// Get executes the request as GET against path.
func (b *RequestBuilder) Get(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "GET", path)
}

// Post executes the request as POST against path.
func (b *RequestBuilder) Post(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "POST", path)
}

// Put executes the request as PUT against path.
func (b *RequestBuilder) Put(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "PUT", path)
}

// Patch executes the request as PATCH against path.
func (b *RequestBuilder) Patch(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "PATCH", path)
}

// Delete executes the request as DELETE against path.
func (b *RequestBuilder) Delete(path string) (*resty.Response, error) {
	return b.executor.execute(b.req, "DELETE", path)
}

// GetBytes executes a GET request and returns raw response bytes without JSON
// unmarshaling. Use for binary responses such as CSV exports or file downloads.
func (b *RequestBuilder) GetBytes(path string) (*resty.Response, []byte, error) {
	return b.executor.executeGetBytes(b.req, path)
}

// mockRequestExecutor backs a RequestBuilder in unit tests.
type mockRequestExecutor struct {
	fn func(method, path string, result any) (*resty.Response, error)
}

func (m *mockRequestExecutor) execute(req *resty.Request, method, path string) (*resty.Response, error) {
	return m.fn(method, path, nil)
}

func (m *mockRequestExecutor) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	resp, err := m.fn("GET", path, nil)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}

// NewMockRequestBuilder returns a RequestBuilder suitable for unit tests.
// The fn callback receives the HTTP method, path, and returns a pre-programmed response.
func NewMockRequestBuilder(fn func(method, path string, result any) (*resty.Response, error)) *RequestBuilder {
	return &RequestBuilder{
		req:      resty.New().R(),
		executor: &mockRequestExecutor{fn: fn},
	}
}
