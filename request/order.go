package request

type OrderRequest struct {
	Id                uint                `json:"id"`
	OrderId           string              `json:"order_id"`
	PaymentId         string              `json:"payment_id"`
	UserId            uint                `json:"user_id"`
	AddressId         uint                `json:"address_id"`
	Remark            string              `json:"remark"`
	Origin            string              `json:"origin"`
	Color             string              `json:"color"`
	Status            int                 `json:"status"`
	RefundStatus      int                 `json:"refund_status"`
	OriginAmount      int64               `json:"origin_amount"`
	Amount            int64               `json:"amount"`
	PaidTime          int64               `json:"paid_time"`
	EndTime           int64               `json:"end_time"`
	DeliverTime       int64               `json:"deliver_time"`
	FinishedTime      int64               `json:"finished_time"`
	DiscountAmount    int64               `json:"discount_amount"` // 可能一个订单支持多个优惠
	CouponCodeId      string              `json:"-"`
	ShareUserId       uint                `json:"share_user_id"`       // 分享者
	ShareAmount       int64               `json:"share_amount"`        // 分销可得金额
	ShareParentAmount int64               `json:"share_parent_amount"` // 分销可得金额
	ShareSelfAmount   int64               `json:"share_self_amount"`   // 本级
	ExpressCompany    string              `json:"express_company"`     // 快递公司
	TrackingNumber    string              `json:"tracking_number"`     // 快递单号
	Address           OrderAddressRequest `json:"address"`
	Details           []OrderDetail       `json:"details"`
}

type OrderDetail struct {
	Id          uint   `json:"id"`
	OrderId     string `json:"order_id"`
	UserId      uint   `json:"user_id"`
	GoodsId     uint   `json:"goods_id"`
	GoodsItemId uint   `json:"goods_item_id"`
	Price       int64  `json:"price"`
	OriginPrice int64  `json:"origin_price"`
	Amount      int64  `json:"amount"`
	RealAmount  int64  `json:"real_amount"` // 实际支付的金额，用于退款的时候进行退款操作
	Quantity    int    `json:"quantity"`
	Status      int    `json:"status"`
}

type OrderRefundRequest struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type PluginPayConfig struct {
	AlipayAppId      string `json:"alipay_app_id"`
	AlipayPrivateKey string `json:"alipay_private_key"`
	AlipayCertPath   string `json:"alipay_cert_path"`

	WeixinAppId     string `json:"weixin_app_id"`
	WeixinAppSecret string `json:"weixin_app_secret"`
	WeixinMchId     string `json:"weixin_mch_id"`
	WeixinApiKey    string `json:"weixin_api_key"`
	WeixinCertPath  string `json:"weixin_cert_path"`
	WeixinKeyPath   string `json:"weixin_key_path"`
}

type OrderAddressRequest struct {
	Id          uint   `json:"id"`
	UserId      uint   `json:"user_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Country     string `json:"country"`
	AddressInfo string `json:"address_info"`
	Postcode    string `json:"postcode"`
	Status      int    `json:"status"`
}

type OrderExportRequest struct {
	Status    string `json:"status"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
