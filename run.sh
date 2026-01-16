#!/bin/zsh

go build -o bookings cmd/web/*.go
./bookings -dbname=booking -dbuser=tjasaspes -cache=false -production=false