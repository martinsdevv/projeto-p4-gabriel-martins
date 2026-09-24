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
