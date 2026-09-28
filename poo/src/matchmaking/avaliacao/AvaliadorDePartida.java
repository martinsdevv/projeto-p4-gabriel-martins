package matchmaking.avaliacao;

import java.util.List;

import matchmaking.dominio.Fila;
import matchmaking.dominio.Partida;

/**
 * Pergunta as restricoes e, se a forma estiver certa, a politica de tolerancia.
 * Os criterios so entram na qualidade; nao decidem sozinhos se a partida existe.
 */
public final class AvaliadorDePartida {
    private final Fila fila;
    private final PoliticaDeTolerancia politica;
    private final List<CriterioQualidade> criterios;
    private final List<Restricao> restricoes;

    public AvaliadorDePartida(Fila fila, PoliticaDeTolerancia politica,
            List<CriterioQualidade> criterios, List<Restricao> restricoes) {
        this.fila = fila;
        this.politica = politica;
        this.criterios = List.copyOf(criterios);
        this.restricoes = List.copyOf(restricoes);
    }

    public static AvaliadorDePartida doContrato(Fila fila) {
        return new AvaliadorDePartida(
                fila,
                PoliticaDeTolerancia.DO_CONTRATO,
                CriterioQualidade.doContrato(),
                RestricoesDoContrato.estruturais());
    }

    public Avaliacao avaliar(Partida partida) {
        for (Restricao restricao : restricoes) {
            if (!restricao.satisfeita(partida, fila)) {
                return Avaliacao.reprovada(restricao.regra());
            }
        }
        Qualidade qualidade = Qualidade.de(partida, criterios);
        if (!politica.aceita(qualidade)) {
            return Avaliacao.reprovada(Regra.R6, qualidade);
        }
        return Avaliacao.aprovada(qualidade);
    }
}
