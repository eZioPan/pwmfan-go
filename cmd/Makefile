# Install pwmfan-go

PREFIX=/usr/local
CONFIG=/etc/pwmfan

main:
	go build -ldflag "-s -w" -o pwmfan

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
	rm -rf pwmfan

uninstall:
	systemctl stop pwmfan
	systemctl disable pwmfan
	rm -rf /etc/systemd/system/pwmfan.service
	rm -rf /usr/local/bin/pwmfan
	rm -rf /etc/pwmfan