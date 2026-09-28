package matchmaking.dominio;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Pedido ja checado para entrar na fila.
 * Ainda nao e um jogador: a fila e quem cria o Jogador e o Grupo.
 */
public final class Inscricao {
    private final String id;
    private final int habilidade;
    private final Posicao preferida;
    private final List<Posicao> alternativas;
    private final int esperaSegundos;
    private final int latenciaMs;
    private final String grupoId;
    private final String regiao;

    public Inscricao(String id, int habilidade, Posicao preferida, List<Posicao> alternativas,
            int esperaSegundos, int latenciaMs, String grupoId, String regiao) {
        this.id = id;
        this.habilidade = habilidade;
        this.preferida = preferida;
        this.alternativas = Collections.unmodifiableList(new ArrayList<Posicao>(alternativas));
        this.esperaSegundos = esperaSegundos;
        this.latenciaMs = latenciaMs;
        this.grupoId = grupoId;
        this.regiao = regiao;
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

    public List<Posicao> alternativas() {
        return alternativas;
    }

    public int esperaSegundos() {
        return esperaSegundos;
    }

    public int latenciaMs() {
        return latenciaMs;
    }

    /** Nulo quando o jogador entrou sozinho. */
    public String grupoId() {
        return grupoId;
    }

    public String regiao() {
        return regiao;
    }
}
