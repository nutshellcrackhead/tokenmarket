#!/bin/sh

rm -rf build/*

go build -o build/api api
go build -o build/auth auth
go build -o build/profile profile
go build -o build/referral referral
go build -o build/wallet wallet
go build -o build/operations operations
go build -o build/deposits deposits
go build -o build/bonus bonus

make migrate