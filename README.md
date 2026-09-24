# P4 - Matchmaking competitivo

**[P4-ETAPA-01] Proposta do problema**  
**[P4-ETAPA-02] Contrato semântico e testes** — ver [`testes/casos.md`](testes/casos.md)  
**[P4-ETAPA-03] Implementação imperativa** — ver [`imperativo/`](imperativo/) e [`imperativo/DECISOES.md`](imperativo/DECISOES.md)

## Etapas

| Tag | Entrega | Onde está |
| --- | --- | --- |
| `[P4-ETAPA-01]` | Proposta e especificação do problema | este README e [`docs/problema.md`](docs/problema.md) |
| `[P4-ETAPA-02]` | Contrato de comportamento e casos de teste | [`testes/casos.md`](testes/casos.md) |
| `[P4-ETAPA-03]` | Solução no paradigma imperativo | [`imperativo/`](imperativo/) |

A Etapa 2 congela **o que** o sistema deve fazer, com casos reutilizáveis nas quatro implementações. A especificação da Etapa 1 permanece como descrição do problema.

## Testes do modelo imperativo

Os 18 casos ficam em [`testes/`](testes/) (`n01.txt` a `n10.txt`, `l01.txt` a `l04.txt`, `i01.txt` a `i04.txt`). Os comandos abaixo são rodados de dentro de `imperativo/`.

No Windows, `go run` grava o executável na pasta temporária e o sistema pode recusar a abertura (`Access is denied`). Quando isso acontecer, use o `.exe` gerado na pasta do projeto.

### Com o Go na máquina

Um caso:

```text
go run . ..\testes\n01.txt
```

Troque o nome do arquivo para ver outro (`l03.txt`, `i01.txt`, etc.).

Todos de uma vez, num relatório:

```text
go run . -relatorio ..\testes
```

A comparação automática com a saída do contrato:

```text
go test .
```

### Com o executável, sem instalar Go

Os dois binários já estão em `imperativo/`. Os comandos são rodados de dentro dessa pasta.

No Windows, o arquivo é [`imperativo.exe`](imperativo/imperativo.exe).

Um caso:

```text
.\imperativo.exe ..\testes\n01.txt
```

Todos de uma vez:

```text
.\imperativo.exe -relatorio ..\testes
```

No Linux, o arquivo é [`imperativo`](imperativo/imperativo). Na primeira vez, libera a execução.

Um caso:

```text
chmod +x imperativo
./imperativo ../testes/n01.txt
```

Todos de uma vez:

```text
./imperativo -relatorio ../testes
```

Quem tiver o Go e quiser gerar de novo:

```text
go build -o imperativo.exe .
```

```text
set GOOS=linux
set GOARCH=amd64
go build -o imperativo .
```

## Sobre o projeto

Este projeto propõe o estudo de um problema de formação de partidas para jogos competitivos, inspirado nos desafios conhecidos do matchmaking de jogos como *League of Legends*.

O objetivo não é reproduzir o algoritmo proprietário de nenhum jogo específico, mas definir um problema próprio, com regras claras, que permita investigar como diferentes paradigmas de programação podem representar e resolver a mesma situação.

Ao longo do projeto, o mesmo problema será reformulado utilizando os paradigmas:

- Programação Imperativa;
- Programação Orientada a Objetos;
- Programação Funcional;
- Programação Lógica.

O projeto segue a proposta da disciplina de manter conceitualmente o mesmo problema durante as diferentes implementações, reformulando a solução de acordo com cada paradigma.

## 1. Descrição do problema

Em um jogo competitivo baseado em partidas entre equipes, jogadores entram em uma fila esperando encontrar uma partida.

Cada jogador possui características relevantes para a formação da partida, como:

- nível de habilidade estimado;
- posições/funções que aceita jogar;
- tempo de espera na fila;
- região;
- latência de conexão;
- possibilidade de estar jogando em grupo com outro jogador.

O sistema de matchmaking deve selecionar jogadores da fila e dividi-los em equipes de forma que a partida seja considerada adequada segundo um conjunto de critérios.

O problema surge porque esses critérios podem entrar em conflito. Uma partida muito equilibrada pode exigir uma espera maior, enquanto uma partida criada rapidamente pode apresentar maior diferença de habilidade ou uma distribuição pior das posições.

Portanto, o problema consiste em encontrar uma formação de partida que represente um bom compromisso entre **qualidade competitiva e tempo de espera**.

## 2. Objetivo

O sistema deverá receber uma fila de jogadores e formar partidas válidas, procurando maximizar a qualidade das partidas segundo critérios previamente definidos.

