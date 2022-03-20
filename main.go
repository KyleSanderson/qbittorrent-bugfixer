package main

import (
	"github.com/KyleSanderson/autobrr/pkg/qbittorrent"
	"strconv"
	"fmt"
	"os"
)

func main() {
	port, _ := strconv.ParseInt(os.Getenv("port"), 10, 16)
	set := qbittorrent.Settings {
		Hostname: os.Getenv("host"),
		Port: uint(port),
		Username: os.Getenv("user"),
		Password: os.Getenv("password"),
	}

	c := qbittorrent.NewClient(set)

	if err := c.Login(); err != nil {
		fmt.Printf("Unable to login: %q\n", err)
		return
	}

	torrents, err := c.GetTorrents()
	if err != nil {
		fmt.Printf("Unable to get Torrents: %q\n", err)
		return
	}

	fixStalled(c, torrents)
	fixForcedUpDown(c, torrents)
	fixErrored(c, torrents)
	fixAutoTMM(c, torrents)
	fixChecking(c, torrents)
	fixMoving(c, torrents)
}

func fixStalled(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.TotalSize == 0 || t.State != qbittorrent.TorrentStateStalledDl {
			continue
		}

		if (float64(t.AmountLeft) / float64(t.TotalSize)) > 0.1 {
			continue
		}

		fmt.Printf("Stalled Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.SetForceStart(hashes, true); err != nil {
		fmt.Printf("Bad forcestart: %q\n", err)
		return
	}
	
	if err := c.Recheck(hashes); err != nil {
		fmt.Printf("Bad recheck: %q\n", err)
		return
	}
}

func fixForcedUpDown(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.State != qbittorrent.TorrentStateForcedDl && t.State != qbittorrent.TorrentStateForcedUp {
			continue
		}

		fmt.Printf("Forced UpDown Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.Resume(hashes); err != nil {
		fmt.Printf("Bad resume: %q\n", err)
		return
	}
}

func fixErrored(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.State != qbittorrent.TorrentStateError {
			continue
		}

		fmt.Printf("Errored Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.Resume(hashes); err != nil {
		fmt.Printf("Bad resume: %q\n", err)
		return
	}

	if err := c.SetForceStart(hashes, true); err != nil {
		fmt.Printf("Bad resume: %q\n", err)
		return
	}
}

func fixAutoTMM(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.AutoManaged == true {
			continue
		}

		fmt.Printf("AutoManaged Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.SetAutoManagement(hashes, true); err != nil {
		fmt.Printf("Bad Management: %q\n", err)
		return
	}
}

func fixChecking(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.State != qbittorrent.TorrentStateCheckingUp && t.State != qbittorrent.TorrentStateCheckingDl {
			continue
		}

		fmt.Printf("Checking Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.Resume(hashes); err != nil {
		fmt.Printf("Bad resume: %q\n", err)
		return
	}

	if err := c.SetForceStart(hashes, true); err != nil {
		fmt.Printf("Bad forcestart: %q\n", err)
		return
	}
}

func fixMoving(c *qbittorrent.Client, torrents []qbittorrent.Torrent) {
	hashes := make([]string, 0, len(torrents))
	for _, t := range torrents {
		if t.State != qbittorrent.TorrentStateMoving {
			continue
		}

		fmt.Printf("Moving Torrent: %q\n", t.Name)
		hashes = append(hashes, t.Hash)
	}

	if err := c.Resume(hashes); err != nil {
		fmt.Printf("Bad resume: %q\n", err)
		return
	}

	if err := c.SetForceStart(hashes, true); err != nil {
		fmt.Printf("Bad forcestart: %q\n", err)
		return
	}
}
