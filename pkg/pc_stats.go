package pcstat

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// page cache status
// Bytes: size of the file (from os.File.Stat())
// Pages: array of booleans: true if cached, false otherwise
type PcStatus struct {
	Name      string    `json:"name,omitempty"`
	Size      int64     `json:"size,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Mtime     time.Time `json:"mtime,omitempty"`
	Pages     int       `json:"pages,omitempty"`
	Cached    int       `json:"cached,omitempty"`
	Uncached  int       `json:"uncached,omitempty"`
	Percent   float64   `json:"percent,omitempty"`
	PPStat    []bool    `json:"pp_stat,omitempty"`
}

func GetPcStatus(fname string) (PcStatus, error) {
	pcs := PcStatus{Name: fname}

	f, err := os.Open(fname)
	if err != nil {
		return pcs, fmt.Errorf("could not open file for read: %v", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return pcs, fmt.Errorf("could not stat the file: %v", err)
	}
	if fi.IsDir() {
		return pcs, errors.New("file is a directory")
	}

	pcs.Size = fi.Size()
	pcs.Timestamp = time.Now()
	pcs.Mtime = fi.ModTime()

	pcs.PPStat, err = fileMincore(f, fi.Size())
	if err != nil {
		return pcs, err
	}
	for _, ok := range pcs.PPStat {
		if ok {
			pcs.Cached++
		}
	}
	pcs.Pages = len(pcs.PPStat)
	pcs.Uncached = pcs.Pages - pcs.Cached
	pcs.Percent = (float64(pcs.Cached) / float64(pcs.Pages)) * 100.00

	return pcs, nil
}
