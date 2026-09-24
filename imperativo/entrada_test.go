package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCasosEntradaInvalidaDaEtapa2(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		codigo   string
		mensagem string
	}{
		{
			nome: "I01",
			entrada: `
p01 | 1500 | Top     | - | 10 | 20 | - | BR
p01 | 1510 | Jungle  | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 | 1500 | Mid     | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | ADC     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | Support | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
p10 | 1500 | Top     | - | 10 | 20 | - | BR
`,
			codigo:   "ID_DUPLICADO",
			mensagem: "identificador repetido na fila: p01",
		},
		{
			nome: "I02",
			entrada: `
p01 | 1500 | Top     | - | 10 | 20 | - | BR
p02 | 1500 | Top     | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 |   -1 | Jungle  | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | Mid     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | ADC     | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
p10 | 1500 | Support | - | 10 | 20 | - | BR
`,
			codigo:   "HABILIDADE_INVALIDA",
			mensagem: "habilidade deve ser maior ou igual a 0 (jogador p04)",
		},
		{
			nome: "I03",
			entrada: `
p01 | 1500 | Top      | - | 10 | 20 | - | BR
p02 | 1500 | Top      | - | 10 | 20 | - | BR
p03 | 1500 | Jungle   | - | 10 | 20 | - | BR
p04 | 1500 | Jungle   | - | 10 | 20 | - | BR
p05 | 1500 | Mid      | - | 10 | 20 | - | BR
p06 | 1500 | Mid      | - | 10 | 20 | - | BR
p07 | 1500 | Atirador | - | 10 | 20 | - | BR
p08 | 1500 | ADC      | - | 10 | 20 | - | BR
p09 | 1500 | Support  | - | 10 | 20 | - | BR
p10 | 1500 | Support  | - | 10 | 20 | - | BR
`,
			codigo:   "POSICAO_DESCONHECIDA",
			mensagem: "posicao nao reconhecida: Atirador (jogador p07)",
		},
		{
			nome: "I04",
			entrada: `
p01 | 1500 | Top     | - | 10 | 20 | G9 | BR
p02 | 1500 | Top     | - | 10 | 20 | G9 | BR
p03 | 1500 | Jungle  | - | 10 | 20 | G9 | BR
p04 | 1500 | Jungle  | - | 10 | 20 | G9 | BR
p05 | 1500 | Mid     | - | 10 | 20 | G9 | BR
p06 | 1500 | Mid     | - | 10 | 20 | G9 | BR
p07 | 1500 | ADC     | - | 10 | 20 | -  | BR
p08 | 1500 | ADC     | - | 10 | 20 | -  | BR
p09 | 1500 | Support | - | 10 | 20 | -  | BR
p10 | 1500 | Support | - | 10 | 20 | -  | BR
`,
			codigo:   "GRUPO_MAIOR_QUE_EQUIPE",
			mensagem: "grupo G9 possui 6 jogadores, acima do tamanho maximo de uma equipe (5)",
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			fila, erro, ok := lerFila(caso.entrada)
			if ok {
				t.Fatalf("a fila foi aceita com %d jogadores", len(fila))
			}
			if erro.Codigo != caso.codigo {
				t.Fatalf("codigo %q, esperava %q", erro.Codigo, caso.codigo)
			}
			if erro.Mensagem != caso.mensagem {
				t.Fatalf("mensagem %q, esperava %q", erro.Mensagem, caso.mensagem)
			}

			var buf bytes.Buffer
			escreverErroEntrada(&buf, erro)
			texto := buf.String()
			if !strings.Contains(texto, "status: ERRO_ENTRADA\n") {
				t.Fatalf("saida sem status: %q", texto)
			}
			if !strings.Contains(texto, "codigo: "+caso.codigo+"\n") {
				t.Fatalf("saida sem codigo: %q", texto)
			}
		})
	}
}

