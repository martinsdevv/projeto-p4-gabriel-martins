package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSaidaDosDezoitoCasos(t *testing.T) {
	casos := []struct {
		arquivo string
		codigo  int
		saida   string
	}{
		{"n01.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p03, Top), (p02, Jungle), (p01, Mid), (p09, ADC), (p05, Support)
equipe_b: (p08, Top), (p07, Jungle), (p06, Mid), (p04, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9993
espera: {min: 30, max: 60, media: 45.0}
amplitude_latencia: 7
fila_restante: []
`},
		{"n02.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p07, ADC), (p09, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 10, max: 10, media: 10.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n03.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p02, Top), (p03, Jungle), (p01, Mid), (p04, ADC), (p05, Support)
equipe_b: (p08, Top), (p07, Jungle), (p06, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 1
desequilibrio_alternativas: 1
qualidade: 9850
espera: {min: 20, max: 20, media: 20.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n04.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 260
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 7400
espera: {min: 10, max: 400, media: 49.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n05.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p02, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p03, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 30, max: 30, media: 30.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n06.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p08, Top), (p12, Jungle), (p01, Mid), (p09, ADC), (p10, Support)
equipe_b: (p11, Top), (p07, Jungle), (p06, Mid), (p04, ADC), (p05, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9994
espera: {min: 5, max: 60, media: 34.0}
amplitude_latencia: 6
fila_restante: [p02, p03]
`},
		{"n07.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 10000
espera: {min: 15, max: 15, media: 15.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n08.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p08, Top), (p03, Jungle), (p05, Mid), (p04, ADC), (p01, Support)
equipe_b: (p02, Top), (p07, Jungle), (p06, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 2
desequilibrio_alternativas: 0
qualidade: 9800
espera: {min: 25, max: 25, media: 25.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"n09.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p05, Jungle), (p07, Mid), (p09, ADC), (p11, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 10000
espera: {min: 20, max: 20, media: 20.0}
amplitude_latencia: 0
fila_restante: [p03]
`},
		{"n10.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p03, Top), (p05, Jungle), (p07, Mid), (p01, ADC), (p02, Support)
equipe_b: (p04, Top), (p06, Jungle), (p08, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9996
espera: {min: 20, max: 40, media: 24.0}
amplitude_latencia: 4
fila_restante: []
`},
		{"l01.txt", 0, `status: NENHUMA_PARTIDA
motivo: JOGADORES_INSUFICIENTES
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09]
`},
		{"l02.txt", 0, `status: NENHUMA_PARTIDA
motivo: COMPOSICAO_IMPOSSIVEL
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09, p10]
`},
		{"l03.txt", 0, `status: NENHUMA_PARTIDA
motivo: DIFERENCA_HABILIDADE
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09, p10]
`},
		{"l04.txt", 0, `status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p07, ADC), (p09, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 0, max: 0, media: 0.0}
amplitude_latencia: 0
fila_restante: []
`},
		{"i01.txt", 1, `status: ERRO_ENTRADA
codigo: ID_DUPLICADO
mensagem: identificador repetido na fila: p01
`},
		{"i02.txt", 1, `status: ERRO_ENTRADA
codigo: HABILIDADE_INVALIDA
mensagem: habilidade deve ser maior ou igual a 0 (jogador p04)
`},
		{"i03.txt", 1, `status: ERRO_ENTRADA
codigo: POSICAO_DESCONHECIDA
mensagem: posicao nao reconhecida: Atirador (jogador p07)
`},
		{"i04.txt", 1, `status: ERRO_ENTRADA
codigo: GRUPO_MAIOR_QUE_EQUIPE
mensagem: grupo G9 possui 6 jogadores, acima do tamanho maximo de uma equipe (5)
`},
	}

	for _, caso := range casos {
		t.Run(caso.arquivo, func(t *testing.T) {
			bytesArquivo, err := os.ReadFile(filepath.Join("..", "testes", caso.arquivo))
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			codigo := executar(string(bytesArquivo), &buf)
			if codigo != caso.codigo {
				t.Fatalf("codigo %d, esperava %d\n%s", codigo, caso.codigo, buf.String())
			}
			if buf.String() != caso.saida {
				t.Fatalf("saida:\n%s\nesperada:\n%s", buf.String(), caso.saida)
			}
		})
	}
}
