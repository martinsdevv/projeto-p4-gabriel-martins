# Contrato semântico e casos de teste

**[P4-ETAPA-02]**

Este documento transforma a especificação da Etapa 1 em um **contrato de comportamento**: descreve o que o sistema deve produzir, não como deve calcular o resultado internamente.

Os mesmos casos deverão ser utilizados nas quatro implementações posteriores (imperativa, orientada a objetos, funcional e lógica). A implementação pode variar; o resultado observável, não.

---

## 1. O que está no contrato e o que não está

### Faz parte do contrato (o que o programa faz)

- formato da entrada;
- quando uma formação é válida;
- como a qualidade de uma formação válida é medida;
- qual formação deve ser escolhida quando várias são válidas;
- formato da saída;
- distinção entre **entrada inválida** e **nenhuma partida possível**.

### Não faz parte do contrato (como o programa faz)

- linguagem, paradigma, estruturas de dados ou algoritmo de busca;
- ordem de exploração das formações candidatas;
- uso de recursão, backtracking, restrições, objetos ou estado mutável;
- qualquer otimização interna, desde que o resultado final obedeça a este contrato.

Os casos desta etapa usam filas com no máximo 12 jogadores, para que a otimalidade do resultado seja verificável em qualquer paradigma.

---

## 2. Modelo da entrada

A entrada é uma **fila ordenada de jogadores**. A ordem da fila é a ordem de chegada e deve ser preservada na saída em `fila_restante`.

Cada jogador possui exatamente os campos abaixo.

| Campo | Tipo | Significado |
| --- | --- | --- |
| `id` | texto | Identificador único na fila. A comparação lexicográfica usa a ordem do conjunto de caracteres ASCII. |
| `habilidade` | inteiro | Habilidade estimada. Deve ser **maior ou igual a 0**. |
| `preferida` | posição | Exatamente uma posição preferida. |
| `alternativas` | lista de posições | Zero ou mais posições alternativas. Não precisa incluir a preferida. Duplicatas são ignoradas. |
| `espera_s` | inteiro | Tempo de espera na fila, em segundos. Deve ser **maior ou igual a 0**. |
| `latencia_ms` | inteiro | Latência, em milissegundos. Deve ser **maior ou igual a 0**. |
| `grupo` | texto ou vazio | Identificador do grupo. Vazio (`-`) significa jogador solo. |
| `regiao` | texto | Região declarada do jogador. |

### Posições válidas

O conjunto fechado de posições é:

```text
Top
Jungle
Mid
ADC
Support
```

A comparação de nomes de posição é **exata e sensível a maiúsculas**. Qualquer outro valor é entrada inválida.

### Posições aceitas por um jogador

```text
aceitas(jogador) = {preferida} ∪ alternativas
```

O jogador só pode ser colocado em uma posição pertencente a `aceitas(jogador)`.

### Região

A região é um dado obrigatório do jogador, conforme a Etapa 1. **Nesta versão do contrato ela não restringe a validade nem altera a qualidade.** Serve para manter o modelo estável; um critério futuro só poderá usá-la se o contrato for revisado explicitamente.

### Parâmetros fixos do sistema

Os parâmetros abaixo são parte do contrato e não variam entre os casos de teste.

```text
JOGADORES_POR_PARTIDA          = 10
JOGADORES_POR_EQUIPE           = 5
EQUIPES_POR_PARTIDA            = 2
LIMITE_BASE_DIFERENCA          = 200
BONUS_ESPERA_DIVISOR           = 4
LIMITE_ABSOLUTO_DIFERENCA      = 400
```

---

## 3. Formação candidata

Uma formação candidata é composta por:

- um subconjunto de **exatamente 10** jogadores da fila, sem repetição;
- uma partição desses 10 jogadores em duas equipes de 5;
- uma atribuição de posição a cada jogador.

### Rotulação canônica das equipes

Depois de escolhidos os 10 jogadores e as duas equipes:

1. **Equipe A** é a equipe que contém o `id` lexicograficamente menor entre os 10 selecionados.
2. **Equipe B** é a outra equipe.
3. Dentro de cada equipe, os jogadores são apresentados na ordem de posições:

```text
Top, Jungle, Mid, ADC, Support
```

Essa rotulação existe apenas para tornar a saída determinística. Não há vantagem competitiva em ser Equipe A ou Equipe B.

---

## 4. Validade de uma formação

Uma formação é **válida** se e somente se todas as restrições abaixo forem verdadeiras.

### R1 — Quantidade de jogadores

A formação contém exatamente 10 jogadores distintos, todos presentes na fila.

### R2 — Quantidade de equipes

Os 10 jogadores estão em duas equipes de exatamente 5 jogadores. Nenhum jogador pertence às duas equipes.

### R3 — Posição aceita

Cada jogador ocupa uma posição em `aceitas(jogador)`.

### R4 — Composição

Cada equipe possui as cinco posições `{Top, Jungle, Mid, ADC, Support}` **exatamente uma vez**.

### R5 — Grupos

Se dois ou mais jogadores selecionados compartilham o mesmo `grupo` não vazio, todos eles devem estar na mesma equipe.

Um grupo de tamanho 1 (ou jogador com `grupo = -`) não impõe restrição extra.

