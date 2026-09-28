package matchmaking;

import matchmaking.busca.FormadorDePartida;
import matchmaking.entrada.LeitorDeFila;
import matchmaking.entrada.Leitura;
import matchmaking.saida.Resultado;

/** Encadeia a leitura da fila e a formacao da partida. */
public final class Aplicacao {
    private final LeitorDeFila leitor;
    private final FormadorDePartida formador;

    public Aplicacao() {
        this(new LeitorDeFila(), new FormadorDePartida());
    }

    public Aplicacao(LeitorDeFila leitor, FormadorDePartida formador) {
        this.leitor = leitor;
        this.formador = formador;
    }

    public Resultado executar(String texto) {
        Leitura leitura = leitor.ler(texto);
        if (!leitura.valida()) {
            return leitura.erro();
        }
        return formador.formar(leitura.fila());
    }
}
