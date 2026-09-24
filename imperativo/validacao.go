package main

import "fmt"

// validarFila aplica as regras de entrada do contrato.
// A primeira violação encerra a função. A ordem é fixa para que
// duas implementações, e os testes, vejam o mesmo código quando
// houver mais de um problema:
//  1. cada jogador, na ordem da fila (habilidade, espera, latência, posições);
//  2. identificadores duplicados;
//  3. grupo com mais jogadores do que cabem em uma equipe.
func validarFila(fila []Jogador) (ErroEntrada, bool) {
	for i := 0; i < len(fila); i++ {
		erro, ok := validarJogador(fila[i])
		if !ok {
			return erro, false
		}
	}

	erro, ok := validarIDs(fila)
	if !ok {
		return erro, false
	}
	return validarGrupos(fila)
}

func validarJogador(jogador Jogador) (ErroEntrada, bool) {
	if jogador.Habilidade < 0 {
		return ErroEntrada{
			Codigo:   "HABILIDADE_INVALIDA",
			Mensagem: fmt.Sprintf("habilidade deve ser maior ou igual a 0 (jogador %s)", jogador.ID),
		}, false
	}
	if jogador.EsperaS < 0 {
		return ErroEntrada{
			Codigo:   "ESPERA_INVALIDA",
			Mensagem: fmt.Sprintf("espera_s deve ser maior ou igual a 0 (jogador %s)", jogador.ID),
		}, false
	}
	if jogador.LatenciaMs < 0 {
		return ErroEntrada{
			Codigo:   "LATENCIA_INVALIDA",
			Mensagem: fmt.Sprintf("latencia_ms deve ser maior ou igual a 0 (jogador %s)", jogador.ID),
		}, false
	}
	if jogador.Preferida == "" {
		return ErroEntrada{
			Codigo:   "PREFERIDA_AUSENTE",
			Mensagem: fmt.Sprintf("posicao preferida ausente (jogador %s)", jogador.ID),
		}, false
	}
	if !posicaoValida(jogador.Preferida) {
		return ErroEntrada{
			Codigo:   "POSICAO_DESCONHECIDA",
			Mensagem: fmt.Sprintf("posicao nao reconhecida: %s (jogador %s)", jogador.Preferida, jogador.ID),
		}, false
	}
	for i := 0; i < len(jogador.Alternativas); i++ {
		if !posicaoValida(jogador.Alternativas[i]) {
			return ErroEntrada{
				Codigo:   "POSICAO_DESCONHECIDA",
				Mensagem: fmt.Sprintf("posicao nao reconhecida: %s (jogador %s)", jogador.Alternativas[i], jogador.ID),
			}, false
		}
	}
	return ErroEntrada{}, true
}

// validarIDs usa um mapa como estado mutável: a cada jogador, consulta
// e depois grava o id já visto. O primeiro id repetido é o reportado.
func validarIDs(fila []Jogador) (ErroEntrada, bool) {
	vistos := map[string]bool{}
	for i := 0; i < len(fila); i++ {
		id := fila[i].ID
		if vistos[id] {
			return ErroEntrada{
				Codigo:   "ID_DUPLICADO",
				Mensagem: "identificador repetido na fila: " + id,
			}, false
		}
		vistos[id] = true
	}
	return ErroEntrada{}, true
}

// validarGrupos conta ocorrências de cada grupo não vazio.
// Grupo vazio (solo) não entra no mapa. Tamanho 1 também é válido:
// a restrição de "permanecer junto" só importa na busca, mais adiante.
func validarGrupos(fila []Jogador) (ErroEntrada, bool) {
	contagem := map[string]int{}
	var ordem []string

	for i := 0; i < len(fila); i++ {
		grupo := fila[i].Grupo
		if grupo == "" {
			continue
		}
		if contagem[grupo] == 0 {
			ordem = append(ordem, grupo)
		}
		contagem[grupo] = contagem[grupo] + 1
	}

	for i := 0; i < len(ordem); i++ {
		grupo := ordem[i]
		if contagem[grupo] > JogadoresPorEquipe {
			return ErroEntrada{
				Codigo: "GRUPO_MAIOR_QUE_EQUIPE",
				Mensagem: fmt.Sprintf(
					"grupo %s possui %d jogadores, acima do tamanho maximo de uma equipe (%d)",
					grupo, contagem[grupo], JogadoresPorEquipe,
				),
			}, false
		}
	}
	return ErroEntrada{}, true
}
