package matchmaking.busca;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;

import matchmaking.avaliacao.Avaliacao;
import matchmaking.avaliacao.AvaliadorDePartida;
import matchmaking.avaliacao.PartidaAvaliada;
import matchmaking.dominio.Equipe;
import matchmaking.dominio.Fila;
import matchmaking.dominio.Jogador;
import matchmaking.dominio.Partida;
import matchmaking.dominio.Posicao;
import matchmaking.saida.NenhumaPartida;
import matchmaking.saida.PartidaFormada;
import matchmaking.saida.Resultado;

/**
 * Procura a partida otima pedindo ao avaliador o que e valido
 * e a propria partida avaliada o que e melhor.
 * A fila nao e reordenada nem alterada.
 */
public final class FormadorDePartida {
    public Resultado formar(Fila fila) {
        if (fila.tamanho() < Partida.jogadoresPorPartida()) {
            return new NenhumaPartida("JOGADORES_INSUFICIENTES", fila.ids());
        }
        AvaliadorDePartida avaliador = AvaliadorDePartida.doContrato(fila);
        Exploracao exploracao = new Exploracao();
        gerarSubconjuntos(fila.jogadores(), Partida.jogadoresPorPartida(), avaliador, exploracao);
        if (!exploracao.achou()) {
            return new NenhumaPartida(exploracao.motivo(), fila.ids());
        }
        PartidaAvaliada melhor = exploracao.melhor();
        return new PartidaFormada(melhor, fila.idsExceto(melhor.partida().jogadores()));
    }

    private void gerarSubconjuntos(List<Jogador> jogadores, int escolher,
            AvaliadorDePartida avaliador, Exploracao exploracao) {
        gerar(jogadores, escolher, 0, new ArrayList<Jogador>(), avaliador, exploracao);
    }

    private void gerar(List<Jogador> jogadores, int escolher, int inicio, List<Jogador> atual,
            AvaliadorDePartida avaliador, Exploracao exploracao) {
        if (atual.size() == escolher) {
            atribuirPosicoes(new ArrayList<Jogador>(atual), avaliador, exploracao);
            return;
        }
        for (int i = inicio; i <= jogadores.size() - (escolher - atual.size()); i++) {
            atual.add(jogadores.get(i));
            gerar(jogadores, escolher, i + 1, atual, avaliador, exploracao);
            atual.remove(atual.size() - 1);
        }
    }

    private void atribuirPosicoes(List<Jogador> subconjunto, AvaliadorDePartida avaliador, Exploracao exploracao) {
        subconjunto.sort(Comparator.comparingInt(jogador -> jogador.posicoesAceitas().size()));
        Posicao[] atribuida = new Posicao[subconjunto.size()];
        int[] contagem = new int[Posicao.porEquipe()];
        atribuir(subconjunto, 0, atribuida, contagem, avaliador, exploracao);
    }

    private void atribuir(List<Jogador> ordem, int passo, Posicao[] atribuida, int[] contagem,
            AvaliadorDePartida avaliador, Exploracao exploracao) {
        if (passo == ordem.size()) {
            exploracao.registrarAtribuicao();
            repartirEntreEquipes(ordem, atribuida, avaliador, exploracao);
            return;
        }
        Jogador jogador = ordem.get(passo);
        for (Posicao posicao : jogador.posicoesAceitas()) {
            int indice = posicao.ordinal();
            if (contagem[indice] >= Partida.quantidadeDeEquipes()) {
                continue;
            }
            contagem[indice]++;
            atribuida[passo] = posicao;
            atribuir(ordem, passo + 1, atribuida, contagem, avaliador, exploracao);
            atribuida[passo] = null;
            contagem[indice]--;
        }
    }

    private void repartirEntreEquipes(List<Jogador> ordem, Posicao[] atribuida,
            AvaliadorDePartida avaliador, Exploracao exploracao) {
        EnumMap<Posicao, List<Jogador>> porPosicao = new EnumMap<Posicao, List<Jogador>>(Posicao.class);
        for (Posicao posicao : Posicao.values()) {
            porPosicao.put(posicao, new ArrayList<Jogador>());
        }
        for (int i = 0; i < ordem.size(); i++) {
            porPosicao.get(atribuida[i]).add(ordem.get(i));
        }
        for (Posicao posicao : Posicao.values()) {
            if (porPosicao.get(posicao).size() != Partida.quantidadeDeEquipes()) {
                return;
            }
        }
        int combinacoes = 1 << Posicao.porEquipe();
        Posicao[] posicoes = Posicao.values();
        for (int mascara = 0; mascara < combinacoes; mascara++) {
            Map<Posicao, Jogador> ladoA = new EnumMap<Posicao, Jogador>(Posicao.class);
            Map<Posicao, Jogador> ladoB = new EnumMap<Posicao, Jogador>(Posicao.class);
            for (int p = 0; p < posicoes.length; p++) {
                List<Jogador> par = porPosicao.get(posicoes[p]);
                boolean primeiroNaA = (mascara & (1 << p)) != 0;
                if (primeiroNaA) {
                    ladoA.put(posicoes[p], par.get(0));
                    ladoB.put(posicoes[p], par.get(1));
                } else {
                    ladoA.put(posicoes[p], par.get(1));
                    ladoB.put(posicoes[p], par.get(0));
                }
            }
            Partida partida = Partida.rotulada(Equipe.montar(ladoA), Equipe.montar(ladoB));
            Avaliacao avaliacao = avaliador.avaliar(partida);
            if (avaliacao.passouDosGrupos()) {
                exploracao.registrarGrupos();
            }
            if (avaliacao.aprovada()) {
                exploracao.considerar(new PartidaAvaliada(partida, avaliacao.qualidade()));
            }
        }
    }

    /** Estado da busca. So o formador enxerga o campeao enquanto explora. */
    private static final class Exploracao {
        private PartidaAvaliada melhor;
        private boolean viuAtribuicao;
        private boolean viuGrupos;

        void registrarAtribuicao() {
            viuAtribuicao = true;
        }

        void registrarGrupos() {
            viuGrupos = true;
        }

        void considerar(PartidaAvaliada candidata) {
            if (melhor == null || candidata.supera(melhor)) {
                melhor = candidata;
            }
        }

        boolean achou() {
            return melhor != null;
        }

        PartidaAvaliada melhor() {
            return melhor;
        }

        String motivo() {
            if (viuGrupos) {
                return "DIFERENCA_HABILIDADE";
            }
            if (viuAtribuicao) {
                return "GRUPO_INCOMPATIVEL";
            }
            return "COMPOSICAO_IMPOSSIVEL";
        }
    }
}
