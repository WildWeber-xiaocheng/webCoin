// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.6

package types

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}

type CodeRequest struct {
	Phone   string `json:"phone,optional"`
	Country string `json:"country,optional"`
}
type CodeResponse struct {
}
