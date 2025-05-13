package dataformats

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/xuri/excelize/v2"
)

type MsExcelDecoder struct {
	config *MsExcelDecoderConfig
}

func NewMsExcelDecoder(config *MsExcelDecoderConfig) *MsExcelDecoder {
	if config == nil {
		config = NewMsExcelDecoderConfig()
	}
	return &MsExcelDecoder{config: config}
}

// Decode implements dataformats.IContentDecoder.
func (m *MsExcelDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsExcelDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *MsExcelDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsExcelDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (d *MsExcelDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if d == nil {
		return nil, fmt.Errorf("MsExcelDecoder is nil")
	}
	f, err := excelize.OpenReader(stream)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_PlainText}
	var worksheetNumber int64 = 0

	for _, sheetName := range f.GetSheetList() {
		worksheetNumber++
		var sb strings.Builder
		if d.config.WithWorksheetNumber {
			sb.WriteString(strings.Replace(d.config.WorksheetNumberTemplate, "{number}", fmt.Sprintf("%d", worksheetNumber), -1))
			sb.WriteString("\n")
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		for _, row := range rows {
			sb.WriteString(d.config.RowPrefix)
			for i, cell := range row {
				if d.config.WithQuotes {
					sb.WriteString("\"")
					if cell == "" {
						sb.WriteString(d.config.BlankCellValue)
					} else {
						sb.WriteString(strings.ReplaceAll(cell, "\"", "\"\""))
					}
					sb.WriteString("\"")
				} else {
					if cell == "" {
						sb.WriteString(d.config.BlankCellValue)
					} else {
						sb.WriteString(cell)
					}
				}

				if i < len(row)-1 {
					sb.WriteString(d.config.ColumnSeparator)
				}
			}
			sb.WriteString(d.config.RowSuffix)
			sb.WriteString("\n")
		}

		if d.config.WithEndOfWorksheetMarker {
			sb.WriteString(strings.Replace(d.config.EndOfWorksheetMarkerTemplate, "{number}", fmt.Sprintf("%d", worksheetNumber), -1))
			sb.WriteString("\n")
		}

		sentencesAreComplete := true
		chunk := dataformats.Chunk{Content: text.NormalizeNewlines(sb.String(), true), Number: worksheetNumber, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)}
		result.Sections = append(result.Sections, chunk)
	}

	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (m *MsExcelDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_MsExcelX)
}

var _ dataformats.IContentDecoder = (*MsExcelDecoder)(nil)
