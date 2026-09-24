# Decisões de implementação — paradigma imperativo

**[P4-ETAPA-03]**

A solução está em Go, a linguagem escolhida na Etapa 1 para o paradigma imperativo. Os tipos são registros. O que o programa faz está em funções, com variáveis, atribuição e laços.

## O que o programa guarda e altera

A fila é um slice. Começa vazia e recebe um jogador a cada linha útil, por `append`.

Durante a validação aparecem dois mapas locais, usados e descartados na própria função:

- `vistos` marca o id que já passou. Se o mesmo id aparece de novo, a entrada é rejeitada.
- `contagem` e `ordem` contam os grupos. `ordem` lembra quem apareceu primeiro, para o laço seguinte percorrer os grupos sempre na mesma sequência.

`Aceitas`, em cada jogador, fica vazio até a fila inteira ser aceita. Aí `completarFila` percorre os índices e `preencherAceitas` reescreve esse campo: a preferida mais as alternativas, sem repetir. A região é guardada, mas nesta versão do contrato não entra na validação nem na qualidade.

Uma formação é um registro com duas equipes de cinco vagas. Cada vaga aponta para um jogador da fila e traz a posição atribuída. `Metricas` nasce zerada. `acumularEquipe` recebe o ponteiro e vai somando habilidade, espera, latência e deslocamentos. Os mapas usados na checagem de ids e de grupos também são locais: se o id já estava em `vistos`, a formação falha em R1; se o grupo já estava na outra equipe, falha em R5.

A busca guarda o seu estado em `estadoBusca` enquanto explora. Lá ficam a formação campeã, as métricas, a lista ordenada de ids e a tupla canônica. `achou` começa falso. A primeira candidata válida ocupa o lugar. As próximas só substituem esses campos se `supera` disser que são melhores. `houvePosicao` e `houveGrupo` são marcados no caminho e só servem quando nenhuma candidata vale, para escolher o motivo. `contagem` e `posDe` sobem e descem na atribuição de posições: o laço marca, a chamada seguinte explora, a volta desfaz.

A fila em si não é reordenada. A ordem de chegada só é relida em `idsNaOrdem`, na hora de montar quem ficou de fora.

## Efeitos colaterais

Escrever `Aceitas` muda o registro que o chamador já tem. Fora isso, avaliar uma formação não mexe na fila nem nas equipes: devolve um `Metricas` preenchido. A fila entra nessa conta só para confirmar que cada selecionado está nela.

A saída escreve no `io.Writer` que recebe. No `main`, esse escritor é a saída padrão. Ler o arquivo, ou a entrada padrão, também é efeito externo, e fica isolado em `lerTexto`.

## Controle e funções

Os laços com índice percorrem linhas, jogadores, alternativas, posições e grupos. Em cada regra de entrada, um `if` ou devolve o erro na hora ou segue. Se a mesma fila tiver mais de um problema, vale o primeiro desta ordem:

1. linha com número de campos diferente de 8, identificador vazio, número ilegível ou região vazia (`CAMPO_OBRIGATORIO`);
2. em cada jogador, na ordem da fila: `HABILIDADE_INVALIDA`, `ESPERA_INVALIDA`, `LATENCIA_INVALIDA`, `PREFERIDA_AUSENTE`, `POSICAO_DESCONHECIDA` (preferida, depois alternativas);
3. `ID_DUPLICADO`;
4. `GRUPO_MAIOR_QUE_EQUIPE`.

Os casos I01–I04 têm uma violação cada, e as mensagens repetem o texto de `testes/casos.md`.

A formação válida passa por R1, R3, R4, R5, o cálculo e R6. O primeiro fracasso volta. R2 não tem laço próprio: o registro já tem cinco vagas por equipe, e jogador repetido nas duas cai em R1.

| Função | O que faz |
| --- | --- |
| `lerTexto` | Lê o arquivo ou a entrada padrão. |
| `lerFila`, `lerLinha` | Separam os oito campos e montam a fila. |
| `validarFila`, `validarJogador`, `validarIDs`, `validarGrupos` | Aplicam as regras de entrada, uma função por regra. |
| `completarFila`, `preencherAceitas` | Preenchem `Aceitas` depois que a fila passou. |
| `avaliarFormacao` | Percorre as restrições e devolve métricas, regra violada e validade. |
| `jogadoresDistintos`, `posicoesAceitas`, `composicaoValida`, `gruposUnidos` | R1, R3, R4 e R5. |
| `calcularMetricas`, `acumularEquipe`, `limiteDiferenca` | Somas, Q e o teto de habilidade. |
| `buscarPartida` | Explora as candidatas e fica com a melhor. |
| `escreverErroEntrada`, `escreverPartida`, `escreverNenhumaPartida` | Gravam o texto do contrato. |
| `executar` | Encadeia leitura, validação, busca e escrita. |

A linha de um jogador é:

```text
id | habilidade | preferida | alternativas | espera_s | latencia_ms | grupo | regiao
```

`-` em alternativas ou em grupo significa vazio. Duplicata em alternativas é ignorada. Nome de posição compara exato, com maiúsculas e minúsculas.

Q usa `|somaA - somaB|` inteiro, como no contrato. A divisão por 5 aparece só em `DiferencaHabilidade`. A espera mexe no limite de R6 e não entra em Q.

## Como a busca escolhe

Há três laços, de fora para dentro:

1. `gerarSubconjuntos` lista combinações de 10 índices.
2. `explorarPosicoes` atribui uma posição aceita, no máximo dois jogadores por posição. Quem aceita menos posições é tentado primeiro.
3. `explorarEquipes` testa as 32 formas de mandar um jogador de cada posição para cada equipe.

`rotularEquipes` troca A e B quando o menor id dos dez está na equipe B. `supera` segue a ordem da Seção 6 do contrato: Q, diferença, deslocados, desequilíbrio, latência, soma de espera, ids e tupla canônica.

Sem partida válida, o motivo sai do que a exploração chegou a ver:

| O que aconteceu | Motivo |
| --- | --- |
| Menos de 10 jogadores | `JOGADORES_INSUFICIENTES` |
| Houve composição com grupos, e todas caíram em R6 | `DIFERENCA_HABILIDADE` |
| Houve atribuição de posições, e o grupo impediu as equipes | `GRUPO_INCOMPATIVEL` |
| Nenhuma atribuição de posições | `COMPOSICAO_IMPOSSIVEL` |

N06 em `testes/casos.md` foi alinhado a essa ordem. `p11` e `p12` entram em equipes opostas (Q 9994, amplitude 6). `p02` e `p03` ficam na fila. O subconjunto de N01 tem Q 9993 e perde.

## Saída

Entrada inválida escreve `ERRO_ENTRADA` e o processo termina com código 1. Fila válida segue para a busca e escreve `PARTIDA_FORMADA` ou `NENHUMA_PARTIDA`, com código 0. A média de espera é `EsperaSoma / 10`, com uma casa decimal. Lista vazia de ids sai como `[]`.

O relatório (`-relatorio`) repete cada caso numa linha. A coluna `espera` é o tempo de fila dos dez selecionados, não o tempo que o programa levou para calcular.

## Por que a solução é predominantemente imperativa

O programa é uma sequência de passos sobre dados que mudam. A fila cresce na leitura. Os mapas de validação são preenchidos e consultados em laços. A busca guarda um campeão e o substitui por atribuição quando aparece uma formação melhor. A saída é escrita no destino que foi passado. Cada função recebe a fila, a formação ou o escritor por parâmetro, e quem chamou decide o próximo passo a partir do retorno.
