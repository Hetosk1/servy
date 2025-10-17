# Servy 

A open-source server configuration tool written in go, i.e. can be used for load balancing, reverse proxy and static site serving and a lot more.

The whole project runs on a `YAML` configuration file which is located at `/etc/servy/servy.yaml`


## Steps for setting up Servy in Ubuntu based distributions

1. Build Code 

    ```
    go built -o servy main.go
    ```

2. Move Code to Binaries 

    ```
    sudo mv servy /usr/local/bin/
    ```

3. Create a `Systemd` Service 

    ```
    sudo touch /etc/systemd/system/servy.service
    sudo touch /etc/servy/servy.yaml
    ```

4. Configure the `systemd` Service

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

5. Start system

    ```
    sudo systemctl daemon-reload
    sudo systemctl enable servy
    sudo systemctl start servy
    sudo systemctl status servy
    ```


## Services Offered currently 
1. Load Balancer
2. Reverse Proxy
3. Static Site Serving 

### Settings up Load Balancer 

```
loadbalancer: 
    service: "on" #on/off 
    port: "80" #enter port number
    servers: 
        - "url1"
        - "url2"
        - "urlN"
```
### Settings up Reverse Proxy 

```
reverseproxy: 
    service: "on" #on/off 
    port: "81" #enter port number
    proxies:
        - proxy: "/pathforproxy"
          address: "addressforproxydivert"
        - proxy: "/pathforproxy"
          address: "addressforproxydivert"
```
