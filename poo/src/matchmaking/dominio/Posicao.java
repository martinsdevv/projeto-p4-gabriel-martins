package matchmaking.dominio;

import java.util.Optional;

/**
 * Conjunto fechado de funcoes da composicao tradicional.
 * A ordem de declaracao e a ordem canonica da saida.
 */
public enum Posicao {
    TOP("Top"),
    JUNGLE("Jungle"),
    MID("Mid"),
    ADC("ADC"),
    SUPPORT("Support");

    private final String nome;

    Posicao(String nome) {
        this.nome = nome;
    }

    public String nome() {
        return nome;
    }

    public static int porEquipe() {
        return values().length;
    }

    public static Optional<Posicao> reconhecer(String nome) {
        for (Posicao posicao : values()) {
            if (posicao.nome.equals(nome)) {
                return Optional.of(posicao);
            }
        }
        return Optional.empty();
    }
}
