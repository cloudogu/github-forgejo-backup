//go:build linux

package disk

import (
	"fmt"
	"golang.org/x/sys/unix"
)

type Usage struct {
	TotalBytes     uint64
	FreeBytes      uint64 // free for unprivileged user (same as df "avail")
	UsedBytes      uint64
	UtilizationPct float64 // df-style: used / (used + avail) * 100
}

func UsageOf(path string) (Usage, error) {

	unix.Sync()

	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return Usage{}, fmt.Errorf("statfs %q: %w", path, err)
	}

	bsize := uint64(st.Bsize)
	blocksUsed := st.Blocks - st.Bfree // blocks used by anyone
	blocksAvail := st.Bavail           // blocks available to non-root
	total := st.Blocks * bsize
	free := blocksAvail * bsize
	used := blocksUsed * bsize

	den := float64(blocksUsed + blocksAvail)
	pct := 0.0
	if den > 0 {
		pct = float64(blocksUsed) / den * 100.0
	}

	return Usage{
		TotalBytes:     total,
		FreeBytes:      free,
		UsedBytes:      used,
		UtilizationPct: pct,
	}, nil

}
