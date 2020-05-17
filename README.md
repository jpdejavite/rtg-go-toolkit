# Rtg go tollkit

A golang toolkit for rtg apps

# unit test with coverage

```
go test -coverprofile coverage.txt ./...
```

# unit test mockgenerator

```
mockgen -source=pkg/config/globalconfig.go -destination=mock/config/globalconfig.go
```