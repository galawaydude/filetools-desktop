package doc

import (
	"context"

	"github.com/galawaydude/filetools-desktop/internal/tool"
)

func init() {
	tool.Register(tool.Define(tool.Config{
		ID:          "doc.wordtopdf",
		Name:        "Word to PDF",
		Description: "Turn a Word document (.docx or .doc) into a PDF. Needs the free LibreOffice program installed.",
		Category:    tool.CategoryDoc,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".docx", ".doc"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if err := cancelledCtx(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.3, "Converting to PDF…")
			out, err := libreConvert(ctx, req.Inputs[0], req.OutDir, "pdf", "pdf")
			if err != nil {
				return tool.Result{}, err
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}}, nil
		},
	}))

	tool.Register(tool.Define(tool.Config{
		ID:          "doc.pdftoword",
		Name:        "PDF to Word",
		Description: "Convert a PDF into an editable Word document. Needs LibreOffice. Results vary — complex layouts may not come out perfectly.",
		Category:    tool.CategoryDoc,
		Input:       tool.InputSingleFile,
		Extensions:  []string{".pdf"},
		Run: func(ctx context.Context, req tool.Request, p tool.Progress) (tool.Result, error) {
			if err := cancelledCtx(ctx); err != nil {
				return tool.Result{}, err
			}
			p.Update(0.3, "Converting to Word…")
			out, err := libreConvert(ctx, req.Inputs[0], req.OutDir, "docx:MS Word 2007 XML:UTF8", "docx")
			if err != nil {
				return tool.Result{}, err
			}
			p.Update(1, "Done")
			return tool.Result{Outputs: []string{out}, Message: "Converted using a best-effort layout. Please check the result — complex PDFs may need tidying up."}, nil
		},
	}))
}
