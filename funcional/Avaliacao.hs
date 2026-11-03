module Avaliacao where

import Tipos
import Data.Maybe (mapMaybe)

-- R5: verifica se nenhum grupo aparece nas duas equipes ao mesmo tempo.
gruposUnidos :: Formacao -> Bool
gruposUnidos f =
    let gruposA = mapMaybe (jGrupo . fst) (equipeA f)
        gruposB = mapMaybe (jGrupo . fst) (equipeB f)
    in all (\g -> notElem g gruposB) gruposA

somaHabilidade :: Equipe -> Int
somaHabilidade eq = sum (map (jHabilidade . fst) eq)

foraPreferida :: Equipe -> Int
foraPreferida eq = length (filter (\(jog, pos) -> pos /= jPreferida jog) eq)

limiteDiferenca :: Int -> Int
limiteDiferenca maxEspera = min limiteAbsoluto (limiteBase + quot maxEspera bonusDivisor)

calcularMetricas :: Formacao -> Metricas
calcularMetricas f =
    let eqA            = equipeA f
        eqB            = equipeB f
        todosJogadores = map fst (eqA ++ eqB)
        somaA          = somaHabilidade eqA
        somaB          = somaHabilidade eqB
        difSoma        = abs (somaA - somaB)
        difHab         = quot difSoma jogadoresPorEquipe
        foraA          = foraPreferida eqA
        foraB          = foraPreferida eqB
        foraTotal      = foraA + foraB
        desequilibrio  = abs (foraA - foraB)
        esperas        = map jEspera todosJogadores
        espMin         = minimum esperas
        espMax         = maximum esperas
        espSoma        = sum esperas
        latencias      = map jLatencia todosJogadores
        amplitude      = maximum latencias - minimum latencias
        qualidade      = 10000 - 2 * difSoma - 100 * foraTotal - 50 * desequilibrio - amplitude
    in Metricas
        { mDiferenca     = difHab
        , mForaPref      = foraTotal
        , mDesequilibrio = desequilibrio
        , mQualidade     = qualidade
        , mAmplitude     = amplitude
        , mEsperaMin     = espMin
        , mEsperaMax     = espMax
        , mEsperaSoma    = espSoma
        }

respeitaHabilidade :: Metricas -> Bool
respeitaHabilidade m = mDiferenca m <= limiteDiferenca (mEsperaMax m)
