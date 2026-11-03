module Main where

import Tipos (Resultado(..))
import Saida (processar, formatarResultado, formatarRelatorio)
import Data.List (isSuffixOf, sortOn)
import System.Directory (doesDirectoryExist, listDirectory)
import System.Environment (getArgs)
import System.Exit (ExitCode(..), exitSuccess, exitWith)
import System.IO (hPutStrLn, stderr)

-- Fronteira com o mundo externo (IO): lê argumentos/arquivos e delega todo o
-- processamento e formatação para as funções puras de Saida, Validacao e Busca.
main :: IO ()
main = do
    args <- getArgs
    case args of
        ["-relatorio"]        -> executarRelatorio "../testes"
        ["-relatorio", pasta] -> executarRelatorio pasta
        ("-relatorio" : _)    -> falhaUso "uso: funcional -relatorio [pasta]"
        [caminho] -> do
            ehPasta <- doesDirectoryExist caminho
            if ehPasta
                then executarRelatorio caminho
                else readFile caminho >>= executarTexto
        [] -> getContents >>= executarTexto
        _  -> falhaUso "uso: funcional [arquivo | -relatorio pasta]"

executarTexto :: String -> IO ()
executarTexto texto = do
    let resultado = processar texto
    putStr (formatarResultado resultado)
    case resultado of
        Erro _ -> exitWith (ExitFailure 1)
        _      -> exitSuccess

executarRelatorio :: FilePath -> IO ()
executarRelatorio pasta = do
    existe <- doesDirectoryExist pasta
    if not existe
        then falhaUso ("nenhum caso em " ++ pasta)
        else do
            arquivos <- listDirectory pasta
            let txts = sortOn chaveArquivo (filter (isSuffixOf ".txt") arquivos)
            if null txts
                then falhaUso ("nenhum caso em " ++ pasta)
                else do
                    casos <- traverse (lerCaso pasta) txts
                    putStr (formatarRelatorio casos)

lerCaso :: FilePath -> FilePath -> IO (String, Resultado)
lerCaso pasta arquivo = do
    conteudo <- readFile (pasta ++ "/" ++ arquivo)
    let nome = take (length arquivo - 4) arquivo
    return (nome, processar conteudo)

chaveArquivo :: String -> (Int, String)
chaveArquivo nome = (faixa nome, nome)
  where
    faixa ('n' : _) = 0
    faixa ('l' : _) = 1
    faixa ('i' : _) = 2
    faixa _         = 3

falhaUso :: String -> IO a
falhaUso msg = do
    hPutStrLn stderr msg
    exitWith (ExitFailure 2)

