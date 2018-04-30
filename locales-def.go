package strongo

import "fmt"

// TODO: This module should be in a dedicate package?

const (
	// LocaleCodeUndefined is undefined locale code
	LocaleCodeUndefined = "UNDEFINED"
	// LocaleCodeEnUS is locale code
	LocaleCodeEnUS = "en-US"
	// LocaleCodeEnUK is locale code
	LocaleCodeEnUK = "en-UK"
	// LocalCodeRuRu is locale code
	LocalCodeRuRu = "ru-RU"
	// LOCALE_ID_ID     = "id-ID"

	// LocaleCodeFaIR is locale code
	LocaleCodeFaIR = "fa-IR"
	// LocaleCodeItIT is locale code
	LocaleCodeItIT = "it-IT"

	// LocaleCodeDeDE is locale code
	LocaleCodeDeDE = "de-DE"
	// LocaleCodeEsES is locale code
	LocaleCodeEsES = "es-ES"
	// LocaleCodeFrFR is locale code
	LocaleCodeFrFR = "fr-FR"
	// LocaleCodePlPL is locale code
	LocaleCodePlPL = "pl-PL"
	// LocaleCodePtPT is locale code
	LocaleCodePtPT = "pt-PT"
	// LocaleCodePtBR is locale code
	LocaleCodePtBR = "pt-BR"

	// LocaleCodeKoKO is locale code
	LocaleCodeKoKO = "ko-KO"
	// LocaleCodeJaJP is locale code
	LocaleCodeJaJP = "ja-JP"
	// LocaleCodeZhCN is locale code
	LocaleCodeZhCN = "zh-CN"
)

//"4. French ",
//"5. Spanish ",
//"6. Italian \xF0\x9F\x87\xAE\xF0\x9F\x87\xB9",

var (
	// LocaleUndefined is undefined locale
	LocaleUndefined = Locale{Code5: LocaleCodeUndefined, NativeTitle: "Undefined", EnglishTitle: "Undefined"}

	// LocaleEnUS is en-US locale
	LocaleEnUS = Locale{Code5: LocaleCodeEnUS, NativeTitle: "English", EnglishTitle: "English", FlagIcon: "🇺🇸"}

	// LocaleEnUK = Locale{Code5: LocaleCodeEnUK, NativeTitle: "English", EnglishTitle: "English", FlagIcon: "🇺🇸"}

	// LocaleRuRu is locale
	LocaleRuRu = Locale{Code5: LocalCodeRuRu, NativeTitle: "Русский", EnglishTitle: "Russian", FlagIcon: "🇷🇺"}

	// LocaleIdId is locale
	//  LocaleIdId = Locale{Code5: LOCALE_ID_ID, NativeTitle: "Indonesian", EnglishTitle: "Indonesian", FlagIcon: ""}

	// LocaleDeDe is locale
	LocaleDeDe = Locale{Code5: LocaleCodeDeDE, NativeTitle: "Deutsch", EnglishTitle: "German", FlagIcon: "🇩🇪"}

	// LocaleEsEs is locale
	LocaleEsEs = Locale{Code5: LocaleCodeEsES, NativeTitle: "Español", EnglishTitle: "Spanish", FlagIcon: "🇪🇸"}

	// LocaleFrFr is locale
	LocaleFrFr = Locale{Code5: LocaleCodeFrFR, NativeTitle: "Français", EnglishTitle: "French", FlagIcon: "🇫🇷"}

	// LocaleItIt is locale
	LocaleItIt = Locale{Code5: LocaleCodeItIT, NativeTitle: "Italiano", EnglishTitle: "Italian", FlagIcon: "🇮🇹"}

	// LocalePlPl is locale
	LocalePlPl = Locale{Code5: LocaleCodePlPL, NativeTitle: "Polszczyzna", EnglishTitle: "Polish", FlagIcon: "🇵🇱"}

	// LocalePtPt is locale
	LocalePtPt = Locale{Code5: LocaleCodePtPT, NativeTitle: "Português (PT)", EnglishTitle: "Portuguese (PT)", FlagIcon: "🇵🇹"}

	// LocalePtBr is locale
	LocalePtBr = Locale{Code5: LocaleCodePtBR, NativeTitle: "Português (BR)", EnglishTitle: "Portuguese (BR)", FlagIcon: "🇧🇷"}

	// LocaleFaIr is locale
	LocaleFaIr = Locale{Code5: LocaleCodeFaIR, IsRtl: true, NativeTitle: "فارسی", EnglishTitle: "Farsi", FlagIcon: "🇮🇷"}

	// LocaleKoKo is locale
	LocaleKoKo = Locale{Code5: LocaleCodeKoKO, NativeTitle: "한국어/조선말", EnglishTitle: "Korean", FlagIcon: "🇰🇷"}

	// LocaleJaJp is locale
	LocaleJaJp = Locale{Code5: LocaleCodeJaJP, NativeTitle: "日本語", EnglishTitle: "Japanese", FlagIcon: "🇯🇵"}

	// LocaleZhCn is locale
	LocaleZhCn = Locale{Code5: LocaleCodeZhCN, NativeTitle: "中文", EnglishTitle: "Chinese", FlagIcon: "🇨🇳"}
)

// LocalesByCode5 map of locales by 5-character code
var LocalesByCode5 = map[string]Locale{
	LocaleCodeEnUS: LocaleEnUS,
	//LocaleCodeEnUK: LocaleEnUK,
	LocalCodeRuRu: LocaleRuRu,
	// LOCALE_ID_ID: LocaleIdId,
	LocaleCodeDeDE: LocaleDeDe,
	LocaleCodeEsES: LocaleEsEs,
	LocaleCodeFrFR: LocaleFrFr,
	LocaleCodeItIT: LocaleItIt,
	LocaleCodePlPL: LocalePlPl,
	LocaleCodePtPT: LocalePtPt,
	LocaleCodePtBR: LocalePtBr,
	LocaleCodeFaIR: LocaleFaIr,
	LocaleCodeKoKO: LocaleKoKo,
	LocaleCodeJaJP: LocaleJaJp,
	LocaleCodeZhCN: LocaleZhCn,
}

// GetLocaleByCode5 returns locale by 5-character code
func GetLocaleByCode5(code5 string) Locale {
	if locale, ok := LocalesByCode5[code5]; ok {
		return locale
	}
	panic(fmt.Sprintf("Unknown locale: [%v]", code5))
}
