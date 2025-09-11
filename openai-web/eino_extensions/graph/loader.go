package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino-ext/components/document/parser/docx"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino-ext/components/document/parser/xlsx"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/futugyousuzu/go-openai-web/models"
)

func getLoaderNode(ctx context.Context, node models.Node) (document.Loader, error) {
	if loader, ok := node.Data["loader"].(string); ok && len(loader) > 0 {
		var loaderParser parser.Parser = nil
		if p, ok := node.Data["parser"].(string); ok && len(p) > 0 {
			switch p {
			case "docx":
				loaderParser, _ = docx.NewDocxParser(ctx, &docx.Config{
					ToSections:      true,
					IncludeComments: true,
					IncludeHeaders:  true,
					IncludeFooters:  true,
					IncludeTables:   true,
				})
			case "pdf":
				loaderParser, _ = pdf.NewPDFParser(ctx, &pdf.Config{ToPages: true})
			case "xlsx":
				loaderParser, _ = xlsx.NewXlsxParser(ctx, &xlsx.Config{})
			}
		}

		switch loader {
		case "url":
			loaderConfig := &url.LoaderConfig{
				Parser: loaderParser,
			}

			return url.NewLoader(ctx, loaderConfig)
		}
	}

	return nil, fmt.Errorf("invalid document loader node: %s", node.ID)
}
