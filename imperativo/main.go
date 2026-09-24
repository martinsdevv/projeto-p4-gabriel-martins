package main

import (
	"fmt"
	"io"
	"os"
)

// Fluxo desta parte: ler a fila, validar, e ou rejeitar a entrada
// ou confirmar que ela está bem formada.
// A busca da partida (PARTIDA_FORMADA / NENHUMA_PARTIDA) entra na
// parte seguinte; até lá, fila válida produz FILA_VALIDA.
func main() {
	texto, err := lerTexto(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	fila, erro, ok := lerFila(texto)
	if !ok {
		escreverErroEntrada(os.Stdout, erro)
		os.Exit(1)
	}

	fmt.Printf("status: FILA_VALIDA\n")
	fmt.Printf("jogadores: %d\n", len(fila))
}

func lerTexto(args []string) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("uso: imperativo [arquivo]")
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
