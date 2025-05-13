package pipeline

const (
	MimeTypes_PlainText           = "text/plain"
	MimeTypes_MarkDown            = "text/markdown"
	MimeTypes_MarkDownOld1        = "text/x-markdown"
	MimeTypes_MarkDownOld2        = "text/plain-markdown"
	MimeTypes_Html                = "text/html"
	MimeTypes_XHTML               = "application/xhtml+xml"
	MimeTypes_XML                 = "application/xml"
	MimeTypes_XML2                = "text/xml"
	MimeTypes_JSONLD              = "application/ld+json"
	MimeTypes_CascadingStyleSheet = "text/css"
	MimeTypes_JavaScript          = "text/javascript"
	MimeTypes_BourneShellScript   = "application/x-sh"

	MimeTypes_ImageBmp  = "image/bmp"
	MimeTypes_ImageGif  = "image/gif"
	MimeTypes_ImageJpeg = "image/jpeg"
	MimeTypes_ImagePng  = "image/png"
	MimeTypes_ImageTiff = "image/tiff"
	MimeTypes_ImageWebP = "image/webp"
	MimeTypes_ImageSVG  = "image/svg+xml"

	MimeTypes_WebPageUrl          = "text/x-uri"
	MimeTypes_TextEmbeddingVector = "float[]"
	MimeTypes_Json                = "application/json"
	MimeTypes_CSVData             = "text/csv"

	MimeTypes_Pdf         = "application/pdf"
	MimeTypes_RTFDocument = "application/rtf"

	MimeTypes_MsWord        = "application/msword"
	MimeTypes_MsWordX       = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MimeTypes_MsPowerPoint  = "application/vnd.ms-powerpoint"
	MimeTypes_MsPowerPointX = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MimeTypes_MsExcel       = "application/vnd.ms-excel"
	MimeTypes_MsExcelX      = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	MimeTypes_OpenDocumentText         = "application/vnd.oasis.opendocument.text"
	MimeTypes_OpenDocumentSpreadsheet  = "application/vnd.oasis.opendocument.spreadsheet"
	MimeTypes_OpenDocumentPresentation = "application/vnd.oasis.opendocument.presentation"
	MimeTypes_ElectronicPublicationZip = "application/epub+zip"

	MimeTypes_AudioAAC      = "audio/aac"
	MimeTypes_AudioMP3      = "audio/mpeg"
	MimeTypes_AudioWaveform = "audio/wav"
	MimeTypes_AudioOGG      = "audio/ogg"
	MimeTypes_AudioOpus     = "audio/opus"
	MimeTypes_AudioWEBM     = "audio/webm"

	MimeTypes_VideoMP4        = "video/mp4"
	MimeTypes_VideoMPEG       = "video/mpeg"
	MimeTypes_VideoOGG        = "video/ogg"
	MimeTypes_VideoOGGGeneric = "application/ogg"
	MimeTypes_VideoWEBM       = "video/webm"

	MimeTypes_ArchiveTar  = "application/x-tar"
	MimeTypes_ArchiveGzip = "application/gzip"
	MimeTypes_ArchiveZip  = "application/zip"
	MimeTypes_ArchiveRar  = "application/vnd.rar"
	MimeTypes_Archive7Zip = "application/x-7z-compressed"
)

const (
	FileExtensions_PlainText = ".txt"
	FileExtensions_MarkDown  = ".md"

	FileExtensions_Htm                 = ".htm"
	FileExtensions_Html                = ".html"
	FileExtensions_XHTML               = ".xhtml"
	FileExtensions_XML                 = ".xml"
	FileExtensions_JSONLD              = ".jsonld"
	FileExtensions_CascadingStyleSheet = ".css"
	FileExtensions_JavaScript          = ".js"
	FileExtensions_BourneShellScript   = ".sh"

	FileExtensions_ImageBmp   = ".bmp"
	FileExtensions_ImageGif   = ".gif"
	FileExtensions_ImageJpeg  = ".jpeg"
	FileExtensions_ImageJpg   = ".jpg"
	FileExtensions_ImagePng   = ".png"
	FileExtensions_ImageTiff  = ".tiff"
	FileExtensions_ImageTiff2 = ".tif"
	FileExtensions_ImageWebP  = ".webp"
	FileExtensions_ImageSVG   = ".svg"

	FileExtensions_WebPageUrl          = ".url"
	FileExtensions_TextEmbeddingVector = ".text_embedding"
	FileExtensions_Json                = ".json"
	FileExtensions_CSVData             = ".csv"

	FileExtensions_Pdf         = ".pdf"
	FileExtensions_RTFDocument = ".rtf"

	FileExtensions_MsWord        = ".doc"
	FileExtensions_MsWordX       = ".docx"
	FileExtensions_MsPowerPoint  = ".ppt"
	FileExtensions_MsPowerPointX = ".pptx"
	FileExtensions_MsExcel       = ".xls"
	FileExtensions_MsExcelX      = ".xlsx"

	FileExtensions_OpenDocumentText         = ".odt"
	FileExtensions_OpenDocumentSpreadsheet  = ".ods"
	FileExtensions_OpenDocumentPresentation = ".odp"
	FileExtensions_ElectronicPublicationZip = ".epub"

	FileExtensions_AudioAAC      = ".aac"
	FileExtensions_AudioMP3      = ".mp3"
	FileExtensions_AudioWaveform = ".wav"
	FileExtensions_AudioOGG      = ".oga"
	FileExtensions_AudioOpus     = ".opus"
	FileExtensions_AudioWEBM     = ".weba"

	FileExtensions_VideoMP4        = ".mp4"
	FileExtensions_VideoMPEG       = ".mpeg"
	FileExtensions_VideoOGG        = ".ogv"
	FileExtensions_VideoOGGGeneric = ".ogx"
	FileExtensions_VideoWEBM       = ".webm"

	FileExtensions_ArchiveTar  = ".tar"
	FileExtensions_ArchiveGzip = ".gz"
	FileExtensions_ArchiveZip  = ".zip"
	FileExtensions_ArchiveRar  = ".rar"
	FileExtensions_Archive7Zip = ".7z"
)

