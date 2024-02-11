package functions

type TechnicalIndicatorsType interface {
	privateTechnicalIndicatorsType()
	String() string
}

type technicalIndicatorsType string

func (c technicalIndicatorsType) privateTechnicalIndicatorsType() {}
func (c technicalIndicatorsType) String() string {
	return string(c)
}

const SMA technicalIndicatorsType = "SMA"
const EMA technicalIndicatorsType = "EMA"
const WMA technicalIndicatorsType = "WMA"
const DEMA technicalIndicatorsType = "DEMA"
const TEMA technicalIndicatorsType = "TEMA"
const TRIMA technicalIndicatorsType = "TRIMA"
const KAMA technicalIndicatorsType = "KAMA"
const MAMA technicalIndicatorsType = "MAMA"
const VWAP technicalIndicatorsType = "VWAP"
const T3 technicalIndicatorsType = "T3"
const MACD technicalIndicatorsType = "MACD"
const MACDEXT technicalIndicatorsType = "MACDEXT"
const STOCH technicalIndicatorsType = "STOCH"
const STOCHF technicalIndicatorsType = "STOCHF"
const RSI technicalIndicatorsType = "RSI"
const STOCHRSI technicalIndicatorsType = "STOCHRSI"
const WILLR technicalIndicatorsType = "WILLR"
const ADX technicalIndicatorsType = "ADX"
const ADXR technicalIndicatorsType = "ADXR"
const APO technicalIndicatorsType = "APO"
const PPO technicalIndicatorsType = "PPO"
const MOM technicalIndicatorsType = "MOM"
const BOP technicalIndicatorsType = "BOP"
const CCI technicalIndicatorsType = "CCI"
const CMO technicalIndicatorsType = "CMO"
const ROC technicalIndicatorsType = "ROC"
const ROCR technicalIndicatorsType = "ROCR"
const AROON technicalIndicatorsType = "AROON"
const AROONOSC technicalIndicatorsType = "AROONOSC"
const MFI technicalIndicatorsType = "MFI"
const TRIX technicalIndicatorsType = "TRIX"
const ULTOSC technicalIndicatorsType = "ULTOSC"
const DX technicalIndicatorsType = "DX"
const MINUS_DI technicalIndicatorsType = "MINUS_DI"
const PLUS_DI technicalIndicatorsType = "PLUS_DI"
const MINUS_DM technicalIndicatorsType = "MINUS_DM"
const PLUS_DM technicalIndicatorsType = "PLUS_DM"
const BBANDS technicalIndicatorsType = "BBANDS"
const MIDPOINT technicalIndicatorsType = "MIDPOINT"
const MIDPRICE technicalIndicatorsType = "MIDPRICE"
const SAR technicalIndicatorsType = "SAR"
const TRANGE technicalIndicatorsType = "TRANGE"
const ATR technicalIndicatorsType = "ATR"
const NATR technicalIndicatorsType = "NATR"
const AD technicalIndicatorsType = "AD"
const ADOSC technicalIndicatorsType = "ADOSC"
const OBV technicalIndicatorsType = "OBV"
const HT_TRENDLINE technicalIndicatorsType = "HT_TRENDLINE"
const HT_SINE technicalIndicatorsType = "HT_SINE"
const HT_TRENDMODE technicalIndicatorsType = "HT_TRENDMODE"
const HT_DCPERIOD technicalIndicatorsType = "HT_DCPERIOD"
const HT_DCPHASE technicalIndicatorsType = "HT_DCPHASE"
const HT_PHASOR technicalIndicatorsType = "HT_PHASOR"