Para cada possível formação, o sistema deverá ser capaz de verificar se ela atende às regras do problema e avaliar sua qualidade.

A solução deverá permitir comparar diferentes estratégias de formação de partidas posteriormente, mantendo os mesmos dados e critérios.

## 3. Entradas

O sistema deverá receber uma fila contendo jogadores.

Cada jogador deverá possuir, no mínimo:

- identificador;
- habilidade estimada;
- uma ou mais posições/funções preferidas;
- tempo de espera;
- região;
- latência;
- informação sobre eventual grupo/dupla.

Também deverão ser fornecidos os parâmetros utilizados para avaliar as partidas, como:

- quantidade de jogadores por partida;
- quantidade de jogadores por equipe;
- tolerância máxima de diferença de habilidade;
- regras de distribuição de posições;
- tolerância relacionada ao tempo de espera;
- critérios utilizados para avaliar uma formação.



## 4. Saídas

Para cada partida formada, o sistema deverá produzir:

- jogadores selecionados;
- divisão dos jogadores entre as equipes;
- posição/função atribuída a cada jogador;
- diferença estimada de habilidade entre as equipes;
- quantidade de jogadores que não estão em sua posição preferida;
- qualidade estimada da partida;
- informações relevantes sobre o tempo de espera.

Caso não seja possível formar uma partida válida, o sistema deverá indicar que nenhum agrupamento adequado foi encontrado dentro das condições atuais.

## 5. Regras do problema

As regras iniciais do problema serão:

1. Uma partida deverá possuir exatamente 10 jogadores.
2. Cada partida deverá possuir duas equipes com 5 jogadores.
3. Cada jogador deverá ocupar uma posição/função compatível com as posições que aceita jogar.
4. Uma equipe não poderá possuir mais de um jogador ocupando a mesma posição principal, considerando a configuração de posições adotada pelo projeto.
5. A diferença de habilidade entre as equipes deverá ser considerada na avaliação da partida.
6. Uma formação com menor diferença de habilidade será considerada melhor que uma formação com maior diferença, quando os demais critérios forem equivalentes.
7. Jogadores que estejam esperando há mais tempo deverão receber maior tolerância na busca por uma partida, evitando que a tentativa de obter uma formação perfeita provoque espera indefinida.
8. A quantidade de jogadores colocados fora de sua posição preferida deverá ser considerada na avaliação da partida.
9. Quando possível, a quantidade de jogadores fora da posição preferida deverá ser equilibrada entre as equipes.
10. A latência deverá ser considerada como um fator de qualidade da partida, sem necessariamente impedir sua formação.
11. Jogadores que pertençam ao mesmo grupo deverão permanecer na mesma equipe, quando essa configuração for permitida pelo modelo.
12. Uma partida deverá ser avaliada por múltiplos critérios, e não somente pela diferença de habilidade.

As regras poderão ser refinadas durante a especificação posterior, desde que o problema conceitual permaneça o mesmo.

## 6. Casos de exemplo



### Caso 1 - Partida perfeitamente equilibrada

**Entrada:** 10 jogadores com habilidades próximas, todos podendo ocupar posições compatíveis e com baixa latência.

**Saída esperada:** uma partida válida com duas equipes de força semelhante e sem necessidade de alterar posições preferidas.

### Caso 2 - Diferença de habilidade

**Entrada:** 10 jogadores, mas com uma distribuição de habilidade que permite formar equipes com diferenças distintas.

**Saída esperada:** o sistema deverá preferir a formação que minimize a diferença de habilidade, desde que as demais regras sejam respeitadas.

### Caso 3 - Posições disputadas

**Entrada:** muitos jogadores preferem a mesma posição e poucas pessoas aceitam outras posições.

**Saída esperada:** o sistema deverá procurar uma distribuição válida, podendo colocar alguns jogadores em uma posição alternativa aceita por eles.

### Caso 4 - Tempo de espera

**Entrada:** um jogador está esperando há muito mais tempo que os demais e não existe uma formação perfeitamente equilibrada utilizando apenas jogadores com habilidades próximas.

**Saída esperada:** o sistema poderá aceitar uma formação menos equilibrada para evitar que o jogador continue esperando indefinidamente, desde que as regras mínimas sejam respeitadas.

### Caso 5 - Grupos

**Entrada:** existem jogadores que estão em grupo e precisam permanecer na mesma equipe.

**Saída esperada:** o sistema deverá considerar essa relação ao montar as equipes e rejeitar formações que quebrem a restrição do grupo.

## 7. Casos-limite



### Caso-limite 1 - Não há jogadores suficientes

