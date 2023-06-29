package datadome

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type JSData struct {
	TimeSinceLastEvent float64 `json:"ttst,omitempty"`
	Ifov               bool    `json:"ifov,omitempty"`
	TimeToDoWebGL      float64 `json:"tagpu,omitempty"`
	Glvd               string  `json:"glvd,omitempty"`
	Glrd               string  `json:"glrd,omitempty"`
	Hc                 int     `json:"hc,omitempty"`
	BrOh               int     `json:"br_oh,omitempty"`
	BrOw               int     `json:"br_ow,omitempty"`
	Ua                 string  `json:"ua,omitempty"`
	Wbd                bool    `json:"wbd,omitempty"`
	Wdif               bool    `json:"wdif,omitempty"`
	Wdifrm             bool    `json:"wdifrm,omitempty"`
	Npmtm              bool    `json:"npmtm,omitempty"`
	BrH                int     `json:"br_h,omitempty"`
	BrW                int     `json:"br_w,omitempty"`
	CookieLength       int     `json:"nddc,omitempty"`
	RsH                int     `json:"rs_h,omitempty"`
	RsW                int     `json:"rs_w,omitempty"`
	RsCd               int     `json:"rs_cd,omitempty"`
	Phe                bool    `json:"phe,omitempty"`
	Nm                 bool    `json:"nm,omitempty"`
	Jsf                bool    `json:"jsf,omitempty"`
	Lg                 string  `json:"lg,omitempty"`
	Pr                 int     `json:"pr,omitempty"`
	ArsH               int     `json:"ars_h,omitempty"`
	ArsW               int     `json:"ars_w,omitempty"`
	TimezoneOffset     int     `json:"tz,omitempty"`
	StrSs              bool    `json:"str_ss,omitempty"`
	StrLs              bool    `json:"str_ls,omitempty"`
	StrIdb             bool    `json:"str_idb,omitempty"`
	StrOdb             bool    `json:"str_odb,omitempty"`
	Plgod              bool    `json:"plgod,omitempty"`
	NavPluginLength    int     `json:"plg,omitempty"`
	Plgne              bool    `json:"plgne,omitempty"`
	Plgre              bool    `json:"plgre,omitempty"`
	Plgof              bool    `json:"plgof,omitempty"`
	Plggt              bool    `json:"plggt,omitempty"`
	Pltod              bool    `json:"pltod,omitempty"`
	Hcovdr             bool    `json:"hcovdr,omitempty"`
	Hcovdr2            bool    `json:"hcovdr2,omitempty"`
	Plovdr             bool    `json:"plovdr,omitempty"`
	Plovdr2            bool    `json:"plovdr2,omitempty"`
	Ftsovdr            bool    `json:"ftsovdr,omitempty"`
	Ftsovdr2           bool    `json:"ftsovdr2,omitempty"`
	Lb                 bool    `json:"lb,omitempty"`
	EvalLength         int     `json:"eva,omitempty"`
	Lo                 bool    `json:"lo,omitempty"`
	TsMtp              int     `json:"ts_mtp,omitempty"`
	TsTec              bool    `json:"ts_tec,omitempty"`
	TsTsa              bool    `json:"ts_tsa,omitempty"`
	Vnd                string  `json:"vnd,omitempty"`
	Bid                string  `json:"bid,omitempty"`
	Mmt                string  `json:"mmt,omitempty"`
	Plu                string  `json:"plu,omitempty"`
	Hdn                bool    `json:"hdn,omitempty"`
	Awe                bool    `json:"awe,omitempty"`
	Geb                bool    `json:"geb,omitempty"`
	Dat                bool    `json:"dat,omitempty"`
	Med                string  `json:"med,omitempty"`
	Aco                string  `json:"aco,omitempty"`
	Acots              bool    `json:"acots,omitempty"`
	Acmp               string  `json:"acmp,omitempty"`
	Acmpts             bool    `json:"acmpts,omitempty"`
	Acw                string  `json:"acw,omitempty"`
	Acwts              bool    `json:"acwts,omitempty"`
	Acma               string  `json:"acma,omitempty"`
	Acmats             bool    `json:"acmats,omitempty"`
	Acaa               string  `json:"acaa,omitempty"`
	Acaats             bool    `json:"acaats,omitempty"`
	Ac3                string  `json:"ac3,omitempty"`
	Ac3ts              bool    `json:"ac3ts,omitempty"`
	Acf                string  `json:"acf,omitempty"`
	Acfts              bool    `json:"acfts,omitempty"`
	Acmp4              string  `json:"acmp4,omitempty"`
	Acmp4ts            bool    `json:"acmp4ts,omitempty"`
	Acmp3              string  `json:"acmp3,omitempty"`
	Acmp3ts            bool    `json:"acmp3ts,omitempty"`
	Acwm               string  `json:"acwm,omitempty"`
	Acwmts             bool    `json:"acwmts,omitempty"`
	Ocpt               bool    `json:"ocpt,omitempty"`
	Vco                string  `json:"vco,omitempty"`
	Vcots              bool    `json:"vcots,omitempty"`
	Vch                string  `json:"vch,omitempty"`
	Vchts              bool    `json:"vchts,omitempty"`
	Vcw                string  `json:"vcw,omitempty"`
	Vcwts              bool    `json:"vcwts,omitempty"`
	Vc3                string  `json:"vc3,omitempty"`
	Vc3ts              bool    `json:"vc3ts,omitempty"`
	Vcmp               string  `json:"vcmp,omitempty"`
	Vcmpts             bool    `json:"vcmpts,omitempty"`
	Vcq                string  `json:"vcq,omitempty"`
	Vcqts              bool    `json:"vcqts,omitempty"`
	Vc1                string  `json:"vc1,omitempty"`
	Vc1ts              bool    `json:"vc1ts,omitempty"`
	Dvm                int     `json:"dvm,omitempty"`
	Sqt                bool    `json:"sqt,omitempty"`
	So                 string  `json:"so,omitempty"`
	Wdw                bool    `json:"wdw,omitempty"`
	Cokys              string  `json:"cokys,omitempty"`
	Ecpc               bool    `json:"ecpc,omitempty"`
	Lgs                bool    `json:"lgs,omitempty"`
	Lgsod              bool    `json:"lgsod,omitempty"`
	Psn                bool    `json:"psn,omitempty"`
	Edp                bool    `json:"edp,omitempty"`
	Addt               bool    `json:"addt,omitempty"`
	Wsdc               bool    `json:"wsdc,omitempty"`
	Ccsr               bool    `json:"ccsr,omitempty"`
	Nuad               bool    `json:"nuad,omitempty"`
	Bcda               bool    `json:"bcda,omitempty"`
	Idn                bool    `json:"idn,omitempty"`
	Capi               bool    `json:"capi,omitempty"`
	Svde               bool    `json:"svde,omitempty"`
	Vpbq               bool    `json:"vpbq,omitempty"`
	Ucdv               bool    `json:"ucdv,omitempty"`
	Spwn               bool    `json:"spwn,omitempty"`
	Emt                bool    `json:"emt,omitempty"`
	Bfr                bool    `json:"bfr,omitempty"`
	Dbov               bool    `json:"dbov,omitempty"`
	Cfpfe              string  `json:"cfpfe,omitempty"`
	Stcfp              string  `json:"stcfp,omitempty"`
	Prm                bool    `json:"prm,omitempty"`
	Tzp                string  `json:"tzp,omitempty"`
	Cvs                bool    `json:"cvs,omitempty"`
	Usb                string  `json:"usb,omitempty"`
	TimeSeconds        int64   `json:"jset,omitempty"`
}

type DDResponse struct {
	Status int    `json:"status"`
	Cookie string `json:"cookie"`
}

type DataDomeSession struct {
	httpClient     tls_client.HttpClient
	Payload        []Param
	GenData        GenerateData
	DataDomeDomain string
}

type GenerateData struct {
	Proxy           string
	DataDomeKey     string
	DataDomeVersion string
	Refer           string
	Cid             string
	UserAgent       string
	Path            string
	Host            string
}

type Param struct {
	key   string
	value string
}
