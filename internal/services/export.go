package services

import (
	"encoding/csv"
	"io"

	"github.com/stjudewashere/seonaut/internal/models"
)

type (
	ExportRepository interface {
		ExportLinks(*models.Crawl) <-chan *models.ExportLink
		ExportExternalLinks(*models.Crawl) <-chan *models.ExportLink
		ExportImages(crawl *models.Crawl) <-chan *models.ExportImage
		ExportScripts(crawl *models.Crawl) <-chan *models.Script
		ExportStyles(crawl *models.Crawl) <-chan *models.Style
		ExportIframes(crawl *models.Crawl) <-chan *models.Iframe
		ExportAudios(crawl *models.Crawl) <-chan *models.Audio
		ExportVideos(crawl *models.Crawl) <-chan *models.ExportVideo
		ExportHreflangs(crawl *models.Crawl) <-chan *models.ExportHreflang
	}

	Exporter struct {
		repository ExportRepository
	}
)

func NewExporter(r ExportRepository) *Exporter {
	return &Exporter{
		repository: r,
	}
}

// Export internal links as a CSV file
func (e *Exporter) ExportLinks(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Destination",
		"Text",
	})

	lStream := e.repository.ExportLinks(crawl)

	for v := range lStream {
		w.Write([]string{
			v.Origin,
			v.Destination,
			v.Text,
		})
	}

	w.Flush()
}

// Export internal links as a CSV file
func (e *Exporter) ExportExternalLinks(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Destination",
		"Text",
	})

	lStream := e.repository.ExportExternalLinks(crawl)

	for v := range lStream {
		w.Write([]string{
			v.Origin,
			v.Destination,
			v.Text,
		})
	}

	w.Flush()
}

// Export all images as a CSV file
func (e *Exporter) ExportImages(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Image URL",
		"Alt",
	})

	iStream := e.repository.ExportImages(crawl)

	for v := range iStream {
		w.Write([]string{
			v.Origin,
			v.Image,
			v.Alt,
		})
	}

	w.Flush()
}

// Export all scripts as a CSV file
func (e *Exporter) ExportScripts(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Script URL",
	})

	sStream := e.repository.ExportScripts(crawl)

	for v := range sStream {
		w.Write([]string{
			v.Origin,
			v.Script,
		})
	}

	w.Flush()
}

// Export all CSS styles as a CSV file
func (e *Exporter) ExportStyles(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Style URL",
	})

	sStream := e.repository.ExportStyles(crawl)

	for v := range sStream {
		w.Write([]string{
			v.Origin,
			v.Style,
		})
	}

	w.Flush()
}

// Export all CSS styles as a CSV file
func (e *Exporter) ExportIframes(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Iframe URL",
	})

	vStream := e.repository.ExportIframes(crawl)

	for v := range vStream {
		w.Write([]string{
			v.Origin,
			v.Iframe,
		})
	}

	w.Flush()
}

// Export all audio as a CSV file
func (e *Exporter) ExportAudios(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Audio URL",
	})

	vStream := e.repository.ExportAudios(crawl)

	for v := range vStream {
		w.Write([]string{
			v.Origin,
			v.Audio,
		})
	}

	w.Flush()
}

// Export all video as a CSV file
func (e *Exporter) ExportVideos(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Video URL",
	})

	vStream := e.repository.ExportVideos(crawl)

	for v := range vStream {
		w.Write([]string{
			v.Origin,
			v.Video,
		})
	}

	w.Flush()
}

// Export all hreflangs as a CSV file
func (e *Exporter) ExportHreflangs(f io.Writer, crawl *models.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Origin Language",
		"Hreflang",
		"Hreflang Language",
	})

	vStream := e.repository.ExportHreflangs(crawl)

	for v := range vStream {
		w.Write([]string{
			v.Origin,
			v.OriginLang,
			v.Hreflang,
			v.HreflangLang,
		})
	}

	w.Flush()
}