var FileExtensionToMimeType = map[string]string{
	FileExtensions_PlainText: MimeTypes_PlainText,
	FileExtensions_MarkDown:  MimeTypes_MarkDown,

	FileExtensions_Htm:                 MimeTypes_Html,
	FileExtensions_Html:                MimeTypes_Html,
	FileExtensions_XHTML:               MimeTypes_XHTML,
	FileExtensions_XML:                 MimeTypes_XML,
	FileExtensions_JSONLD:              MimeTypes_JSONLD,
	FileExtensions_CascadingStyleSheet: MimeTypes_CascadingStyleSheet,
	FileExtensions_JavaScript:          MimeTypes_JavaScript,
	FileExtensions_BourneShellScript:   MimeTypes_BourneShellScript,

	FileExtensions_ImageBmp:   MimeTypes_ImageBmp,
	FileExtensions_ImageGif:   MimeTypes_ImageGif,
	FileExtensions_ImageJpeg:  MimeTypes_ImageJpeg,
	FileExtensions_ImageJpg:   MimeTypes_ImageJpeg,
	FileExtensions_ImagePng:   MimeTypes_ImagePng,
	FileExtensions_ImageTiff:  MimeTypes_ImageTiff,
	FileExtensions_ImageTiff2: MimeTypes_ImageTiff,
	FileExtensions_ImageWebP:  MimeTypes_ImageWebP,
	FileExtensions_ImageSVG:   MimeTypes_ImageSVG,

	FileExtensions_WebPageUrl:          MimeTypes_WebPageUrl,
	FileExtensions_TextEmbeddingVector: MimeTypes_TextEmbeddingVector,
	FileExtensions_Json:                MimeTypes_Json,
	FileExtensions_CSVData:             MimeTypes_CSVData,

	FileExtensions_Pdf:         MimeTypes_Pdf,
	FileExtensions_RTFDocument: MimeTypes_RTFDocument,

	FileExtensions_MsWord:        MimeTypes_MsWord,
	FileExtensions_MsWordX:       MimeTypes_MsWordX,
	FileExtensions_MsPowerPoint:  MimeTypes_MsPowerPoint,
	FileExtensions_MsPowerPointX: MimeTypes_MsPowerPointX,
	FileExtensions_MsExcel:       MimeTypes_MsExcel,
	FileExtensions_MsExcelX:      MimeTypes_MsExcelX,

	FileExtensions_OpenDocumentText:         MimeTypes_OpenDocumentText,
	FileExtensions_OpenDocumentSpreadsheet:  MimeTypes_OpenDocumentSpreadsheet,
	FileExtensions_OpenDocumentPresentation: MimeTypes_OpenDocumentPresentation,
	FileExtensions_ElectronicPublicationZip: MimeTypes_ElectronicPublicationZip,

	FileExtensions_AudioAAC:      MimeTypes_AudioAAC,
	FileExtensions_AudioMP3:      MimeTypes_AudioMP3,
	FileExtensions_AudioWaveform: MimeTypes_AudioWaveform,
	FileExtensions_AudioOGG:      MimeTypes_AudioOGG,
	FileExtensions_AudioOpus:     MimeTypes_AudioOpus,
	FileExtensions_AudioWEBM:     MimeTypes_AudioWEBM,

	FileExtensions_VideoMP4:        MimeTypes_VideoMP4,
	FileExtensions_VideoMPEG:       MimeTypes_VideoMPEG,
	FileExtensions_VideoOGG:        MimeTypes_VideoOGG,
	FileExtensions_VideoOGGGeneric: MimeTypes_VideoOGGGeneric,
	FileExtensions_VideoWEBM:       MimeTypes_VideoWEBM,

	FileExtensions_ArchiveTar:  MimeTypes_ArchiveTar,
	FileExtensions_ArchiveGzip: MimeTypes_ArchiveGzip,
	FileExtensions_ArchiveZip:  MimeTypes_ArchiveZip,
	FileExtensions_ArchiveRar:  MimeTypes_ArchiveRar,
	FileExtensions_Archive7Zip: MimeTypes_Archive7Zip,
}
