package abt

import (
	bencode "code.google.com/p/bencode-go"
	"io"
	"os"
)

type ABTFile struct {
	Sections SectionList
	Source   ABTFileSource
}

type SectionList []SectionData

type SectionData struct {
	Path     string
	Origin   string
	Start    int64
	Size     int64
	Filesize int64
}

type ABTFileSource interface {
	io.Reader
	io.ReaderAt
}

func getSectionList(input io.Reader) (list SectionList, err error) {

	list = make(SectionList, 0, 5)
	err = bencode.Unmarshal(input, &list)

	return
}

func NewABTFile(source ABTFileSource) (abtfile *ABTFile, err error) {
	sectionList, err := getSectionList(source)
	if err != nil {
		return
	}

	abtfile = &ABTFile{sectionList, source}

	return
}

func OpenABTFile(path *string) (abtfile *ABTFile, err error) {
	file, err := os.Open(*path)
	if err != nil {
		return
	}

	abtfile, err = NewABTFile(file)

	return
}

func (sections *SectionList) Write(target io.Writer) (err error) {

	return bencode.Marshal(target, *sections)
}

func (abtfile *ABTFile) Write(target io.Writer) (err error) {

	err = abtfile.Sections.Write(target)

	return
}
