package main

// ResultadoBusca é o que a busca deixa para a saída.
// Formou falso guarda o motivo de NENHUMA_PARTIDA. A fila restante
// preserva a ordem original dos jogadores que não entraram.
type ResultadoBusca struct {
	Formou       bool
	Motivo       string
	Formacao     Formacao
	Metricas     Metricas
	FilaRestante []string
}

// estadoBusca é o estado mutável da exploração. A cada candidata válida,
// considerar atualiza Formacao e Metricas se a nova supera a guardada.
type estadoBusca struct {
	achou        bool
	formacao     Formacao
	metricas     Metricas
	ids          []string
	tupla        [10]string
	houvePosicao bool
	houveGrupo   bool
}

// buscarPartida percorre subconjuntos de 10 jogadores, atribui posições
// e reparte equipes. O melhor candidato fica em estadoBusca.
func buscarPartida(fila []Jogador) ResultadoBusca {
	if len(fila) < JogadoresPorPartida {
		return ResultadoBusca{
			Motivo:       "JOGADORES_INSUFICIENTES",
			FilaRestante: idsNaOrdem(fila, nil),
		}
	}

	var estado estadoBusca
	gerarSubconjuntos(len(fila), JogadoresPorPartida, func(escolhidos []int) {
		explorarSubconjunto(fila, escolhidos, &estado)
	})

	if !estado.achou {
		motivo := "COMPOSICAO_IMPOSSIVEL"
		if estado.houveGrupo {
			motivo = "DIFERENCA_HABILIDADE"
		} else if estado.houvePosicao {
			motivo = "GRUPO_INCOMPATIVEL"
		}
		return ResultadoBusca{
			Motivo:       motivo,
			FilaRestante: idsNaOrdem(fila, nil),
		}
	}

	selecionados := map[string]bool{}
	for i := 0; i < len(estado.ids); i++ {
		selecionados[estado.ids[i]] = true
	}
	return ResultadoBusca{
		Formou:       true,
		Formacao:     estado.formacao,
		Metricas:     estado.metricas,
		FilaRestante: idsNaOrdem(fila, selecionados),
	}
}

func explorarSubconjunto(fila []Jogador, escolhidos []int, estado *estadoBusca) {
	ordem := ordenarPorRigidez(fila, escolhidos)
	posDe := make([]int, len(fila))
	for i := 0; i < len(posDe); i++ {
		posDe[i] = -1
	}
	var contagem [5]int
	explorarPosicoes(fila, ordem, 0, posDe, contagem[:], estado)
}

// ordenarPorRigidez copia os índices e os ordena por quem aceita menos
// posições. O laço de atribuição encontra cedo um beco sem saída.
func ordenarPorRigidez(fila []Jogador, escolhidos []int) []int {
	ordem := make([]int, len(escolhidos))
	for i := 0; i < len(escolhidos); i++ {
		ordem[i] = escolhidos[i]
	}
	for i := 1; i < len(ordem); i++ {
		atual := ordem[i]
		j := i
		for j > 0 && len(fila[ordem[j-1]].Aceitas) > len(fila[atual].Aceitas) {
			ordem[j] = ordem[j-1]
			j--
		}
		ordem[j] = atual
	}
	return ordem
}

func explorarPosicoes(fila []Jogador, ordem []int, passo int, posDe []int, contagem []int, estado *estadoBusca) {
	if passo == len(ordem) {
		estado.houvePosicao = true
		explorarEquipes(fila, ordem, posDe, estado)
		return
	}

	indice := ordem[passo]
	jogador := &fila[indice]
	for a := 0; a < len(jogador.Aceitas); a++ {
		posicao := indicePosicao(jogador.Aceitas[a])
		if posicao < 0 || contagem[posicao] >= 2 {
			continue
		}
		contagem[posicao] = contagem[posicao] + 1
		posDe[indice] = posicao
		explorarPosicoes(fila, ordem, passo+1, posDe, contagem, estado)
		contagem[posicao] = contagem[posicao] - 1
		posDe[indice] = -1
	}
}

func explorarEquipes(fila []Jogador, ordem []int, posDe []int, estado *estadoBusca) {
	var pares [5][2]int
	var cheios [5]int
	for i := 0; i < len(ordem); i++ {
		indice := ordem[i]
		posicao := posDe[indice]
		pares[posicao][cheios[posicao]] = indice
		cheios[posicao] = cheios[posicao] + 1
	}

	for mascara := 0; mascara < 32; mascara++ {
		var formacao Formacao
		for p := 0; p < 5; p++ {
			primeiroNaA := mascara&(1<<p) != 0
			indiceA := pares[p][0]
			indiceB := pares[p][1]
			if !primeiroNaA {
				indiceA, indiceB = indiceB, indiceA
			}
			formacao.EquipeA[p] = Atribuicao{Jogador: &fila[indiceA], Posicao: posicoes[p]}
			formacao.EquipeB[p] = Atribuicao{Jogador: &fila[indiceB], Posicao: posicoes[p]}
		}
		rotularEquipes(&formacao)

		_, regra, ok := avaliarFormacao(fila, formacao)
		if regra == "R6" || ok {
			estado.houveGrupo = true
		}
		if !ok {
			continue
		}
		considerar(estado, formacao, fila)
	}
}

