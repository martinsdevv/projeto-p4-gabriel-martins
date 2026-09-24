package main

// avaliarFormacao percorre a formação e aplica R1, R3, R4, R5 e R6.
// R2 está na forma do registro: cada equipe é um arranjo de cinco vagas,
// e o laço de jogadores distintos recusa quem aparecer nas duas.
// Se a formação é válida, as métricas incluem Q. Se cai em R6, as
// métricas vêm preenchidas mesmo assim, para dar para ver o limite que estourou.
func avaliarFormacao(fila []Jogador, formacao Formacao) (Metricas, string, bool) {
	if !jogadoresDistintos(fila, formacao) {
		return Metricas{}, "R1", false
	}
	if !posicoesAceitas(formacao) {
		return Metricas{}, "R3", false
	}
	if !composicaoValida(formacao.EquipeA) || !composicaoValida(formacao.EquipeB) {
		return Metricas{}, "R4", false
	}
	if !gruposUnidos(formacao) {
		return Metricas{}, "R5", false
	}

	metricas := calcularMetricas(formacao)
	if metricas.DiferencaHabilidade > metricas.LimiteDiferenca {
		return metricas, "R6", false
	}
	return metricas, "", true
}

func jogadoresDistintos(fila []Jogador, formacao Formacao) bool {
	vistos := map[string]bool{}
	equipes := [2][5]Atribuicao{formacao.EquipeA, formacao.EquipeB}
	for e := 0; e < len(equipes); e++ {
		for i := 0; i < len(equipes[e]); i++ {
			jogador := equipes[e][i].Jogador
			if jogador == nil || !idNaFila(fila, jogador.ID) {
				return false
			}
			if vistos[jogador.ID] {
				return false
			}
			vistos[jogador.ID] = true
		}
	}
	return len(vistos) == JogadoresPorPartida
}

func idNaFila(fila []Jogador, id string) bool {
	for i := 0; i < len(fila); i++ {
		if fila[i].ID == id {
			return true
		}
	}
	return false
}

func posicoesAceitas(formacao Formacao) bool {
	equipes := [2][5]Atribuicao{formacao.EquipeA, formacao.EquipeB}
	for e := 0; e < len(equipes); e++ {
		for i := 0; i < len(equipes[e]); i++ {
			vaga := equipes[e][i]
			if !contem(vaga.Jogador.Aceitas, vaga.Posicao) {
				return false
			}
		}
	}
	return true
}

func composicaoValida(equipe [5]Atribuicao) bool {
	var contagem [5]int
	for i := 0; i < len(equipe); i++ {
		indice := indicePosicao(equipe[i].Posicao)
		if indice < 0 {
			return false
		}
		contagem[indice] = contagem[indice] + 1
	}
	for i := 0; i < len(contagem); i++ {
		if contagem[i] != 1 {
			return false
		}
	}
	return true
}

func indicePosicao(nome string) int {
	for i := 0; i < len(posicoes); i++ {
		if posicoes[i] == nome {
			return i
		}
	}
	return -1
}

// gruposUnidos grava, para cada grupo não vazio, a equipe em que ele
// apareceu. A segunda ocorrência em outra equipe torna a formação inválida.
func gruposUnidos(formacao Formacao) bool {
	equipeDoGrupo := map[string]int{}
	equipes := [2][5]Atribuicao{formacao.EquipeA, formacao.EquipeB}
	for e := 0; e < len(equipes); e++ {
		for i := 0; i < len(equipes[e]); i++ {
			grupo := equipes[e][i].Jogador.Grupo
			if grupo == "" {
				continue
			}
			equipe, jaVisto := equipeDoGrupo[grupo]
			if jaVisto && equipe != e {
				return false
			}
			equipeDoGrupo[grupo] = e
		}
	}
	return true
}

// calcularMetricas zera um registro e o preenche em uma passagem por equipe.
// O tempo de espera entra só no limite de R6, não em Q.
func calcularMetricas(formacao Formacao) Metricas {
	var m Metricas
	primeiro := formacao.EquipeA[0].Jogador
	m.EsperaMin = primeiro.EsperaS
	m.EsperaMax = primeiro.EsperaS
	latMin := primeiro.LatenciaMs
	latMax := primeiro.LatenciaMs

	acumularEquipe(&m, formacao.EquipeA, true, &latMin, &latMax)
	acumularEquipe(&m, formacao.EquipeB, false, &latMin, &latMax)

	diferencaSoma := m.SomaA - m.SomaB
	if diferencaSoma < 0 {
		diferencaSoma = -diferencaSoma
	}
	m.DiferencaHabilidade = diferencaSoma / JogadoresPorEquipe
	m.ForaPreferencia = m.ForaPrefA + m.ForaPrefB
	m.DesequilibrioAlt = m.ForaPrefA - m.ForaPrefB
	if m.DesequilibrioAlt < 0 {
		m.DesequilibrioAlt = -m.DesequilibrioAlt
	}
	m.AmplitudeLatencia = latMax - latMin
	m.LimiteDiferenca = limiteDiferenca(m.EsperaMax)
	m.Qualidade = 10000 - 2*diferencaSoma - 100*m.ForaPreferencia - 50*m.DesequilibrioAlt - m.AmplitudeLatencia
	return m
}

func acumularEquipe(m *Metricas, equipe [5]Atribuicao, equipeA bool, latMin *int, latMax *int) {
	for i := 0; i < len(equipe); i++ {
		jogador := equipe[i].Jogador
		if equipeA {
			m.SomaA = m.SomaA + jogador.Habilidade
		} else {
			m.SomaB = m.SomaB + jogador.Habilidade
		}
		if equipe[i].Posicao != jogador.Preferida {
			if equipeA {
				m.ForaPrefA = m.ForaPrefA + 1
			} else {
				m.ForaPrefB = m.ForaPrefB + 1
			}
		}
		m.EsperaSoma = m.EsperaSoma + jogador.EsperaS
		if jogador.EsperaS < m.EsperaMin {
			m.EsperaMin = jogador.EsperaS
		}
		if jogador.EsperaS > m.EsperaMax {
			m.EsperaMax = jogador.EsperaS
		}
		if jogador.LatenciaMs < *latMin {
			*latMin = jogador.LatenciaMs
		}
		if jogador.LatenciaMs > *latMax {
			*latMax = jogador.LatenciaMs
		}
	}
}

func limiteDiferenca(maxEspera int) int {
	limite := LimiteBaseDiferenca + maxEspera/BonusEsperaDivisor
	if limite > LimiteAbsolutoDiferenca {
		limite = LimiteAbsolutoDiferenca
	}
	return limite
}
