module main
go 1.19

require "accounts" v1.0.0
replace "accounts" => "./accounts"

require "dictionary" v1.0.0
replace "dictionary" => "./dictionary"

require "errorMsgs" v1.0.0
replace "errorMsgs" => "./errorMsgs"