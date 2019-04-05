package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

const (
	MIN_FILE_SIZE     = 1           // bytes
	MAX_FILE_SIZE     = 1024 * 1024 // bytes
	IMAGE_TYPES       = "image/(gif|p?jpeg|(x-)?png)"
	ACCEPT_FILE_TYPES = IMAGE_TYPES
)

const (
	SUCCESS           = 0
	INVALID_FILE_SIZE = 1
	INVALID_FILE_TYPE = 2
)

var (
	imageTypes      = regexp.MustCompile(IMAGE_TYPES)
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type FileInfo struct {
	Name      string
	Type      string
	Size      int64
	ErrorCode int8
	ErrorMsg  string
}

type RspJson struct {
	RetCode        int8   `json:"retCode"`
	ErrMsg         string `json:"errMsg"`
	SrcImgPath     string `json:"srcImgPath"`
	InitSegImgPath string `json:"initSegImgPath"`
	DstImgPath     string `json:"dstImgPath"`
}

func (fi *FileInfo) ValidateType() (valid bool) {
	if imageTypes.MatchString(fi.Type) {
		return true
	}
	fi.ErrorCode = INVALID_FILE_TYPE
	fi.ErrorMsg = "File type not allowed"

	return false
}

func (fi *FileInfo) ValidateSize() (valid bool) {
	if fi.Size < MIN_FILE_SIZE {
		fi.ErrorCode = INVALID_FILE_SIZE
		fi.ErrorMsg = "File is too small"
	} else if fi.Size > MAX_FILE_SIZE {
		fi.ErrorCode = INVALID_FILE_SIZE
		fi.ErrorMsg = "File is too big"
	} else {
		return true
	}
	return false
}

func handleUpload(r *http.Request, p *multipart.Part) (fi *FileInfo) {
	fi = &FileInfo{
		Name:      p.FileName(),
		Type:      p.Header.Get("Content-Type"),
		ErrorCode: SUCCESS,
		ErrorMsg:  "",
	}

	if !fi.ValidateType() {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			fi.ErrorMsg = rec.(error).Error()
		}
	}()

	newFile, err := os.Create("upload_files/" + fi.Name)
	if err != nil {
		log.Fatal(err)
	}

	lr := &io.LimitedReader{R: p, N: MAX_FILE_SIZE + 1}
	_, err = io.Copy(newFile, lr)
	check(err)

	fi.Size = MAX_FILE_SIZE + 1 - lr.N
	if !fi.ValidateSize() {
		return
	}

	err = newFile.Sync()
	return
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<24)) // Copy max: 16 MiB
	return b.String()
}

func handleUploads(r *http.Request) (fileInfo *FileInfo) {
	fileInfo = new(FileInfo)
	mr, err := r.MultipartReader()
	check(err)

	r.Form, err = url.ParseQuery(r.URL.RawQuery)
	check(err)

	part, err := mr.NextPart()
	check(err)

	if name := part.FormName(); name != "" {
		if part.FileName() != "" {
			fileInfo = handleUpload(r, part)
		} else {
			r.Form[name] = append(r.Form[name], getFormValue(part))
		}
	}
	return
}