### R6 — Tolerância de habilidade

Sejam `somaA` e `somaB` as somas das habilidades das duas equipes.

```text
diferenca_habilidade = |somaA - somaB| / 5
max_espera           = máximo de espera_s entre os 10 selecionados
limite_diferenca     = mínimo(
                         LIMITE_ABSOLUTO_DIFERENCA,
                         LIMITE_BASE_DIFERENCA + (max_espera // BONUS_ESPERA_DIVISOR)
                       )
```

`//` é a divisão inteira truncada em direção a zero (como em `10 // 4 = 2`).

A formação só é válida se:

```text
diferenca_habilidade <= limite_diferenca
```

O tempo de espera **não aumenta a qualidade** da partida. Ele apenas amplia a tolerância máxima de diferença de habilidade, evitando espera indefinida por uma formação perfeita.

---

## 5. Qualidade de uma formação válida

A qualidade só é definida para formações válidas. Quanto maior o valor, melhor a formação.

```text
fora_preferencia     = quantidade de jogadores cuja posição atribuída ≠ preferida
fora_pref_A          = quantidade de jogadores fora da preferida na Equipe A
fora_pref_B          = quantidade de jogadores fora da preferida na Equipe B
desequilibrio_alt    = |fora_pref_A - fora_pref_B|
amplitude_latencia   = max(latencia_ms) - min(latencia_ms)   entre os 10 selecionados

Q = 10000
    - 2 * |somaA - somaB|
    - 100 * fora_preferencia
    - 50 * desequilibrio_alt
    - amplitude_latencia
```

A fórmula usa `|somaA - somaB|` em vez da média para permanecer em aritmética inteira. Como cada equipe tem 5 jogadores, isso é equivalente a penalizar `10 * diferenca_habilidade`.

Interpretação dos pesos, alinhada à Etapa 1:

| Critério | Efeito em Q | Leitura |
| --- | --- | --- |
| C1 equilíbrio de habilidade | `-2 * \|somaA - somaB\|` | Menor diferença entre equipes é melhor. |
| C2 posições preferidas | `-100` por jogador deslocado | Usar alternativa só vale a pena se reduzir bastante o desequilíbrio de habilidade. |
| C3 equilíbrio de alternativas | `-50` por unidade de desequilíbrio | Se o número de deslocados for o mesmo, é melhor distribuí-los entre as equipes. |
| C4 latência | `-1` por milissegundo de amplitude | Critério de qualidade, não restrição absoluta. |
| C5 tempo de espera | não entra em Q | Atua somente em R6, relaxando o limite de habilidade. |

---

## 6. Escolha da formação

Se não existir nenhuma formação válida, o sistema **não cria partida**.

Se existirem uma ou mais formações válidas, o sistema deve devolver **uma formação ótima** segundo a ordem total abaixo. O primeiro critério que desempatar decide.

1. maior `Q`;
2. menor `diferenca_habilidade`;
3. menor `fora_preferencia`;
4. menor `desequilibrio_alt`;
5. menor `amplitude_latencia`;
6. maior soma de `espera_s` dos 10 selecionados (atender quem já espera);
7. conjunto dos 10 `id` lexicograficamente menor, comparando as listas de IDs em ordem crescente;
8. atribuição canônica lexicograficamente menor, comparando a tupla:

```text
(id_Top_A, id_Jungle_A, id_Mid_A, id_ADC_A, id_Support_A,
 id_Top_B, id_Jungle_B, id_Mid_B, id_ADC_B, id_Support_B)
```

Essa ordem existe para que as quatro implementações produzam a **mesma** partida quando várias formações empatam em qualidade.

---

## 7. Saídas

Há três classes de resultado. Elas não se misturam.

### 7.1 Entrada inválida

A fila não obedece ao modelo da Seção 2. O sistema **não tenta** formar partida.

```text
status: ERRO_ENTRADA
codigo: <codigo>
mensagem: <texto livre descritivo>
```

Códigos previstos:

| Código | Condição |
| --- | --- |
| `ID_DUPLICADO` | Dois ou mais jogadores com o mesmo `id`. |
| `HABILIDADE_INVALIDA` | `habilidade < 0`. |
| `ESPERA_INVALIDA` | `espera_s < 0`. |
| `LATENCIA_INVALIDA` | `latencia_ms < 0`. |
| `POSICAO_DESCONHECIDA` | `preferida` ou alguma alternativa fora do conjunto de posições. |
| `PREFERIDA_AUSENTE` | Jogador sem posição preferida. |
| `GRUPO_MAIOR_QUE_EQUIPE` | Algum grupo não vazio contém 6 ou mais jogadores na fila. |
| `CAMPO_OBRIGATORIO` | Falta algum campo obrigatório. |

Se várias violações existirem na mesma entrada, o sistema pode reportar qualquer uma delas. Os casos de teste abaixo contém **uma única** violação cada, para o código ser verificável.

### 7.2 Nenhuma partida

A entrada é bem formada, mas não existe formação válida.

```text
status: NENHUMA_PARTIDA
motivo: <motivo>
fila_restante: <ids na ordem original da fila>
```

Motivos previstos:

