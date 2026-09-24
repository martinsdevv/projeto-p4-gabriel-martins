// Registros de dados da fila e parâmetros fixos do contrato da Etapa 2.
// O comportamento (leitura, validação, busca) fica em subprogramas;
// estes tipos só guardam valores.
package main

// Parâmetros do contrato. São os mesmos em todos os casos de teste.
const (
	JogadoresPorPartida     = 10
	JogadoresPorEquipe      = 5
	EquipesPorPartida       = 2
	LimiteBaseDiferenca     = 200
	BonusEsperaDivisor      = 4
	LimiteAbsolutoDiferenca = 400
)

// posicoes é a ordem canônica da composição e da saída.
// O laço de validação e, depois, o de busca percorrem este slice.
var posicoes = []string{"Top", "Jungle", "Mid", "ADC", "Support"}

// Jogador é um registro da fila.
// Aceitas começa vazio e é preenchido por completarFila depois que a
// fila inteira passa na validação.
type Jogador struct {
	ID           string
	Habilidade   int
	Preferida    string
	Alternativas []string
	EsperaS      int
	LatenciaMs   int
	Grupo        string
	Regiao       string
	Aceitas      []string
}

// ErroEntrada é o resultado observável de uma fila que viola o modelo.
type ErroEntrada struct {
	Codigo   string
	Mensagem string
}

// Atribuicao é um jogador já colocado numa posição.
// Jogador aponta para o registro que está na fila.
type Atribuicao struct {
	Jogador *Jogador
	Posicao string
}

// Formacao é uma candidata: duas equipes de cinco, na ordem
// Top, Jungle, Mid, ADC, Support. A busca é quem preenche as vagas;
// avaliarFormacao só confere o registro que já chegou montado.
type Formacao struct {
	EquipeA [5]Atribuicao
	EquipeB [5]Atribuicao
}

// Metricas guarda os acumuladores da avaliação. avaliarFormacao
// preenche os campos por atribuição, na ordem das regras.
type Metricas struct {
	SomaA               int
	SomaB               int
	DiferencaHabilidade int
	ForaPreferencia     int
	ForaPrefA           int
	ForaPrefB           int
	DesequilibrioAlt    int
	AmplitudeLatencia   int
	Qualidade           int
	EsperaMin           int
	EsperaMax           int
	EsperaSoma          int
	LimiteDiferenca     int
}

func posicaoValida(nome string) bool {
	for i := 0; i < len(posicoes); i++ {
		if posicoes[i] == nome {
			return true
		}
	}
	return false
}

func contem(lista []string, valor string) bool {
	for i := 0; i < len(lista); i++ {
		if lista[i] == valor {
			return true
		}
	}
	return false
}
