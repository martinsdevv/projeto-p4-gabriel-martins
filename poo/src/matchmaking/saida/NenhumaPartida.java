package matchmaking.saida;

import java.util.List;

public final class NenhumaPartida implements Resultado {
    private final String motivo;
    private final List<String> filaRestante;

    public NenhumaPartida(String motivo, List<String> filaRestante) {
        this.motivo = motivo;
        this.filaRestante = List.copyOf(filaRestante);
    }

    @Override
    public String texto() {
        return "status: NENHUMA_PARTIDA\n"
                + "motivo: " + motivo + "\n"
                + "fila_restante: " + TextoDaFila.ids(filaRestante) + "\n";
    }

    @Override
    public int codigoDeSaida() {
        return 0;
    }

    @Override
    public String status() {
        return "NENHUMA_PARTIDA";
    }

    @Override
    public String detalhe() {
        return motivo;
    }

    @Override
    public String textoEspera() {
        return "-";
    }
}
