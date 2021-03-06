package extract

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mholt/archiver"
	"github.com/yonkornilov/snowcapper/config"
	"github.com/yonkornilov/snowcapper/context"
)

func Run(c *context.Context, b config.Binary, downloadPath string) (binaryPath string, err error) {
	binaryPath = b.GetBinaryPath()

	if c.IsDryRun {
		fmt.Printf("DRY-RUN: Extracting %s from %s\n", b.Format, downloadPath)
		err, extractedPath := extract(c, b.Format, downloadPath)
		if err != nil {
			return "", err
		}
		fmt.Printf("DRY-RUN: Copying %s to %s\n", getExtractedBinaryPath(b, extractedPath), binaryPath)
		fmt.Printf("DRY-RUN: Removing %s\n", extractedPath)
		fmt.Printf("DRY-RUN: Removing %s\n", downloadPath)

		return binaryPath, nil
	}
	fmt.Printf("Extracting %s from %s\n", b.Format, downloadPath)
	err, extractedPath := extract(c, b.Format, downloadPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Copying %s to %s\n", getExtractedBinaryPath(b, extractedPath), binaryPath)
	err = copyToTarget(getExtractedBinaryPath(b, extractedPath), binaryPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Removing %s\n", extractedPath)
	err = os.RemoveAll(extractedPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Removing %s\n", downloadPath)
	err = os.RemoveAll(b.Src)
	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

func extract(c *context.Context, archiveType string, src string) (error, string) {
	extractedPath := getExtractedPath(archiveType, src)
	var err error
	if c.IsDryRun {
		if archiver.SupportedFormats[getArchiverFormat(archiveType)] == nil {
			err = errors.New(fmt.Sprintf("Error: 'Type' must be one of: %s", archiver.SupportedFormats))
		}
		if err != nil {
			return err, ""
		}
		return nil, extractedPath
	}
	if archiver.SupportedFormats[getArchiverFormat(archiveType)] == nil {
		err = errors.New(fmt.Sprintf("Error: 'Type' must be one of: %s", archiver.SupportedFormats))
	}
	if err != nil {
		return err, ""
	}
	err = archiver.SupportedFormats[getArchiverFormat(archiveType)].Open(src, extractedPath)
	if err != nil {
		return err, ""
	}

	return nil, extractedPath
}

func copyToTarget(src string, target string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}

func getArchiverFormat(archiveType string) string {
	if archiveType == "tar" {
		return "Tar"
	}
	if archiveType == "zip" {
		return "Zip"
	}
	if archiveType == "rar" {
		return "Rar"
	}
	if archiveType == "tar.xz" {
		return "TarXZ"
	}
	archiverFormat := strings.Replace(archiveType, "tar.", "", 1)
	archiverFormat = strings.Title(archiverFormat)
	return "Tar" + archiverFormat
}

func getExtractedPath(archiveType string, src string) string {
	extractedPath := strings.Replace(src, "."+archiveType, "", -1)
	return extractedPath
}

func getExtractedBinaryPath(b config.Binary, extractedPath string) string {
	return extractedPath + "/" + b.Name
}
