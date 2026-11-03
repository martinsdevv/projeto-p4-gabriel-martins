# Como o raciocínio mudou ao passar para o paradigma funcional

**[P4-ETAPA-06]**

O problema continua o mesmo congelado na Etapa 2: dada uma fila de jogadores, formar a partida ótima segundo o contrato semântico ou indicar o motivo pelo qual nenhuma partida existe (ou por que a entrada é inválida). A linguagem utilizada é **Haskell**, conforme previsto na Etapa 1. A saída para os 18 casos de teste (`n01` a `n10`, `l01` a `l04`, `i01` a `i04`) é idêntica à das soluções imperativa (Go) e orientada a objetos (Java).

O que mudou profundamente foi a forma de decompor e raciocinar sobre o problema.

---

## 1. Reflexão obrigatória: de "Que estado devo modificar?" para "Que transformação deve ser aplicada aos dados?"

### Como o problema era pensado antes (Imperativo e OO)

Na versão **imperativa (`imperativo/`)**, o raciocínio girava em torno de **variáveis e estruturas mutáveis ao longo do tempo**:
- Em `entrada.go`, um slice `fila` começava vazio e crescia com `append` a cada iteração do laço `for`. Depois da validação, `completarFila` percorria a fila com ponteiros (`&fila[i]`) para mutar o campo `Aceitas` de cada registro.
- Em `validacao.go`, mapas mutáveis (`vistos := map[string]bool{}` e `contagem := map[string]int{}`) eram alimentados passo a passo para detectar duplicatas e grupos maiores que uma equipe.
- Em `busca.go`, toda a exploração dependia de um registro mutável `estadoBusca` passado por ponteiro (`&estado`). O algoritmo de backtracking (`explorarPosicoes`) escrevia `contagem[posicao]++` e `posDe[indice] = posicao`, chamava a recursão e depois desfazia a mutação (`contagem[posicao]--`, `posDe[indice] = -1`). Duas flags booleanas (`houvePosicao` e `houveGrupo`) eram alteradas como efeito colateral durante a travessia para decidir o motivo de falha caso nenhuma partida fosse encontrada.
- Em `formacao.go`, `calcularMetricas` criava um registro zerado `var m Metricas` e passava ponteiros (`&m`, `&latMin`, `&latMax`) para `acumularEquipe` somar habilidades e atualizar mínimos/máximos.

Na versão **orientada a objetos (`poo/`)**, os dados do domínio passaram a ser encapsulados em objetos imutáveis (`Jogador`, `Equipe`, `Partida`, `Qualidade`), mas a **coordenação da busca** ainda preservava um objeto mutável interno (`FormadorDePartida.Exploracao`), que guardava a melhor `PartidaAvaliada` encontrada até o momento e atualizava as marcas booleanas `houveAtribuicaoDePosicoes` e `houveFormacaoComGruposUnidos` conforme o formador visitava cada combinação.

### Como o raciocínio mudou no paradigma funcional (`funcional/`)

Em Haskell, **não há atribuição destrutiva, variáveis mutáveis, laços `for` nem estado intermediário compartilhado**. A pergunta deixou de ser *"qual variável ou objeto acumulador devo atualizar neste passo?"* e passou a ser *"qual é o tipo de entrada, qual é o tipo de saída e que composição de funções transforma um no outro?"*

O programa inteiro no núcleo puro (`Saida.executarPuro`) tornou-se literalmente uma composição matemática `String -> String`:

$$\text{Texto Bruto} \xrightarrow{\text{lerFila}} \text{Either ErroEntrada [Jogador]} \xrightarrow{\text{either Erro buscarPartida}} \text{Resultado} \xrightarrow{\text{formatarResultado}} \text{Texto de Saída}$$

Cada etapa que antes exigia mutação de estado virou uma transformação de valores imutáveis:

1. **Do preenchimento in-place (`completarFila`) para função pura sob demanda (`aceitas`):**
   Em vez de um campo `Aceitas` que nasce vazio e precisa ser preenchido por efeito colateral depois da leitura, `aceitas :: Jogador -> [Posicao]` (`Tipos.hs`) é uma função pura `nub (jPreferida j : jAlternativas j)`. Sempre que aplicada ao mesmo jogador, devolve a mesma lista (transparência referencial).
