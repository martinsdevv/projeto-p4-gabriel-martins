module Busca where

import Tipos
import Avaliacao
import Data.List (minimumBy, sort, sortOn)
import Data.Ord (Down(..), comparing)

subconjuntos :: Int -> [a] -> [[a]]
subconjuntos 0 _        = [[]]
subconjuntos _ []       = []
subconjuntos k (x : xs) = map (x :) (subconjuntos (k - 1) xs) ++ subconjuntos k xs

type ParesPorPosicao = [(Posicao, (Jogador, Jogador))]

atribuirPosicoes :: [Jogador] -> [ParesPorPosicao]
atribuirPosicoes dezJogadores =
    let ordenados = sortOn (length . aceitas) dezJogadores
    in map agruparPorPosicao (buscar [] ordenados)
  where
    buscar :: [Vaga] -> [Jogador] -> [[Vaga]]
    buscar escolhidas [] = [escolhidas]
    buscar escolhidas (j : resto) =
        [ resultado
        | pos <- aceitas j
        , vagasNaPosicao pos escolhidas < 2
        , resultado <- buscar ((j, pos) : escolhidas) resto
        ]

    vagasNaPosicao :: Posicao -> [Vaga] -> Int
    vagasNaPosicao pos vagas = length (filter (\(_, p) -> p == pos) vagas)

    agruparPorPosicao :: [Vaga] -> ParesPorPosicao
    agruparPorPosicao vagas =
        [ (pos, extrairPar pos vagas) | pos <- todasPosicoes ]

    extrairPar :: Posicao -> [Vaga] -> (Jogador, Jogador)
    extrairPar pos vagas =
        case [ jog | (jog, p) <- vagas, p == pos ] of
            [j1, j2] -> (j1, j2)
            _        -> error "impossivel: cada posicao tem exatamente 2 jogadores"

gerarFormacoes :: ParesPorPosicao -> [Formacao]
gerarFormacoes pares = map (uncurry rotularCanonicamente) (dividir pares)
  where
    dividir :: ParesPorPosicao -> [(Equipe, Equipe)]
    dividir [] = [([], [])]
    dividir ((pos, (j1, j2)) : resto) =
        [ ((jA, pos) : eqA, (jB, pos) : eqB)
        | (eqA, eqB) <- dividir resto
        , (jA, jB)   <- [(j1, j2), (j2, j1)]
        ]

rotularCanonicamente :: Equipe -> Equipe -> Formacao
rotularCanonicamente eq1 eq2
    | minimum (map (jId . fst) eq1) < minimum (map (jId . fst) eq2) = Formacao eq1 eq2
    | otherwise                                                     = Formacao eq2 eq1

chaveEscolha :: (Formacao, Metricas) -> (Down Int, Int, Int, Int, Int, Down Int, [String], [String])
chaveEscolha (f, m) =
    let idsCanonicos = map (jId . fst) (equipeA f ++ equipeB f)
        idsOrdenados = sort idsCanonicos
    in ( Down (mQualidade m)
       , mDiferenca m
       , mForaPref m
       , mDesequilibrio m
       , mAmplitude m
       , Down (mEsperaSoma m)
       , idsOrdenados
       , idsCanonicos
       )

buscarPartida :: [Jogador] -> Resultado
buscarPartida fila
    | length fila < jogadoresPorPartida =
        Nenhuma JogadoresInsuficientes (map jId fila)
    | null todasFormacoes =
        Nenhuma ComposicaoImpossivel (map jId fila)
    | null comGruposUnidos =
        Nenhuma GrupoIncompativel (map jId fila)
    | null validas =
        Nenhuma DiferencaHabilidade (map jId fila)
    | otherwise =
        let (melhorF, melhorM) = minimumBy (comparing chaveEscolha) validas
            idsEscolhidos      = map (jId . fst) (equipeA melhorF ++ equipeB melhorF)
            restantes          = [ jId j | j <- fila, notElem (jId j) idsEscolhidos ]
        in Formada melhorF melhorM restantes
  where
    todasAtribuicoes = concatMap atribuirPosicoes (subconjuntos jogadoresPorPartida fila)
    todasFormacoes   = concatMap gerarFormacoes todasAtribuicoes
    comGruposUnidos  = filter gruposUnidos todasFormacoes
    avaliadas        = map (\f -> (f, calcularMetricas f)) comGruposUnidos
    validas          = filter (\(_, m) -> respeitaHabilidade m) avaliadas
