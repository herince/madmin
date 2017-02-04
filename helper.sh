#!/bin/bash

go fmt ./

for i in `seq 1 10`
do
	echo
done

# compile scss
sass --update static/css/scss/main.scss:static/css/main.css

# build server
go build github.com/herince/madmin/app

# run server
go run madmin.go