2. **De flags de erro (`bool`) e saídas precoces (`if !ok { return }`) para o tipo algébrico `Either`:**
   A leitura e a validação (`Entrada.hs` e `Validacao.hs`) não alteram variáveis de erro nem retornam tuplas `(Jogador, ErroEntrada, bool)`. Elas transformam linhas em `Either ErroEntrada Jogador` e compõem a fila inteira com `traverse lerLinha`, encadeando `validarIDs` e `validarGrupos` pela mônada `Either`. Se qualquer etapa resulta em `Left erro`, o pipeline propaga o erro automaticamente; se todas resultam em `Right`, a fila validada segue adiante.
3. **Do `estadoBusca` mutável e backtracking in-place para pipeline de listas:**
   Em `Busca.hs`, não existe acumulador do "melhor até agora" nem flags `houvePosicao`/`houveGrupo`. A busca é expressa como uma sequência de transformações sobre listas de candidatos:
   - `subconjuntos 10 fila` transforma `[Jogador]` em `[[Jogador]]`;
   - `concatMap atribuirPosicoes` transforma os subconjuntos em atribuições de vagas válidas (`[ParesPorPosicao]`);
   - `concatMap gerarFormacoes` transforma cada atribuição nas 32 partições rotuladas (`[Formacao]`);
   - `filter gruposUnidos` retém apenas as formações que satisfazem R5;
   - `map (\f -> (f, calcularMetricas f))` anexa as métricas a cada formação;
   - `filter (\(_, m) -> respeitaHabilidade m)` retém apenas as que satisfazem R6.

   O motivo de `NENHUMA_PARTIDA` deixou de depender de flags mutadas durante a busca: ele emerge diretamente de verificar qual estágio do pipeline produziu uma lista vazia (`null todasFormacoes`, `null comGruposUnidos` ou `null validas`).
4. **De 60 linhas de `if` de desempate (`supera`) para uma transformação em tupla (`chaveEscolha`):**
   Os 8 critérios de desempate da Seção 6 do contrato foram expressos como uma função `chaveEscolha` que transforma `(Formacao, Metricas)` na tupla `(Down Int, Int, Int, Int, Int, Down Int, [String], [String])`. Como tuplas e listas em Haskell possuem ordenação lexicográfica derivada automaticamente, a escolha da formação ótima reduz-se a `minimumBy (comparing chaveEscolha) validas`.

---

## 2. Identificação dos elementos exigidos pelo enunciado

### 2.1 Dados imutáveis (`Tipos.hs`)

Todos os tipos de dados da solução são algébricos (ADTs) e imutáveis por definição em Haskell:
- **Tipos soma (uniões discriminadas):**
  - `Posicao = Top | Jungle | Mid | ADC | Support`
  - `Motivo = JogadoresInsuficientes | ComposicaoImpossivel | DiferencaHabilidade | GrupoIncompativel`
  - `Resultado = Erro ErroEntrada | Nenhuma Motivo [String] | Formada Formacao Metricas [String]`
  - Uso de `Maybe String` para representar `jGrupo` (`Nothing` para jogador solo, `Just g` para grupo), eliminando o valor sentinela `""` ou `null`.
- **Registros imutáveis (tipos produto):** `Jogador`, `Formacao`, `Metricas` e `ErroEntrada`.

### 2.2 Funções puras vs. Funções impuras e fronteira com o mundo externo

A separação entre núcleo funcional puro e operações de entrada/saída é absoluta e verificável pelo sistema de tipos do Haskell (`IO`):

| Módulo | Natureza | Funções e papel |
| --- | --- | --- |
| `Tipos.hs` | **100% Puro** | Definições de tipos imutáveis, constantes do contrato e `aceitas :: Jogador -> [Posicao]`. |
| `Entrada.hs` | **100% Puro** | `trim`, `separarPor`, `lerPosicao`, `lerInteiro`, `validarPosicao`, `lerAlternativas`, `lerGrupo`, `exigir`, `lerLinha`, `lerLinhas`. |
| `Validacao.hs` | **100% Puro** | `validarIDs`, `validarGrupos`, `lerFila`. |
| `Avaliacao.hs` | **100% Puro** | `gruposUnidos`, `somaHabilidade`, `foraPreferida`, `limiteDiferenca`, `calcularMetricas`, `respeitaHabilidade`. |
| `Busca.hs` | **100% Puro** | `subconjuntos`, `atribuirPosicoes`, `gerarFormacoes`, `rotularCanonicamente`, `chaveEscolha`, `buscarPartida`. |
| `Saida.hs` | **100% Puro** | `processar`, `executarPuro`, `formatarResultado`, `formatarEspera`, `formatarMotivo`, `formatarEquipe`, `formatarIds`, `formatarMediaDez`, `formatarRelatorio`. |
| `Main.hs` | **Impuro (`IO`)** | **Único ponto de interação com o mundo externo:** `main`, `executarTexto`, `executarRelatorio`, `lerCaso`, `falhaUso` (usam `getArgs`, `readFile`, `getContents`, `doesDirectoryExist`, `listDirectory`, `putStr`, `hPutStrLn stderr`, `exitWith`). Até a auxiliar `chaveArquivo` em `Main.hs` é pura. |