func TestFilaBemFormadaPreencheAceitas(t *testing.T) {
	entrada := `
p01 | 1500 | Mid | Mid, Top, Top | 60 | 20 | - | BR
p02 | 1510 | Jungle | - | 30 | 22 | G1 | BR
`
	fila, erro, ok := lerFila(entrada)
	if !ok {
		t.Fatalf("fila rejeitada: %s %s", erro.Codigo, erro.Mensagem)
	}
	if len(fila) != 2 {
		t.Fatalf("quantidade %d", len(fila))
	}
	if fila[0].Habilidade != 1500 || fila[0].EsperaS != 60 || fila[0].LatenciaMs != 20 {
		t.Fatalf("campos numericos de p01: %+v", fila[0])
	}
	if fila[0].Grupo != "" || fila[0].Regiao != "BR" || fila[0].Preferida != "Mid" {
		t.Fatalf("campos textuais de p01: %+v", fila[0])
	}
	if len(fila[0].Alternativas) != 2 || fila[0].Alternativas[0] != "Mid" || fila[0].Alternativas[1] != "Top" {
		t.Fatalf("alternativas %+v", fila[0].Alternativas)
	}
	if len(fila[0].Aceitas) != 2 || fila[0].Aceitas[0] != "Mid" || fila[0].Aceitas[1] != "Top" {
		t.Fatalf("aceitas %+v", fila[0].Aceitas)
	}
	if fila[1].Grupo != "G1" || len(fila[1].Aceitas) != 1 || fila[1].Aceitas[0] != "Jungle" {
		t.Fatalf("p02 %+v", fila[1])
	}
}

func TestGrupoDeCincoEFilaCurtaSaoEntradaValida(t *testing.T) {
	entrada := `
p01 | 1500 | Top | - | 10 | 20 | G1 | BR
p02 | 1500 | Jungle | - | 10 | 20 | G1 | BR
p03 | 1500 | Mid | - | 10 | 20 | G1 | BR
p04 | 1500 | ADC | - | 10 | 20 | G1 | BR
p05 | 1500 | Support | - | 10 | 20 | G1 | BR
`
	fila, erro, ok := lerFila(entrada)
	if !ok {
		t.Fatalf("grupo de 5 foi rejeitado: %s %s", erro.Codigo, erro.Mensagem)
	}
	if len(fila) != 5 {
		t.Fatalf("quantidade %d", len(fila))
	}
}

func TestOrdemDasViolacoes(t *testing.T) {
	entrada := `
p01 | -1 | Top | - | 10 | 20 | - | BR
p01 | 1500 | Jungle | - | 10 | 20 | - | BR
`
	_, erro, ok := lerFila(entrada)
	if ok {
		t.Fatal("fila aceita")
	}
	if erro.Codigo != "HABILIDADE_INVALIDA" {
		t.Fatalf("codigo %s", erro.Codigo)
	}
}

func TestEsperaLatenciaEAlternativaInvalidas(t *testing.T) {
	_, erro, ok := lerFila("p01 | 1500 | Top | - | -3 | 20 | - | BR\n")
	if ok || erro.Codigo != "ESPERA_INVALIDA" {
		t.Fatalf("espera: ok=%v erro=%+v", ok, erro)
	}

	_, erro, ok = lerFila("p01 | 1500 | Top | - | 10 | -5 | - | BR\n")
	if ok || erro.Codigo != "LATENCIA_INVALIDA" {
		t.Fatalf("latencia: ok=%v erro=%+v", ok, erro)
	}

	_, erro, ok = lerFila("p01 | 1500 | Mid | Atirador | 10 | 20 | - | BR\n")
	if ok || erro.Codigo != "POSICAO_DESCONHECIDA" || erro.Mensagem != "posicao nao reconhecida: Atirador (jogador p01)" {
		t.Fatalf("alternativa: ok=%v erro=%+v", ok, erro)
	}

	_, erro, ok = lerFila("p01 | 1500 |  | - | 10 | 20 | - | BR\n")
	if ok || erro.Codigo != "PREFERIDA_AUSENTE" {
		t.Fatalf("preferida: ok=%v erro=%+v", ok, erro)
	}
}
