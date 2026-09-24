package main

import (
	"fmt"
	"io"
)

// escreverErroEntrada grava o formato da Seção 7.1 do contrato.
// Efeito colateral: escreve no io.Writer recebido. Em main, esse
// escritor é a saída padrão. Não há fila_restante neste resultado.
func escreverErroEntrada(saida io.Writer, erro ErroEntrada) {
	fmt.Fprintf(saida, "status: ERRO_ENTRADA\n")
	fmt.Fprintf(saida, "codigo: %s\n", erro.Codigo)
	fmt.Fprintf(saida, "mensagem: %s\n", erro.Mensagem)
}
