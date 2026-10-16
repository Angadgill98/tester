package main

type Response struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Args   []any  `json:"args"`
}