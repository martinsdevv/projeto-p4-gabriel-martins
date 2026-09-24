package main

import "testing"

func idsEquipe(equipe [5]Atribuicao) []string {
	ids := make([]string, 5)
	for i := 0; i < 5; i++ {
		ids[i] = equipe[i].Jogador.ID
		if equipe[i].Posicao != posicoes[i] {
			ids[i] = ids[i] + "@" + equipe[i].Posicao
		}
	}
	return ids
}

func exigirPartida(t *testing.T, nome string, equipeA []string, equipeB []string, restante []string, qualidade int) {
	t.Helper()
	fila := filaArquivo(t, nome)
	resultado := buscarPartida(fila)
	if !resultado.Formou {
		t.Fatalf("%s motivo %s", nome, resultado.Motivo)
	}
	obtidaA := idsEquipe(resultado.Formacao.EquipeA)
	obtidaB := idsEquipe(resultado.Formacao.EquipeB)
	if !idsIguais(obtidaA, equipeA) || !idsIguais(obtidaB, equipeB) {
		t.Fatalf("%s\nA %v\nB %v", nome, obtidaA, obtidaB)
	}
	if !idsIguais(resultado.FilaRestante, restante) && !(len(resultado.FilaRestante) == 0 && len(restante) == 0) {
		t.Fatalf("%s restante %v", nome, resultado.FilaRestante)
	}
	if resultado.Metricas.Qualidade != qualidade {
		t.Fatalf("%s Q %d", nome, resultado.Metricas.Qualidade)
	}
}

func exigirNenhuma(t *testing.T, nome string, motivo string, restante []string) {
	t.Helper()
	fila := filaArquivo(t, nome)
	resultado := buscarPartida(fila)
	if resultado.Formou {
		t.Fatalf("%s formou partida", nome)
	}
	if resultado.Motivo != motivo {
		t.Fatalf("%s motivo %s", nome, resultado.Motivo)
	}
	if !idsIguais(resultado.FilaRestante, restante) {
		t.Fatalf("%s restante %v", nome, resultado.FilaRestante)
	}
}

func TestBuscaCasosNormais(t *testing.T) {
	exigirPartida(t, "n01.txt",
		[]string{"p03", "p02", "p01", "p09", "p05"},
		[]string{"p08", "p07", "p06", "p04", "p10"},
		nil, 9993)
	exigirPartida(t, "n02.txt",
		[]string{"p01", "p03", "p05", "p08", "p10"},
		[]string{"p02", "p04", "p06", "p07", "p09"},
		nil, 8000)
	exigirPartida(t, "n03.txt",
		[]string{"p02", "p03", "p01", "p04", "p05"},
		[]string{"p08", "p07", "p06", "p09", "p10"},
		nil, 9850)
	exigirPartida(t, "n04.txt",
		[]string{"p01", "p03", "p05", "p07", "p09"},
		[]string{"p02", "p04", "p06", "p08", "p10"},
		nil, 7400)
	exigirPartida(t, "n05.txt",
		[]string{"p01", "p02", "p05", "p07", "p09"},
		[]string{"p03", "p04", "p06", "p08", "p10"},
		nil, 8000)
	exigirPartida(t, "n06.txt",
		[]string{"p08", "p12", "p01", "p09", "p10"},
		[]string{"p11", "p07", "p06", "p04", "p05"},
		[]string{"p02", "p03"}, 9994)
	exigirPartida(t, "n07.txt",
		[]string{"p01", "p03", "p05", "p07", "p09"},
		[]string{"p02", "p04", "p06", "p08", "p10"},
		nil, 10000)
	exigirPartida(t, "n08.txt",
		[]string{"p08", "p03", "p05", "p04", "p01"},
		[]string{"p02", "p07", "p06", "p09", "p10"},
		nil, 9800)
	exigirPartida(t, "n09.txt",
		[]string{"p01", "p04", "p06", "p08", "p10"},
		[]string{"p02", "p05", "p07", "p09", "p11"},
		[]string{"p03"}, 10000)
	exigirPartida(t, "n10.txt",
		[]string{"p03", "p05", "p07", "p01", "p02"},
		[]string{"p04", "p06", "p08", "p09", "p10"},
		nil, 9996)
}

func TestBuscaCasosLimite(t *testing.T) {
	exigirNenhuma(t, "l01.txt", "JOGADORES_INSUFICIENTES",
		[]string{"p01", "p02", "p03", "p04", "p05", "p06", "p07", "p08", "p09"})
	exigirNenhuma(t, "l02.txt", "COMPOSICAO_IMPOSSIVEL",
		[]string{"p01", "p02", "p03", "p04", "p05", "p06", "p07", "p08", "p09", "p10"})
	exigirNenhuma(t, "l03.txt", "DIFERENCA_HABILIDADE",
		[]string{"p01", "p02", "p03", "p04", "p05", "p06", "p07", "p08", "p09", "p10"})
	exigirPartida(t, "l04.txt",
		[]string{"p01", "p03", "p05", "p08", "p10"},
		[]string{"p02", "p04", "p06", "p07", "p09"},
		nil, 8000)
}
