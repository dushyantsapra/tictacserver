### Make docker image

```bash
docker build -t tictacserver .
```

### Run 
```bash
docker run -it -p 8080:8080 tictacserver
```

### Start new game
```bash
curl -X GET localhost:8080/game
```

### Make turn
```bash
curl -d '{"player":"X","row":1,"column":2}' localhost:8080/game/move
```

### delete game
```bash
curl -X DELETE localhost:8080/game
```
