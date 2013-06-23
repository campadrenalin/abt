package abt

import (
	bencode "code.google.com/p/bencode-go"
	"errors"
	"io"
    "os"
)

type SectionData struct {
	Path     string
	Start    int64
	Size     int64
	Filesize int64
}

func parseSectionData(data interface{}) (sdstruct SectionData, err error) {
    // Convert interface{} to map
    toMap, ok := data.(map[string]interface{})
    if !ok {
        err = errors.New("Item in section map was not a map")
        return
    }

    // Ensure all necessary fields are present
    fields := [...]string {"path","start","size","filesize"}
    for _, fieldname := range fields {
        _, ok := toMap[fieldname]
        if !ok {
            err = errors.New("Missing field: " + fieldname)
        }
    }

    // Extract fields
    raw_path     := toMap["path"]
    raw_start    := toMap["start"]
    raw_size     := toMap["size"]
    raw_filesize := toMap["filesize"]

    path,  ok_path  := raw_path.(string)
    start, ok_start := raw_start.(int64)
    size,  ok_size  := raw_size.(int64)
    filesize, ok_fsize := raw_filesize.(int64)

    if !(ok_path && ok_start && ok_size && ok_fsize) {
        err = errors.New("A field was the wrong type")
        return
    }

    sdstruct.Path  = path
    sdstruct.Start = start
    sdstruct.Size  = size
    sdstruct.Filesize  = filesize

    return
}

func getSectionList(input io.Reader) (data []SectionData, err error) {
	var decoded interface{}
	decoded, err = bencode.Decode(input)
	if err != nil {
		return
	}
	toList, ok := decoded.([]interface{})
	if !ok {
		err = errors.New("Could not parse section map as list")
        return
	}

    len := cap(toList)
    data = make([]SectionData, len, len)

	for index, value := range toList {
        var sdstruct SectionData

		sdstruct, err = parseSectionData(value)
        if err != nil {
            return
        }
		data[index]  = sdstruct
	}
	return
}

type ABTFileSource interface {
    io.Reader
    io.ReaderAt
}

type ABTFile struct {
    SectionList []SectionData
    Source ABTFileSource
}

func NewABTFile (source ABTFileSource) (abtfile *ABTFile, err error) {
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
