.PHONY: make install clean
make:
	@go build -o prometheus-racsgpfs-exporter main.go

install: make
	@cp prometheus-racsgpfs-exporter /usr/local/sbin/prometheus-racsgpfs-exporter
	@cp deploy/systemd/prometheus-racsgpfs-exporter.service /etc/systemd/system/prometheus-racsgpfs-exporter.service

clean:
	rm -f prometheus-racsgpfs-exporter
	rm -f /etc/systemd/system/prometheus-racsgpfs-exporter.service

