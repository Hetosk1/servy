Build Code
`go built -o servy main.go`

Move Code to Binaries 
`sudo mv my-lb /usr/local/bin/`

Create a Systemd Service 
`sudo touch /etc/systemd/system/servy.service`

```
[Unit]
Description=Go Load Balancer
After=network.target

[Service]
ExecStart=/usr/local/bin/my-lb
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

Start system
```
sudo systemctl daemon-reload
sudo systemctl enable servy
sudo systemctl start servy
sudo systemctl status servy
```

