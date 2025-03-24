package dataformats

import (
	"golang.org/x/text/language"
)

type MsExcelDecoderConfig struct {
	WithWorksheetNumber          bool
	WithEndOfWorksheetMarker     bool
	WithQuotes                   bool
	WorksheetNumberTemplate      string
	EndOfWorksheetMarkerTemplate string
	RowPrefix                    string
	ColumnSeparator              string
	RowSuffix                    string
	BlankCellValue               string
	BooleanTrueValue             string
	BooleanFalseValue            string
	TimeSpanFormat               string
	TimeSpanProvider             language.Tag
	DateFormat                   string
	DateFormatProvider           language.Tag
}

func NewMsExcelDecoderConfig() *MsExcelDecoderConfig {
	return &MsExcelDecoderConfig{
		WithWorksheetNumber:          true,
		WithEndOfWorksheetMarker:     false,
		WithQuotes:                   true,
		WorksheetNumberTemplate:      "\n# Worksheet {number}\n",
		EndOfWorksheetMarkerTemplate: "\n# End of worksheet {number}",
		RowPrefix:                    "",
		ColumnSeparator:              ", ",
		RowSuffix:                    "",
		BlankCellValue:               "",
		BooleanTrueValue:             "TRUE",
		BooleanFalseValue:            "FALSE",
		TimeSpanFormat:               "15:04:05",
		TimeSpanProvider:             language.English,
		DateFormat:                   "2006-01-02",
		DateFormatProvider:           language.English,
	}
}

type MsPowerPointDecoderConfig struct {
	SlideNumberTemplate      string
	EndOfSlideMarkerTemplate string
	WithSlideNumber          bool
	WithEndOfSlideMarker     bool
	SkipHiddenSlides         bool
}

func NewMsPowerPointDecoderConfig() *MsPowerPointDecoderConfig {
	return &MsPowerPointDecoderConfig{
		SlideNumberTemplate:      "# Slide {number}",
		EndOfSlideMarkerTemplate: "# End of slide {number}",
		WithSlideNumber:          true,
		WithEndOfSlideMarker:     false,
		SkipHiddenSlides:         true,
	}
}
