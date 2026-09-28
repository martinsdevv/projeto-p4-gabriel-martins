package matchmaking.avaliacao;

import java.util.List;

import matchmaking.dominio.Partida;

/**
 * Criterio que entra em Q.
 * A penalidade e sempre peso vezes a medida; cada subclasse so diz o que medir.
 * O tempo de espera nao e subclasse: ele nao entra em Q.
 */
public abstract class CriterioQualidade {
    private final String nome;
    private final int peso;

    protected CriterioQualidade(String nome, int peso) {
        if (peso < 0) {
            throw new IllegalArgumentException("peso negativo");
        }
        this.nome = nome;
        this.peso = peso;
    }

    public final int penalidade(Partida partida) {
        return peso * medir(partida);
    }

    protected abstract int medir(Partida partida);

    public final String nome() {
        return nome;
    }

    public final int peso() {
        return peso;
    }

    public static List<CriterioQualidade> doContrato() {
        return List.of(
                new EquilibrioDeHabilidade(),
                new PosicoesForaDaPreferida(),
                new DesequilibrioDeAlternativas(),
                new AmplitudeDeLatencia());
    }

    /** C1. A medida e a diferenca das somas, nao a media. */
    public static final class EquilibrioDeHabilidade extends CriterioQualidade {
        public EquilibrioDeHabilidade() {
            super("equilibrio de habilidade", 2);
        }

        @Override
        protected int medir(Partida partida) {
            return partida.diferencaDeSomas();
        }
    }

    /** C2. Cada jogador fora da preferida. */
    public static final class PosicoesForaDaPreferida extends CriterioQualidade {
        public PosicoesForaDaPreferida() {
            super("posicoes fora da preferida", 100);
        }

        @Override
        protected int medir(Partida partida) {
            return partida.foraDaPreferida();
        }
    }

    /** C3. Modulo da diferenca de deslocados entre as equipes. */
    public static final class DesequilibrioDeAlternativas extends CriterioQualidade {
        public DesequilibrioDeAlternativas() {
            super("desequilibrio de alternativas", 50);
        }

        @Override
        protected int medir(Partida partida) {
            return partida.desequilibrioDeAlternativas();
        }
    }

    /** C4. Amplitude de latencia, em milissegundos. */
    public static final class AmplitudeDeLatencia extends CriterioQualidade {
        public AmplitudeDeLatencia() {
            super("amplitude de latencia", 1);
        }

        @Override
        protected int medir(Partida partida) {
            return partida.amplitudeDeLatencia();
        }
    }
}
