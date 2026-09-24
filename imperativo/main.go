package main

import (
	"fmt"
	"io"
	"os"
)

// Fluxo: ler a fila, validar e, se a entrada for válida, buscar a partida.
// A escrita na saída padrão é o efeito final dessa sequência.
func main() {
	args := os.Args[1:]
	if pasta, ok, err := pastaRelatorio(args); ok || err != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
		if err = escreverRelatorio(os.Stdout, pasta); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
		return
	}

	texto, err := lerTexto(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}
	os.Exit(executar(texto, os.Stdout))
}

// executar devolve 1 quando a entrada é inválida e 0 quando há um
// resultado do contrato, com ou sem partida.
func executar(texto string, saida io.Writer) int {
	fila, erro, ok := lerFila(texto)
	if !ok {
		escreverErroEntrada(saida, erro)
		return 1
	}
	resultado := buscarPartida(fila)
	if !resultado.Formou {
		escreverNenhumaPartida(saida, resultado)
		return 0
	}
	escreverPartida(saida, resultado)
	return 0
}

func lerTexto(args []string) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("uso: imperativo [arquivo | -relatorio pasta]")
	}
	if len(args) == 1 {
		bytes, err := os.ReadFile(args[0])
		if err != nil {
			return "", fmt.Errorf("falha ao ler arquivo: %s", err.Error())
		}
		return string(bytes), nil
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("falha ao ler entrada: %s", err.Error())
	}
	return string(bytes), nil
}
