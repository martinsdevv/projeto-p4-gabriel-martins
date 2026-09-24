package main

import (
	"fmt"
	"strconv"
	"strings"
)

// lerFila percorre o texto linha a linha, acumula jogadores em um slice
// e, se a validação passar, preenche Aceitas em cada registro.
// O slice fila é o estado desta etapa: cresce por append a cada linha útil.
func lerFila(texto string) ([]Jogador, ErroEntrada, bool) {
	var fila []Jogador
	linhas := strings.Split(texto, "\n")
	numero := 0

	for i := 0; i < len(linhas); i++ {
		linha := strings.TrimSpace(linhas[i])
		if linha == "" {
			continue
		}
		numero++
		jogador, erro, ok := lerLinha(linha, numero)
		if !ok {
			return nil, erro, false
		}
		fila = append(fila, jogador)
	}

	erro, ok := validarFila(fila)
	if !ok {
		return nil, erro, false
	}
	completarFila(fila)
	return fila, ErroEntrada{}, true
}

// lerLinha separa os oito campos do contrato.
// Ainda não aplica as regras de domínio: isso fica em validarFila,
// para o fluxo ficar em duas passagens visíveis (ler, depois validar).
func lerLinha(linha string, numero int) (Jogador, ErroEntrada, bool) {
	partes := strings.Split(linha, "|")
	if len(partes) != 8 {
		return Jogador{}, ErroEntrada{
			Codigo:   "CAMPO_OBRIGATORIO",
			Mensagem: fmt.Sprintf("linha %d: campo obrigatorio ausente", numero),
		}, false
	}

	for i := 0; i < len(partes); i++ {
		partes[i] = strings.TrimSpace(partes[i])
	}

	id := partes[0]
	if id == "" {
		return Jogador{}, ErroEntrada{
			Codigo:   "CAMPO_OBRIGATORIO",
			Mensagem: fmt.Sprintf("linha %d: identificador ausente", numero),
		}, false
	}

	habilidade, okHab := inteiroCampo(partes[1])
	espera, okEspera := inteiroCampo(partes[4])
	latencia, okLatencia := inteiroCampo(partes[5])
	if !okHab || !okEspera || !okLatencia {
		return Jogador{}, ErroEntrada{
			Codigo:   "CAMPO_OBRIGATORIO",
			Mensagem: fmt.Sprintf("linha %d: campo numerico ilegivel (jogador %s)", numero, id),
		}, false
	}

	if partes[7] == "" {
		return Jogador{}, ErroEntrada{
			Codigo:   "CAMPO_OBRIGATORIO",
			Mensagem: fmt.Sprintf("linha %d: regiao ausente (jogador %s)", numero, id),
		}, false
	}

	jogador := Jogador{
		ID:           id,
		Habilidade:   habilidade,
		Preferida:    partes[2],
		Alternativas: lerAlternativas(partes[3]),
		EsperaS:      espera,
		LatenciaMs:   latencia,
		Grupo:        lerGrupo(partes[6]),
		Regiao:       partes[7],
	}
	return jogador, ErroEntrada{}, true
}

// lerAlternativas ignora vazio, "-" e duplicatas.
// Se a preferida aparecer de novo na lista, ela permanece aqui;
// preencherAceitas evita repetir essa posição em Aceitas.
func lerAlternativas(campo string) []string {
	if campo == "" || campo == "-" {
		return nil
	}
	pedacos := strings.Split(campo, ",")
	var alternativas []string
	for i := 0; i < len(pedacos); i++ {
		nome := strings.TrimSpace(pedacos[i])
		if nome == "" || nome == "-" {
			continue
		}
		if contem(alternativas, nome) {
			continue
		}
		alternativas = append(alternativas, nome)
	}
	return alternativas
}

func lerGrupo(campo string) string {
	if campo == "" || campo == "-" {
		return ""
	}
	return campo
}

func inteiroCampo(campo string) (int, bool) {
	valor, err := strconv.Atoi(campo)
	if err != nil {
		return 0, false
	}
	return valor, true
}

// completarFila altera cada jogador da fila já validada.
// Efeito: o campo Aceitas de cada registro passa a ser
// {preferida} ∪ alternativas, sem repetição.
func completarFila(fila []Jogador) {
	for i := 0; i < len(fila); i++ {
		preencherAceitas(&fila[i])
	}
}

func preencherAceitas(jogador *Jogador) {
	jogador.Aceitas = nil
	jogador.Aceitas = append(jogador.Aceitas, jogador.Preferida)
	for i := 0; i < len(jogador.Alternativas); i++ {
		alt := jogador.Alternativas[i]
		if alt == jogador.Preferida || contem(jogador.Aceitas, alt) {
			continue
		}
		jogador.Aceitas = append(jogador.Aceitas, alt)
	}
}
