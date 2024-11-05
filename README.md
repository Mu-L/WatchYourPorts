[![Main-Docker](https://github.com/aceberg/watchyourports/actions/workflows/main-docker.yml/badge.svg)](https://github.com/aceberg/watchyourports/actions/workflows/main-docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aceberg/watchyourports)](https://goreportcard.com/report/github.com/aceberg/watchyourports)
[![Maintainability](https://api.codeclimate.com/v1/badges/e8f67994120fc7936aeb/maintainability)](https://codeclimate.com/github/aceberg/WatchYourPorts/maintainability)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/aceberg/watchyourports)

<h1><a href="https://github.com/aceberg/watchyourports">
    <img src="https://raw.githubusercontent.com/aceberg/watchyourports/main/assets/logo.png" width="35" />
</a>WatchYourPorts</h1>

Open ports inventory for local servers. Exports data to InfluxDB2/Grafana 

- [Quick start](https://github.com/aceberg/watchyourports#quick-start)
- [Login](https://github.com/aceberg/watchyourports#login)
- [Import ports from Docker](https://github.com/aceberg/watchyourports#import-ports-from-docker)
- [Config](https://github.com/aceberg/watchyourports#config)
- [Options](https://github.com/aceberg/watchyourports#options)
- [Local network only](https://github.com/aceberg/watchyourports#local-network-only)
- [API](https://github.com/aceberg/watchyourports#api)
- [Thanks](https://github.com/aceberg/watchyourports#thanks)


![Screenshot](https://raw.githubusercontent.com/aceberg/WatchYourPorts/main/assets/Screenshot1.png)   
<details>
  <summary>More screenshots</summary>
  <img src="https://raw.githubusercontent.com/aceberg/WatchYourPorts/main/assets/Screenshot2.png">
  <img src="https://raw.githubusercontent.com/aceberg/WatchYourPorts/main/assets/Screenshot3.png">
</details> 

## Quick start

```sh
docker run --name wyp \
-e "TZ=Asia/Novosibirsk" \
-v ~/.dockerdata/WatchYourPorts:/data/WatchYourPorts \
-p 8853:8853 \
aceberg/watchyourports
```
Or use [docker-compose.yml](docker-compose.yml)


## Auth
You can limit access to WYP with [ForAuth](https://github.com/aceberg/ForAuth). Here is an example: [docker-compose-auth.yml](docker-compose-auth.yml)   
Also, SSO tools like Authelia should work.

## Import ports from Docker
1. Run [docker-export.sh](configs/docker-export.sh) on a server, where Docker is installed. `$ADDR` is IP or domain name of the server, without `http(s)://` prefix. It will be used to ping ports.
```sh
./docker-export.sh $ADDR
```   
2. Paste the output to `hosts.yaml` file in WatchYourPorts config dir
3. You can add as many servers to `hosts.yaml`, as you want


## Config


Configuration can be done through `config.yaml` file or GUI, or environment variables

| Variable  | Description | Default |
| --------  | ----------- | ------- |
| HOST | Listen address | 0.0.0.0 |
| PORT   | Port for web GUI | 8853 |
| THEME | Any theme name from https://bootswatch.com in lowcase or [additional](https://github.com/aceberg/aceberg-bootswatch-fork) | grass |
| COLOR | Background color: light or dark | dark |
| TIMEOUT | How often watched ports are scanned (minutes) | 10 |
| HIST_TRIM | How many port states are saved in memory and displayed | 90 |
| TZ | Set your timezone for correct time | "" |

### InfluxDB2 config
This config matches Grafana's config for InfluxDB data source

| Variable  | Description | Default | Example |
| --------  | ----------- | ------- | ------- |
| INFLUX_ENABLE | Enable export to InfluxDB2 | false | true |
| INFLUX_SKIP_TLS | Skip TLS Verify | false | true |
| INFLUX_ADDR | Address:port of InfluxDB2 server | | https://192.168.2.3:8086/ |
| INFLUX_BUCKET | InfluxDB2 bucket | | test |
| INFLUX_ORG | InfluxDB2 org | | home |
| INFLUX_TOKEN | Secret token, generated by InfluxDB2 | | |

## Options

| Key  | Description | Default | 
| --------  | ----------- | ------- | 
| -d | Path to config dir | /data/WatchYourPorts | 
| -n | Path to local JS and Themes ([node-bootstrap](https://github.com/aceberg/my-dockerfiles/tree/main/node-bootstrap)) | "" | 

## Local network only
By default, this app pulls themes, icons and fonts from the internet. But, in some cases, it may be useful to have an independent from global network setup. I created a separate [image](https://github.com/aceberg/my-dockerfiles/tree/main/node-bootstrap) with all necessary modules and fonts.    
```sh
docker run --name node-bootstrap       \
    -v ~/.dockerdata/icons:/app/icons  \ # For local images
    -p 8850:8850                       \
    aceberg/node-bootstrap
```
```sh
docker run --name wyp \
    -v ~/.dockerdata/WatchYourPorts:/data/WatchYourPorts \
    -p 8853:8853 \
    aceberg/watchyourports -n "http://$YOUR_IP:8850"
```
Or use [docker-compose](docker-compose-local.yml)

## API
```http
GET /api/all
```
Returns all data about saved addresses in `json`.
<details>
  <summary>Response example</summary>
  
```json
{
    "192.168.2.2": {
        "Name": "SomeAddrName",
        "Addr": "192.168.2.2",
        "PortMap": {},  // All saved ports will be here
        "Total": 0,
        "Watching": 0,
        "Online": 0,
        "Offline": 0
    },
}
```
</details><br>   

```http
GET /api/history
```
All history data from memory.
<details>
  <summary>Response example</summary>
  
```json
{
"192.168.2.3:8849": {
        "Name": "OS",
        "Addr": "192.168.2.3",
        "Port": 8849,
        "PortName": "MiniBoard",
        "State": [
            {
                "Date": "2024-06-28 22:42:45",
                "State": true
            },
            {
                "Date": "2024-06-28 22:52:45",
                "State": true
            }
        ],
        "NowState": true
    },
}
```
</details><br> 

```http
GET /api/port/:addr
```
Returns current PortMap for `addr`.    
<details>
  <summary>Request example</summary>

```bash
curl http://0.0.0.0:8853/api/port/192.168.2.2
```
</details>
<details>
  <summary>Response example</summary>
  
```json
{
    "8850": {
        "Name": "node-bootstrap",
        "Port": 8850,
        "State": true,
        "Watch": true
    },
    "8851": {
        "Name": "Exercise Diary",
        "Port": 8851,
        "State": true,
        "Watch": true
    },

}
```
</details><br>  

```http
GET /api/port/:addr/:port
```
Gets state of one port
<details>
  <summary>Request example</summary>

```bash
curl http://0.0.0.0:8853/api/port/192.168.2.2/8844
```
</details>
<details>
  <summary>Response example</summary>
  
```json
{
    "Name": "git-syr",
    "Port": 8844,
    "State": true,
    "Watch": true
}
```
</details><br>  

## Thanks
- All go packages listed in [dependencies](https://github.com/aceberg/watchyourports/network/dependencies)
- [Bootstrap](https://getbootstrap.com/)
- Themes: [Free themes for Bootstrap](https://bootswatch.com)
- Favicon and logo: [Flaticon](https://www.flaticon.com/icons/)