package telescope

import (
	"github.com/sirupsen/logrus"
	"time"
)

// ClientRequest @Bean
type ClientRequest struct {
	Method  string `json:"method"`
	Uri     string `json:"uri"`
	Headers struct {
		ContentLength string `json:"content-length"`
		UserAgent     string `json:"user-agent"`
		ContentType   string `json:"content-type"`
		Host          string `json:"host"`
	} `json:"headers"`
	Payload struct {
		Request struct {
			Head struct {
				Appid       string    `json:"appid"`
				Version     string    `json:"version"`
				Reserve     string    `json:"reserve"`
				SignType    string    `json:"sign_type"`
				RequestTime time.Time `json:"request_time"`
			} `json:"head"`
			Body struct {
				RequestId     string    `json:"request_id"`
				BrandCode     string    `json:"brand_code"`
				StoreSn       string    `json:"store_sn"`
				WorkstationSn string    `json:"workstation_sn"`
				CheckSn       int64     `json:"check_sn"`
				Scene         int       `json:"scene"`
				SalesTime     time.Time `json:"sales_time"`
				Amount        int       `json:"amount"`
				Currency      string    `json:"currency"`
				Subject       string    `json:"subject"`
				Operator      string    `json:"operator"`
				Customer      int       `json:"customer"`
				IndustryCode  string    `json:"industry_code"`
				PosInfo       string    `json:"pos_info"`
				NotifyUrl     string    `json:"notify_url"`
			} `json:"body"`
		} `json:"request"`
		Signature string `json:"signature"`
	} `json:"payload"`
	ResponseStatus  int `json:"response_status"`
	ResponseHeaders struct {
		Date                string `json:"date"`
		ContentType         string `json:"content-type"`
		ContentLength       string `json:"content-length"`
		Connection          string `json:"connection"`
		SetCookie           string `json:"set-cookie"`
		Vary                string `json:"vary"`
		CacheControl        string `json:"cache-control"`
		Pragma              string `json:"pragma"`
		Expires             string `json:"expires"`
		XContentTypeOptions string `json:"x-content-type-options"`
		XFrameOptions       string `json:"x-frame-options"`
		XXssProtection      string `json:"x-xss-protection"`
		ReferrerPolicy      string `json:"referrer-policy"`
	} `json:"response_headers"`
	Response struct {
		Response struct {
			Head struct {
				Appid        string    `json:"appid"`
				Version      string    `json:"version"`
				Reserve      string    `json:"reserve"`
				SignType     string    `json:"sign_type"`
				ResponseTime time.Time `json:"response_time"`
			} `json:"head"`
			Body struct {
				ResultCode  string `json:"result_code"`
				BizResponse struct {
					ResultCode string `json:"result_code"`
					Data       struct {
						BrandCode     string `json:"brand_code"`
						StoreSn       string `json:"store_sn"`
						WorkstationSn string `json:"workstation_sn"`
						CheckSn       string `json:"check_sn"`
						OrderSn       string `json:"order_sn"`
						OrderToken    string `json:"order_token"`
						OrderSource   int    `json:"order_source"`
						RequestId     string `json:"request_id"`
					} `json:"data"`
				} `json:"biz_response"`
			} `json:"body"`
		} `json:"response"`
		Signature string `json:"signature"`
	} `json:"response"`
	Hostname string `json:"hostname"`
	User     struct {
		Id    int         `json:"id"`
		Name  interface{} `json:"name"`
		Email interface{} `json:"email"`
	} `json:"user"`
}

func (b ClientRequest) BindType() string {
	return "client_request"
}

func (b ClientRequest) Handler(entry *logrus.Entry) (*entries, []tag) {
	return &entries{}, nil
}
