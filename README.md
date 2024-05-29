# tradeEngine

## Done

- gin restful and jwt
- mysql implemnt
- logger implemt
- services handle table Member, Wallet, StockInfo, Stock, SellOrder, BuyOrder

## Todo

- update flow after add/edit/delete sell order
- update flow after add/edit/delete buy order
- thread handle order match
- socket.io or websocket
- documant like swagger

## How to run

1. install docker
2. bash exec `docker-compose.exe -f cfg\mysql\mysql.yml up -d`
3. bash exec `go build main.go` or run `vscode` to launch project
