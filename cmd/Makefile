# Install pwmfan-go

PREFIX=/usr/local
CONFIG=/etc/pwmfan

main:
	go build -v -o pwmfan -ldflags "-s -w"

install:
	mkdir -p /etc/pwmfan/
	cp pwmfan $(PREFIX)/bin/pwmfan
	cp config.json $(CONFIG)/config.json
	cp systemd/pwmfan.service /etc/systemd/system/pwmfan.service
	chmod 0755 $(PREFIX)/bin/pwmfan
	chmod 0644 $(CONFIG)/config.json
	chmod 0644 /etc/systemd/system/pwmfan.service
	systemctl daemon-reload
	systemctl start pwmfan
	systemctl enable pwmfan

clean:
	rm pwmfan

uninstall:
	systemctl stop pwmfan
	systemctl disable pwmfan
	rm /etc/systemd/system/pwmfan.service
	rm /usr/local/bin/pwmfan
	rm -r /etc/pwmfan
