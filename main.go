package main

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: add home USR filesets

var err error

type Collector struct {
	sizeGB  *prometheus.Desc
	quotaGB *prometheus.Desc
	inodes  *prometheus.Desc
}

type FilesetInfo struct {
	name    string
	sizeGB  float64
	quotaGB float64
	inodes  float64
	fstype  string
}

func main() {

	// set up logging
	lvl := slog.LevelInfo
	_, found := os.LookupEnv("RACSGPFS_EXPORTER_DEBUG")
	if found {
		lvl = slog.LevelDebug
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(l)
	slog.Debug("debug logging enabled")

	listenAddress, found := os.LookupEnv("RACSGPFS_EXPORTER_LISTEN_ADDRESS")
	if !found {
		listenAddress = ":8030"
	}

	r := prometheus.NewRegistry()

	r.MustRegister(NewCollector())

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func NewCollector() *Collector {
	labels := []string{"name", "fstype"}
	return &Collector{
		sizeGB:  prometheus.NewDesc("racsgpfs_size_gb", "Current fileset size in GB", labels, nil),
		quotaGB: prometheus.NewDesc("racsgpfs_quota_gb", "Current fileset quota in GB", labels, nil),
		inodes:  prometheus.NewDesc("racsgpfs_inodes", "Current fileset inode count", labels, nil),
	}
}

func (ac *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ac.sizeGB
	ch <- ac.inodes
}

func (ac *Collector) Collect(ch chan<- prometheus.Metric) {
	vis, err := ParseGPFS()
	if err != nil {
		slog.Error("failed to parse gpfs data", "error", err)
		return
	}
	for _, vi := range vis {
		ch <- prometheus.MustNewConstMetric(ac.sizeGB, prometheus.GaugeValue, vi.sizeGB, vi.name, vi.fstype)
		ch <- prometheus.MustNewConstMetric(ac.quotaGB, prometheus.GaugeValue, vi.quotaGB, vi.name, vi.fstype)
		ch <- prometheus.MustNewConstMetric(ac.inodes, prometheus.GaugeValue, vi.inodes, vi.name, vi.fstype)
	}
}

func ParseGPFS() ([]*FilesetInfo, error) {
	var fsInfos []*FilesetInfo
	cmd := exec.Command("sh", "-c", `/usr/lpp/mmfs/bin/mmrepquota --block-size g fs1 | grep -v GRP | grep -v "Block Limits" | grep -v "in_doubt" | awk '{printf "%s|%s|%s|%s|%s\n", $3, $1, $4, $5, $10}'`)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	data := buf.String()
	for _, line := range strings.Split(data, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		info := parseLine(line)
		if info == nil {
			continue
		}
		fsInfos = append(fsInfos, info)
	}
	return fsInfos, nil
}

func parseLine(l string) *FilesetInfo {
	elems := strings.Split(l, "|")
	fstype := elems[0]
	name := elems[1]
	sizeGBstr := elems[2]
	sizeGB, _ := strconv.ParseFloat(sizeGBstr, 64)
	quotaGBstr := elems[3]
	quotaGB, _ := strconv.ParseFloat(quotaGBstr, 64)
	inodesStr := elems[4]
	inodes, _ := strconv.ParseFloat(inodesStr, 64)
	return &FilesetInfo{name: name, sizeGB: sizeGB, quotaGB: quotaGB, inodes: inodes, fstype: fstype}
}