| Motivo | Quando usar |
| --- | --- |
| `JOGADORES_INSUFICIENTES` | Menos de 10 jogadores na fila. |
| `COMPOSICAO_IMPOSSIVEL` | Não há como atribuir posições respeitando R3 e R4, mesmo ignorando habilidade. |
| `DIFERENCA_HABILIDADE` | Existem composições de posições e grupos, mas nenhuma cabe em R6. |
| `GRUPO_INCOMPATIVEL` | A restrição de grupo impede qualquer formação que satisfaria as demais regras. |

### 7.3 Partida formada

```text
status: PARTIDA_FORMADA
equipe_a: lista de 5 pares (id, posicao) na ordem Top, Jungle, Mid, ADC, Support
equipe_b: lista de 5 pares (id, posicao) na mesma ordem
diferenca_habilidade: |somaA - somaB| / 5
jogadores_fora_preferencia: inteiro
desequilibrio_alternativas: inteiro
qualidade: Q
espera: {min, max, media} em segundos, apenas dos 10 selecionados
amplitude_latencia: inteiro
fila_restante: ids dos não selecionados, na ordem original da fila
```

`media` de espera é a média aritmética das 10 esperas, podendo ter uma casa decimal (`45.0`, `30.5`, etc.).

---

## 8. Convenções dos casos de teste

Cada caso possui identificador, entrada, saída esperada e descrição.

Notação compacta de um jogador:

```text
id | habilidade | preferida | alternativas | espera_s | latencia_ms | grupo | regiao
```

Lista vazia de alternativas é escrita `-`. Grupo vazio é escrito `-`.

A menos que o caso diga o contrário, `regiao = BR`.

Resumo:

| Identificador | Classe | Intenção |
| --- | --- | --- |
| N01 | Normal | Partida equilibrada com todos nas preferidas |
| N02 | Normal | Escolher o menor desequilíbrio de habilidade |
| N03 | Normal | Uso obrigatório de uma posição alternativa |
| N04 | Normal | Tempo de espera amplia a tolerância de habilidade |
| N05 | Normal | Grupo permanece na mesma equipe |
| N06 | Normal | Fila com mais de 10 jogadores: selecionar o melhor subconjunto |
| N07 | Normal | Não usar alternativa se existe formação só com preferidas |
| N08 | Normal | Distribuir jogadores deslocados entre as duas equipes |
| N09 | Normal | Latência como critério de qualidade na seleção |
| N10 | Normal | Dupla ADC + Support na mesma equipe |
| L01 | Limite | Menos de 10 jogadores |
| L02 | Limite | Todos aceitam apenas Mid |
| L03 | Limite | Diferença de habilidade acima da tolerância, espera baixa |
| L04 | Limite | Diferença exatamente igual ao limite, com espera zero |
| I01 | Entrada inválida | Identificadores duplicados |
| I02 | Entrada inválida | Habilidade negativa |
| I03 | Entrada inválida | Posição desconhecida |
| I04 | Entrada inválida | Grupo com mais de 5 jogadores |

---

## 9. Casos normais

### N01 — Partida equilibrada com todos nas preferidas

**Descrição:** Dez jogadores com habilidades próximas, duas pessoas por posição, todos podendo ocupar apenas a posição preferida. Existe formação válida com diferença de habilidade zero. O sistema deve devolvê-la, sem deslocar ninguém.

**Entrada:**

```text
p01 | 1500 | Mid     | - | 60 | 20 | - | BR
p02 | 1510 | Jungle  | - | 60 | 22 | - | BR
p03 | 1490 | Top     | - | 60 | 18 | - | BR
p04 | 1520 | ADC     | - | 60 | 25 | - | BR
p05 | 1505 | Support | - | 60 | 21 | - | BR
p06 | 1500 | Mid     | - | 30 | 19 | - | BR
p07 | 1495 | Jungle  | - | 30 | 23 | - | BR
p08 | 1510 | Top     | - | 30 | 20 | - | BR
p09 | 1515 | ADC     | - | 30 | 24 | - | BR
p10 | 1495 | Support | - | 30 | 22 | - | BR
```

**Cálculo relevante:**

```text
Equipe A (contém p01): p03 Top 1490, p02 Jungle 1510, p01 Mid 1500, p09 ADC 1515, p05 Support 1505
somaA = 7520
Equipe B: p08 Top 1510, p07 Jungle 1495, p06 Mid 1500, p04 ADC 1520, p10 Support 1495
somaB = 7520
diferenca_habilidade = 0
fora_preferencia = 0
amplitude_latencia = 25 - 18 = 7
Q = 10000 - 7 = 9993
```

