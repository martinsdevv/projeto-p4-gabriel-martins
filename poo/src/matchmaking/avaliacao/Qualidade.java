package matchmaking.avaliacao;

import java.util.List;

import matchmaking.dominio.Jogador;
import matchmaking.dominio.Partida;
import matchmaking.dominio.Posicao;

/** Valor de uma partida valida, mais os campos que desempatam a escolha. */
public final class Qualidade {
    private static final int BASE = 10000;

    private final int diferencaHabilidade;
    private final int foraPreferencia;
    private final int desequilibrioAlternativas;
    private final int amplitudeLatencia;
    private final int valor;
    private final int esperaMin;
    private final int esperaMax;
    private final int esperaSoma;

    private Qualidade(int diferencaHabilidade, int foraPreferencia, int desequilibrioAlternativas,
            int amplitudeLatencia, int valor, int esperaMin, int esperaMax, int esperaSoma) {
        this.diferencaHabilidade = diferencaHabilidade;
        this.foraPreferencia = foraPreferencia;
        this.desequilibrioAlternativas = desequilibrioAlternativas;
        this.amplitudeLatencia = amplitudeLatencia;
        this.valor = valor;
        this.esperaMin = esperaMin;
        this.esperaMax = esperaMax;
        this.esperaSoma = esperaSoma;
    }

    public static Qualidade de(Partida partida, List<CriterioQualidade> criterios) {
        int penalidade = 0;
        for (CriterioQualidade criterio : criterios) {
            penalidade += criterio.penalidade(partida);
        }
        int esperaMin = Integer.MAX_VALUE;
        int esperaMax = Integer.MIN_VALUE;
        int esperaSoma = 0;
        for (Jogador jogador : partida.jogadores()) {
            int espera = jogador.esperaSegundos();
            esperaMin = Math.min(esperaMin, espera);
            esperaMax = Math.max(esperaMax, espera);
            esperaSoma += espera;
        }
        int diferenca = partida.diferencaDeSomas() / Posicao.porEquipe();
        return new Qualidade(
                diferenca,
                partida.foraDaPreferida(),
                partida.desequilibrioDeAlternativas(),
                partida.amplitudeDeLatencia(),
                BASE - penalidade,
                esperaMin,
                esperaMax,
                esperaSoma);
    }

    /**
     * Positivo quando esta qualidade e preferivel a outra nos criterios 1 a 6.
     * Zero quando o desempate fica para os ids e para a tupla canonica.
     */
    public int preferenciaSobre(Qualidade outra) {
        if (valor != outra.valor) {
            return Integer.compare(valor, outra.valor);
        }
        if (diferencaHabilidade != outra.diferencaHabilidade) {
            return Integer.compare(outra.diferencaHabilidade, diferencaHabilidade);
        }
        if (foraPreferencia != outra.foraPreferencia) {
            return Integer.compare(outra.foraPreferencia, foraPreferencia);
        }
        if (desequilibrioAlternativas != outra.desequilibrioAlternativas) {
            return Integer.compare(outra.desequilibrioAlternativas, desequilibrioAlternativas);
        }
        if (amplitudeLatencia != outra.amplitudeLatencia) {
            return Integer.compare(outra.amplitudeLatencia, amplitudeLatencia);
        }
        return Integer.compare(esperaSoma, outra.esperaSoma);
    }

    public int diferencaHabilidade() {
        return diferencaHabilidade;
    }

    public int foraPreferencia() {
        return foraPreferencia;
    }

    public int desequilibrioAlternativas() {
        return desequilibrioAlternativas;
    }

    public int amplitudeLatencia() {
        return amplitudeLatencia;
    }

    public int valor() {
        return valor;
    }

    public int esperaMin() {
        return esperaMin;
    }

    public int esperaMax() {
        return esperaMax;
    }

    public int esperaSoma() {
        return esperaSoma;
    }
}
