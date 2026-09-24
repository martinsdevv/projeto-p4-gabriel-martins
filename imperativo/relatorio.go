package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// pastaRelatorio reconhece `imperativo -relatorio [pasta]` e também
// um único argumento que já é um diretório.
func pastaRelatorio(args []string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, nil
	}
	if args[0] == "-relatorio" {
		if len(args) > 2 {
			return "", false, fmt.Errorf("uso: imperativo -relatorio [pasta]")
		}
		if len(args) == 2 {
			return args[1], true, nil
		}
		return filepath.Join("..", "testes"), true, nil
	}
	if len(args) == 1 {
		info, err := os.Stat(args[0])
		if err == nil && info.IsDir() {
			return args[0], true, nil
		}
	}
	return "", false, nil
}

// escreverRelatorio percorre os .txt da pasta, executa cada fila e
// imprime uma linha por caso. O escritor recebe o relatório pronto.
func escreverRelatorio(saida io.Writer, pasta string) error {
	entradas, err := filepath.Glob(filepath.Join(pasta, "*.txt"))
	if err != nil {
		return err
	}
	if len(entradas) == 0 {
		return fmt.Errorf("nenhum caso em %s", pasta)
	}
	sort.Slice(entradas, func(i, j int) bool {
		return nomeAntes(entradas[i], entradas[j])
	})

	var linhas []string
	partidas := 0
	nenhuma := 0
	erros := 0
	for i := 0; i < len(entradas); i++ {
		texto, err := os.ReadFile(entradas[i])
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		executar(string(texto), &buf)
		saidaCaso := buf.String()
		status, detalhe := resumir(saidaCaso)
		switch status {
		case "PARTIDA_FORMADA":
			partidas++
		case "NENHUMA_PARTIDA":
			nenhuma++
		case "ERRO_ENTRADA":
			erros++
		}
		nome := strings.TrimSuffix(filepath.Base(entradas[i]), ".txt")
		linhas = append(linhas, fmt.Sprintf("%-6s %-18s %-24s %s", nome, status, detalhe, textoEspera(saidaCaso)))
	}

	fmt.Fprintf(saida, "%-6s %-18s %-24s %s\n", "caso", "status", "detalhe", "espera")
	for i := 0; i < len(linhas); i++ {
		fmt.Fprintln(saida, linhas[i])
	}
	fmt.Fprintf(saida, "\n%d casos: %d partidas, %d sem partida, %d erros de entrada\n",
		len(linhas), partidas, nenhuma, erros)
	return nil
}

func nomeAntes(caminhoA, caminhoB string) bool {
	a := filepath.Base(caminhoA)
	b := filepath.Base(caminhoB)
	ordem := func(nome string) int {
		if strings.HasPrefix(nome, "n") {
			return 0
		}
		if strings.HasPrefix(nome, "l") {
			return 1
		}
		if strings.HasPrefix(nome, "i") {
			return 2
		}
		return 3
	}
	if ordem(a) != ordem(b) {
		return ordem(a) < ordem(b)
	}
	return a < b
}

func resumir(saida string) (string, string) {
	status := valorCampo(saida, "status")
	switch status {
	case "PARTIDA_FORMADA":
		return status, "qualidade " + valorCampo(saida, "qualidade")
	case "NENHUMA_PARTIDA":
		return status, valorCampo(saida, "motivo")
	case "ERRO_ENTRADA":
		return status, valorCampo(saida, "codigo")
	default:
		return status, ""
	}
}

// textoEspera é o tempo de fila dos dez selecionados. Sem partida, não há formação.
func textoEspera(saida string) string {
	espera := valorCampo(saida, "espera")
	if espera == "" {
		return "-"
	}
	return espera
}

func valorCampo(texto, campo string) string {
	prefixo := campo + ": "
	linhas := strings.Split(texto, "\n")
	for i := 0; i < len(linhas); i++ {
		if strings.HasPrefix(linhas[i], prefixo) {
			return strings.TrimPrefix(linhas[i], prefixo)
		}
	}
	return ""
}
