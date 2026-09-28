package matchmaking.avaliacao;

/**
 * Ate onde a espera dos selecionados amplia a diferenca de habilidade.
 * Nao altera Q. So decide se a formacao ainda cabe em R6.
 */
public final class PoliticaDeTolerancia {
    public static final PoliticaDeTolerancia DO_CONTRATO = new PoliticaDeTolerancia(200, 4, 400);

    private final int limiteBase;
    private final int divisorBonus;
    private final int limiteAbsoluto;

    public PoliticaDeTolerancia(int limiteBase, int divisorBonus, int limiteAbsoluto) {
        this.limiteBase = limiteBase;
        this.divisorBonus = divisorBonus;
        this.limiteAbsoluto = limiteAbsoluto;
    }

    public int limite(int esperaMaxima) {
        int ampliado = limiteBase + esperaMaxima / divisorBonus;
        return Math.min(limiteAbsoluto, ampliado);
    }

    public boolean aceita(Qualidade qualidade) {
        return qualidade.diferencaHabilidade() <= limite(qualidade.esperaMax());
    }
}
