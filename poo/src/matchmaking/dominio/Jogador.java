package matchmaking.dominio;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

/**
 * Jogador pronto para o matchmaking.
 * Nao existe com habilidade negativa nem com posicao fora do conjunto:
 * esses casos sao recusados antes, como erro de entrada.
 */
public final class Jogador {
    private final String id;
    private final int habilidade;
    private final Posicao preferida;
    private final List<Posicao> aceitas;
    private final int esperaSegundos;
    private final int latenciaMs;
    private final Grupo grupo;
    private final String regiao;

    private Jogador(String id, int habilidade, Posicao preferida, List<Posicao> aceitas,
            int esperaSegundos, int latenciaMs, Grupo grupo, String regiao) {
        if (id == null || id.isEmpty()) {
            throw new IllegalArgumentException("jogador sem identificador");
        }
        if (habilidade < 0 || esperaSegundos < 0 || latenciaMs < 0) {
            throw new IllegalArgumentException("atributo numerico invalido para " + id);
        }
        if (preferida == null || aceitas.isEmpty() || !aceitas.contains(preferida)) {
            throw new IllegalArgumentException("posicoes invalidas para " + id);
        }
        if (regiao == null || regiao.isEmpty()) {
            throw new IllegalArgumentException("regiao ausente para " + id);
        }
        this.id = id;
        this.habilidade = habilidade;
        this.preferida = preferida;
        this.aceitas = List.copyOf(aceitas);
        this.esperaSegundos = esperaSegundos;
        this.latenciaMs = latenciaMs;
        this.grupo = grupo;
        this.regiao = regiao;
    }

    static Jogador aPartirDe(Inscricao inscricao, Grupo grupo) {
        List<Posicao> aceitas = new ArrayList<Posicao>();
        aceitas.add(inscricao.preferida());
        for (Posicao alternativa : inscricao.alternativas()) {
            if (!aceitas.contains(alternativa)) {
                aceitas.add(alternativa);
            }
        }
        return new Jogador(
                inscricao.id(),
                inscricao.habilidade(),
                inscricao.preferida(),
                aceitas,
                inscricao.esperaSegundos(),
                inscricao.latenciaMs(),
                grupo,
                inscricao.regiao());
    }

    public String id() {
        return id;
    }

    public int habilidade() {
        return habilidade;
    }

    public Posicao preferida() {
        return preferida;
    }

    public boolean prefere(Posicao posicao) {
        return preferida == posicao;
    }

    public boolean aceita(Posicao posicao) {
        return aceitas.contains(posicao);
    }

    /** Preferida primeiro, depois as alternativas, sem repetir. */
    public List<Posicao> posicoesAceitas() {
        return aceitas;
    }

    public int esperaSegundos() {
        return esperaSegundos;
    }

    public int latenciaMs() {
        return latenciaMs;
    }

    /**
     * Ausente quando o jogador e solo.
     * Estar em grupo e uma relacao, nao um subtipo de jogador.
     */
    public Optional<Grupo> grupo() {
        return Optional.ofNullable(grupo);
    }

    /** Guardada no modelo. Nesta versao do contrato nenhuma regra a consulta. */
    public String regiao() {
        return regiao;
    }

    @Override
    public boolean equals(Object outro) {
        if (this == outro) {
            return true;
        }
        if (!(outro instanceof Jogador)) {
            return false;
        }
        return id.equals(((Jogador) outro).id);
    }

    @Override
    public int hashCode() {
        return id.hashCode();
    }
}
