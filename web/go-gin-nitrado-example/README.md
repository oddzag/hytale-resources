# Example Go app to interface with Nitrado's WebServer plugin

This example depends on already having [Nitrado's WebServer and Query plugin configured](https://github.com/oddzag/hytale-resources/blob/main/guides/how-to-nitrado-webserver-query-plugin.md), and a running [Hytale server](https://github.com/oddzag/hytale-resources/blob/main/guides/how-to-dedicated-server-debian-13.md).

### To set-up:

- Clone the repo
- Initialize the go project - and install dependencies
- Create the `.env` file
- Start the app. Rename it from `go-gin-nitrado-example` if you'd like

```
git clone https://github.com/oddzag/go-gin-nitrado-example.git
cd go-gin-nitrado-example
go mod go-gin-nitrado-example
go mod tidy
nano .env
```

`go-gin-nitrado-example/.env`
```env
HYTALE_QUERY_URL=http://192.168.1.10:5523/Nitrado/Query
HYTALE_QUERY_USERNAME=serviceAccount.oddzag
HYTALE_QUERY_PASSWORD=ThisIsATestPassword
HYTALE_QUERY_PORT=8080
```
*these credentials are from [this guide](https://github.com/oddzag/hytale-resources/blob/main/guides/how-to-nitrado-webserver-query-plugin.md), change them to match <u>your</u> setup*


Now start the server and in your browser, navigate to http://192.168.1.10:8080/query and you should see the server's `json` payload from Nitrado's Query endpoint 
```
go run main.go
```
```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /query                    --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2026/01/26 - 15:09:51 | 200 |  167.473658ms |      192.168.1.10 | GET      "/query"
```

If you have `ufw` set up, you may need to open port `8080` to access your `/query` endpoint