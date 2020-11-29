.PHONY: nchat
nchat:
	GOOS=linux CGO_ENABLED=0 go build -o nchat ./server/*.go

.PHONY: nclient
nclient:
	GOOS=linux CGO_ENABLED=0 go build -o nclient ./client/*.go
