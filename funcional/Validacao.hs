module Validacao where

import Tipos
import Entrada (lerLinhas)
import Data.List (nub)
import Data.Maybe (mapMaybe)

validarIDs :: [Jogador] -> Either ErroEntrada ()
validarIDs fila = verificar [] fila
  where
    verificar :: [String] -> [Jogador] -> Either ErroEntrada ()
    verificar _ [] = Right ()
    verificar vistos (j : resto)
        | elem (jId j) vistos = Left (ErroEntrada "ID_DUPLICADO"
            ("identificador repetido na fila: " ++ jId j))
        | otherwise = verificar (jId j : vistos) resto

validarGrupos :: [Jogador] -> Either ErroEntrada ()
validarGrupos fila =
    let todosGrupos = mapMaybe jGrupo fila
        gruposUnicos = nub todosGrupos
        contar g = (g, length (filter (== g) todosGrupos))
        acimaDoLimite = filter (\(_, qtd) -> qtd > jogadoresPorEquipe) (map contar gruposUnicos)
    in case acimaDoLimite of
        ((g, qtd) : _) -> Left (ErroEntrada "GRUPO_MAIOR_QUE_EQUIPE"
            ("grupo " ++ g ++ " possui " ++ show qtd ++ " jogadores, acima do tamanho maximo de uma equipe (" ++ show jogadoresPorEquipe ++ ")"))
        [] -> Right ()

lerFila :: String -> Either ErroEntrada [Jogador]
lerFila texto = do
    fila <- lerLinhas texto
    validarIDs fila
    validarGrupos fila
    Right fila