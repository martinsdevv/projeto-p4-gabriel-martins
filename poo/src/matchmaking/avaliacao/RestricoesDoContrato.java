package matchmaking.avaliacao;

import java.util.HashSet;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import matchmaking.dominio.Equipe;
import matchmaking.dominio.Fila;
import matchmaking.dominio.Grupo;
import matchmaking.dominio.Jogador;
import matchmaking.dominio.Partida;
import matchmaking.dominio.Posicao;
import matchmaking.dominio.Vaga;

/** R1 a R5, na ordem em que o avaliador as consulta. */
public final class RestricoesDoContrato {
    private RestricoesDoContrato() {
    }

    public static List<Restricao> estruturais() {
        return List.of(
                new JogadoresDistintos(),
                new DuasEquipes(),
                new PosicaoAceita(),
                new ComposicaoCompleta(),
                new GruposNaMesmaEquipe());
    }

    static final class JogadoresDistintos implements Restricao {
        @Override
        public Regra regra() {
            return Regra.R1;
        }

        @Override
        public boolean satisfeita(Partida partida, Fila fila) {
            Set<String> vistos = new HashSet<String>();
            for (Jogador jogador : partida.jogadores()) {
                if (!fila.contem(jogador) || !vistos.add(jogador.id())) {
                    return false;
                }
            }
            return vistos.size() == Partida.jogadoresPorPartida();
        }
    }

    static final class DuasEquipes implements Restricao {
        @Override
        public Regra regra() {
            return Regra.R2;
        }

        @Override
        public boolean satisfeita(Partida partida, Fila fila) {
            return partida.equipeA().tamanho() == Posicao.porEquipe()
                    && partida.equipeB().tamanho() == Posicao.porEquipe();
        }
    }

    static final class PosicaoAceita implements Restricao {
        @Override
        public Regra regra() {
            return Regra.R3;
        }

        @Override
        public boolean satisfeita(Partida partida, Fila fila) {
            return aceita(partida.equipeA()) && aceita(partida.equipeB());
        }

        private boolean aceita(Equipe equipe) {
            for (Vaga vaga : equipe.vagasNaOrdem()) {
                if (!vaga.jogador().aceita(vaga.posicao())) {
                    return false;
                }
            }
            return true;
        }
    }

    static final class ComposicaoCompleta implements Restricao {
        @Override
        public Regra regra() {
            return Regra.R4;
        }

        @Override
        public boolean satisfeita(Partida partida, Fila fila) {
            return partida.equipeA().composicaoCompleta() && partida.equipeB().composicaoCompleta();
        }
    }

    static final class GruposNaMesmaEquipe implements Restricao {
        @Override
        public Regra regra() {
            return Regra.R5;
        }

        @Override
        public boolean satisfeita(Partida partida, Fila fila) {
            Map<Grupo, Equipe> lado = new IdentityHashMap<Grupo, Equipe>();
            return marcar(partida.equipeA(), lado) && marcar(partida.equipeB(), lado);
        }

        private boolean marcar(Equipe equipe, Map<Grupo, Equipe> lado) {
            for (Vaga vaga : equipe.vagasNaOrdem()) {
                if (vaga.jogador().grupo().isEmpty()) {
                    continue;
                }
                Grupo grupo = vaga.jogador().grupo().get();
                Equipe anterior = lado.get(grupo);
                if (anterior != null && anterior != equipe) {
                    return false;
                }
                lado.put(grupo, equipe);
            }
            return true;
        }
    }
}
