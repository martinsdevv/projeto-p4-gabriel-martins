package main

import (
	"fmt"
	"io"
)

// escreverErroEntrada grava o formato da Seção 7.1 do contrato.
// Efeito colateral: escreve no io.Writer recebido. Em main, esse
// escritor é a saída padrão. Não há fila_restante neste resultado.
func escreverErroEntrada(saida io.Writer, erro ErroEntrada) {
	fmt.Fprintf(saida, "status: ERRO_ENTRADA\n")
	fmt.Fprintf(saida, "codigo: %s\n", erro.Codigo)
	fmt.Fprintf(saida, "mensagem: %s\n", erro.Mensagem)
}

// escreverPartida grava a Seção 7.3. A média de espera sai com uma casa decimal.
func escreverPartida(saida io.Writer, resultado ResultadoBusca) {
	m := resultado.Metricas
	media := float64(m.EsperaSoma) / float64(JogadoresPorPartida)
	fmt.Fprintf(saida, "status: PARTIDA_FORMADA\n")
	fmt.Fprintf(saida, "equipe_a: %s\n", textoEquipe(resultado.Formacao.EquipeA))
	fmt.Fprintf(saida, "equipe_b: %s\n", textoEquipe(resultado.Formacao.EquipeB))
	fmt.Fprintf(saida, "diferenca_habilidade: %d\n", m.DiferencaHabilidade)
	fmt.Fprintf(saida, "jogadores_fora_preferencia: %d\n", m.ForaPreferencia)
	fmt.Fprintf(saida, "desequilibrio_alternativas: %d\n", m.DesequilibrioAlt)
	fmt.Fprintf(saida, "qualidade: %d\n", m.Qualidade)
	fmt.Fprintf(saida, "espera: {min: %d, max: %d, media: %.1f}\n", m.EsperaMin, m.EsperaMax, media)
	fmt.Fprintf(saida, "amplitude_latencia: %d\n", m.AmplitudeLatencia)
	fmt.Fprintf(saida, "fila_restante: %s\n", textoIds(resultado.FilaRestante))
}

// escreverNenhumaPartida grava a Seção 7.2. A fila restante mantém a ordem de chegada.
func escreverNenhumaPartida(saida io.Writer, resultado ResultadoBusca) {
	fmt.Fprintf(saida, "status: NENHUMA_PARTIDA\n")
	fmt.Fprintf(saida, "motivo: %s\n", resultado.Motivo)
	fmt.Fprintf(saida, "fila_restante: %s\n", textoIds(resultado.FilaRestante))
}

func textoEquipe(equipe [5]Atribuicao) string {
	texto := ""
	for i := 0; i < len(equipe); i++ {
		if i > 0 {
			texto = texto + ", "
		}
		texto = texto + "(" + equipe[i].Jogador.ID + ", " + equipe[i].Posicao + ")"
	}
	return texto
}

func textoIds(ids []string) string {
	if len(ids) == 0 {
		return "[]"
	}
	texto := "["
	for i := 0; i < len(ids); i++ {
		if i > 0 {
			texto = texto + ", "
		}
		texto = texto + ids[i]
	}
	return texto + "]"
}
