package matchmaking.avaliacao;

public final class Avaliacao {
    private final Regra violada;
    private final Qualidade qualidade;

    private Avaliacao(Regra violada, Qualidade qualidade) {
        this.violada = violada;
        this.qualidade = qualidade;
    }

    public static Avaliacao reprovada(Regra regra) {
        return new Avaliacao(regra, null);
    }

    public static Avaliacao reprovada(Regra regra, Qualidade qualidade) {
        return new Avaliacao(regra, qualidade);
    }

    public static Avaliacao aprovada(Qualidade qualidade) {
        return new Avaliacao(null, qualidade);
    }

    public boolean aprovada() {
        return violada == null;
    }

    /** Passou de R5: ou a partida vale, ou so a tolerancia de habilidade barrou. */
    public boolean passouDosGrupos() {
        return aprovada() || violada == Regra.R6;
    }

    public Qualidade qualidade() {
        return qualidade;
    }
}
