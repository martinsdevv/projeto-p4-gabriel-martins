package matchmaking.dominio;

/** Um jogador colocado numa funcao. A vaga nao existe fora da equipe. */
public final class Vaga {
    private final Jogador jogador;
    private final Posicao posicao;

    public Vaga(Jogador jogador, Posicao posicao) {
        if (jogador == null || posicao == null) {
            throw new IllegalArgumentException("vaga incompleta");
        }
        this.jogador = jogador;
        this.posicao = posicao;
    }

    public Jogador jogador() {
        return jogador;
    }

    public Posicao posicao() {
        return posicao;
    }

    public boolean deslocada() {
        return !jogador.prefere(posicao);
    }
}
