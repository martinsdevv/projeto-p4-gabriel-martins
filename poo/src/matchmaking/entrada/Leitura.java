package matchmaking.entrada;

import matchmaking.dominio.Fila;
import matchmaking.saida.ErroEntrada;

/** Ou uma fila pronta, ou o erro que impediu de monta-la. */
public final class Leitura {
    private final Fila fila;
    private final ErroEntrada erro;

    private Leitura(Fila fila, ErroEntrada erro) {
        this.fila = fila;
        this.erro = erro;
    }

    public static Leitura ok(Fila fila) {
        return new Leitura(fila, null);
    }

    public static Leitura rejeitada(ErroEntrada erro) {
        return new Leitura(null, erro);
    }

    public boolean valida() {
        return erro == null;
    }

    public Fila fila() {
        return fila;
    }

    public ErroEntrada erro() {
        return erro;
    }
}
