module Saida where

import Tipos
import Validacao (lerFila)
import Busca (buscarPartida)
import Data.List (intercalate)

-- Núcleo puro completo: transforma o texto bruto da fila no Resultado do domínio.
processar :: String -> Resultado
processar = either Erro buscarPartida . lerFila

-- Composição pura de ponta a ponta: String de entrada -> String de saída.
executarPuro :: String -> String
executarPuro = formatarResultado . processar

formatarResultado :: Resultado -> String
formatarResultado (Erro err) = unlines
    [ "status: ERRO_ENTRADA"
    , "codigo: " ++ codigo err
    , "mensagem: " ++ mensagem err
    ]
formatarResultado (Nenhuma mot restantes) = unlines
    [ "status: NENHUMA_PARTIDA"
    , "motivo: " ++ formatarMotivo mot
    , "fila_restante: " ++ formatarIds restantes
    ]
formatarResultado (Formada form met restantes) = unlines
    [ "status: PARTIDA_FORMADA"
    , "equipe_a: " ++ formatarEquipe (equipeA form)
    , "equipe_b: " ++ formatarEquipe (equipeB form)
    , "diferenca_habilidade: " ++ show (mDiferenca met)
    , "jogadores_fora_preferencia: " ++ show (mForaPref met)
    , "desequilibrio_alternativas: " ++ show (mDesequilibrio met)
    , "qualidade: " ++ show (mQualidade met)
    , "espera: " ++ formatarEspera met
    , "amplitude_latencia: " ++ show (mAmplitude met)
    , "fila_restante: " ++ formatarIds restantes
    ]

formatarEspera :: Metricas -> String
formatarEspera met =
    "{min: " ++ show (mEsperaMin met)
    ++ ", max: " ++ show (mEsperaMax met)
    ++ ", media: " ++ formatarMediaDez (mEsperaSoma met) ++ "}"

formatarMotivo :: Motivo -> String
formatarMotivo JogadoresInsuficientes = "JOGADORES_INSUFICIENTES"
formatarMotivo ComposicaoImpossivel   = "COMPOSICAO_IMPOSSIVEL"
formatarMotivo DiferencaHabilidade    = "DIFERENCA_HABILIDADE"
formatarMotivo GrupoIncompativel      = "GRUPO_INCOMPATIVEL"

formatarEquipe :: Equipe -> String
formatarEquipe eq =
    intercalate ", " [ "(" ++ jId jog ++ ", " ++ show pos ++ ")" | (jog, pos) <- eq ]

formatarIds :: [String] -> String
formatarIds ids = "[" ++ intercalate ", " ids ++ "]"

formatarMediaDez :: Int -> String
formatarMediaDez soma =
    show (quot soma jogadoresPorPartida) ++ "." ++ show (rem soma jogadoresPorPartida)

-- Formatação pura do relatório de todos os casos de teste
formatarRelatorio :: [(String, Resultado)] -> String
formatarRelatorio casos =
    let cabecalho = formatarLinhaTabela "caso" "status" "detalhe" "espera"
        linhas    = map (uncurry formatarCasoRelatorio) casos
        partidas  = length [ () | (_, Formada _ _ _) <- casos ]
        nenhuma   = length [ () | (_, Nenhuma _ _)   <- casos ]
        erros     = length [ () | (_, Erro _)        <- casos ]
        rodape    = show (length casos) ++ " casos: "
                    ++ show partidas ++ " partidas, "
                    ++ show nenhuma ++ " sem partida, "
                    ++ show erros ++ " erros de entrada"
    in unlines (cabecalho : linhas ++ ["", rodape])

formatarCasoRelatorio :: String -> Resultado -> String
formatarCasoRelatorio nome (Formada _ met _) =
    formatarLinhaTabela nome "PARTIDA_FORMADA" ("qualidade " ++ show (mQualidade met)) (formatarEspera met)
formatarCasoRelatorio nome (Nenhuma mot _) =
    formatarLinhaTabela nome "NENHUMA_PARTIDA" (formatarMotivo mot) "-"
formatarCasoRelatorio nome (Erro err) =
    formatarLinhaTabela nome "ERRO_ENTRADA" (codigo err) "-"

formatarLinhaTabela :: String -> String -> String -> String -> String
formatarLinhaTabela c1 c2 c3 c4 =
    alinharEsquerda 6 c1 ++ " "
    ++ alinharEsquerda 18 c2 ++ " "
    ++ alinharEsquerda 24 c3 ++ " "
    ++ c4

alinharEsquerda :: Int -> String -> String
alinharEsquerda largura texto =
    texto ++ replicate (max 0 (largura - length texto)) ' '
