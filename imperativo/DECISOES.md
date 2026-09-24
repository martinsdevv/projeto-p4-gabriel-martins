# Decisões de implementação - paradigma imperativo

**[P4-ETAPA-03]**

Linguagem: Go, conforme a escolha registrada na Etapa 1. Os tipos são registros e os procedimentos são funções. O fluxo está escrito com variáveis, atribuição e laços.

## Partes


| Parte | Conteúdo                                                            | Estado   |
| ----- | ------------------------------------------------------------------- | -------- |
| 1     | Modelo da fila, leitura e validação da entrada                      | feita    |
| 2     | Formação, restrições R1–R6 e qualidade Q                            | pendente |
| 3     | Busca da formação ótima, com estado do melhor candidato             | pendente |
| 4     | Saídas `PARTIDA_FORMADA` e `NENHUMA_PARTIDA`                        | pendente |
| 5     | Execução dos 18 casos de `testes/casos.md` e fechamento deste texto | pendente |




## Parte 1 - estado, efeitos e organização



### Estados mantidos

- `fila`, um slice que começa vazio e recebe um jogador a cada linha não vazia (`append`).
- `vistos`, mapa preenchido durante a varredura de identificadores. A escrita `vistos[id] = true` registra que aquele id já apareceu.
- `contagem` e `ordem`, preenchidos durante a varredura dos grupos. `ordem` guarda a primeira aparição de cada grupo para o laço seguinte percorrer os grupos de forma estável.
- `Aceitas`, em cada `Jogador`. Permanece vazio até `completarFila`. Só é preenchido se a fila inteira for válida.

A região é armazenada e não participa da validação nem, nesta versão do contrato, da qualidade.

### Operações que modificam estado

- `lerFila` aumenta `fila`.
- `validarIDs` e `validarGrupos` escrevem nos mapas locais e descartam esses mapas ao retornar.
- `preencherAceitas` recebe um ponteiro e reatribui `Aceitas` do registro apontado. `completarFila` percorre a fila e chama esse procedimento em cada índice.



### Efeitos colaterais

- A mutação de `Jogador.Aceitas` é o efeito sobre o dado do chamador.
- `escreverErroEntrada` escreve no `io.Writer` recebido como parâmetro. No `main`, esse escritor é `os.Stdout`.

A leitura do arquivo ou da entrada padrão também é efeito externo, isolado em `lerTexto`.

### Estruturas de controle

Laços `for` com índice percorrem linhas, jogadores, alternativas, posições e grupos. Em cada regra, um `if` decide se a função retorna na hora com o erro ou segue para a próxima verificação. A ordem desses retornos é a ordem do contrato prático desta implementação (ver abaixo).

### Subprogramas


| Função                                          | Papel                                                            |
| ----------------------------------------------- | ---------------------------------------------------------------- |
| `lerTexto`                                      | Obtém o texto bruto (arquivo ou entrada padrão).                 |
| `lerFila`                                       | Conduz o fluxo: ler linhas, validar, completar.                  |
| `lerLinha`                                      | Separa os oito campos de uma linha.                              |
| `validarFila`                                   | Aplica as regras de entrada na ordem fixa.                       |
| `validarJogador`, `validarIDs`, `validarGrupos` | Uma regra por subprograma, com os dados recebidos por parâmetro. |
| `completarFila`, `preencherAceitas`             | Atualizam `Aceitas` depois da validação.                         |
| `escreverErroEntrada`                           | Produz o texto de `ERRO_ENTRADA`.                                |


`main` só encadeia esses passos.

### Por que esta parte é imperativa

O programa descreve uma sequência de passos sobre dados que mudam: a fila cresce, os mapas de controle são atualizados, e `Aceitas` é escrito no registro já existente. Cada subprograma recebe explicitamente o que precisa (a fila, um jogador, o escritor) e o chamador decide o que fazer com o retorno antes de seguir. O resultado observável depende dessa ordem de execução.

### Ordem da primeira violação

Quando a mesma fila tem mais de um problema, o código reportado é o primeiro desta sequência:

1. linha com número de campos diferente de 8, identificador vazio, número ilegível ou região vazia (`CAMPO_OBRIGATORIO`);
2. em cada jogador, na ordem da fila: `HABILIDADE_INVALIDA`, `ESPERA_INVALIDA`, `LATENCIA_INVALIDA`, `PREFERIDA_AUSENTE`, `POSICAO_DESCONHECIDA` (preferida, depois alternativas);
3. `ID_DUPLICADO`;
4. `GRUPO_MAIOR_QUE_EQUIPE`.

Os casos I01–I04 têm uma violação cada. As mensagens desses quatro casos reproduzem o texto de `testes/casos.md`.

### Saída provisória

Entrada inválida já sai no formato do contrato. Entrada válida, nesta parte, imprime `status: FILA_VALIDA` e a quantidade de jogadores. Esse status não faz parte do contrato; as partes 3 e 4 o substituem por `PARTIDA_FORMADA` ou `NENHUMA_PARTIDA`.

### Leitura da linha

```text
id | habilidade | preferida | alternativas | espera_s | latencia_ms | grupo | regiao
```

`-` em alternativas ou em grupo significa vazio. Duplicatas em alternativas são ignoradas. A comparação de posição é exata, com distinção de maiúsculas.