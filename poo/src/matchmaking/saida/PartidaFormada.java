package matchmaking.saida;

import java.util.List;
import java.util.Locale;

import matchmaking.avaliacao.PartidaAvaliada;
import matchmaking.avaliacao.Qualidade;
import matchmaking.dominio.Partida;

public final class PartidaFormada implements Resultado {
    private final Partida partida;
    private final Qualidade qualidade;
    private final List<String> filaRestante;

    public PartidaFormada(PartidaAvaliada avaliada, List<String> filaRestante) {
        this.partida = avaliada.partida();
        this.qualidade = avaliada.qualidade();
        this.filaRestante = List.copyOf(filaRestante);
    }

    @Override
    public String texto() {
        return "status: PARTIDA_FORMADA\n"
                + "equipe_a: " + partida.equipeA().texto() + "\n"
                + "equipe_b: " + partida.equipeB().texto() + "\n"
                + "diferenca_habilidade: " + qualidade.diferencaHabilidade() + "\n"
                + "jogadores_fora_preferencia: " + qualidade.foraPreferencia() + "\n"
                + "desequilibrio_alternativas: " + qualidade.desequilibrioAlternativas() + "\n"
                + "qualidade: " + qualidade.valor() + "\n"
                + "espera: " + textoEspera() + "\n"
                + "amplitude_latencia: " + qualidade.amplitudeLatencia() + "\n"
                + "fila_restante: " + TextoDaFila.ids(filaRestante) + "\n";
    }

    @Override
    public int codigoDeSaida() {
        return 0;
    }

    @Override
    public String status() {
        return "PARTIDA_FORMADA";
    }

    @Override
    public String detalhe() {
        return "qualidade " + qualidade.valor();
    }

    @Override
    public String textoEspera() {
        double media = qualidade.esperaSoma() / (double) Partida.jogadoresPorPartida();
        return String.format(Locale.US, "{min: %d, max: %d, media: %.1f}",
                qualidade.esperaMin(), qualidade.esperaMax(), media);
    }
}
