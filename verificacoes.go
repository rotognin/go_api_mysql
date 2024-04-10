package main

import (
	"strconv"
)

func checarID(id string) string {
	var retorno string
	retorno = ""

	numero, err := strconv.Atoi(id)
	if err != nil {
		retorno += " ID informado não é numérico."
	}

	if numero == 0 {
		retorno += " ID deve ser maior que zero."
	}

	return retorno
}

func checarTitle(title string) string {
	var retorno string
	retorno = ""

	if title == "" {
		retorno += " O título não pode estar em branco"
	}

	return retorno
}

func checarArtist(artist string) string {
	var retorno string
	retorno = ""

	if artist == "" {
		retorno += " O artista deve ser informado"
	}

	return retorno
}

func checarPrice(price float64) string {
	var retorno string
	retorno = ""

	if price == 0 {
		retorno += " O preço deve ser informado"
	}

	return retorno
}