package meilisearch

// Locate is localized attributes for index in search
//
// https://www.meilisearch.com/docs/reference/api/settings#localized-attributes-object
type Locate string

const (
	EPO Locate = "epo" // EPO is Esperanto
	ENG Locate = "eng" // ENG is English
	RUS Locate = "rus" // RUS is Russian
	CMN Locate = "cmn" // CMN is Mandarin Chinese
	SPA Locate = "spa" // SPA is Spanish
	POR Locate = "por" // POR is Portuguese
	ITA Locate = "ita" // ITA is Italian
	BEN Locate = "ben" // BEN is Bengali
	FRA Locate = "fra" // FRA is French
	DEU Locate = "deu" // DEU is German
	UKR Locate = "ukr" // UKR is Ukrainian
	KAT Locate = "kat" // KAT is Georgian
	ARA Locate = "ara" // ARA is Arabic
	HIN Locate = "hin" // HIN is Hindi
	JPN Locate = "jpn" // JPN is Japanese
	HEB Locate = "heb" // HEB is Hebrew
	YID Locate = "yid" // YID is Yiddish
	POL Locate = "pol" // POL is Polish
	AMH Locate = "amh" // AMH is Amharic
	JAV Locate = "jav" // JAV is Javanese
	KOR Locate = "kor" // KOR is Korean
	NOB Locate = "nob" // NOB is Norwegian Bokmål
	DAN Locate = "dan" // DAN is Danish
	SWE Locate = "swe" // SWE is Swedish
	FIN Locate = "fin" // FIN is Finnish
	TUR Locate = "tur" // TUR is Turkish
	NLD Locate = "nld" // NLD is Dutch
	HUN Locate = "hun" // HUN is Hungarian
	CES Locate = "ces" // CES is Czech
	ELL Locate = "ell" // ELL is Greek
	BUL Locate = "bul" // BUL is Bulgarian
	BEL Locate = "bel" // BEL is Belarusian
	MAR Locate = "mar" // MAR is Marathi
	KAN Locate = "kan" // KAN is Kannada
	RON Locate = "ron" // RON is Romanian
	SLV Locate = "slv" // SLV is Slovenian
	HRV Locate = "hrv" // HRV is Croatian
	SRP Locate = "srp" // SRP is Serbian
	MKD Locate = "mkd" // MKD is Macedonian
	LIT Locate = "lit" // LIT is Lithuanian
	LAV Locate = "lav" // LAV is Latvian
	EST Locate = "est" // EST is Estonian
	TAM Locate = "tam" // TAM is Tamil
	VIE Locate = "vie" // VIE is Vietnamese
	URD Locate = "urd" // URD is Urdu
	THA Locate = "tha" // THA is Thai
	GUJ Locate = "guj" // GUJ is Gujarati
	UZB Locate = "uzb" // UZB is Uzbek
	PAN Locate = "pan" // PAN is Punjabi
	AZE Locate = "aze" // AZE is Azerbaijani
	IND Locate = "ind" // IND is Indonesian
	TEL Locate = "tel" // TEL is Telugu
	PES Locate = "pes" // PES is Persian
	MAL Locate = "mal" // MAL is Malayalam
	ORI Locate = "ori" // ORI is Odia
	MYA Locate = "mya" // MYA is Burmese
	NEP Locate = "nep" // NEP is Nepali
	SIN Locate = "sin" // SIN is Sinhala
	KHM Locate = "khm" // KHM is Khmer
	TUK Locate = "tuk" // TUK is Turkmen
	AKA Locate = "aka" // AKA is Akan
	ZUL Locate = "zul" // ZUL is Zulu
	SNA Locate = "sna" // SNA is Shona
	AFR Locate = "afr" // AFR is Afrikaans
	LAT Locate = "lat" // LAT is Latin
	SLK Locate = "slk" // SLK is Slovak
	CAT Locate = "cat" // CAT is Catalan
	TGL Locate = "tgl" // TGL is Tagalog
	HYE Locate = "hye" // HYE is Armenian
)

func (l Locate) String() string {
	return string(l)
}
