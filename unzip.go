package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	errBadFile = errors.New("Bad file")
)

func unzip(f *tb.File, bot *tb.Bot) (string, error) {
	id := uuid.NewString()
	path := "tmp/" + id + ".zip"
	dest := "tmp/" + id
	err := bot.Download(f, path)
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return "", err
	}

	// Check for zip bombs
	var size uint64 = 0
	for _, v := range r.File {
		size += v.UncompressedSize64
	}
	if size > maxSize {
		return "", errBadFile
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return "", err
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return errBadFile
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return dest, nil
}
