module Entrada where

import Tipos
import Data.Char (isSpace)
import Data.List (dropWhileEnd, nub)
import Text.Read (readMaybe)

trim :: String -> String
trim = dropWhile isSpace . dropWhileEnd isSpace

separarPor :: Char -> String -> [String]
separarPor c texto = 
    case break (== c) texto of 
        (antes, []) -> [antes]
        (antes, _ : resto) -> antes : separarPor c resto

lerPosicao :: String -> Maybe Posicao
lerPosicao "Top" = Just Top
lerPosicao "Jungle" = Just Jungle
lerPosicao "Mid" = Just Mid
lerPosicao "ADC" = Just ADC
lerPosicao "Support" = Just Support
lerPosicao _ = Nothing

lerInteiro :: Int -> String -> String -> Either ErroEntrada Int
lerInteiro numLinha idJog texto =
    case readMaybe texto of
        Just n -> Right n
        Nothing -> Left (ErroEntrada "CAMPO_OBRIGATORIO"
            ("linha " ++ show numLinha ++ ": campo numerico ilegivel (jogador " ++ idJog ++ ")"))
    
validarPosicao :: String -> String -> Either ErroEntrada Posicao
validarPosicao idJog nome =
    case lerPosicao nome of
        Just pos -> Right pos
        Nothing -> Left (ErroEntrada "POSICAO_DESCONHECIDA"
            ("posicao nao reconhecida: " ++ nome ++ " (jogador " ++ idJog ++ ")"))

lerAlternativas :: String -> String -> Either ErroEntrada [Posicao]
lerAlternativas idJog campo
    | campo == "" || campo == "-" = Right []
    | otherwise =
        let nomes = filter (\s -> s /= "" && s /= "-") (map trim (separarPor ',' campo))
        in fmap nub (traverse (validarPosicao idJog) nomes)

lerGrupo :: String -> Maybe String
lerGrupo campo 
    | campo == "" || campo == "-" = Nothing
    | otherwise = Just campo

exigir :: Bool -> ErroEntrada -> Either ErroEntrada ()
exigir True _ = Right ()
exigir False erro = Left erro

lerLinha :: (Int, String) -> Either ErroEntrada Jogador
lerLinha (num, linha) = 
    case map trim (separarPor '|' linha) of
        [idJog, sHab, sPref, sAlt, sEsp, sLat, sGrp, regiao] -> do
            exigir (idJog /= "") (ErroEntrada "CAMPO_OBRIGATORIO"
                ("linha " ++ show num ++ ": identificador ausente"))
            hab <- lerInteiro num idJog sHab
            esp <- lerInteiro num idJog sEsp
            lat <- lerInteiro num idJog sLat
            exigir (regiao /= "") (ErroEntrada "CAMPO_OBRIGATORIO"
                ("linha " ++ show num ++ ": regiao ausente (jogador " ++ idJog ++ ")"))
            exigir (hab >= 0) (ErroEntrada "HABILIDADE_INVALIDA"
                ("habilidade deve ser maior ou igual a 0 (jogador " ++ idJog ++ ")"))
            exigir (esp >= 0) (ErroEntrada "ESPERA_INVALIDA"
                ("espera_s deve ser maior ou igual a 0 (jogador " ++ idJog ++ ")"))
            exigir (lat >= 0) (ErroEntrada "LATENCIA_INVALIDA"
                ("latencia_ms deve ser maior ou igual a 0 (jogador " ++ idJog ++ ")"))
            exigir (sPref /= "") (ErroEntrada "PREFERIDA_AUSENTE"
                ("posicao preferida ausente (jogador " ++ idJog ++ ")"))
            
            pref <- validarPosicao idJog sPref
            alts <- lerAlternativas idJog sAlt
            Right (Jogador idJog hab pref alts esp lat (lerGrupo sGrp) regiao)
        _ -> Left (ErroEntrada "CAMPO_OBRIGATORIO"
            ("linha " ++ show num ++ ": campo obrigatorio ausente"))

lerLinhas :: String -> Either ErroEntrada [Jogador]
lerLinhas texto = 
    let linhasUteis = filter (/= "") (map trim (lines texto))
        numeradas = zip [1..] linhasUteis
    in traverse lerLinha numeradas