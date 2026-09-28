# Como o modelo mudou ao passar do paradigma imperativo para o orientado a objetos

**[P4-ETAPA-04]**

O problema continua o mesmo da Etapa 2: dada uma fila, formar a partida ótima do contrato, ou explicar por que nenhuma existe. A linguagem é Java, a escolha da Etapa 1 para este paradigma. A saída dos 18 casos é a mesma da implementação em Go. O que mudou foi o modelo.

Na versão imperativa, um jogador é um registro. A fila é um slice que cresce com `append`. Mapas locais marcam ids e grupos, são preenchidos no laço e somem quando a função retorna. `Aceitas` nasce vazio e é reescrito depois da validação. A busca guarda o campeão em `estadoBusca` e o substitui por atribuição. As regras são funções que recebem esses registros e devolvem um código.

Aqui o mesmo domínio virou objetos que colaboram. O texto lido ainda não é um jogador. O jogador só passa a existir quando a entrada já é válida, e a partida não é um arranjo de ponteiros que uma função percorre: é um objeto composto por duas equipes, cada uma composta por cinco vagas. Quem decide se a formação vale, e quem decide se ela é melhor que outra, são esses objetos. O formador só explora.

A exploração continua sendo um procedimento. O contrato exige a formação ótima, e a fila dos testes cabe numa busca exaustiva. Isso não voltou a ser o modelo imperativo: o laço não calcula Q, não aplica R5 e não compara a tupla. Ele pergunta ao avaliador se a partida vale e à partida avaliada se ela supera a campeã.

## Representação do estado

No imperativo, o estado útil está espalhado em variáveis que mudam. A fila muda de tamanho na leitura. O campo `Aceitas` muda no registro que o chamador já tem. Os mapas de validação são estado temporário da função. `estadoBusca` começa com `achou` falso e vai recebendo formação, métricas, ids e tupla. `houvePosicao` e `houveGrupo` são marcas deixadas no caminho para escolher o motivo quando nada presta.

No modelo orientado a objetos, a fila, o jogador, o grupo, a equipe, a partida e a qualidade não são reescritos depois de criados. A ordem de chegada fica na `Fila`. `idsExceto` lê essa ordem para montar `fila_restante`; ninguém reordena a fila para isso.

O que ainda muda durante o processo é o estado da busca, e ele ficou escondido. `FormadorDePartida.Exploracao` guarda a melhor `PartidaAvaliada` e duas marcas: se alguma atribuição de posições existiu, e se alguma formação passou dos grupos. Quando `formar` termina, esse objeto some. Quem chamou recebe um `Resultado`, não as marcas.

Há uma fronteira que o registro único do Go não tinha. `LeitorDeFila` acumula um `Registro` privado, que ainda é texto interpretado: posição é string, grupo é string ou nulo, habilidade pode ser negativa. As regras de entrada olham essa lista. Só depois `Inscricao` vira `Jogador`, e o construtor do jogador recusa habilidade, espera ou latência negativas e posição fora do conjunto. Um jogador inválido não é um objeto pela metade. É um erro de entrada, e a busca nem começa.

O grupo também deixou de ser um mapa refeito a cada checagem. Jogadores do mesmo grupo apontam para o mesmo objeto `Grupo`. A restrição compara essa identidade. Um jogador solo não tem grupo: a ausência é o caso solo, não um valor `"-"` guardado para o laço seguinte tratar.

## Responsabilidades

Cada conceito responde pelo que é dele.

O `Jogador` sabe se aceita uma posição e se aquela posição é a preferida. A lista aceita já nasce com a preferida na frente e sem repetir alternativa. Ninguém de fora preenche esse campo depois.

O `Grupo` é o conjunto de quem entrou junto. O tamanho máximo é conferido na entrada, antes do grupo existir. Na formação, a pergunta é outra: os membros selecionados caíram na mesma equipe?

A `Fila` é o agregado. Ela guarda a ordem de chegada, diz se um jogador está nela e devolve os ids de quem ficou de fora.

A `Vaga` é um jogador numa função, e sabe se está deslocado. A `Equipe` soma habilidade, conta deslocados e se apresenta na ordem Top, Jungle, Mid, ADC, Support. A `Partida` junta as duas equipes, rotula A e B pelo menor id dos dez e expõe as medidas: diferença das somas, jogadores fora da preferida, desequilíbrio e amplitude de latência.

`CriterioQualidade` transforma medida em penalidade. `Qualidade` soma as penalidades, guarda Q e os campos de desempate, e sabe se uma qualidade é preferível a outra nos seis primeiros critérios do contrato. `PartidaAvaliada` completa a ordem total com a lista de ids e a tupla canônica. A comparação não é mais uma função solta recebendo slices.

`PoliticaDeTolerancia` é R6. A espera máxima amplia o limite e não entra em Q. Por isso ela não é um critério de qualidade.

Cada `Restricao` é uma regra estrutural, de R1 a R5. O `AvaliadorDePartida` percorre essa lista e, se a forma estiver certa, pergunta à política se a diferença cabe. O `FormadorDePartida` gera subconjuntos, atribui posições e reparte as equipes. Ele não contém a fórmula.

Os três resultados implementam `Resultado` e cada um escreve o próprio texto. `Main` não escolhe o formato com um `if` em cima de um status. No imperativo, o relatório relê o texto que acabou de gravar para achar `qualidade` ou `motivo`. Aqui o relatório pergunta `status`, `detalhe` e `textoEspera` ao objeto.

## Relacionamento entre componentes

A partida é composta por duas equipes. Sem as equipes não há partida, e uma equipe não é compartilhada entre partidas. A equipe é composta por cinco vagas. A vaga não circula sozinha: ela nasce dentro da equipe.