Existem menos de 10 jogadores disponíveis ou jogadores suficientes, mas não existe uma combinação que satisfaça as regras.

**Comportamento esperado:** nenhuma partida deverá ser criada.

### Caso-limite 2 - Concentração extrema em uma posição

Todos ou quase todos os jogadores disponíveis escolhem a mesma posição como preferência e não aceitam posições alternativas.

**Comportamento esperado:** o sistema deverá reconhecer que não é possível formar uma composição válida, em vez de atribuir posições que os jogadores não aceitam.

### Caso-limite 3 - Habilidades muito diferentes

Os únicos jogadores disponíveis apresentam diferenças de habilidade muito grandes.

**Comportamento esperado:** o sistema deverá avaliar se a expansão da tolerância permitida pelo tempo de espera possibilita uma partida. Caso contrário, deverá manter os jogadores na fila.

## 8. Restrições e fora do escopo

O projeto não pretende:

- reproduzir o algoritmo real de matchmaking de *League of Legends*;
- utilizar dados internos ou proprietários de empresas de jogos;
- conectar-se à fila real de nenhum jogo;
- criar um cliente ou servidor de jogo;
- simular a partida depois que as equipes forem formadas;
- implementar comunicação em tempo real entre jogadores;
- desenvolver uma interface gráfica como parte essencial do problema;
- determinar o MMR real de jogadores;
- provar que o algoritmo desenvolvido é superior ao matchmaking real de qualquer jogo.

O foco será o problema abstrato de formação e avaliação de partidas a partir de uma fila de jogadores e de critérios definidos.

## 9. Principais conceitos do domínio

- Jogador
- Fila
- Partida
- Equipe
- Posição/função
- Habilidade estimada
- Tempo de espera
- Latência
- Grupo/dupla
- Formação de partida
- Restrição
- Critério de qualidade
- Equilíbrio entre equipes
- Posição preferida
- Posição alternativa



## 10. Adequação aos quatro paradigmas

O problema foi escolhido por possuir regras, estado, relações e critérios de avaliação que podem ser representados de maneiras significativamente diferentes.

### Programação Imperativa

O problema possui um estado mutável representado pela fila de jogadores. Uma solução imperativa pode manipular essa fila passo a passo, selecionar jogadores, testar formações, modificar o estado e construir partidas.

### Programação Orientada a Objetos

O domínio possui conceitos naturalmente identificáveis, como jogador, equipe, partida, posição e regra de matchmaking. Uma solução orientada a objetos pode organizar o comportamento e os dados em entidades que colaboram para formar e avaliar partidas.

### Programação Funcional

A formação de uma partida pode ser vista como uma sequência de transformações sobre dados: gerar candidatos, filtrar formações inválidas, calcular métricas, ordenar ou comparar resultados e selecionar uma formação. Isso permite explorar funções puras, composição e ausência de estado compartilhado.

### Programação Lógica

O problema possui diversas relações e restrições. Uma solução lógica pode expressar as condições que uma formação precisa satisfazer e buscar atribuições de jogadores às equipes que atendam às regras, permitindo uma abordagem diferente da especificação explícita do procedimento de busca.

## 11. Linguagens inicialmente consideradas

As linguagens são apenas uma escolha inicial e poderão ser revisadas conforme o desenvolvimento.


| Paradigma           | Linguagem considerada | Justificativa inicial                                                                                                          |
| ------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| Imperativo          | Go                    | Favorece um estilo procedural com estado mutável, structs e controle de fluxo explícito, adequado à manipulação passo a passo da fila e à construção das partidas. |
| Orientado a Objetos | Java                  | Possui suporte consolidado ao paradigma orientado a objetos e permite modelar claramente os conceitos do domínio.              |
| Funcional           | Haskell               | É uma linguagem funcional que permite explorar diretamente funções puras, composição, imutabilidade e avaliação de expressões. |
| Lógico              | Prolog                | É adequada à representação de fatos, relações e restrições, permitindo explorar diretamente a natureza lógica do problema.     |


A escolha definitiva das linguagens deverá considerar também as exigências da disciplina e o que for necessário para demonstrar adequadamente cada paradigma.

## Relação com o objetivo do P4

O problema deverá permanecer conceitualmente o mesmo durante o projeto. O que será alterado entre as etapas será a forma de modelar e expressar a solução de acordo com cada paradigma.

A intenção é comparar, ao final, como cada paradigma trata:

- representação dos dados;
- representação do estado;
- controle do fluxo;
- decomposição da solução;
- regras e restrições;
- efeitos colaterais;
- abstração;
- expressão do algoritmo e das relações.