### 2.3 Composição funcional (`.`)

Exemplos concretos de composição de funções (`.`) no código:
- `Entrada.trim = dropWhile isSpace . dropWhileEnd isSpace` — compõe a remoção de espaços no final e no início em estilo *point-free*.
- `Avaliacao.gruposUnidos`: `mapMaybe (jGrupo . fst) (equipeA f)` — compõe o seletor de tupla `fst` com o seletor de registro `jGrupo`.
- `Avaliacao.somaHabilidade`: `sum (map (jHabilidade . fst) eq)` — extrai o jogador da vaga e sua habilidade em uma única função composta.
- `Busca.atribuirPosicoes`: `sortOn (length . aceitas) dezJogadores` — ordena jogadores compondo `aceitas` com `length`.
- `Saida.processar = either Erro buscarPartida . lerFila` — compõe a leitura/validação com o despacho para busca ou erro.
- `Saida.executarPuro = formatarResultado . processar` — compõe todo o núcleo do domínio com a formatação textual.

### 2.4 Funções de primeira classe e de alta ordem

Funções são passadas como argumentos, retornadas e combinadas em todos os módulos:
- `map` e `fmap` (`Entrada.hs`, `Avaliacao.hs`, `Busca.hs`, `Saida.hs`) para transformar listas e valores dentro de `Either`;
- `filter` (`Entrada.lerAlternativas`, `Entrada.lerLinhas`, `Validacao.validarGrupos`, `Avaliacao.foraPreferida`, `Busca.buscarPartida`) recebendo predicados e funções anônimas (`\s -> ...`, `\(_, m) -> ...`);
- `traverse` (`Entrada.lerAlternativas`, `Entrada.lerLinhas`, `Main.executarRelatorio`) para aplicar uma função com efeito (`Either` ou `IO`) a cada item de uma lista e inverter a estrutura (`[Either e a] -> Either e [a]`);
- `all` (`Avaliacao.gruposUnidos`) para verificar universalidade sobre os grupos da equipe;
- `concatMap` (`Busca.buscarPartida`) para encadear geradores de listas;
- `sortOn` e `minimumBy (comparing chaveEscolha)` (`Busca.hs`), onde `comparing :: Ord a => (b -> a) -> b -> b -> Ordering` é uma função de alta ordem que recebe `chaveEscolha` e devolve um comparador passado para `minimumBy`;
- `uncurry` (`Busca.gerarFormacoes`, `Saida.formatarRelatorio`) que transforma uma função curried `a -> b -> c` em uma função sobre pares `(a, b) -> c`;
- `either` (`Saida.processar`) que recebe duas funções (`Erro` e `buscarPartida`) para eliminar o tipo `Either`.

### 2.5 Recursão

Onde laços imperativos com índices mutáveis eram usados em Go e Java, a versão em Haskell utiliza recursão estrutural sobre listas:
- `Entrada.separarPor`: divide uma string em torno de um delimitador `,` ou `|` usando `break (== c)` e recursão na cauda restante.
- `Validacao.validarIDs` (`verificar`): percorre recursivamente a lista de jogadores acumulando os IDs já vistos em uma lista imutável `jId j : vistos`.
- `Busca.subconjuntos`: gera todas as combinações de $k$ elementos de uma lista pela recorrência $\binom{n}{k} = \binom{n-1}{k-1} + \binom{n-1}{k}$.
- `Busca.atribuirPosicoes` (`buscar`): explora recursivamente a árvore de atribuição de posições com compreensão de listas, passando `((j, pos) : escolhidas)` para o próximo nível sem mutar o nível atual.
- `Busca.gerarFormacoes` (`dividir`): gera recursivamente as $2^5 = 32$ partições entre Equipe A e Equipe B.

---

## 3. Execução

Os comandos abaixo são executados de dentro de `funcional/`:

No Windows:
```text
.\rodar.bat ..\testes\n01.txt
.\rodar.bat -relatorio ..\testes
```

No Linux:
```text
chmod +x rodar.sh
./rodar.sh ../testes/n01.txt
./rodar.sh -relatorio ../testes
```

