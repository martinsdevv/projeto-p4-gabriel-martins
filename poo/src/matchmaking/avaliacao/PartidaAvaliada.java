package matchmaking.avaliacao;

import java.util.List;

import matchmaking.dominio.Partida;

/** Partida que ja passou no avaliador, pronta para disputar o lugar de melhor. */
public final class PartidaAvaliada {
    private final Partida partida;
    private final Qualidade qualidade;

    public PartidaAvaliada(Partida partida, Qualidade qualidade) {
        this.partida = partida;
        this.qualidade = qualidade;
    }

    public Partida partida() {
        return partida;
    }

    public Qualidade qualidade() {
        return qualidade;
    }

    public boolean supera(PartidaAvaliada outra) {
        int preferencia = qualidade.preferenciaSobre(outra.qualidade);
        if (preferencia != 0) {
            return preferencia > 0;
        }
        int ids = comparar(partida.idsOrdenados(), outra.partida.idsOrdenados());
        if (ids != 0) {
            return ids < 0;
        }
        return comparar(partida.tuplaCanonica(), outra.partida.tuplaCanonica()) < 0;
    }

    private static int comparar(List<String> atual, List<String> outra) {
        int n = Math.min(atual.size(), outra.size());
        for (int i = 0; i < n; i++) {
            int diferenca = atual.get(i).compareTo(outra.get(i));
            if (diferenca != 0) {
                return diferenca;
            }
        }
        return Integer.compare(atual.size(), outra.size());
    }
}