A fila agrega jogadores e grupos. Ela é a dona dos dois. O jogador agrega o grupo, e vários jogadores apontam para o mesmo objeto. O grupo não é parte exclusiva de um jogador, então não faz sentido embuti-lo como campos copiados nem transformá-lo em subclasse.

A vaga referencia o jogador que continua na fila. Selecionar alguém para uma partida não o tira do objeto fila; a fila restante é calculada por quem não entrou na partida vencedora.

O avaliador é composto pela política, pela lista de critérios e pela lista de restrições. O formador usa um avaliador. Ele não herda o avaliador, porque buscar não é um tipo de avaliar. São tarefas diferentes, ligadas por colaboração.

Equipe A e Equipe B são o mesmo tipo. `Partida.rotulada` troca as referências quando o menor id está do outro lado. Duas subclasses `EquipeA` e `EquipeB` teriam os mesmos métodos de soma, deslocados e texto. A diferença entre elas é o lugar que ocupam na partida, não o comportamento.

## Reutilização

`Posicao` é o conjunto fechado, na ordem canônica. A leitura recusa nome fora dele. A equipe exige uma vaga de cada. A saída imprime nessa ordem. No Go, isso era um slice global de strings e um laço `posicaoValida` repetido por quem precisasse.

`Equipe.texto` serve para os dois lados. `CriterioQualidade.penalidade` é o mesmo cálculo para as quatro medidas; o que muda é `medir`. `PartidaAvaliada.supera` é o único lugar da ordem total da Seção 6. A `Aplicacao` que resolve um arquivo é a mesma que o relatório chama para cada caso.

`PoliticaDeTolerancia.DO_CONTRATO` concentra 200, 4 e 400. Outra política pode ser passada ao construtor do avaliador sem copiar a busca nem a fórmula de Q.

## Encapsulamento

No imperativo, o pacote `main` enxerga todos os campos. Qualquer função pode ler `Regiao`, reescrever `Aceitas` ou mexer nas métricas pelo ponteiro que `acumularEquipe` recebe. O estado da busca também é um registro aberto para as funções que o recebem.

Aqui o campo não é a interface. O jogador expõe `aceita` e `prefere`, não uma lista para o chamador completar. `Grupo.adicionar` e `fechar` ficam no pacote do domínio: depois que a fila fecha o grupo, a lista de membros não cresce. `Exploracao` é classe privada do formador. O registro lido do arquivo é classe privada do leitor.

A região continua no jogador, porque o contrato a exige, e nenhum critério a consulta. Isso é de propósito. No registro público, qualquer laço novo podia passar a usá-la sem que o modelo dissesse que uma regra nova tinha entrado. Aqui, usá-la significa um critério ou uma restrição que peça `regiao()` com esse nome.

A qualidade é montada de uma vez em `Qualidade.de` e depois só é lida. Não há um acumulador zerado que funções vão preenchendo por atribuição.

## Extensão do sistema

Um critério novo de qualidade é uma subclasse de `CriterioQualidade` que implementa `medir`, registrada em `doContrato`. O laço que soma penalidades não muda, e o peso continua na superclasse. Foi assim que C1 a C4 entraram. A espera não entrou por aí: ela não aumenta Q, só relaxa R6, então vive na política.

Uma restrição estrutural nova implementa `Restricao` e entra na lista de `RestricoesDoContrato`. O formador não ganha mais um `if`. Uma tolerância diferente é outro objeto `PoliticaDeTolerancia`, com outra base, outro divisor ou outro teto.

Há um limite honesto. A busca assume duas equipes: para cada função, um bit escolhe qual dos dois jogadores vai para um lado, e `rotulada` define quem é A. Uma partida de três equipes não sai de uma subclasse. Ela exigiria outro formador. A composição tradicional também está no enum `Posicao`. Trocar o conjunto de funções muda o enum, não um parâmetro escondido. O contrato desta etapa é esse jogo, com essas cinco funções e essas duas equipes. O que o modelo deixa aberto é regra, peso e tolerância. O formato da partida, não.

## Herança

A herança que ficou é a de `CriterioQualidade`. As quatro subclasses definem a medida. O método `penalidade` é final e faz `peso * medir`. Esse é o trecho comum de C1 a C4: coisas diferentes, cobradas do mesmo jeito. Estender Q é escrever outra medida, não copiar a conta.

Não usei herança em três lugares onde ela só cumpriria o enunciado.

Jogador solo e jogador de grupo não são subclasses. Entrar num grupo não muda o tipo do jogador. Muda a ligação com um `Grupo`. Vários jogadores compartilham o mesmo grupo, e um grupo de uma pessoa não impõe equipe. Isso é agregação.

Equipe A e Equipe B não são subclasses. O comportamento é o mesmo. A partida compõe duas equipes e a rotulação canônica só escolhe qual referência fica em A.

`PartidaFormada`, `NenhumaPartida` e `ErroEntrada` não herdam um resultado abstrato. Não compartilham campos nem um algoritmo. Compartilham o contrato `Resultado`: texto, código de saída, status, detalhe e espera. `Main` fala com a interface. Uma superclasse só para os três herdarem `texto()` vazio não diria nada sobre o domínio.

## Execução

Os comandos são rodados de dentro de `poo/`. O script compila e executa.

No Windows:

```text
.\rodar.bat ..\testes\n01.txt
.\rodar.bat -relatorio ..\testes
```

No Linux, na primeira vez, libera a execução:

```text
chmod +x rodar.sh
./rodar.sh ../testes/n01.txt
./rodar.sh -relatorio ../testes
```

Entrada inválida termina com código 1. Partida formada ou nenhuma partida termina com código 0. Falha de uso ou de leitura termina com código 2.
