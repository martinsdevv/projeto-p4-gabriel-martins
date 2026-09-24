package main

import (
	"os"
	"path/filepath"
	"testing"
)

func filaArquivo(t *testing.T, nome string) []Jogador {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Join("..", "testes", nome))
	if err != nil {
		t.Fatal(err)
	}
	fila, erro, ok := lerFila(string(bytes))
	if !ok {
		t.Fatalf("fila %s rejeitada: %s %s", nome, erro.Codigo, erro.Mensagem)
	}
	return fila
}

func vaga(fila []Jogador, id string, posicao string) Atribuicao {
	for i := 0; i < len(fila); i++ {
		if fila[i].ID == id {
			return Atribuicao{Jogador: &fila[i], Posicao: posicao}
		}
	}
	return Atribuicao{Posicao: posicao}
}

func TestN01FormacaoEquilibrada(t *testing.T) {
	fila := filaArquivo(t, "n01.txt")
	formacao := Formacao{
		EquipeA: [5]Atribuicao{
			vaga(fila, "p03", "Top"),
			vaga(fila, "p02", "Jungle"),
			vaga(fila, "p01", "Mid"),
			vaga(fila, "p09", "ADC"),
			vaga(fila, "p05", "Support"),
		},
		EquipeB: [5]Atribuicao{
			vaga(fila, "p08", "Top"),
			vaga(fila, "p07", "Jungle"),
			vaga(fila, "p06", "Mid"),
			vaga(fila, "p04", "ADC"),
			vaga(fila, "p10", "Support"),
		},
	}
	m, regra, ok := avaliarFormacao(fila, formacao)
	if !ok || regra != "" {
		t.Fatalf("regra %s", regra)
	}
	if m.DiferencaHabilidade != 0 || m.Qualidade != 9993 || m.AmplitudeLatencia != 7 {
		t.Fatalf("metricas %+v", m)
	}
	if m.EsperaMin != 30 || m.EsperaMax != 60 || m.EsperaSoma != 450 {
		t.Fatalf("espera %+v", m)
	}
	if m.SomaA != 7520 || m.SomaB != 7520 {
		t.Fatalf("somas %d %d", m.SomaA, m.SomaB)
	}
}

func TestN04EsperaAmpliaLimiteEL03Recusa(t *testing.T) {
	formacaoDe := func(fila []Jogador) Formacao {
		return Formacao{
			EquipeA: [5]Atribuicao{
				vaga(fila, "p01", "Top"),
				vaga(fila, "p03", "Jungle"),
				vaga(fila, "p05", "Mid"),
				vaga(fila, "p07", "ADC"),
				vaga(fila, "p09", "Support"),
			},
			EquipeB: [5]Atribuicao{
				vaga(fila, "p02", "Top"),
				vaga(fila, "p04", "Jungle"),
				vaga(fila, "p06", "Mid"),
				vaga(fila, "p08", "ADC"),
				vaga(fila, "p10", "Support"),
			},
		}
	}

	filaN04 := filaArquivo(t, "n04.txt")
	m, _, ok := avaliarFormacao(filaN04, formacaoDe(filaN04))
	if !ok || m.DiferencaHabilidade != 260 || m.LimiteDiferenca != 300 || m.Qualidade != 7400 {
		t.Fatalf("n04 ok=%v metricas %+v", ok, m)
	}

	filaL03 := filaArquivo(t, "l03.txt")
	m, regra, ok := avaliarFormacao(filaL03, formacaoDe(filaL03))
	if ok || regra != "R6" || m.LimiteDiferenca != 202 || m.DiferencaHabilidade != 260 {
		t.Fatalf("l03 ok=%v regra=%s metricas %+v", ok, regra, m)
	}
}

func TestN05GrupoSeparadoInvalido(t *testing.T) {
	fila := filaArquivo(t, "n05.txt")
	junta := Formacao{
		EquipeA: [5]Atribuicao{
			vaga(fila, "p01", "Top"),
			vaga(fila, "p02", "Jungle"),
			vaga(fila, "p05", "Mid"),
			vaga(fila, "p07", "ADC"),
			vaga(fila, "p09", "Support"),
		},
		EquipeB: [5]Atribuicao{
			vaga(fila, "p03", "Top"),
			vaga(fila, "p04", "Jungle"),
			vaga(fila, "p06", "Mid"),
			vaga(fila, "p08", "ADC"),
			vaga(fila, "p10", "Support"),
		},
	}
	m, _, ok := avaliarFormacao(fila, junta)
	if !ok || m.DiferencaHabilidade != 200 || m.Qualidade != 8000 || m.LimiteDiferenca != 207 {
		t.Fatalf("junta ok=%v %+v", ok, m)
	}

	separada := junta
	separada.EquipeA[1] = vaga(fila, "p04", "Jungle")
	separada.EquipeB[1] = vaga(fila, "p02", "Jungle")
	_, regra, ok := avaliarFormacao(fila, separada)
	if ok || regra != "R5" {
		t.Fatalf("separada ok=%v regra=%s", ok, regra)
	}
}

func TestN03AlternativaContaNaQualidade(t *testing.T) {
	fila := filaArquivo(t, "n03.txt")
	formacao := Formacao{
		EquipeA: [5]Atribuicao{
			vaga(fila, "p02", "Top"),
			vaga(fila, "p03", "Jungle"),
			vaga(fila, "p01", "Mid"),
			vaga(fila, "p04", "ADC"),
			vaga(fila, "p05", "Support"),
		},
		EquipeB: [5]Atribuicao{
			vaga(fila, "p08", "Top"),
			vaga(fila, "p07", "Jungle"),
			vaga(fila, "p06", "Mid"),
			vaga(fila, "p09", "ADC"),
			vaga(fila, "p10", "Support"),
		},
	}
	m, _, ok := avaliarFormacao(fila, formacao)
	if !ok || m.ForaPreferencia != 1 || m.DesequilibrioAlt != 1 || m.Qualidade != 9850 {
		t.Fatalf("ok=%v %+v", ok, m)
	}

	formacao.EquipeA[0].Posicao = "Mid"
	_, regra, ok := avaliarFormacao(fila, formacao)
	if ok || regra != "R4" {
		t.Fatalf("composicao ok=%v regra=%s", ok, regra)
	}
}

func TestPosicaoRecusadaEJogadorRepetido(t *testing.T) {
	fila := filaArquivo(t, "n01.txt")
	formacao := Formacao{
		EquipeA: [5]Atribuicao{
			vaga(fila, "p03", "Top"),
			vaga(fila, "p02", "Jungle"),
			vaga(fila, "p01", "Mid"),
			vaga(fila, "p09", "ADC"),
			vaga(fila, "p05", "Support"),
		},
		EquipeB: [5]Atribuicao{
			vaga(fila, "p08", "Top"),
			vaga(fila, "p07", "Jungle"),
			vaga(fila, "p06", "Mid"),
			vaga(fila, "p04", "ADC"),
			vaga(fila, "p10", "Support"),
		},
	}
	formacao.EquipeA[2].Posicao = "Top"
	_, regra, ok := avaliarFormacao(fila, formacao)
	if ok || regra != "R3" {
		t.Fatalf("posicao ok=%v regra=%s", ok, regra)
	}

	formacao.EquipeA[2].Posicao = "Mid"
	formacao.EquipeB[2] = formacao.EquipeA[2]
	_, regra, ok = avaliarFormacao(fila, formacao)
	if ok || regra != "R1" {
		t.Fatalf("repetido ok=%v regra=%s", ok, regra)
	}
}
