package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRelatorioDosCasos(t *testing.T) {
	var buf bytes.Buffer
	err := escreverRelatorio(&buf, "../testes")
	if err != nil {
		t.Fatal(err)
	}
	texto := buf.String()
	if !strings.Contains(texto, "18 casos: 11 partidas, 3 sem partida, 4 erros de entrada") {
		t.Fatalf("resumo:\n%s", texto)
	}
	if !strings.Contains(texto, "n06") || !strings.Contains(texto, "qualidade 9994") {
		t.Fatalf("n06:\n%s", texto)
	}
	if !strings.Contains(texto, "l03") || !strings.Contains(texto, "DIFERENCA_HABILIDADE") {
		t.Fatalf("l03:\n%s", texto)
	}
	if !strings.Contains(texto, "i01") || !strings.Contains(texto, "ID_DUPLICADO") {
		t.Fatalf("i01:\n%s", texto)
	}
	if !strings.Contains(texto, "media: 45.0") || !strings.Contains(texto, "media: 34.0") {
		t.Fatalf("espera:\n%s", texto)
	}
	if strings.Contains(texto, "tempo total") {
		t.Fatal("o relatorio nao deve medir o tempo de execucao")
	}
}