// rotularEquipes troca A e B quando o menor id dos dez está na equipe B.
// A ordem das posições dentro da equipe já é a canônica.
func rotularEquipes(formacao *Formacao) {
	menor := formacao.EquipeA[0].Jogador.ID
	menorNaA := true
	for i := 0; i < 5; i++ {
		idA := formacao.EquipeA[i].Jogador.ID
		idB := formacao.EquipeB[i].Jogador.ID
		if idA < menor {
			menor = idA
			menorNaA = true
		}
		if idB < menor {
			menor = idB
			menorNaA = false
		}
	}
	if !menorNaA {
		formacao.EquipeA, formacao.EquipeB = formacao.EquipeB, formacao.EquipeA
	}
}

func considerar(estado *estadoBusca, formacao Formacao, fila []Jogador) {
	metricas, _, ok := avaliarFormacao(fila, formacao)
	if !ok {
		return
	}
	ids := idsDaFormacao(formacao)
	tupla := tuplaCanonico(formacao)
	if estado.achou && !supera(metricas, ids, tupla, estado) {
		return
	}
	estado.achou = true
	estado.formacao = formacao
	estado.metricas = metricas
	estado.ids = ids
	estado.tupla = tupla
}

func supera(metricas Metricas, ids []string, tupla [10]string, estado *estadoBusca) bool {
	if metricas.Qualidade != estado.metricas.Qualidade {
		return metricas.Qualidade > estado.metricas.Qualidade
	}
	if metricas.DiferencaHabilidade != estado.metricas.DiferencaHabilidade {
		return metricas.DiferencaHabilidade < estado.metricas.DiferencaHabilidade
	}
	if metricas.ForaPreferencia != estado.metricas.ForaPreferencia {
		return metricas.ForaPreferencia < estado.metricas.ForaPreferencia
	}
	if metricas.DesequilibrioAlt != estado.metricas.DesequilibrioAlt {
		return metricas.DesequilibrioAlt < estado.metricas.DesequilibrioAlt
	}
	if metricas.AmplitudeLatencia != estado.metricas.AmplitudeLatencia {
		return metricas.AmplitudeLatencia < estado.metricas.AmplitudeLatencia
	}
	if metricas.EsperaSoma != estado.metricas.EsperaSoma {
		return metricas.EsperaSoma > estado.metricas.EsperaSoma
	}
	if !idsIguais(ids, estado.ids) {
		return listaMenor(ids, estado.ids)
	}
	return tuplaMenor(tupla, estado.tupla)
}

func idsDaFormacao(formacao Formacao) []string {
	ids := make([]string, 0, 10)
	for i := 0; i < 5; i++ {
		ids = append(ids, formacao.EquipeA[i].Jogador.ID)
		ids = append(ids, formacao.EquipeB[i].Jogador.ID)
	}
	for i := 1; i < len(ids); i++ {
		atual := ids[i]
		j := i
		for j > 0 && ids[j-1] > atual {
			ids[j] = ids[j-1]
			j--
		}
		ids[j] = atual
	}
	return ids
}

func tuplaCanonico(formacao Formacao) [10]string {
	var tupla [10]string
	for i := 0; i < 5; i++ {
		tupla[i] = formacao.EquipeA[i].Jogador.ID
		tupla[i+5] = formacao.EquipeB[i].Jogador.ID
	}
	return tupla
}

func idsIguais(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func listaMenor(a []string, b []string) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func tuplaMenor(a [10]string, b [10]string) bool {
	for i := 0; i < 10; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func idsNaOrdem(fila []Jogador, selecionados map[string]bool) []string {
	var restante []string
	for i := 0; i < len(fila); i++ {
		if selecionados[fila[i].ID] {
			continue
		}
		restante = append(restante, fila[i].ID)
	}
	return restante
}

func gerarSubconjuntos(n int, escolher int, visitar func([]int)) {
	if n < escolher {
		return
	}
	indice := make([]int, escolher)
	for i := 0; i < escolher; i++ {
		indice[i] = i
	}
	for {
		copia := make([]int, escolher)
		for i := 0; i < escolher; i++ {
			copia[i] = indice[i]
		}
		visitar(copia)

		i := escolher - 1
		for i >= 0 && indice[i] == n-escolher+i {
			i--
		}
		if i < 0 {
			return
		}
		indice[i] = indice[i] + 1
		for j := i + 1; j < escolher; j++ {
			indice[j] = indice[j-1] + 1
		}
	}
}