Há outra formação também com diferença 0 (`p01` com `p08, p07, p04, p10`). A ordem total da Seção 6 escolhe a atribuição cujo primeiro `id` de Top na Equipe A é `p03`, menor que `p08`.

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p03, Top), (p02, Jungle), (p01, Mid), (p09, ADC), (p05, Support)
equipe_b: (p08, Top), (p07, Jungle), (p06, Mid), (p04, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9993
espera: {min: 30, max: 60, media: 45.0}
amplitude_latencia: 7
fila_restante: []
```

---

### N02 — Preferir a menor diferença de habilidade

**Descrição:** Há dois jogadores por posição. Cinco jogadores têm habilidade 2000 e cinco têm 1000. Formações 5–0 e 4–1 violam R6. A única família válida é 3 jogadores fortes + 2 fracos versus 2 fortes + 3 fracos, com diferença 200. O sistema deve escolher essa família, não uma formação mais desequilibrada.

**Entrada:**

```text
p01 | 2000 | Top     | - | 10 | 20 | - | BR
p02 | 1000 | Top     | - | 10 | 20 | - | BR
p03 | 2000 | Jungle  | - | 10 | 20 | - | BR
p04 | 1000 | Jungle  | - | 10 | 20 | - | BR
p05 | 2000 | Mid     | - | 10 | 20 | - | BR
p06 | 1000 | Mid     | - | 10 | 20 | - | BR
p07 | 2000 | ADC     | - | 10 | 20 | - | BR
p08 | 1000 | ADC     | - | 10 | 20 | - | BR
p09 | 2000 | Support | - | 10 | 20 | - | BR
p10 | 1000 | Support | - | 10 | 20 | - | BR
```

**Cálculo relevante:**

```text
limite_diferenca = min(400, 200 + (10 // 4)) = 202
5 fortes vs 5 fracos → diferença 1000 → inválida
4 fortes vs 1 forte  → diferença 600  → inválida
3 fortes vs 2 fortes → diferença 200  → válida

Equipe A contém p01 (Top 2000). A atribuição canônica de menor tupla é:
A: p01 Top 2000, p03 Jungle 2000, p05 Mid 2000, p08 ADC 1000, p10 Support 1000  → soma 8000
B: p02 Top 1000, p04 Jungle 1000, p06 Mid 1000, p07 ADC 2000, p09 Support 2000  → soma 7000
Q = 10000 - 2 * 1000 = 8000
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p07, ADC), (p09, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 10, max: 10, media: 10.0}
amplitude_latencia: 0
fila_restante: []
```

---

### N03 — Posição alternativa obrigatória

**Descrição:** Há mais pretendentes a Mid do que vagas de Mid, e apenas um jogador tem Top como preferida. A única formação válida coloca `p02` em Top, que é alternativa dele. A partida continua válida porque Top está em `aceitas(p02)`.

**Entrada:**

```text
p01 | 1500 | Mid     | Support | 20 | 20 | - | BR
p02 | 1500 | Mid     | Top     | 20 | 20 | - | BR
p03 | 1500 | Jungle  | -       | 20 | 20 | - | BR
p04 | 1500 | ADC     | -       | 20 | 20 | - | BR
p05 | 1500 | Support | -       | 20 | 20 | - | BR
p06 | 1500 | Mid     | -       | 20 | 20 | - | BR
p07 | 1500 | Jungle  | -       | 20 | 20 | - | BR
p08 | 1500 | Top     | -       | 20 | 20 | - | BR
p09 | 1500 | ADC     | -       | 20 | 20 | - | BR
p10 | 1500 | Support | -       | 20 | 20 | - | BR
```

**Cálculo relevante:**

```text
Tops possíveis: p08 (preferida) e p02 (alternativa). São necessários 2 Tops.
Portanto p02 deve jogar Top. p01 permanece Mid.
fora_preferencia = 1 (apenas p02)
desequilibrio_alt = 1 (um deslocado em uma equipe, zero na outra)
Q = 10000 - 100 - 50 = 9850
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p02, Top), (p03, Jungle), (p01, Mid), (p04, ADC), (p05, Support)
equipe_b: (p08, Top), (p07, Jungle), (p06, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 1
desequilibrio_alternativas: 1
qualidade: 9850
espera: {min: 20, max: 20, media: 20.0}
amplitude_latencia: 0
fila_restante: []
```

---

### N04 — Tempo de espera amplia a tolerância

**Descrição:** Um jogador muito mais habilidoso espera há 400 segundos. A diferença resultante é 260, acima do limite-base 200, mas dentro do limite relaxado pelo tempo de espera. O sistema deve formar a partida.

**Entrada:**

```text
p01 | 2800 | Top     | - | 400 | 20 | - | BR
p02 | 1500 | Top     | - |  10 | 20 | - | BR
p03 | 1500 | Jungle  | - |  10 | 20 | - | BR
p04 | 1500 | Jungle  | - |  10 | 20 | - | BR
p05 | 1500 | Mid     | - |  10 | 20 | - | BR
p06 | 1500 | Mid     | - |  10 | 20 | - | BR
p07 | 1500 | ADC     | - |  10 | 20 | - | BR
p08 | 1500 | ADC     | - |  10 | 20 | - | BR
p09 | 1500 | Support | - |  10 | 20 | - | BR
p10 | 1500 | Support | - |  10 | 20 | - | BR
```

**Cálculo relevante:**

```text
Equipe de p01: 2800 + 4 * 1500 = 8800  → média 1760
Outra equipe: 5 * 1500 = 7500          → média 1500
diferenca_habilidade = 260
max_espera = 400
limite_diferenca = min(400, 200 + (400 // 4)) = 300
260 <= 300 → válida
Q = 10000 - 2 * 1300 = 7400
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 260
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 7400
espera: {min: 10, max: 400, media: 49.0}
amplitude_latencia: 0
fila_restante: []
```

---

### N05 — Jogadores do mesmo grupo na mesma equipe

**Descrição:** `p01` e `p02` entram juntos e são os dois jogadores mais fortes. Separá-los produziria diferença 0, mas violaria R5. O sistema deve mantê-los juntos, mesmo com diferença 200.

**Entrada:**

```text
p01 | 2000 | Top     | - | 30 | 20 | G1 | BR
p02 | 2000 | Jungle  | - | 30 | 20 | G1 | BR
p03 | 1500 | Top     | - | 30 | 20 | -  | BR
p04 | 1500 | Jungle  | - | 30 | 20 | -  | BR
p05 | 1500 | Mid     | - | 30 | 20 | -  | BR
p06 | 1500 | Mid     | - | 30 | 20 | -  | BR
p07 | 1500 | ADC     | - | 30 | 20 | -  | BR
p08 | 1500 | ADC     | - | 30 | 20 | -  | BR
p09 | 1500 | Support | - | 30 | 20 | -  | BR
p10 | 1500 | Support | - | 30 | 20 | -  | BR
```

**Cálculo relevante:**

```text
p01 e p02 obrigatoriamente na mesma equipe
Equipe A: 2000 + 2000 + 1500 + 1500 + 1500 = 8500 → média 1700
Equipe B: 1500 * 5 = 7500                        → média 1500
diferenca_habilidade = 200
limite_diferenca = 200 + (30 // 4) = 207
200 <= 207 → válida
Q = 10000 - 2 * 1000 = 8000
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p02, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p03, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 30, max: 30, media: 30.0}
amplitude_latencia: 0
fila_restante: []
```

A Equipe A contém `p01` e, por R5, também `p02`. Qualquer saída que separe `p01` de `p02` é incorreta.

---

### N06 — Fila com 12 jogadores: escolher o melhor subconjunto

**Descrição:** A fila tem 12 jogadores. Dez deles formam a partida equilibrada de N01. Os outros dois têm habilidade muito alta e piorariam a diferença. O sistema deve selecionar os dez primeiros e deixar os dois fortes na fila.

**Entrada:** os 10 jogadores de N01, nesta ordem, seguidos de:

```text
p11 | 3000 | Top     | - | 5 | 20 | - | BR
p12 | 3000 | Jungle  | - | 5 | 20 | - | BR
```

**Cálculo relevante:**

Incluir `p11` no lugar de um Top ~1500 eleva `|somaA - somaB|` em cerca de 1500 pontos e reduz Q muito abaixo de 9993. O subconjunto `{p01,...,p10}` é ótimo e reproduz N01.

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p03, Top), (p02, Jungle), (p01, Mid), (p09, ADC), (p05, Support)
equipe_b: (p08, Top), (p07, Jungle), (p06, Mid), (p04, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9993
espera: {min: 30, max: 60, media: 45.0}
amplitude_latencia: 7
fila_restante: [p11, p12]
```

---

### N07 — Não usar alternativa se existe formação só com preferidas

**Descrição:** Todos os jogadores aceitam também outras posições, mas já existe cobertura completa das preferidas. Deslocar alguém reduziria Q em pelo menos 100 sem ganho de equilíbrio. O sistema deve manter `fora_preferencia = 0`.

**Entrada:**

```text
p01 | 1500 | Top     | Mid, Jungle     | 15 | 20 | - | BR
p02 | 1500 | Top     | ADC, Support    | 15 | 20 | - | BR
p03 | 1500 | Jungle  | Top, Mid        | 15 | 20 | - | BR
p04 | 1500 | Jungle  | Support, ADC    | 15 | 20 | - | BR
p05 | 1500 | Mid     | Top, Jungle     | 15 | 20 | - | BR
p06 | 1500 | Mid     | ADC, Support    | 15 | 20 | - | BR
p07 | 1500 | ADC     | Mid, Support    | 15 | 20 | - | BR
p08 | 1500 | ADC     | Top, Jungle     | 15 | 20 | - | BR
p09 | 1500 | Support | ADC, Mid        | 15 | 20 | - | BR
p10 | 1500 | Support | Top, Jungle     | 15 | 20 | - | BR
```

**Cálculo relevante:**

```text
Formação só com preferidas: diferença 0, fora_preferencia 0, Q = 10000
Qualquer deslocamento: Q <= 9900
A atribuição canônica coloca em cada posição o menor id disponível daquela preferida.
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p07, ADC), (p09, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 10000
espera: {min: 15, max: 15, media: 15.0}
amplitude_latencia: 0
fila_restante: []
```

---

### N08 — Equilíbrio de jogadores fora da preferida entre as equipes

**Descrição:** Dois jogadores precisam usar alternativa: `p02` em Top e `p01` em Support. Eles ocupam posições diferentes, então podem ficar na mesma equipe ou em equipes opostas. Com `fora_preferencia` igual a 2 nos dois casos, C3 exige a distribuição 1–1 (`desequilibrio_alt = 0`).

**Entrada:**

```text
p01 | 1500 | Mid     | Support | 25 | 20 | - | BR
p02 | 1500 | Mid     | Top     | 25 | 20 | - | BR
p03 | 1500 | Jungle  | -       | 25 | 20 | - | BR
p04 | 1500 | ADC     | -       | 25 | 20 | - | BR
p05 | 1500 | Mid     | -       | 25 | 20 | - | BR
p06 | 1500 | Mid     | -       | 25 | 20 | - | BR
p07 | 1500 | Jungle  | -       | 25 | 20 | - | BR
p08 | 1500 | Top     | -       | 25 | 20 | - | BR
p09 | 1500 | ADC     | -       | 25 | 20 | - | BR
p10 | 1500 | Support | -       | 25 | 20 | - | BR
```

**Cálculo relevante:**

```text
Tops possíveis: p08 e p02 → p02 em Top
Supports possíveis: p10 e p01 → p01 em Support
Mids restantes: p05 e p06
fora_preferencia = 2 sempre

Mesma equipe: desequilibrio_alt = 2 → Q = 10000 - 200 - 100 = 9700
Equipes opostas: desequilibrio_alt = 0 → Q = 10000 - 200 = 9800  ← escolhida

Equipe A contém p01 (Support). Logo p02 deve ir para a Equipe B.
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p08, Top), (p03, Jungle), (p05, Mid), (p04, ADC), (p01, Support)
equipe_b: (p02, Top), (p07, Jungle), (p06, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 2
desequilibrio_alternativas: 0
qualidade: 9800
espera: {min: 25, max: 25, media: 25.0}
amplitude_latencia: 0
fila_restante: []
```

---

### N09 — Latência como critério de qualidade

**Descrição:** Onze jogadores, todos com habilidade 1500. Há três Tops; um deles tem latência 200 ms e os demais 20 ms. Incluir o Top de alta latência não muda equilíbrio nem posições, mas aumenta `amplitude_latencia` em 180 e reduz Q. O sistema deve deixá-lo na fila.

**Entrada:**

```text
p01 | 1500 | Top     | - | 20 |  20 | - | BR
p02 | 1500 | Top     | - | 20 |  20 | - | BR
p03 | 1500 | Top     | - | 20 | 200 | - | BR
p04 | 1500 | Jungle  | - | 20 |  20 | - | BR
p05 | 1500 | Jungle  | - | 20 |  20 | - | BR
p06 | 1500 | Mid     | - | 20 |  20 | - | BR
p07 | 1500 | Mid     | - | 20 |  20 | - | BR
p08 | 1500 | ADC     | - | 20 |  20 | - | BR
p09 | 1500 | ADC     | - | 20 |  20 | - | BR
p10 | 1500 | Support | - | 20 |  20 | - | BR
p11 | 1500 | Support | - | 20 |  20 | - | BR
```

**Cálculo relevante:**

```text
Subconjunto sem p03: amplitude_latencia = 0, Q = 10000
Subconjunto com p03: amplitude_latencia = 180, Q = 9820
```

Entre os dois Tops de 20 ms (`p01` e `p02`), a atribuição canônica usa `p01` na Equipe A.

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p04, Jungle), (p06, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p05, Jungle), (p07, Mid), (p09, ADC), (p11, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 10000
espera: {min: 20, max: 20, media: 20.0}
amplitude_latencia: 0
fila_restante: [p03]
```

---

### N10 — Dupla ADC + Support na mesma equipe

**Descrição:** Uma dupla entra na fila como ADC e Support. O restante da fila cobre as outras posições com habilidade homogênea. A dupla deve permanecer junta. Separá-la seria tão equilibrada quanto, mas violaria R5.

**Entrada:**

```text
p01 | 1510 | ADC     | - | 40 | 18 | G2 | BR
p02 | 1490 | Support | - | 40 | 22 | G2 | BR
p03 | 1500 | Top     | - | 20 | 20 | -  | BR
p04 | 1500 | Top     | - | 20 | 20 | -  | BR
p05 | 1500 | Jungle  | - | 20 | 20 | -  | BR
p06 | 1500 | Jungle  | - | 20 | 20 | -  | BR
p07 | 1500 | Mid     | - | 20 | 20 | -  | BR
p08 | 1500 | Mid     | - | 20 | 20 | -  | BR
p09 | 1500 | ADC     | - | 20 | 20 | -  | BR
p10 | 1500 | Support | - | 20 | 20 | -  | BR
```

**Cálculo relevante:**

```text
Equipe A contém p01, portanto também p02.
somaA = 1510 + 1490 + 1500 + 1500 + 1500 = 7500
somaB = 7500
diferenca_habilidade = 0
amplitude_latencia = 22 - 18 = 4
Q = 10000 - 4 = 9996
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p03, Top), (p05, Jungle), (p07, Mid), (p01, ADC), (p02, Support)
equipe_b: (p04, Top), (p06, Jungle), (p08, Mid), (p09, ADC), (p10, Support)
diferenca_habilidade: 0
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 9996
espera: {min: 20, max: 40, media: 24.0}
amplitude_latencia: 4
fila_restante: []
```

---

## 10. Casos-limite

### L01 — Jogadores insuficientes

**Descrição:** A fila é bem formada, mas tem 9 jogadores. Não é possível satisfazer R1. Nenhuma partida deve ser criada e a fila permanece intacta.

**Entrada:**

```text
p01 | 1500 | Top     | - | 10 | 20 | - | BR
p02 | 1500 | Top     | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 | 1500 | Jungle  | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | Mid     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | ADC     | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
```

**Saída esperada:**

```text
status: NENHUMA_PARTIDA
motivo: JOGADORES_INSUFICIENTES
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09]
```

O mesmo motivo aplica-se a qualquer fila bem formada com 0 a 9 jogadores.

---

### L02 — Concentração extrema em uma posição

**Descrição:** Dez jogadores aceitam exclusivamente Mid. Não há como montar Top, Jungle, ADC e Support sem violar R3. O sistema deve reconhecer a impossibilidade de composição, e não atribuir posições recusadas.

**Entrada:**

```text
p01 | 1500 | Mid | - | 50 | 20 | - | BR
p02 | 1500 | Mid | - | 50 | 20 | - | BR
p03 | 1500 | Mid | - | 50 | 20 | - | BR
p04 | 1500 | Mid | - | 50 | 20 | - | BR
p05 | 1500 | Mid | - | 50 | 20 | - | BR
p06 | 1500 | Mid | - | 50 | 20 | - | BR
p07 | 1500 | Mid | - | 50 | 20 | - | BR
p08 | 1500 | Mid | - | 50 | 20 | - | BR
p09 | 1500 | Mid | - | 50 | 20 | - | BR
p10 | 1500 | Mid | - | 50 | 20 | - | BR
```

**Saída esperada:**

```text
status: NENHUMA_PARTIDA
motivo: COMPOSICAO_IMPOSSIVEL
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09, p10]
```

---

### L03 — Habilidade muito diferente e espera insuficiente

**Descrição:** Mesma distribuição de N04 (um jogador 2800 e nove 1500), porém todos com espera baixa. A diferença 260 ultrapassa o limite 202. Sem relaxamento suficiente, ninguém deve ser agrupado.

**Entrada:**

```text
p01 | 2800 | Top     | - | 10 | 20 | - | BR
p02 | 1500 | Top     | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 | 1500 | Jungle  | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | Mid     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | ADC     | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
p10 | 1500 | Support | - | 10 | 20 | - | BR
```

**Cálculo relevante:**

```text
diferenca_habilidade = 260
limite_diferenca = min(400, 200 + (10 // 4)) = 202
260 > 202 → inválida
Não há outra partição: o 2800 estará em alguma equipe de 5.
```

**Saída esperada:**

```text
status: NENHUMA_PARTIDA
motivo: DIFERENCA_HABILIDADE
fila_restante: [p01, p02, p03, p04, p05, p06, p07, p08, p09, p10]
```

Este caso é o complemento de N04: a espera é o único fator que muda o resultado.

---

### L04 — Diferença exatamente no limite, espera zero

**Descrição:** Fronteira de R6. A melhor formação válida tem diferença exatamente 200 e `max_espera = 0`, portanto `limite_diferenca = 200`. Igualdade deve ser aceita (`<=`, não `<`).

**Entrada:** igual a N02, porém todos com `espera_s = 0`.

```text
p01 | 2000 | Top     | - | 0 | 20 | - | BR
p02 | 1000 | Top     | - | 0 | 20 | - | BR
p03 | 2000 | Jungle  | - | 0 | 20 | - | BR
p04 | 1000 | Jungle  | - | 0 | 20 | - | BR
p05 | 2000 | Mid     | - | 0 | 20 | - | BR
p06 | 1000 | Mid     | - | 0 | 20 | - | BR
p07 | 2000 | ADC     | - | 0 | 20 | - | BR
p08 | 1000 | ADC     | - | 0 | 20 | - | BR
p09 | 2000 | Support | - | 0 | 20 | - | BR
p10 | 1000 | Support | - | 0 | 20 | - | BR
```

**Cálculo relevante:**

```text
limite_diferenca = min(400, 200 + (0 // 4)) = 200
diferença da formação 3–2 = 200 → válida
diferença 4–1 = 600 → inválida
```

**Saída esperada:**

```text
status: PARTIDA_FORMADA
equipe_a: (p01, Top), (p03, Jungle), (p05, Mid), (p08, ADC), (p10, Support)
equipe_b: (p02, Top), (p04, Jungle), (p06, Mid), (p07, ADC), (p09, Support)
diferenca_habilidade: 200
jogadores_fora_preferencia: 0
desequilibrio_alternativas: 0
qualidade: 8000
espera: {min: 0, max: 0, media: 0.0}
amplitude_latencia: 0
fila_restante: []
```

---

## 11. Casos de entrada inválida

Estes casos **não** são “partida impossível”. A entrada viola o modelo e deve ser rejeitada antes da busca.

### I01 — Identificadores duplicados

**Descrição:** Dois registros usam o `id` `p01`. Sem unicidade, a fila não é uma entrada bem formada.

**Entrada:**

```text
p01 | 1500 | Top     | - | 10 | 20 | - | BR
p01 | 1510 | Jungle  | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 | 1500 | Mid     | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | ADC     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | Support | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
p10 | 1500 | Top     | - | 10 | 20 | - | BR
```

**Saída esperada:**

```text
status: ERRO_ENTRADA
codigo: ID_DUPLICADO
mensagem: identificador repetido na fila: p01
```

Não há `fila_restante` nem tentativa de formação.

---

### I02 — Habilidade negativa

**Descrição:** Habilidade é um valor estimado não negativo. `p04` com `-1` torna a entrada inválida, mesmo que os demais jogadores estejam corretos.

**Entrada:**

```text
p01 | 1500 | Top     | - | 10 | 20 | - | BR
p02 | 1500 | Top     | - | 10 | 20 | - | BR
p03 | 1500 | Jungle  | - | 10 | 20 | - | BR
p04 |   -1 | Jungle  | - | 10 | 20 | - | BR
p05 | 1500 | Mid     | - | 10 | 20 | - | BR
p06 | 1500 | Mid     | - | 10 | 20 | - | BR
p07 | 1500 | ADC     | - | 10 | 20 | - | BR
p08 | 1500 | ADC     | - | 10 | 20 | - | BR
p09 | 1500 | Support | - | 10 | 20 | - | BR
p10 | 1500 | Support | - | 10 | 20 | - | BR
```

**Saída esperada:**

```text
status: ERRO_ENTRADA
codigo: HABILIDADE_INVALIDA
mensagem: habilidade deve ser maior ou igual a 0 (jogador p04)
```

---

### I03 — Posição desconhecida

**Descrição:** `p07` declara preferida `Atirador`, que não pertence ao conjunto fechado de posições. O nome próximo de ADC não deve ser aceito: a comparação é exata.

**Entrada:**

```text
p01 | 1500 | Top      | - | 10 | 20 | - | BR
p02 | 1500 | Top      | - | 10 | 20 | - | BR
p03 | 1500 | Jungle   | - | 10 | 20 | - | BR
p04 | 1500 | Jungle   | - | 10 | 20 | - | BR
p05 | 1500 | Mid      | - | 10 | 20 | - | BR
p06 | 1500 | Mid      | - | 10 | 20 | - | BR
p07 | 1500 | Atirador | - | 10 | 20 | - | BR
p08 | 1500 | ADC      | - | 10 | 20 | - | BR
p09 | 1500 | Support  | - | 10 | 20 | - | BR
p10 | 1500 | Support  | - | 10 | 20 | - | BR
```

**Saída esperada:**

```text
status: ERRO_ENTRADA
codigo: POSICAO_DESCONHECIDA
mensagem: posicao nao reconhecida: Atirador (jogador p07)
```

O mesmo código aplica-se se a posição inválida aparecer apenas em `alternativas`.

---

### I04 — Grupo maior do que uma equipe

**Descrição:** Seis jogadores compartilham o grupo `G9`. R2 exige equipes de 5 e R5 exige que o grupo fique junto; as duas regras não podem ser verdadeiras ao mesmo tempo. Essa configuração é estruturalmente inválida e deve ser rejeitada na validação da entrada, não tratada como busca sem resultado.

**Entrada:**

```text
p01 | 1500 | Top     | - | 10 | 20 | G9 | BR
p02 | 1500 | Top     | - | 10 | 20 | G9 | BR
p03 | 1500 | Jungle  | - | 10 | 20 | G9 | BR
p04 | 1500 | Jungle  | - | 10 | 20 | G9 | BR
p05 | 1500 | Mid     | - | 10 | 20 | G9 | BR
p06 | 1500 | Mid     | - | 10 | 20 | G9 | BR
p07 | 1500 | ADC     | - | 10 | 20 | -  | BR
p08 | 1500 | ADC     | - | 10 | 20 | -  | BR
p09 | 1500 | Support | - | 10 | 20 | -  | BR
p10 | 1500 | Support | - | 10 | 20 | -  | BR
```

**Saída esperada:**

```text
status: ERRO_ENTRADA
codigo: GRUPO_MAIOR_QUE_EQUIPE
mensagem: grupo G9 possui 6 jogadores, acima do tamanho maximo de uma equipe (5)
```

---

## 12. Relação com a Etapa 1

| Item da Etapa 1 | Caso que congela o comportamento |
| --- | --- |
| Caso 1 — partida equilibrada | N01 |
| Caso 2 — menor diferença de habilidade | N02 |
| Caso 3 — posições disputadas | N03, N07 |
| Caso 4 — tempo de espera | N04, L03 |
| Caso 5 — grupos | N05, N10, I04 |
| Caso-limite 1 — poucos jogadores | L01 |
| Caso-limite 2 — todos na mesma posição | L02 |
| Caso-limite 3 — habilidades muito diferentes | L03, L04, N04 |
| C3 equilíbrio de alternativas | N08 |
| C4 latência | N09 |
| Fila maior que uma partida | N06 |

---

## 13. Como as próximas implementações devem usar estes casos

1. Ler a fila no formato da Seção 2 (o encapsulamento em cada linguagem é livre).
2. Validar a entrada. Se inválida, produzir `ERRO_ENTRADA` com o código do caso.
3. Se válida, produzir `PARTIDA_FORMADA` ou `NENHUMA_PARTIDA` conforme as Seções 4 a 7.
4. Comparar a saída com a **saída esperada** deste arquivo. Equipes, métricas e fila restante precisam coincidir.

Um programa que encontra *alguma* partida válida, mas não a formação ótima da Seção 6, **não** satisfaz o contrato.
