module Tipos where

import Data.List (nub)

jogadoresPorPartida, jogadoresPorEquipe :: Int
jogadoresPorPartida = 10
jogadoresPorEquipe = 5

limiteBase, bonusDivisor, limiteAbsoluto :: Int
limiteBase = 200
bonusDivisor = 4
limiteAbsoluto = 400

-- `data` cria um tipo novo. Cada nome após o `=` ou do `|` é um contrutor, que é um valor possível
-- `deriving` pede ao compilador implementações automáticas: 
-- * `show` permite converter para texto (show ADC == "ADC")
-- * `Ord` segue a ordem da declaração, ou seja, Top < Jungle < ...
-- * `Enum` e `Bounded` permitem escrever `[minBound .. maxBound]`
data Posicao 
    = Top 
    | Jungle 
    | Mid 
    | ADC 
    | Support 
    deriving (Show, Eq, Ord, Enum, Bounded)

data Motivo 
    = JogadoresInsuficientes 
    | ComposicaoImpossivel 
    | DiferencaHabilidade 
    | GrupoIncompativel 
    deriving (Show, Eq)

-- a linha com `::` é a assinatura do tipo ("todasPosicoes é uma lista de Posicao")
todasPosicoes :: [Posicao]
todasPosicoes = [minBound .. maxBound]

data Jogador = Jogador {
    jId :: String,
    jHabilidade :: Int,
    jPreferida :: Posicao,
    jAlternativas :: [Posicao],
    jEspera :: Int,
    jLatencia :: Int,
    -- `Maybe String` é porque a ausência faz parte do tipo. `Nothing` ou `Just "G1"`
    jGrupo :: Maybe String,
    jRegiao :: String
} deriving (Show, Eq)

type Vaga = (Jogador, Posicao)
type Equipe = [Vaga]

data Formacao = Formacao {
    equipeA :: Equipe,
    equipeB :: Equipe
} deriving Show

data Metricas = Metricas {
    mDiferenca :: Int,
    mForaPref :: Int,
    mDesequilibrio :: Int,
    mQualidade :: Int,
    mAmplitude :: Int,
    mEsperaMin :: Int,
    mEsperaMax :: Int,
    mEsperaSoma :: Int
} deriving Show

data ErroEntrada = ErroEntrada {
    codigo :: String,
    mensagem :: String
} deriving Show

data Resultado
    = Erro ErroEntrada
    | Nenhuma Motivo [String] --motivo, fila_restante
    | Formada Formacao Metricas [String] --formação, métricas, fila_restante
    deriving Show

aceitas :: Jogador -> [Posicao]
aceitas j = nub(jPreferida j : jAlternativas j)