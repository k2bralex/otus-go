package app

import (
	"errors"
	"github.com/cheggaaa/pb"
	"io"
	"os"
)

func Run(src, trg *string, lim, off *int64) error {
	fileFrom, err := os.Open(*src)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	if _, err = os.Stat(*trg); err == nil {
		return errors.New("file exist already")
	}

	fileTo, err := os.OpenFile(*trg, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	if *lim == 0 {
		if err = limitHandler(src, lim); err != nil {
			return err
		}
	}

	buf := make([]byte, *lim)

	if *off != 0 {
		if _, err = fileFrom.Seek(*off, io.SeekStart); err != nil {
			return err
		}
	}

	bar := pb.StartNew(int(*lim))
	defer bar.Finish()

	for *off < *lim {
		read, err := fileFrom.Read(buf)
		*off += int64(read)
		if err != nil && err == io.EOF {
			return err
		}

		bar.Add(read)

		if read == 0 {
			break
		}

		if _, err = fileTo.Write(buf[:read]); err != nil {
			return err
		}
	}

	return nil
}

func limitHandler(path *string, lim *int64) error {
	info, err := os.Stat(*path)
	if err != nil {
		return err
	}
	*lim = info.Size()
	return nil
}
