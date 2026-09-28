package matchmaking.dominio;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/** Duas equipes. A rotulacao canonica coloca na A o menor id dos dez. */
public final class Partida {
    private final Equipe equipeA;
    private final Equipe equipeB;

    private Partida(Equipe equipeA, Equipe equipeB) {
        this.equipeA = equipeA;
        this.equipeB = equipeB;
    }

    public static int quantidadeDeEquipes() {
        return 2;
    }

    public static int jogadoresPorPartida() {
        return quantidadeDeEquipes() * Posicao.porEquipe();
    }

    public static Partida rotulada(Equipe uma, Equipe outra) {
        if (outra.menorId().compareTo(uma.menorId()) < 0) {
            return new Partida(outra, uma);
        }
        return new Partida(uma, outra);
    }

    public Equipe equipeA() {
        return equipeA;
    }

    public Equipe equipeB() {
        return equipeB;
    }

    public List<Jogador> jogadores() {
        List<Jogador> todos = new ArrayList<Jogador>(jogadoresPorPartida());
        todos.addAll(equipeA.jogadores());
        todos.addAll(equipeB.jogadores());
        return todos;
    }

    public int diferencaDeSomas() {
        return Math.abs(equipeA.somaHabilidade() - equipeB.somaHabilidade());
    }

    public int foraDaPreferida() {
        return equipeA.deslocados() + equipeB.deslocados();
    }

    public int desequilibrioDeAlternativas() {
        return Math.abs(equipeA.deslocados() - equipeB.deslocados());
    }

    public int amplitudeDeLatencia() {
        int minima = Integer.MAX_VALUE;
        int maxima = Integer.MIN_VALUE;
        for (Jogador jogador : jogadores()) {
            minima = Math.min(minima, jogador.latenciaMs());
            maxima = Math.max(maxima, jogador.latenciaMs());
        }
        return maxima - minima;
    }

    public List<String> idsOrdenados() {
        List<String> ids = new ArrayList<String>();
        for (Jogador jogador : jogadores()) {
            ids.add(jogador.id());
        }
        Collections.sort(ids);
        return ids;
    }

    public List<String> tuplaCanonica() {
        List<String> tupla = new ArrayList<String>(jogadoresPorPartida());
        for (Posicao posicao : Posicao.values()) {
            tupla.add(equipeA.jogadorEm(posicao).id());
        }
        for (Posicao posicao : Posicao.values()) {
            tupla.add(equipeB.jogadorEm(posicao).id());
        }
        return tupla;
    }
}
