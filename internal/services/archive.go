package services

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stjudewashere/seonaut/internal/archiver"
	"github.com/stjudewashere/seonaut/internal/models"
)

type ArchiveService struct {
	ArchiveDir string
}

func NewArchiveService(ad string) *ArchiveService {
	return &ArchiveService{
		ArchiveDir: ad,
	}
}

// ArchiveProject returns an archiver for the specified project.
// It returns an error if the archiver couldn't be created.
func (s *ArchiveService) GetArchiveWriter(p *models.Project) (*archiver.Writer, error) {
	return archiver.NewArchiver(s.getArchiveFile(p))
}

// ReadArchive reads an URLs WACZ record from a project's archive.
func (s *ArchiveService) ReadArchiveRecord(p *models.Project, urlStr string) *models.ArchiveRecord {
	waczPath := s.getArchiveFile(p)
	reader := archiver.NewReader(waczPath)

	content := reader.ReadArchive(urlStr)

	record := &models.ArchiveRecord{}
	index := strings.Index(content, "\r\n\r\n")
	if index != -1 {
		record.Headers = content[:index]
		record.Body = strings.TrimSpace(content[index+1:])
	} else {
		record.Headers = content
	}

	return record
}

// ArchiveExists checks if a wacz file exists for the current project.
// It returns true if it exists, otherwise it returns false.
func (s *ArchiveService) ArchiveExists(p *models.Project) bool {
	file := s.getArchiveFile(p)
	_, err := os.Stat(file)
	return err == nil
}

// DeleteArchive removes the wacz archive file for a given project.
// It checks if the file exists before removing it.
func (s *ArchiveService) DeleteArchive(p *models.Project) {
	if !s.ArchiveExists(p) {
		return
	}

	file := s.getArchiveFile(p)
	os.Remove(file)

	// Check if the archive directory is empty and remove it.
	dir := filepath.Dir(file)
	d, err := os.Open(dir)
	if err != nil {
		log.Printf("failed to open archive dir %s: %v", dir, err)
		return
	}

	_, err = d.ReadDir(1)
	if err == nil {
		return // dir is not empty.
	}

	err = os.Remove(dir)
	if err != nil {
		log.Printf("failed to remove empty archive dir %s: %v", dir, err)
	}
}

// GetArchiveFilePath returns the project's wacz file path if it exists,
// otherwise it returns an error.
func (s *ArchiveService) GetArchiveFilePath(p *models.Project) (string, error) {
	if !s.ArchiveExists(p) {
		return "", errors.New("WACZ archive file does not exist")
	}

	file := s.getArchiveFile(p)
	return file, nil
}

// getArchiveFile returns a string with the path to the project's WACZ file.
func (s *ArchiveService) getArchiveFile(p *models.Project) string {
	return s.ArchiveDir + "/" + strconv.FormatInt(p.Id, 10) + "/" + p.Host + ".wacz"
}
